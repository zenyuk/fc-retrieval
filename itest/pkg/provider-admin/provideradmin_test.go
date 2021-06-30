/*
Package provider_admin - set of integration tests, specific to Retrieval Providers
*/
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
	"errors"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ConsenSys/fc-retrieval/client/pkg/fcrclient"
	"github.com/ConsenSys/fc-retrieval/common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrregistermgr"
	"github.com/ConsenSys/fc-retrieval/common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval/common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval/common/pkg/register"
	"github.com/ConsenSys/fc-retrieval/gateway-admin/pkg/fcrgatewayadmin"
	"github.com/ConsenSys/fc-retrieval/itest/config"
	"github.com/ConsenSys/fc-retrieval/itest/pkg/util"
	"github.com/ConsenSys/fc-retrieval/provider-admin/pkg/fcrprovideradmin"
)

// Test the Provider Admin API.
var providerTestProviderConfig = config.NewConfig(".env.provider")
var gatewayConfigGatewayConfig = config.NewConfig(".env.gateway")

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

	// Get env
	rgEnv := util.GetEnvMap("../../.env.register")
	gwEnv := util.GetEnvMap("../../.env.gateway")
	pvEnv := util.GetEnvMap("../../.env.provider")

	// Create shared net
	ctx := context.Background()
	network, networkName := util.CreateNetwork(ctx)

	// Start redis
	redisContainer := util.StartRedis(ctx, networkName, true)

	// Start register
	registerContainer := util.StartRegister(ctx, networkName, util.ColorYellow, rgEnv, true)

	// Start gateway
	gatewayContainer := util.StartGateway(ctx, "gateway", networkName, util.ColorBlue, gwEnv, true)

	// Start provider
	providerContainer := util.StartProvider(ctx, "provider", networkName, util.ColorPurple, pvEnv, true)

	// Start itest
	done := make(chan bool)
	itestContainer := util.StartItest(ctx, networkName, util.ColorGreen, "", "", done, true, "")

	// Block until done.
	if <-done {
		logging.Info("Tests passed, shutdown...")
	} else {
		logging.Error("Tests failed, shutdown...")
	}

	if err := itestContainer.Terminate(ctx); err != nil {
		logging.Error("error while terminating test container: %s", err.Error())
	}
	if err := providerContainer.Terminate(ctx); err != nil {
		logging.Error("error while terminating test container: %s", err.Error())
	}
	if err := gatewayContainer.Terminate(ctx); err != nil {
		logging.Error("error while terminating test container: %s", err.Error())
	}
	if err := registerContainer.Terminate(ctx); err != nil {
		logging.Error("error while terminating test container: %s", err.Error())
	}
	if err := redisContainer.Terminate(ctx); err != nil {
		logging.Error("error while terminating test container: %s", err.Error())
	}
	if err := (*network).Remove(ctx); err != nil {
		logging.Error("error while terminating test container network: %s", err.Error())
	}
}

func TestGetProviderAdminVersion(t *testing.T) {
	versionInfo := fcrclient.GetVersion()
	// Verify that the client version is an integer number.
	actualVersion, err := strconv.Atoi(versionInfo.Version)
	if err != nil {
		panic(err)
	}

	assert.GreaterOrEqualf(t, actualVersion, 1, "version must be 1 or more")
}

func TestInitProviderAdminNoRetrievalKey(t *testing.T) {
	logging.Info("/*******************************************************/")
	logging.Info("/*      Start TestInitProviderAdminNoRetrievalKey	     */")
	logging.Info("/*******************************************************/")
	logging.Error("Wait two seconds for the provider to deploy and be ready for requests")

	// Main key used across the chain
	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		logging.Error("error generating blockchain key: %s", err.Error())
		t.FailNow()
	}

	// Create and start register manager
	var rm = fcrregistermgr.NewFCRRegisterMgr(providerTestProviderConfig.GetString("REGISTER_API_URL"), true, true, 2*time.Second)
	if err := rm.Start(); err != nil {
		logging.Error("error starting Register Manager: %s", err.Error())
		t.FailNow()
	}
	defer rm.ShutdownAndWait()

	// Configure gateway admin
	gConfBuilder := fcrgatewayadmin.CreateSettings()
	gConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	gConfBuilder.SetRegisterURL(providerTestProviderConfig.GetString("REGISTER_API_URL"))
	gConf := gConfBuilder.Build()
	gwAdmin := fcrgatewayadmin.NewFilecoinRetrievalGatewayAdmin(*gConf)

	// Configure provider admin
	pConfBuilder := fcrprovideradmin.CreateSettings()
	pConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	pConfBuilder.SetRegisterURL(providerTestProviderConfig.GetString("REGISTER_API_URL"))
	conf := pConfBuilder.Build()
	pAdmin := fcrprovideradmin.NewFilecoinRetrievalProviderAdmin(*conf)

	// Configure client
	clientConfBuilder := fcrclient.CreateSettings()
	clientConfBuilder.SetEstablishmentTTL(101)
	clientConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	clientConfBuilder.SetRegisterURL(providerTestProviderConfig.GetString("REGISTER_API_URL"))
	clientConf := clientConfBuilder.Build()

	// Create client
	client, err := fcrclient.NewFilecoinRetrievalClient(*clientConf, rm)
	if err != nil {
		logging.Error("error creating retrieval client: %s", err.Error())
		t.FailNow()
	}

	// Initialise gateway
	gatewayRootPubKey, gatewayRetrievalPubKey, gatewayRetrievalPrivateKey, err := generateKeys()
	gatewayID, err := nodeid.NewNodeIDFromHexString("ebc134a429ba7dc4811bf64ccb67057f5bd57ca4676800e2f71731cbcc5eb518")
	if err != nil {
		logging.Error("error generating gateway id")
		os.Exit(1)
	}
	gatewayRegistrar := register.NewGatewayRegister(
		gatewayID.ToString(),
		gatewayConfigGatewayConfig.GetString("GATEWAY_ADDRESS"),
		gatewayRootPubKey,
		gatewayRetrievalPubKey,
		gatewayConfigGatewayConfig.GetString("GATEWAY_REGION_CODE"),
		gatewayConfigGatewayConfig.GetString("NETWORK_INFO_GATEWAY"),
		gatewayConfigGatewayConfig.GetString("NETWORK_INFO_PROVIDER"),
		gatewayConfigGatewayConfig.GetString("NETWORK_INFO_CLIENT"),
		gatewayConfigGatewayConfig.GetString("NETWORK_INFO_ADMIN"),
	)
	err = gwAdmin.InitialiseGateway(gatewayRegistrar, gatewayRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1))
	if err != nil {
		panic(err)
	}

	if err = rm.RegisterGateway(gatewayRegistrar); err != nil {
		t.Errorf("can't register gateway")
	}
	// Add the gateways to the passive list
	added := client.AddGatewaysToUse([]*nodeid.NodeID{gatewayID})
	if !assert.Equal(t, 1, added, "32 gateways should be added") {
		t.FailNow()
	}
	// Make the gateways active, this involves doing an establishment
	addedActive := client.AddActiveGateways([]*nodeid.NodeID{gatewayID})
	if !assert.Equal(t, 1, addedActive, "32 gateways should be activated") {
		t.FailNow()
	}

	// Initialise provider
	providerRootPubKey, providerRetrievalPubKey, providerRetrievalPrivateKey, err := generateKeys()
	providerID, err := nodeid.NewNodeIDFromHexString("ebc134a429ba7dc4811bf64ccb67057f5bd57ca4676800e2f71731cbcc5eb518")
	if err != nil {
		logging.Error("can't generate provider ID: %s", err.Error())
		os.Exit(1)
	}
	providerRegistrar := register.NewProviderRegister(
		providerID.ToString(),
		providerTestProviderConfig.GetString("PROVIDER_ADDRESS"),
		providerRootPubKey,
		providerRetrievalPubKey,
		providerTestProviderConfig.GetString("PROVIDER_REGION_CODE"),
		providerTestProviderConfig.GetString("NETWORK_INFO_GATEWAY"),
		providerTestProviderConfig.GetString("NETWORK_INFO_CLIENT"),
		providerTestProviderConfig.GetString("NETWORK_INFO_ADMIN"),
	)
	// Initialise the provider using provider admin
	err = pAdmin.InitialiseProvider(providerRegistrar, providerRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1))
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}
	// Enroll the provider in the Register srv.
	if err := rm.RegisterProvider(providerRegistrar); err != nil {
		logging.Error("error registering provider: %s", err.Error())
		t.FailNow()
	}

	// Generate random cid offer
	contentID1 := cid.NewRandomContentID()
	contentID2 := cid.NewRandomContentID()
	pieceCIDs := []cid.ContentID{*contentID1, *contentID2}
	expiryDate := time.Now().Local().Add(time.Hour * time.Duration(24)).Unix()

	// Force update
	err = pAdmin.ForceUpdate(providerID)
	if err != nil {
		panic(err)
	}
	err = gwAdmin.ForceUpdate(gatewayID)
	if err != nil {
		panic(err)
	}

	// Publish Group CID
	err = pAdmin.PublishGroupCID(providerID, pieceCIDs, 42, expiryDate, 42)
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}

	// Get all offers
	var gatewayIDs []nodeid.NodeID
	gatewayIDs = make([]nodeid.NodeID, 0)
	logging.Info("Get all offers")
	_, cidgroupInfo, err := pAdmin.GetGroupCIDOffer(providerID, gatewayIDs)
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}
	logging.Info("Get all offers: %d", len(cidgroupInfo))
	assert.GreaterOrEqual(t, len(cidgroupInfo), 1, "Get all offers should be found")

	// Get offers by gatewayIDs real
	gateways := rm.GetAllGateways()
	if gateways == nil {
		logging.ErrorAndPanic("expecting list of gateways, got nil")
	}
	logging.Info("Registered gateways: %+v", gateways)
	realNodeID, err := nodeid.NewNodeIDFromHexString("ebc134a429ba7dc4811bf64ccb67057f5bd57ca4676800e2f71731cbcc5eb518")
	if err != nil {
		logging.Error("can't generate node ID: %s", err.Error())
		os.Exit(1)
	}
	gatewayIDs = append(gatewayIDs, *realNodeID) // Add a gateway
	logging.Info("Get offers by real gatewayID=%s", realNodeID.ToString())
	_, cidgroupInfo, err = pAdmin.GetGroupCIDOffer(providerID, gatewayIDs)
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}
	logging.Info("Get offers by real gatewayID=%s: %d", realNodeID.ToString(), len(cidgroupInfo))
	assert.GreaterOrEqual(t, len(cidgroupInfo), 1, "Get offers by gatewayIDs real should be found")

	// Get offers by gatewayIDs fake
	fakeNodeID, _ := nodeid.NewNodeIDFromHexString("101112131415161718191A1B1C1D1E1F202122232425262728292A2B2C2DFA43")
	gatewayIDs[0] = *fakeNodeID
	logging.Info("Get offers by fake gatewayID=%s", fakeNodeID.ToString())
	_, cidgroupInfo, err = pAdmin.GetGroupCIDOffer(providerID, gatewayIDs)
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

// Helper function to generate set of keys
func generateKeys() (rootPubKey string, retrievalPubKey string, retrievalPrivateKey *fcrcrypto.KeyPair, err error) {
	rootKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		return "", "", nil, fmt.Errorf("error generating blockchain key: %s", err.Error())
	}
	if rootKey == nil {
		return "", "", nil, errors.New("error generating blockchain key")
	}

	rootPubKey, err = rootKey.EncodePublicKey()
	if err != nil {
		return "", "", nil, fmt.Errorf("error encoding public key: %s", err.Error())
	}

	retrievalPrivateKey, err = fcrcrypto.GenerateRetrievalV1KeyPair()
	if err != nil {
		return "", "", nil, fmt.Errorf("error generating retrieval key: %s", err.Error())
	}
	if retrievalPrivateKey == nil {
		return "", "", nil, errors.New("error generating retrieval key")
	}

	retrievalPubKey, err = retrievalPrivateKey.EncodePublicKey()
	if err != nil {
		return "", "", nil, fmt.Errorf("error encoding retrieval pub key: %s", err.Error())
	}
	return
}
