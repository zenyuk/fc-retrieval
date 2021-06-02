package provider_admin

/*
 * Copyright 2021 ConsenSys Software Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ConsenSys/fc-retrieval-client/pkg/fcrclient"
	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/register"
	"github.com/ConsenSys/fc-retrieval-gateway-admin/pkg/fcrgatewayadmin"
	"github.com/ConsenSys/fc-retrieval-provider-admin/pkg/fcrprovideradmin"

	"github.com/ConsenSys/fc-retrieval-itest/config"
	"github.com/ConsenSys/fc-retrieval-itest/pkg/util"
)

// Test the Provider Admin API.
var providerTest_providerConfig = config.NewConfig(".env.provider")
var gatewayConfig_gatewayConfig = config.NewConfig(".env.gateway")

func TestMain(m *testing.M) {
	// Need to make sure this env is not set in host machine
	itestEnv := os.Getenv("ITEST_CALLING_FROM_CONTAINER")

	if itestEnv != "" {
		// Env is set, we are calling from docker container
		m.Run()
		return
	}
	// Env is not set, we are calling from host
	// We need a redis, a register, a gateway and a provider
	tag := util.GetCurrentBranch()

	// Get env
	rgEnv := util.GetEnvMap("../../.env.register")
	gwEnv := util.GetEnvMap("../../.env.gateway")
	pvEnv := util.GetEnvMap("../../.env.provider")

	// Create shared net
	ctx := context.Background()
	network, networkName := util.CreateNetwork(ctx)
	defer (*network).Remove(ctx)

	// Start redis
	redisContainer := util.StartRedis(ctx, networkName, true)
	defer redisContainer.Terminate(ctx)

	// Start register
	registerContainer := util.StartRegister(ctx, tag, networkName, util.ColorYellow, rgEnv, true)
	defer registerContainer.Terminate(ctx)

	// Start gateway
	gatewayContainer := util.StartGateway(ctx, "gateway", tag, networkName, util.ColorBlue, gwEnv, true)
	defer gatewayContainer.Terminate(ctx)

	// Start provider
	providerContainer := util.StartProvider(ctx, "provider", tag, networkName, util.ColorPurple, pvEnv, true)
	defer providerContainer.Terminate(ctx)

	// Start itest
	done := make(chan bool)
	itestContainer := util.StartItest(ctx, tag, networkName, util.ColorGreen, "", "", done, true)
	defer itestContainer.Terminate(ctx)

	// Block until done.
	if <-done {
		logging.Info("Tests passed, shutdown...")
	} else {
		logging.Fatal("Tests failed, shutdown...")
	}
}

func TestGetProviderAdminVersion(t *testing.T) {
	versionInfo := fcrclient.GetVersion()
	// Verify that the client version is an integer number.
	_, err := strconv.Atoi(versionInfo.Version)
	if err != nil {
		panic(err)
	}

	// The version must be 1 or more.
	assert.LessOrEqual(t, 0, 0)
}

func TestInitProviderAdminNoRetrievalKey(t *testing.T) {
	logging.Info("/*******************************************************/")
	logging.Info("/*      Start TestInitProviderAdminNoRetrievalKey	     */")
	logging.Info("/*******************************************************/")
	logging.Error("Wait two seconds for the provider to deploy and be ready for requests")

	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}

	// Gateway
	gwConfBuilder := fcrgatewayadmin.CreateSettings()
	gwConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	gwConfBuilder.SetRegisterURL(providerTest_providerConfig.GetString("REGISTER_API_URL"))
	gwConf := gwConfBuilder.Build()

	gwAdmin := fcrgatewayadmin.NewFilecoinRetrievalGatewayAdmin(*gwConf)

	gatewayRootKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	logging.Info("gatewayRootKey: %v", gatewayRootKey)
	if err != nil {
		panic(err)
	}
	gatewayRootSigningKey, err := gatewayRootKey.EncodePublicKey()
	logging.Info("gatewayRootSigningKey: %s", gatewayRootSigningKey)
	if err != nil {
		panic(err)
	}
	gatewayRetrievalPrivateKey, err := fcrcrypto.GenerateRetrievalV1KeyPair()
	logging.Info("gatewayRetrievalPrivateKey: %v", gatewayRetrievalPrivateKey)
	if err != nil {
		panic(err)
	}
	gatewayRetrievalSigningKey, err := gatewayRetrievalPrivateKey.EncodePublicKey()
	logging.Info("gatewayRetrievalSigningKey: %s", gatewayRetrievalSigningKey)
	if err != nil {
		panic(err)
	}

	// gatewayRootSigningKey := "0104d799bc7141b058b4c9d819ba8d8fa1e87b2ee9132f5b59d3a91edcd72c08cd64d2fd44f99f8d4a0159a65a0c8c0409f646712793ab4fb7b6151654b6e00ca69f"
	// gatewayRetrievalSigningKey := "01041ee440cab4f5e92803e29de7079d317a332b206b21df612fe0d1c34b585df4f44180aa9a75e4c95116ac341256333d7356d42704be43efd8828293ef013d9139"
	// gatewayID, err := nodeid.NewRandomNodeID()
	gatewayID, err := nodeid.NewNodeIDFromHexString("ebc134a429ba7dc4811bf64ccb67057f5bd57ca4676800e2f71731cbcc5eb518")
	gatewayRegister := &register.GatewayRegister{
		NodeID:              gatewayID.ToString(),
		Address:             gatewayConfig_gatewayConfig.GetString("GATEWAY_ADDRESS"),
		RootSigningKey:      gatewayRootSigningKey,
		SigningKey:          gatewayRetrievalSigningKey,
		RegionCode:          gatewayConfig_gatewayConfig.GetString("GATEWAY_REGION_CODE"),
		NetworkInfoGateway:  gatewayConfig_gatewayConfig.GetString("NETWORK_INFO_GATEWAY"),
		NetworkInfoProvider: gatewayConfig_gatewayConfig.GetString("NETWORK_INFO_PROVIDER"),
		NetworkInfoClient:   gatewayConfig_gatewayConfig.GetString("NETWORK_INFO_CLIENT"),
		NetworkInfoAdmin:    gatewayConfig_gatewayConfig.GetString("NETWORK_INFO_ADMIN"),
	}

	err = gwAdmin.InitialiseGateway(gatewayRegister, gatewayRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1))
	if err != nil {
		panic(err)
	}

	// Provider
	confBuilder := fcrprovideradmin.CreateSettings()
	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	confBuilder.SetRegisterURL(providerTest_providerConfig.GetString("REGISTER_API_URL"))
	conf := confBuilder.Build()

	pvadmin := fcrprovideradmin.NewFilecoinRetrievalProviderAdmin(*conf)

	providerRootKey, err := fcrcrypto.GenerateRetrievalV1KeyPair()
	logging.Info("providerRootKey: %v", providerRootKey)
	if err != nil {
		panic(err)
	}
	providerRootSigningKey, err := providerRootKey.EncodePublicKey()
	logging.Info("providerRootSigningKey: %s", providerRootSigningKey)
	if err != nil {
		panic(err)
	}

	providerPrivKey, err := fcrcrypto.GenerateRetrievalV1KeyPair()
	logging.Info("providerPrivKey: %v", providerPrivKey)
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}
	providerSigningKey, err := providerPrivKey.EncodePublicKey()
	logging.Info("providerSigningKey: %s", providerSigningKey)
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}
	providerID, err := nodeid.NewNodeIDFromHexString("ebc134a429ba7dc4811bf64ccb67057f5bd57ca4676800e2f71731cbcc5eb518")
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}

	providerRegister := &register.ProviderRegister{
		NodeID:             providerID.ToString(),
		Address:            providerTest_providerConfig.GetString("PROVIDER_ADDRESS"),
		RootSigningKey:     gatewayRootSigningKey,
		SigningKey:         gatewayRetrievalSigningKey,
		RegionCode:         providerTest_providerConfig.GetString("PROVIDER_REGION_CODE"),
		NetworkInfoGateway: providerTest_providerConfig.GetString("NETWORK_INFO_GATEWAY"),
		NetworkInfoClient:  providerTest_providerConfig.GetString("NETWORK_INFO_CLIENT"),
		NetworkInfoAdmin:   providerTest_providerConfig.GetString("NETWORK_INFO_ADMIN"),
	}

	// Initialise provider
	err = pvadmin.InitialiseProvider(providerRegister, gatewayRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1))
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}

	// Generate random cid offer
	contentID1 := cid.NewRandomContentID()
	contentID2 := cid.NewRandomContentID()
	pieceCIDs := []cid.ContentID{*contentID1, *contentID2}
	expiryDate := time.Now().Local().Add(time.Hour * time.Duration(24)).Unix()

	// Force update
	err = pvadmin.ForceUpdate(providerID)
	if err != nil {
		panic(err)
	}
	err = gwAdmin.ForceUpdate(gatewayID)
	if err != nil {
		panic(err)
	}

	// Publish Group CID
	err = pvadmin.PublishGroupCID(providerID, pieceCIDs, 42, expiryDate, 42)
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}

	// Get all offers
	var gatewayIDs []nodeid.NodeID
	gatewayIDs = make([]nodeid.NodeID, 0)
	logging.Info("Get all offers")
	_, cidgroupInfo, err := pvadmin.GetGroupCIDOffer(providerID, gatewayIDs)
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}
	logging.Info("Get all offers: %d", len(cidgroupInfo))
	assert.GreaterOrEqual(t, len(cidgroupInfo), 1, "Get all offers should be found")

	// Get offers by gatewayIDs real
	gateways, err := register.GetRegisteredGateways(providerTest_providerConfig.GetString("REGISTER_API_URL"))
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}
	logging.Info("Registered gateways: %+v", gateways)
	realNodeID, err := nodeid.NewNodeIDFromHexString("ebc134a429ba7dc4811bf64ccb67057f5bd57ca4676800e2f71731cbcc5eb518")
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}
	gatewayIDs = append(gatewayIDs, *realNodeID) // Add a gateway
	logging.Info("Get offers by real gatewayID=%s", realNodeID.ToString())
	_, cidgroupInfo, err = pvadmin.GetGroupCIDOffer(providerID, gatewayIDs)
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}
	logging.Info("Get offers by real gatewayID=%s: %d", realNodeID.ToString(), len(cidgroupInfo))
	assert.GreaterOrEqual(t, len(cidgroupInfo), 1, "Get offers by gatewayIDs real should be found")

	// Get offers by gatewayIDs fake
	fakeNodeID, _ := nodeid.NewNodeIDFromHexString("101112131415161718191A1B1C1D1E1F202122232425262728292A2B2C2DFA43")
	gatewayIDs[0] = *fakeNodeID
	logging.Info("Get offers by fake gatewayID=%s", fakeNodeID.ToString())
	_, cidgroupInfo, err = pvadmin.GetGroupCIDOffer(providerID, gatewayIDs)
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}
	logging.Info("Get offers by fake gatewayID=%s: %d", fakeNodeID.ToString(), len(cidgroupInfo))
	assert.Equal(t, 0, len(cidgroupInfo), "Get offers by gatewayIDs fake should be empty")

	// The version must be 1 or more.
	assert.LessOrEqual(t, 1, 1)

	logging.Info("/*******************************************************/")
	logging.Info("/*      End TestInitProviderAdminNoRetrievalKey	       */")
	logging.Info("/*******************************************************/")
}
