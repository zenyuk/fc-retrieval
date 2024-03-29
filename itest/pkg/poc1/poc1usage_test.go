/*
Package poc1 - set of end-to-end tests, designed to demonstrate functionality required for Proof of Concept stage 1.
*/
package poc1

import (
	"context"
	"errors"
	"fmt"
	"os"
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
	var rm = fcrregistermgr.NewFCRRegisterMgr(gatewayConfig.GetString("REGISTER_API_URL"), false, true, 10*time.Second)
	if err := rm.Start(); err != nil {
		logging.Error("error starting Register Manager: %s", err.Error())
		t.FailNow()
	}
	defer rm.ShutdownAndWait()

	// Configure gateway admin
	confBuilder := fcrgatewayadmin.CreateSettings()
	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	confBuilder.SetRegisterURL(gatewayConfig.GetString("REGISTER_API_URL"))
	conf := confBuilder.Build()
	gwAdmin := fcrgatewayadmin.NewFilecoinRetrievalGatewayAdmin(*conf)

	// Configure gateway register
	gatewayRootPubKey, gatewayRetrievalPubKey, gatewayRetrievalPrivateKey, err := generateKeys()
	if err != nil {
		logging.Error("can't generate key pairs %s", err.Error())
		t.FailNow()
	}
	gatewayID := nodeid.NewRandomNodeID()
	gatewayRegistrar := register.NewGatewayRegister(
		gatewayID.ToString(),
		gatewayConfig.GetString("GATEWAY_ADDRESS"),
		gatewayRootPubKey,
		gatewayRetrievalPubKey,
		gatewayConfig.GetString("GATEWAY_REGION_CODE"),
		gatewayConfig.GetString("NETWORK_INFO_GATEWAY"),
		gatewayConfig.GetString("NETWORK_INFO_PROVIDER"),
		gatewayConfig.GetString("NETWORK_INFO_CLIENT"),
		gatewayConfig.GetString("NETWORK_INFO_ADMIN"),
	)
	// Initialise the gateway using gateway admin
	err = gwAdmin.InitialiseGateway(gatewayRegistrar, gatewayRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1))
	if err != nil {
		logging.Error("error initialising gateway: %s", err.Error())
		t.FailNow()
	}
	// Enroll the gateway in the Register srv.
	if err := rm.RegisterGateway(gatewayRegistrar); err != nil {
		logging.Error("error registering gateway: %s", err.Error())
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
	var rm = fcrregistermgr.NewFCRRegisterMgr(gatewayConfig.GetString("REGISTER_API_URL"), true, true, 2*time.Second)
	if err := rm.Start(); err != nil {
		logging.Error("error starting Register Manager: %s", err.Error())
		t.FailNow()
	}
	defer rm.ShutdownAndWait()

	// Configure provider admin
	confBuilder := fcrprovideradmin.CreateSettings()
	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	confBuilder.SetRegisterURL(providerConfig.GetString("REGISTER_API_URL"))
	conf := confBuilder.Build()
	pAdmin := fcrprovideradmin.NewFilecoinRetrievalProviderAdmin(*conf)

	// Configure provider register
	providerRootPubKey, providerRetrievalPubKey, providerRetrievalPrivateKey, err := generateKeys()
	if err != nil {
		logging.Error("can't generate key pairs %s", err.Error())
		t.FailNow()
	}
	providerID := nodeid.NewRandomNodeID()
	providerRegistrar := register.NewProviderRegister(
		providerID.ToString(),
		providerConfig.GetString("PROVIDER_ADDRESS"),
		providerRootPubKey,
		providerRetrievalPubKey,
		providerConfig.GetString("PROVIDER_REGION_CODE"),
		providerConfig.GetString("NETWORK_INFO_GATEWAY"),
		providerConfig.GetString("NETWORK_INFO_CLIENT"),
		providerConfig.GetString("NETWORK_INFO_ADMIN"),
	)

	// Initialise the provider using provider admin
	err = pAdmin.InitialiseProvider(providerRegistrar, providerRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1))
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
	var rm = fcrregistermgr.NewFCRRegisterMgr(gatewayConfig.GetString("REGISTER_API_URL"), true, true, time.Second)
	if err := rm.Start(); err != nil {
		logging.Error("error starting Register Manager: %s", err.Error())
		t.FailNow()
	}
	defer rm.ShutdownAndWait()

	// Configure gateway admin
	gConfBuilder := fcrgatewayadmin.CreateSettings()
	gConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	gConfBuilder.SetRegisterURL(gatewayConfig.GetString("REGISTER_API_URL"))
	gConf := gConfBuilder.Build()
	gwAdmin := fcrgatewayadmin.NewFilecoinRetrievalGatewayAdmin(*gConf)

	// Configure provider admin
	pConfBuilder := fcrprovideradmin.CreateSettings()
	pConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	pConfBuilder.SetRegisterURL(providerConfig.GetString("REGISTER_API_URL"))
	pConf := pConfBuilder.Build()
	pAdmin := fcrprovideradmin.NewFilecoinRetrievalProviderAdmin(*pConf)

	// Configure gateway register
	gatewayRootPubKey, gatewayRetrievalPubKey, gatewayRetrievalPrivateKey, err := generateKeys()
	if err != nil {
		logging.Error("can't generate key pairs %s", err.Error())
		t.FailNow()
	}
	gatewayID := nodeid.NewRandomNodeID()
	gatewayRegistrar := register.NewGatewayRegister(
		gatewayID.ToString(),
		gatewayConfig.GetString("GATEWAY_ADDRESS"),
		gatewayRootPubKey,
		gatewayRetrievalPubKey,
		gatewayConfig.GetString("GATEWAY_REGION_CODE"),
		gatewayConfig.GetString("NETWORK_INFO_GATEWAY"),
		gatewayConfig.GetString("NETWORK_INFO_PROVIDER"),
		gatewayConfig.GetString("NETWORK_INFO_CLIENT"),
		gatewayConfig.GetString("NETWORK_INFO_ADMIN"),
	)
	err = gwAdmin.InitialiseGateway(gatewayRegistrar, gatewayRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1))
	if err != nil {
		logging.Error("error initialising gateway: %s", err.Error())
		t.FailNow()
	}
	if err := rm.RegisterGateway(gatewayRegistrar); err != nil {
		logging.Error("error registering gateway: %s", err.Error())
		t.FailNow()
	}

	// Configure provider register
	providerRootPubKey, providerRetrievalPubKey, providerRetrievalPrivateKey, err := generateKeys()
	if err != nil {
		logging.Error("can't generate key pairs %s", err.Error())
		t.FailNow()
	}
	providerID := nodeid.NewRandomNodeID()
	providerRegistrar := register.NewProviderRegister(
		providerID.ToString(),
		providerConfig.GetString("PROVIDER_ADDRESS"),
		providerRootPubKey,
		providerRetrievalPubKey,
		providerConfig.GetString("PROVIDER_REGION_CODE"),
		providerConfig.GetString("NETWORK_INFO_GATEWAY"),
		providerConfig.GetString("NETWORK_INFO_CLIENT"),
		providerConfig.GetString("NETWORK_INFO_ADMIN"),
	)
	err = pAdmin.InitialiseProvider(providerRegistrar, providerRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1))
	if err != nil {
		logging.Error("error initialising provider: %s", err.Error())
		t.FailNow()
	}
	if err := rm.RegisterProvider(providerRegistrar); err != nil {
		logging.Error("error registering provider: %s", err.Error())
		t.FailNow()
	}

	// Force provider and gateway to update
	err = pAdmin.ForceUpdate(providerID)
	if err != nil {
		logging.Error("error forcing update provider: %s", err.Error())
		t.FailNow()
	}
	err = gwAdmin.ForceUpdate(gatewayID)
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
	err = pAdmin.PublishGroupCID(providerID, pieceCIDs, 42, expiryDate, 42)
	if err != nil {
		logging.Error("error publishing group CID: %s", err.Error())
		t.FailNow()
	}

	// Test get all offers
	gatewayIDs := make([]nodeid.NodeID, 0)
	logging.Info("Get all offers")
	_, cidgroupInfo, err := pAdmin.GetGroupCIDOffer(providerID, gatewayIDs)
	if err != nil {
		logging.Error("error getting group CID offer: %s", err.Error())
		t.FailNow()
	}
	if !assert.GreaterOrEqual(t, len(cidgroupInfo), 1, "Offers should be found") {
		t.FailNow()
	}
	// Add a gateway and  verify
	gatewayIDs = append(gatewayIDs, *gatewayID)
	logging.Info("Get offers by gatewayID=%s", gatewayID.ToString())
	_, cidgroupInfo, err = pAdmin.GetGroupCIDOffer(providerID, gatewayIDs)
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
	_, cidgroupInfo, err = pAdmin.GetGroupCIDOffer(providerID, gatewayIDs)
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
	var rm = fcrregistermgr.NewFCRRegisterMgr(gatewayConfig.GetString("REGISTER_API_URL"), false, false, 2*time.Second)
	if err := rm.Start(); err != nil {
		logging.Error("error starting Register Manager: %s", err.Error())
		t.FailNow()
	}
	defer rm.ShutdownAndWait()

	// Configure client
	confBuilder := fcrclient.CreateSettings()
	confBuilder.SetEstablishmentTTL(101)
	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	confBuilder.SetRegisterURL(gatewayConfig.GetString("REGISTER_API_URL"))
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
	var rm = fcrregistermgr.NewFCRRegisterMgr(gatewayConfig.GetString("REGISTER_API_URL"), true, true, time.Second)
	if err := rm.Start(); err != nil {
		logging.Error("error starting Register Manager: %s", err.Error())
	}
	defer rm.ShutdownAndWait()

	// Configure gateway admin
	gConfBuilder := fcrgatewayadmin.CreateSettings()
	gConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	gConfBuilder.SetRegisterURL(gatewayConfig.GetString("REGISTER_API_URL"))
	gConf := gConfBuilder.Build()
	gwAdmin := fcrgatewayadmin.NewFilecoinRetrievalGatewayAdmin(*gConf)

	// Configure client
	confBuilder := fcrclient.CreateSettings()
	confBuilder.SetEstablishmentTTL(101)
	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	confBuilder.SetRegisterURL(gatewayConfig.GetString("REGISTER_API_URL"))
	conf := confBuilder.Build()
	client, err := fcrclient.NewFilecoinRetrievalClient(*conf, rm)
	if err != nil {
		logging.Error("error creating retrieval client: %s", err.Error())
		t.FailNow()
	}

	// Configure gateway register
	gatewayRootPubKey, gatewayRetrievalPubKey, gatewayRetrievalPrivateKey, err := generateKeys()
	if err != nil {
		logging.Error("can't generate key pairs %s", err.Error())
		t.FailNow()
	}
	gatewayID := nodeid.NewRandomNodeID()
	gatewayRegistrar := register.NewGatewayRegister(
		gatewayID.ToString(),
		gatewayConfig.GetString("GATEWAY_ADDRESS"),
		gatewayRootPubKey,
		gatewayRetrievalPubKey,
		gatewayConfig.GetString("GATEWAY_REGION_CODE"),
		gatewayConfig.GetString("NETWORK_INFO_GATEWAY"),
		gatewayConfig.GetString("NETWORK_INFO_PROVIDER"),
		gatewayConfig.GetString("NETWORK_INFO_CLIENT"),
		gatewayConfig.GetString("NETWORK_INFO_ADMIN"),
	)
	err = gwAdmin.InitialiseGateway(gatewayRegistrar, gatewayRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1))
	if err != nil {
		logging.Error("error initialising gateway: %s", err.Error())
		t.FailNow()
	}
	if err := rm.RegisterGateway(gatewayRegistrar); err != nil {
		logging.Error("error registering gateway: %s", err.Error())
		t.FailNow()
	}

	// Add a gateway to the passive list
	added := client.AddGatewaysToUse([]*nodeid.NodeID{gatewayID})
	if !assert.Equal(t, 1, added, "One gateway should be added") {
		t.FailNow()
	}
	// Make the gateway active, this involves doing an establishment
	added = client.AddActiveGateways([]*nodeid.NodeID{gatewayID})
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
	var rm = fcrregistermgr.NewFCRRegisterMgr(gatewayConfig.GetString("REGISTER_API_URL"), true, true, time.Second)
	if err := rm.Start(); err != nil {
		logging.Error("error starting Register Manager: %s", err.Error())
	}
	defer rm.ShutdownAndWait()

	// Configure gateway admin
	gConfBuilder := fcrgatewayadmin.CreateSettings()
	gConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	gConfBuilder.SetRegisterURL(gatewayConfig.GetString("REGISTER_API_URL"))
	gConf := gConfBuilder.Build()
	gwAdmin := fcrgatewayadmin.NewFilecoinRetrievalGatewayAdmin(*gConf)

	// Configure provider admin
	pConfBuilder := fcrprovideradmin.CreateSettings()
	pConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	pConfBuilder.SetRegisterURL(providerConfig.GetString("REGISTER_API_URL"))
	pConf := pConfBuilder.Build()
	pAdmin := fcrprovideradmin.NewFilecoinRetrievalProviderAdmin(*pConf)

	// Configure gateway register
	gatewayRootPubKey, gatewayRetrievalPubKey, gatewayRetrievalPrivateKey, err := generateKeys()
	if err != nil {
		logging.Error("can't generate key pairs %s", err.Error())
		t.FailNow()
	}
	gatewayID := nodeid.NewRandomNodeID()
	gatewayRegistrar := register.NewGatewayRegister(
		gatewayID.ToString(),
		gatewayConfig.GetString("GATEWAY_ADDRESS"),
		gatewayRootPubKey,
		gatewayRetrievalPubKey,
		gatewayConfig.GetString("GATEWAY_REGION_CODE"),
		gatewayConfig.GetString("NETWORK_INFO_GATEWAY"),
		gatewayConfig.GetString("NETWORK_INFO_PROVIDER"),
		gatewayConfig.GetString("NETWORK_INFO_CLIENT"),
		gatewayConfig.GetString("NETWORK_INFO_ADMIN"),
	)
	err = gwAdmin.InitialiseGateway(gatewayRegistrar, gatewayRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1))
	if err != nil {
		logging.Error("error initialising gateway: %s", err.Error())
		t.FailNow()
	}
	if err := rm.RegisterGateway(gatewayRegistrar); err != nil {
		logging.Error("error registering gateway: %s", err.Error())
		t.FailNow()
	}

	// Configure provider register
	providerRootPubKey, providerRetrievalPubKey, providerRetrievalPrivateKey, err := generateKeys()
	if err != nil {
		logging.Error("can't generate key pairs %s", err.Error())
		t.FailNow()
	}
	providerID := nodeid.NewRandomNodeID()
	providerRegistrar := register.NewProviderRegister(
		providerID.ToString(),
		providerConfig.GetString("PROVIDER_ADDRESS"),
		providerRootPubKey,
		providerRetrievalPubKey,
		providerConfig.GetString("PROVIDER_REGION_CODE"),
		providerConfig.GetString("NETWORK_INFO_GATEWAY"),
		providerConfig.GetString("NETWORK_INFO_CLIENT"),
		providerConfig.GetString("NETWORK_INFO_ADMIN"),
	)
	err = pAdmin.InitialiseProvider(providerRegistrar, providerRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1))
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
	clientConfBuilder.SetRegisterURL(gatewayConfig.GetString("REGISTER_API_URL"))
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
	added = client.AddActiveGateways([]*nodeid.NodeID{gatewayID})
	assert.Equal(t, 1, added, "One gateway should be added")

	// Force provider and gateway to update
	err = pAdmin.ForceUpdate(providerID)
	if err != nil {
		logging.Error("error forcing update provider: %s", err.Error())
		t.FailNow()
	}
	err = gwAdmin.ForceUpdate(gatewayID)
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
	err = pAdmin.PublishGroupCID(providerID, pieceCIDs, 42, expiryDate, 42)
	if err != nil {
		logging.Error("error publishing group CID: %s", err.Error())
		t.FailNow()
	}

	// Find published offers
	offers, err := client.FindOffersStandardDiscovery(&(pieceCIDs[0]), gatewayID)
	if err != nil {
		logging.Error("error finding offer via standard discovery: %s", err.Error())
		t.FailNow()
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0.") {
		t.FailNow()
	}

	offers, err = client.FindOffersStandardDiscovery(&(pieceCIDs[1]), gatewayID)
	if err != nil {
		logging.Error("error finding offer via standard discovery: %s", err.Error())
		t.FailNow()
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 1.") {
		t.FailNow()
	}

	offers, err = client.FindOffersStandardDiscovery(&(pieceCIDs[2]), gatewayID)
	if err != nil {
		logging.Error("error finding offer via standard discovery: %s", err.Error())
		t.FailNow()
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 2.") {
		t.FailNow()
	}

	// Negative test using random content ID
	randomCID := cid.NewRandomContentID()
	offers, err = client.FindOffersStandardDiscovery(randomCID, gatewayID)
	if err != nil {
		logging.Error("error finding offer for random content ID via standard discovery: %s", err.Error())
		t.FailNow()
	}
	assert.Equal(t, 0, len(offers), "Shouldn't find any offer with random cid.")

	logging.Info("/*******************************************************/")
	logging.Info("/*        End TestClientStdContentDiscover     	     */")
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
