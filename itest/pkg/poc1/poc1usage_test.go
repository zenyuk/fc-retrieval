/*
Package poc1 - set of end-to-end tests, designed to demonstrate functionality required for Proof of Concept stage 1.
*/
package poc1

import (
	"context"
	"os"
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

var gatewayConfig = config.NewConfig(".env.gateway")
var providerConfig = config.NewConfig(".env.provider")
var gwAdmin *fcrgatewayadmin.FilecoinRetrievalGatewayAdmin
var pAdmin *fcrprovideradmin.FilecoinRetrievalProviderAdmin
var client *fcrclient.FilecoinRetrievalClient
var gwID *nodeid.NodeID
var pID *nodeid.NodeID
var testCIDs []cid.ContentID

func TestMain(m *testing.M) {
	// Need to make sure this env is not set in host machine
	itestEnv := os.Getenv("ITEST_CALLING_FROM_CONTAINER")

	if itestEnv != "" {
		// Env is set, we are calling from docker container
		// This logging should be only called after all tests finished.
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

func TestInitialiseGateway(t *testing.T) {
	logging.Info("/*******************************************************/")
	logging.Info("/*             Start TestInitialiseGateway	         */")
	logging.Info("/*******************************************************/")
	logging.Error("Wait two seconds for the gateway to deploy and be ready for requests")

	time.Sleep(2 * time.Second)

	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		panic(err)
	}

	confBuilder := fcrgatewayadmin.CreateSettings()
	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	confBuilder.SetRegisterURL(gatewayConfig.GetString("REGISTER_API_URL"))
	conf := confBuilder.Build()
	gwAdmin = fcrgatewayadmin.NewFilecoinRetrievalGatewayAdmin(*conf)

	gatewayRootKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		panic(err)
	}
	gatewayRootSigningKey, err := gatewayRootKey.EncodePublicKey()
	if err != nil {
		panic(err)
	}
	gatewayRetrievalPrivateKey, err := fcrcrypto.GenerateRetrievalV1KeyPair()
	if err != nil {
		panic(err)
	}
	gatewayRetrievalSigningKey, err := gatewayRetrievalPrivateKey.EncodePublicKey()
	if err != nil {
		panic(err)
	}
	gatewayID := nodeid.NewRandomNodeID()
	if err != nil {
		panic(err)
	}
	gwID = gatewayID

	gatewayRegister := &register.GatewayRegister{
		NodeID:              gatewayID.ToString(),
		Address:             gatewayConfig.GetString("GATEWAY_ADDRESS"),
		RootSigningKey:      gatewayRootSigningKey,
		SigningKey:          gatewayRetrievalSigningKey,
		RegionCode:          gatewayConfig.GetString("GATEWAY_REGION_CODE"),
		NetworkInfoGateway:  gatewayConfig.GetString("NETWORK_INFO_GATEWAY"),
		NetworkInfoProvider: gatewayConfig.GetString("NETWORK_INFO_PROVIDER"),
		NetworkInfoClient:   gatewayConfig.GetString("NETWORK_INFO_CLIENT"),
		NetworkInfoAdmin:    gatewayConfig.GetString("NETWORK_INFO_ADMIN"),
	}

	err = gwAdmin.InitialiseGateway(gatewayRegister, gatewayRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1))
	if err != nil {
		panic(err)
	}

	logging.Info("/*******************************************************/")
	logging.Info("/*               End TestInitialiseGateway	         */")
	logging.Info("/*******************************************************/")
}

func TestInitialiseProvider(t *testing.T) {
	logging.Info("/*******************************************************/")
	logging.Info("/*             Start TestInitialiseProvider	         */")
	logging.Info("/*******************************************************/")
	logging.Error("Wait two seconds for the provider to deploy and be ready for requests")
	time.Sleep(2 * time.Second)

	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		panic(err)
	}

	confBuilder := fcrprovideradmin.CreateSettings()
	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	confBuilder.SetRegisterURL(providerConfig.GetString("REGISTER_API_URL"))
	conf := confBuilder.Build()
	pAdmin = fcrprovideradmin.NewFilecoinRetrievalProviderAdmin(*conf)

	providerRootKey, err := fcrcrypto.GenerateRetrievalV1KeyPair()
	if err != nil {
		panic(err)
	}
	providerRootSigningKey, err := providerRootKey.EncodePublicKey()
	if err != nil {
		panic(err)
	}
	providerPrivKey, err := fcrcrypto.GenerateRetrievalV1KeyPair()
	if err != nil {
		panic(err)
	}
	providerSigningKey, err := providerPrivKey.EncodePublicKey()
	if err != nil {
		panic(err)
	}
	providerID := nodeid.NewRandomNodeID()
	pID = providerID

	providerRegister := &register.ProviderRegister{
		NodeID:             providerID.ToString(),
		Address:            providerConfig.GetString("PROVIDER_ADDRESS"),
		RootSigningKey:     providerRootSigningKey,
		SigningKey:         providerSigningKey,
		RegionCode:         providerConfig.GetString("PROVIDER_REGION_CODE"),
		NetworkInfoGateway: providerConfig.GetString("NETWORK_INFO_GATEWAY"),
		NetworkInfoClient:  providerConfig.GetString("NETWORK_INFO_CLIENT"),
		NetworkInfoAdmin:   providerConfig.GetString("NETWORK_INFO_ADMIN"),
	}

	// Initialise provider
	err = pAdmin.InitialiseProvider(providerRegister, providerPrivKey, fcrcrypto.DecodeKeyVersion(1))
	if err != nil {
		panic(err)
	}

	logging.Info("/*******************************************************/")
	logging.Info("/*              End TestInitialiseProvider	         */")
	logging.Info("/*******************************************************/")
}

func TestPublishGroupCID(t *testing.T) {
	logging.Info("/*******************************************************/")
	logging.Info("/*      Start TestProviderPublishGroupCIDOffer	     */")
	logging.Info("/*******************************************************/")

	// Generate random group cid offer
	contentID1 := cid.NewRandomContentID()
	contentID2 := cid.NewRandomContentID()
	contentID3 := cid.NewRandomContentID()
	pieceCIDs := []cid.ContentID{*contentID1, *contentID2, *contentID3}
	expiryDate := time.Now().Local().Add(time.Hour * time.Duration(24)).Unix()

	// Set global variable
	testCIDs = pieceCIDs

	// Force update
	err := pAdmin.ForceUpdate(pID)
	if err != nil {
		panic(err)
	}
	err = gwAdmin.ForceUpdate(gwID)
	if err != nil {
		panic(err)
	}

	// Publish Group CID
	err = pAdmin.PublishGroupCID(pID, pieceCIDs, 42, expiryDate, 42)
	if err != nil {
		panic(err)
	}

	// Test get all offers
	gatewayIDs := make([]nodeid.NodeID, 0)
	logging.Info("Get all offers")
	_, cidgroupInfo, err := pAdmin.GetGroupCIDOffer(pID, gatewayIDs)
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}
	if !assert.GreaterOrEqual(t, len(cidgroupInfo), 1, "Offers should be found") {
		return
	}

	// Get offers by gatewayIDs real
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}
	gatewayIDs = append(gatewayIDs, *gwID) // Add a gateway
	logging.Info("Get offers by gatewayID=%s", gwID.ToString())
	_, cidgroupInfo, err = pAdmin.GetGroupCIDOffer(pID, gatewayIDs)
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}
	if !assert.GreaterOrEqual(t, len(cidgroupInfo), 1, "Offers should be found") {
		return
	}

	// Get offers by gatewayIDs fake
	fakeNodeID, _ := nodeid.NewNodeIDFromHexString("101112131415161718191A1B1C1D1E1F202122232425262728292A2B2C2DFA43")
	gatewayIDs[0] = *fakeNodeID
	logging.Info("Get offers by gatewayID=%s", fakeNodeID.ToString())
	_, cidgroupInfo, err = pAdmin.GetGroupCIDOffer(pID, gatewayIDs)
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}
	if !assert.Equal(t, 0, len(cidgroupInfo), "Offers should be empty") {
		return
	}

	logging.Info("/*******************************************************/")
	logging.Info("/*       End TestProviderPublishGroupCIDOffer	         */")
	logging.Info("/*******************************************************/")
}

func TestInitClient(t *testing.T) {
	logging.Info("/*******************************************************/")
	logging.Info("/*                Start TestInitClient        	     */")
	logging.Info("/*******************************************************/")

	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		panic(err)
	}

	confBuilder := fcrclient.CreateSettings()
	confBuilder.SetEstablishmentTTL(101)
	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	confBuilder.SetRegisterURL(gatewayConfig.GetString("REGISTER_API_URL"))
	conf := confBuilder.Build()
	client, err = fcrclient.NewFilecoinRetrievalClient(*conf)
	if err != nil {
		t.Fatal(err)
	}

	logging.Info("/*******************************************************/")
	logging.Info("/*                 End TestInitClient      	         */")
	logging.Info("/*******************************************************/")
}

func TestClientAddGateway(t *testing.T) {
	logging.Info("/*******************************************************/")
	logging.Info("/*             Start TestClientAddGateway     	     */")
	logging.Info("/*******************************************************/")

	// Add a gateway to use
	added := client.AddGatewaysToUse([]*nodeid.NodeID{gwID})
	if !assert.Equal(t, 1, added, "One gateway should be added") {
		return
	}

	// Make the gateway to active, this involves doing an establishment
	added = client.AddActiveGateways([]*nodeid.NodeID{gwID})
	if !assert.Equal(t, 1, added, "One gateway should be added") {
		return
	}

	logging.Info("/*******************************************************/")
	logging.Info("/*              End TestClientAddGateway      	     */")
	logging.Info("/*******************************************************/")
}

func TestClientStdContentDiscover(t *testing.T) {
	logging.Info("/*******************************************************/")
	logging.Info("/*        Start TestClientStdContentDiscover     	     */")
	logging.Info("/*******************************************************/")

	offers, err := client.FindOffersStandardDiscovery(&(testCIDs[0]), gwID)
	if err != nil {
		panic(err)
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0.") {
		return
	}

	offers, err = client.FindOffersStandardDiscovery(&(testCIDs[1]), gwID)
	if err != nil {
		panic(err)
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 1.") {
		return
	}

	offers, err = client.FindOffersStandardDiscovery(&(testCIDs[2]), gwID)
	if err != nil {
		panic(err)
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 2.") {
		return
	}

	randomCID := cid.NewRandomContentID()
	offers, err = client.FindOffersStandardDiscovery(randomCID, gwID)
	if err != nil {
		panic(err)
	}
	if !assert.Equal(t, 0, len(offers), "Shouldn't find any offer with random cid.") {
		return
	}

	logging.Info("/*******************************************************/")
	logging.Info("/*        End TestClientStdContentDiscover     	     */")
	logging.Info("/*******************************************************/")
}
