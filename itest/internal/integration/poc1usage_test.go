package integration

import (
	"testing"
	"time"

	"github.com/ConsenSys/fc-retrieval-client/pkg/fcrclient"
	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-gateway-admin/pkg/fcrgatewayadmin"
	"github.com/ConsenSys/fc-retrieval-itest/config"
	"github.com/ConsenSys/fc-retrieval-provider-admin/pkg/fcrprovideradmin"
	"github.com/ConsenSys/fc-retrieval-register/pkg/register"
	"github.com/stretchr/testify/assert"
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
var gwAdmin *fcrgatewayadmin.FilecoinRetrievalGatewayAdminClient
var pAdmin *fcrprovideradmin.FilecoinRetrievalProviderAdminClient
var client *fcrclient.FilecoinRetrievalClient
var gwID *nodeid.NodeID
var pID *nodeid.NodeID
var testCIDs []cid.ContentID

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
	confBuilder.SetEstablishmentTTL(101)
	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	confBuilder.SetRegisterURL(gatewayConfig.GetString("REGISTER_API_URL"))
	conf := confBuilder.Build()
	gwAdmin = fcrgatewayadmin.NewFilecoinRetrievalGatewayAdminClient(*conf)

	gatewayRootKey, err := fcrgatewayadmin.CreateKey()
	if err != nil {
		panic(err)
	}
	gatewayRootSigningKey, err := gatewayRootKey.EncodePublicKey()
	if err != nil {
		panic(err)
	}
	gatewayRetrievalPrivateKey, err := fcrgatewayadmin.CreateKey()
	if err != nil {
		panic(err)
	}
	gatewayRetrievalSigningKey, err := gatewayRetrievalPrivateKey.EncodePublicKey()
	if err != nil {
		panic(err)
	}
	gatewayID, err := nodeid.NewRandomNodeID()
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

	err = gwAdmin.InitializeGateway(gatewayRegister, gatewayRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1))
	if err != nil {
		panic(err)
	}
	logging.Info("Wait five seconds for the gateway to initialise")
	time.Sleep(5 * time.Second)

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
	pAdmin = fcrprovideradmin.InitFilecoinRetrievalProviderAdminClient(*conf)

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
	providerID, err := nodeid.NewRandomNodeID()
	if err != nil {
		panic(err)
	}
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
	logging.Info("Wait five seconds for the provider to initialise")
	time.Sleep(5 * time.Second)

	logging.Info("/*******************************************************/")
	logging.Info("/*              End TestInitialiseProvider	         */")
	logging.Info("/*******************************************************/")
}

func TestPublishGroupCID(t *testing.T) {
	logging.Info("/*******************************************************/")
	logging.Info("/*      Start TestProviderPublishGroupCIDOffer	     */")
	logging.Info("/*******************************************************/")

	// Generate random group cid offer
	contentID1, _ := cid.NewRandomContentID()
	contentID2, _ := cid.NewRandomContentID()
	contentID3, _ := cid.NewRandomContentID()
	pieceCIDs := []cid.ContentID{*contentID1, *contentID2, *contentID3}
	expiryDate := time.Now().Local().Add(time.Hour * time.Duration(24)).Unix()

	// Set global variable
	testCIDs = pieceCIDs

	// Publish Group CID
	err := pAdmin.PublishGroupCID(pID, pieceCIDs, 42, expiryDate, 42)
	if err != nil {
		panic(err)
	}
	logging.Info("Wait 3 seconds")
	time.Sleep(3 * time.Second)

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
	fakeNodeID, _ := nodeid.NewNodeIDFromString("101112131415161718191A1B1C1D1E1F202122232425262728292A2B2C2DFA43")
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
	client = fcrclient.NewFilecoinRetrievalClient(*conf)

	client.RefreshLatestGatewayInfo()
	client.RefreshLatestProviderInfo()

	logging.Info("/*******************************************************/")
	logging.Info("/*                 End TestInitClient      	         */")
	logging.Info("/*******************************************************/")
}

func TestClientAddGateway(t *testing.T) {
	logging.Info("/*******************************************************/")
	logging.Info("/*             Start TestClientAddGateway     	     */")
	logging.Info("/*******************************************************/")

	// Add a gateway, this involves doing an establishment
	added := client.AddGateways([]*nodeid.NodeID{gwID})
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

	offers, err := client.FindOffersStandardDiscovery(&(testCIDs[0]))
	if err != nil {
		panic(err)
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0.") {
		return
	}

	offers, err = client.FindOffersStandardDiscovery(&(testCIDs[1]))
	if err != nil {
		panic(err)
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 1.") {
		return
	}

	offers, err = client.FindOffersStandardDiscovery(&(testCIDs[2]))
	if err != nil {
		panic(err)
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 2.") {
		return
	}

	randomCID, _ := cid.NewRandomContentID()
	offers, err = client.FindOffersStandardDiscovery(randomCID)
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
