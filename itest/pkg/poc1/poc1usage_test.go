/*
Package poc1 - set of end-to-end tests, designed to demonstrate functionality required for Proof of Concept stage 1.
*/
package poc1

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"

	"github.com/ConsenSys/fc-retrieval/client/pkg/fcrclient"
	"github.com/ConsenSys/fc-retrieval/common/pkg/cid"
	cr "github.com/ConsenSys/fc-retrieval/itest/pkg/util/crypto-facade"
	tc "github.com/ConsenSys/fc-retrieval/itest/pkg/util/test-containers"

	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrregistermgr"
	"github.com/ConsenSys/fc-retrieval/common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval/common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval/common/pkg/register"
	"github.com/ConsenSys/fc-retrieval/gateway-admin/pkg/fcrgatewayadmin"
	"github.com/ConsenSys/fc-retrieval/itest/config"
	"github.com/ConsenSys/fc-retrieval/provider-admin/pkg/fcrprovideradmin"
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
var registerConfig = config.NewConfig(".env.register")
var containers tc.AllContainers

func TestMain(m *testing.M) {
	const testName = "poc1"
	ctx := context.Background()
	var network *testcontainers.Network
	var err error
	containers, network, err = tc.StartContainers(ctx, 1, 1, testName, true, gatewayConfig, providerConfig, registerConfig)
	if err != nil {
		logging.Error("%s failed, container starting error: %s", testName, err.Error())
		tc.StopContainers(ctx, testName, containers, network)
		os.Exit(1)
	}
	defer tc.StopContainers(ctx, testName, containers, network)
	m.Run()
}

func TestInitialiseGateway(t *testing.T) {
	logging.Info("/*******************************************************/")
	logging.Info("/*             Start TestInitialiseGateway	         */")
	logging.Info("/*******************************************************/")

	// Main key used across the chain
	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		logging.Error("error generating blockchain key: %s", err.Error())
		t.FailNow()
	}

	// Create and start register manager
	registerApiEndpoint := "http://" + containers.Register.GetRegisterHostApiEndpoint()
	var rm = fcrregistermgr.NewFCRRegisterMgr(registerApiEndpoint, false, true, 10*time.Second)
	if err := rm.Start(); err != nil {
		logging.Error("error starting Register Manager: %s", err.Error())
		t.FailNow()
	}
	defer rm.ShutdownAndWait()

	// Configure gateway admin
	confBuilder := fcrgatewayadmin.CreateSettings()
	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	confBuilder.SetRegisterURL("http://" + containers.Register.GetRegisterHostApiEndpoint())
	conf := confBuilder.Build()
	gwAdmin := fcrgatewayadmin.NewFilecoinRetrievalGatewayAdmin(*conf)

	// Configure gateway register
	gatewayRootPubKey, gatewayRetrievalPubKey, gatewayRetrievalPrivateKey, err := cr.GenerateKeys()
	if err != nil {
		logging.Error("can't generate key pairs %s", err.Error())
		t.FailNow()
	}
	gatewayID := nodeid.NewRandomNodeID()
	gatewayName := "gateway-0"
	_, _, _, gatewayAdminApiEndpoint := containers.Gateways[gatewayName].GetGatewayHostApiEndpoints()
	gatewayRegistrar := register.NewGatewayRegister(
		gatewayID.ToString(),
		gatewayConfig.GetString("GATEWAY_ADDRESS"),
		gatewayRootPubKey,
		gatewayRetrievalPubKey,
		gatewayConfig.GetString("GATEWAY_REGION_CODE"),
		gatewayName+":"+gatewayConfig.GetString("BIND_GATEWAY_API"),
		gatewayName+":"+gatewayConfig.GetString("BIND_PROVIDER_API"),
		gatewayName+":"+gatewayConfig.GetString("BIND_REST_API"),
		gatewayName+":"+gatewayConfig.GetString("BIND_ADMIN_API"),
	)
	// Initialise the gateway using gateway admin
	err = gwAdmin.InitialiseGateway(gatewayAdminApiEndpoint, gatewayRegistrar, gatewayRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1))
	if err != nil {
		logging.Error("gateway initialising error: %s", err.Error())
		t.FailNow()
	}
	// Enroll the gateway in the Register srv.
	if err := rm.RegisterGateway(gatewayRegistrar); err != nil {
		logging.Error("gateway registering error: %s", err.Error())
		t.FailNow()
	}

	logging.Info("/*******************************************************/")
	logging.Info("/*               End TestInitialiseGateway	         */")
	logging.Info("/*******************************************************/")
}

func TestInitialiseProvider(t *testing.T) {
	logging.Info("/*******************************************************/")
	logging.Info("/*             Start TestInitialiseProvider	         */")
	logging.Info("/*******************************************************/")

	// Main key used across the chain
	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		logging.Error("error generating blockchain key: %s", err.Error())
		t.FailNow()
	}

	// Create and start register manager
	registerApiEndpoint := "http://" + containers.Register.GetRegisterHostApiEndpoint()
	var rm = fcrregistermgr.NewFCRRegisterMgr(registerApiEndpoint, true, true, 2*time.Second)
	if err := rm.Start(); err != nil {
		logging.Error("error starting Register Manager: %s", err.Error())
		t.FailNow()
	}
	defer rm.ShutdownAndWait()

	// Configure provider admin
	confBuilder := fcrprovideradmin.CreateSettings()
	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	confBuilder.SetRegisterURL(registerApiEndpoint)
	conf := confBuilder.Build()
	pAdmin := fcrprovideradmin.NewFilecoinRetrievalProviderAdmin(*conf)

	// Configure provider register
	providerRootPubKey, providerRetrievalPubKey, providerRetrievalPrivateKey, err := cr.GenerateKeys()
	if err != nil {
		logging.Error("can't generate key pairs %s", err.Error())
		t.FailNow()
	}
	providerID := nodeid.NewRandomNodeID()
	providerName := "provider-0"
	_, _, adminApiEndpoint := containers.Providers[providerName].GetProviderHostApiEndpoints()
	providerRegistrar := register.NewProviderRegister(
		providerID.ToString(),
		providerConfig.GetString("PROVIDER_ADDRESS"),
		providerRootPubKey,
		providerRetrievalPubKey,
		providerConfig.GetString("PROVIDER_REGION_CODE"),
		providerName+":"+providerConfig.GetString("BIND_GATEWAY_API"),
		providerName+":"+providerConfig.GetString("BIND_REST_API"),
		providerName+":"+providerConfig.GetString("BIND_ADMIN_API"),
	)

	// Initialise the provider using provider admin
	err = pAdmin.InitialiseProvider(adminApiEndpoint, providerRegistrar, providerRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1))
	if err != nil {
		logging.Error("error initialising provider: %s", err.Error())
		t.FailNow()
	}

	// Enroll the provider in the Register srv.
	if err := rm.RegisterProvider(providerRegistrar); err != nil {
		logging.Error("error registering provider: %s", err.Error())
		t.FailNow()
	}

	logging.Info("/*******************************************************/")
	logging.Info("/*              End TestInitialiseProvider	         */")
	logging.Info("/*******************************************************/")
}

func TestPublishGroupCID(t *testing.T) {
	logging.Info("/*******************************************************/")
	logging.Info("/*      Start TestProviderPublishGroupCIDOffer	     */")
	logging.Info("/*******************************************************/")

	// Main key used across the chain
	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		logging.Error("error generating blockchain key: %s", err.Error())
		t.FailNow()
	}

	// Create and start register manager
	registerApiEndpoint := "http://" + containers.Register.GetRegisterHostApiEndpoint()
	var rm = fcrregistermgr.NewFCRRegisterMgr(registerApiEndpoint, true, true, time.Second)
	if err := rm.Start(); err != nil {
		logging.Error("error starting Register Manager: %s", err.Error())
		t.FailNow()
	}
	defer rm.ShutdownAndWait()

	// Configure gateway admin
	gConfBuilder := fcrgatewayadmin.CreateSettings()
	gConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	gConfBuilder.SetRegisterURL(registerApiEndpoint)
	gConf := gConfBuilder.Build()
	gwAdmin := fcrgatewayadmin.NewFilecoinRetrievalGatewayAdmin(*gConf)

	// Configure provider admin
	pConfBuilder := fcrprovideradmin.CreateSettings()
	pConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	pConfBuilder.SetRegisterURL(registerApiEndpoint)
	pConf := pConfBuilder.Build()
	pAdmin := fcrprovideradmin.NewFilecoinRetrievalProviderAdmin(*pConf)

	// Configure gateway register
	gatewayRootPubKey, gatewayRetrievalPubKey, gatewayRetrievalPrivateKey, err := cr.GenerateKeys()
	if err != nil {
		logging.Error("can't generate key pairs %s", err.Error())
		t.FailNow()
	}
	gatewayID := nodeid.NewRandomNodeID()
	gatewayName := "gateway-0"
	_, _, _, gatewayAdminApiEndpoint := containers.Gateways[gatewayName].GetGatewayHostApiEndpoints()
	gatewayRegistrar := register.NewGatewayRegister(
		gatewayID.ToString(),
		gatewayConfig.GetString("GATEWAY_ADDRESS"),
		gatewayRootPubKey,
		gatewayRetrievalPubKey,
		gatewayConfig.GetString("GATEWAY_REGION_CODE"),
		gatewayName+":"+gatewayConfig.GetString("BIND_GATEWAY_API"),
		gatewayName+":"+gatewayConfig.GetString("BIND_PROVIDER_API"),
		gatewayName+":"+gatewayConfig.GetString("BIND_REST_API"),
		gatewayName+":"+gatewayConfig.GetString("BIND_ADMIN_API"),
	)
	err = gwAdmin.InitialiseGateway(gatewayAdminApiEndpoint, gatewayRegistrar, gatewayRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1))
	if err != nil {
		logging.Error("gateway initialising error: %s", err.Error())
		t.FailNow()
	}
	if err := rm.RegisterGateway(gatewayRegistrar); err != nil {
		logging.Error("gateway registering error: %s", err.Error())
		t.FailNow()
	}

	// Configure provider register
	providerRootPubKey, providerRetrievalPubKey, providerRetrievalPrivateKey, err := cr.GenerateKeys()
	if err != nil {
		logging.Error("can't generate key pairs %s", err.Error())
		t.FailNow()
	}
	providerID := nodeid.NewRandomNodeID()
	providerName := "provider-0"
	_, _, providerAdminApiEndpoint := containers.Providers[providerName].GetProviderHostApiEndpoints()
	providerRegistrar := register.NewProviderRegister(
		providerID.ToString(),
		providerConfig.GetString("PROVIDER_ADDRESS"),
		providerRootPubKey,
		providerRetrievalPubKey,
		providerConfig.GetString("PROVIDER_REGION_CODE"),
		providerName+":"+providerConfig.GetString("BIND_GATEWAY_API"),
		providerName+":"+providerConfig.GetString("BIND_REST_API"),
		providerName+":"+providerConfig.GetString("BIND_ADMIN_API"),
	)
	err = pAdmin.InitialiseProvider(providerAdminApiEndpoint, providerRegistrar, providerRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1))
	if err != nil {
		logging.Error("error initialising provider: %s", err.Error())
		t.FailNow()
	}
	if err := rm.RegisterProvider(providerRegistrar); err != nil {
		logging.Error("error registering provider: %s", err.Error())
		t.FailNow()
	}

	// Force provider and gateway to update
	err = pAdmin.ForceUpdate(providerAdminApiEndpoint, providerID)
	if err != nil {
		logging.Error("error forcing update provider: %s", err.Error())
		t.FailNow()
	}
	err = gwAdmin.ForceUpdate(gatewayAdminApiEndpoint, gatewayID)
	if err != nil {
		logging.Error("error forcing update gateway: %s", err.Error())
		t.FailNow()
	}

	// Generate random group cid offer
	contentID1 := cid.NewRandomContentID()
	contentID2 := cid.NewRandomContentID()
	contentID3 := cid.NewRandomContentID()
	pieceCIDs := []cid.ContentID{*contentID1, *contentID2, *contentID3}
	expiryDate := time.Now().Local().Add(time.Hour * time.Duration(24)).Unix()

	// Publish Group CID
	err = pAdmin.PublishGroupCID(providerAdminApiEndpoint, providerID, pieceCIDs, 42, expiryDate, 42)
	if err != nil {
		logging.Error("error publishing group CID: %s", err.Error())
		t.FailNow()
	}

	// Test get all offers
	gatewayIDs := make([]nodeid.NodeID, 0)
	logging.Info("Get all offers")
	_, cidgroupInfo, err := pAdmin.GetGroupCIDOffer(providerAdminApiEndpoint, providerID, gatewayIDs)
	if err != nil {
		logging.Error("error getting group CID offer: %s", err.Error())
		t.FailNow()
	}
	fmt.Printf("<<<<<<<<<<<<< cidgroupInfo: %#v", cidgroupInfo)
	if !assert.GreaterOrEqual(t, len(cidgroupInfo), 1, "Offers should be found") {
		t.FailNow()
	}
	// Add a gateway and  verify
	gatewayIDs = append(gatewayIDs, *gatewayID)
	logging.Info("Get offers by gatewayID=%s", gatewayID.ToString())
	_, cidgroupInfo, err = pAdmin.GetGroupCIDOffer(providerAdminApiEndpoint, providerID, gatewayIDs)
	if err != nil {
		logging.Error("error getting group CID offer: %s", err.Error())
		t.FailNow()
	}
	if !assert.GreaterOrEqual(t, len(cidgroupInfo), 1, "Offers should be found") {
		t.FailNow()
	}

	// Negative test using fake content ID
	fakeNodeID, _ := nodeid.NewNodeIDFromHexString("101112131415161718191A1B1C1D1E1F202122232425262728292A2B2C2DFA43")
	gatewayIDs[0] = *fakeNodeID
	logging.Info("Get offers by gatewayID=%s", fakeNodeID.ToString())
	_, cidgroupInfo, err = pAdmin.GetGroupCIDOffer(providerAdminApiEndpoint, providerID, gatewayIDs)
	if err != nil {
		logging.Error("error getting group CID offer: %s", err.Error())
		t.FailNow()
	}
	assert.Equal(t, 0, len(cidgroupInfo), "Offers should be empty")

	logging.Info("/*******************************************************/")
	logging.Info("/*       End TestProviderPublishGroupCIDOffer	         */")
	logging.Info("/*******************************************************/")
}

func TestInitClient(t *testing.T) {
	logging.Info("/*******************************************************/")
	logging.Info("/*                Start TestInitClient        	     */")
	logging.Info("/*******************************************************/")

	// Main key used across the chain
	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		logging.Error("error generating blockchain key: %s", err.Error())
		t.FailNow()
	}

	// Create and start register manager
	registerApiEndpoint := "http://" + containers.Register.GetRegisterHostApiEndpoint()
	var rm = fcrregistermgr.NewFCRRegisterMgr(registerApiEndpoint, false, false, 2*time.Second)
	if err := rm.Start(); err != nil {
		logging.Error("error starting Register Manager: %s", err.Error())
		t.FailNow()
	}
	defer rm.ShutdownAndWait()

	// Configure client
	confBuilder := fcrclient.CreateSettings()
	confBuilder.SetEstablishmentTTL(101)
	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	confBuilder.SetRegisterURL(registerApiEndpoint)
	conf := confBuilder.Build()

	// Create client and verify
	_, err = fcrclient.NewFilecoinRetrievalClient(*conf, rm)
	if err != nil {
		logging.Error("error creating retrieval client: %s", err.Error())
		t.FailNow()
	}

	logging.Info("/*******************************************************/")
	logging.Info("/*                 End TestInitClient      	         */")
	logging.Info("/*******************************************************/")
}

func TestClientAddGateway(t *testing.T) {
	logging.Info("/*******************************************************/")
	logging.Info("/*             Start TestClientAddGateway     	     */")
	logging.Info("/*******************************************************/")

	// Main key used across the chain
	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		logging.Error("error generating blockchain key: %s", err.Error())
		t.FailNow()
	}

	// Create and start register manager
	registerApiEndpoint := "http://" + containers.Register.GetRegisterHostApiEndpoint()
	var rm = fcrregistermgr.NewFCRRegisterMgr(registerApiEndpoint, true, true, time.Second)
	if err := rm.Start(); err != nil {
		logging.Error("error starting Register Manager: %s", err.Error())
	}
	defer rm.ShutdownAndWait()

	// Configure gateway admin
	gConfBuilder := fcrgatewayadmin.CreateSettings()
	gConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	gConfBuilder.SetRegisterURL(registerApiEndpoint)
	gConf := gConfBuilder.Build()
	gwAdmin := fcrgatewayadmin.NewFilecoinRetrievalGatewayAdmin(*gConf)

	// Configure client
	confBuilder := fcrclient.CreateSettings()
	confBuilder.SetEstablishmentTTL(101)
	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	confBuilder.SetRegisterURL(registerApiEndpoint)
	conf := confBuilder.Build()
	client, err := fcrclient.NewFilecoinRetrievalClient(*conf, rm)
	if err != nil {
		logging.Error("error creating retrieval client: %s", err.Error())
		t.FailNow()
	}

	// Configure gateway register
	gatewayRootPubKey, gatewayRetrievalPubKey, gatewayRetrievalPrivateKey, err := cr.GenerateKeys()
	if err != nil {
		logging.Error("can't generate key pairs %s", err.Error())
		t.FailNow()
	}
	gatewayID := nodeid.NewRandomNodeID()
	gatewayName := "gateway-0"
	_, _, gatewayClientApiEndpoint, gatewayAdminApiEndpoint := containers.Gateways[gatewayName].GetGatewayHostApiEndpoints()
	gatewayRegistrar := register.NewGatewayRegister(
		gatewayID.ToString(),
		gatewayConfig.GetString("GATEWAY_ADDRESS"),
		gatewayRootPubKey,
		gatewayRetrievalPubKey,
		gatewayConfig.GetString("GATEWAY_REGION_CODE"),
		gatewayName+":"+gatewayConfig.GetString("BIND_GATEWAY_API"),
		gatewayName+":"+gatewayConfig.GetString("BIND_PROVIDER_API"),
		gatewayName+":"+gatewayConfig.GetString("BIND_REST_API"),
		gatewayName+":"+gatewayConfig.GetString("BIND_ADMIN_API"),
	)
	err = gwAdmin.InitialiseGateway(gatewayAdminApiEndpoint, gatewayRegistrar, gatewayRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1))
	if err != nil {
		logging.Error("gateway initialising error: %s", err.Error())
		t.FailNow()
	}
	if err := rm.RegisterGateway(gatewayRegistrar); err != nil {
		logging.Error("gateway registering error: %s", err.Error())
		t.FailNow()
	}

	// Add a gateway to the passive list
	added := client.AddGatewaysToUse([]*nodeid.NodeID{gatewayID})
	if !assert.Equal(t, 1, added, "One gateway should be added") {
		t.FailNow()
	}
	// Make the gateway active, this involves doing an establishment
	added = client.AddActiveGateways(gatewayClientApiEndpoint, []*nodeid.NodeID{gatewayID})
	assert.Equal(t, 1, added, "One gateway should be added")

	logging.Info("/*******************************************************/")
	logging.Info("/*              End TestClientAddGateway      	     */")
	logging.Info("/*******************************************************/")
}

func TestClientStdContentDiscover(t *testing.T) {
	logging.Info("/*******************************************************/")
	logging.Info("/*        Start TestClientStdContentDiscover     	     */")
	logging.Info("/*******************************************************/")

	// Main key used across the chain
	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		logging.Error("error generating blockchain key: %s", err.Error())
		t.FailNow()
	}

	// Create and start register manager
	registerApiEndpoint := "http://" + containers.Register.GetRegisterHostApiEndpoint()
	var rm = fcrregistermgr.NewFCRRegisterMgr(registerApiEndpoint, true, true, time.Second)
	if err := rm.Start(); err != nil {
		logging.Error("error starting Register Manager: %s", err.Error())
	}
	defer rm.ShutdownAndWait()

	// Configure gateway admin
	gConfBuilder := fcrgatewayadmin.CreateSettings()
	gConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	gConfBuilder.SetRegisterURL(registerApiEndpoint)
	gConf := gConfBuilder.Build()
	gwAdmin := fcrgatewayadmin.NewFilecoinRetrievalGatewayAdmin(*gConf)

	// Configure provider admin
	pConfBuilder := fcrprovideradmin.CreateSettings()
	pConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	pConfBuilder.SetRegisterURL(registerApiEndpoint)
	pConf := pConfBuilder.Build()
	pAdmin := fcrprovideradmin.NewFilecoinRetrievalProviderAdmin(*pConf)

	// Configure gateway register
	gatewayRootPubKey, gatewayRetrievalPubKey, gatewayRetrievalPrivateKey, err := cr.GenerateKeys()
	if err != nil {
		logging.Error("can't generate key pairs %s", err.Error())
		t.FailNow()
	}
	gatewayID := nodeid.NewRandomNodeID()
	gatewayName := "gateway-0"
	_, _, gatewayClientApiEndpoint, gatewayAdminApiEndpoint := containers.Gateways[gatewayName].GetGatewayHostApiEndpoints()
	gatewayRegistrar := register.NewGatewayRegister(
		gatewayID.ToString(),
		gatewayConfig.GetString("GATEWAY_ADDRESS"),
		gatewayRootPubKey,
		gatewayRetrievalPubKey,
		gatewayConfig.GetString("GATEWAY_REGION_CODE"),
		gatewayName+":"+gatewayConfig.GetString("BIND_GATEWAY_API"),
		gatewayName+":"+gatewayConfig.GetString("BIND_PROVIDER_API"),
		gatewayName+":"+gatewayConfig.GetString("BIND_REST_API"),
		gatewayName+":"+gatewayConfig.GetString("BIND_ADMIN_API"),
	)
	err = gwAdmin.InitialiseGateway(gatewayAdminApiEndpoint, gatewayRegistrar, gatewayRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1))
	if err != nil {
		logging.Error("gateway initialising error: %s", err.Error())
		t.FailNow()
	}
	if err := rm.RegisterGateway(gatewayRegistrar); err != nil {
		logging.Error("gateway registering error: %s", err.Error())
		t.FailNow()
	}

	// Configure provider register
	providerRootPubKey, providerRetrievalPubKey, providerRetrievalPrivateKey, err := cr.GenerateKeys()
	if err != nil {
		logging.Error("can't generate key pairs %s", err.Error())
		t.FailNow()
	}
	providerID := nodeid.NewRandomNodeID()
	providerName := "provider-0"
	_, _, providerAdminApiEndpoint := containers.Providers[providerName].GetProviderHostApiEndpoints()
	providerRegistrar := register.NewProviderRegister(
		providerID.ToString(),
		providerConfig.GetString("PROVIDER_ADDRESS"),
		providerRootPubKey,
		providerRetrievalPubKey,
		providerConfig.GetString("PROVIDER_REGION_CODE"),
		providerName+":"+providerConfig.GetString("BIND_GATEWAY_API"),
		providerName+":"+providerConfig.GetString("BIND_REST_API"),
		providerName+":"+providerConfig.GetString("BIND_ADMIN_API"),
	)
	err = pAdmin.InitialiseProvider(providerAdminApiEndpoint, providerRegistrar, providerRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1))
	if err != nil {
		logging.Error("error initialising provider: %s", err.Error())
		t.FailNow()
	}
	if err := rm.RegisterProvider(providerRegistrar); err != nil {
		logging.Error("error registering provider: %s", err.Error())
		t.FailNow()
	}

	// Configure client
	clientConfBuilder := fcrclient.CreateSettings()
	clientConfBuilder.SetEstablishmentTTL(101)
	clientConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	clientConfBuilder.SetRegisterURL(registerApiEndpoint)
	clientConf := clientConfBuilder.Build()
	client, err := fcrclient.NewFilecoinRetrievalClient(*clientConf, rm)
	if err != nil {
		logging.Error("error creating retrieval client: %s", err.Error())
		t.FailNow()
	}
	// Add a gateway to the passive list
	added := client.AddGatewaysToUse([]*nodeid.NodeID{gatewayID})
	if !assert.Equal(t, 1, added, "One gateway should be added") {
		t.FailNow()
	}
	// Make the gateway active, this involves doing an establishment
	added = client.AddActiveGateways(gatewayClientApiEndpoint, []*nodeid.NodeID{gatewayID})
	assert.Equal(t, 1, added, "One gateway should be added")

	// Force provider and gateway to update
	err = pAdmin.ForceUpdate(providerAdminApiEndpoint, providerID)
	if err != nil {
		logging.Error("error forcing update provider: %s", err.Error())
		t.FailNow()
	}
	err = gwAdmin.ForceUpdate(gatewayAdminApiEndpoint, gatewayID)
	if err != nil {
		logging.Error("error forcing update gateway: %s", err.Error())
		t.FailNow()
	}

	// Generate random group cid offer
	contentID1 := cid.NewRandomContentID()
	contentID2 := cid.NewRandomContentID()
	contentID3 := cid.NewRandomContentID()
	pieceCIDs := []cid.ContentID{*contentID1, *contentID2, *contentID3}
	expiryDate := time.Now().Local().Add(time.Hour * time.Duration(24)).Unix()

	// Publish Group CID
	err = pAdmin.PublishGroupCID(providerAdminApiEndpoint, providerID, pieceCIDs, 42, expiryDate, 42)
	if err != nil {
		logging.Error("error publishing group CID: %s", err.Error())
		t.FailNow()
	}

	// Find published offers
	offers, err := client.FindOffersStandardDiscovery(gatewayClientApiEndpoint, &(pieceCIDs[0]), gatewayID)
	if err != nil {
		logging.Error("error finding offer via standard discovery: %s", err.Error())
		t.FailNow()
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0.") {
		t.FailNow()
	}

	offers, err = client.FindOffersStandardDiscovery(gatewayClientApiEndpoint, &(pieceCIDs[1]), gatewayID)
	if err != nil {
		logging.Error("error finding offer via standard discovery: %s", err.Error())
		t.FailNow()
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 1.") {
		t.FailNow()
	}

	offers, err = client.FindOffersStandardDiscovery(gatewayClientApiEndpoint, &(pieceCIDs[2]), gatewayID)
	if err != nil {
		logging.Error("error finding offer via standard discovery: %s", err.Error())
		t.FailNow()
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 2.") {
		t.FailNow()
	}

	// Negative test using random content ID
	randomCID := cid.NewRandomContentID()
	offers, err = client.FindOffersStandardDiscovery(gatewayClientApiEndpoint, randomCID, gatewayID)
	if err != nil {
		logging.Error("error finding offer for random content ID via standard discovery: %s", err.Error())
		t.FailNow()
	}
	assert.Equal(t, 0, len(offers), "Shouldn't find any offer with random cid.")

	logging.Info("/*******************************************************/")
	logging.Info("/*        End TestClientStdContentDiscover     	     */")
	logging.Info("/*******************************************************/")
}
