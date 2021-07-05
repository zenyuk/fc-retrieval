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
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"

	"github.com/ConsenSys/fc-retrieval/client/pkg/fcrclient"
	"github.com/ConsenSys/fc-retrieval/common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrregistermgr"
	"github.com/ConsenSys/fc-retrieval/common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval/common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval/common/pkg/register"
	"github.com/ConsenSys/fc-retrieval/gateway-admin/pkg/fcrgatewayadmin"
	"github.com/ConsenSys/fc-retrieval/itest/config"
	cr "github.com/ConsenSys/fc-retrieval/itest/pkg/util/crypto-facade"
	tc "github.com/ConsenSys/fc-retrieval/itest/pkg/util/test-containers"
	"github.com/ConsenSys/fc-retrieval/provider-admin/pkg/fcrprovideradmin"
)

// Test the Provider Admin API.
var providerConfig = config.NewConfig(".env.provider")
var gatewayConfig = config.NewConfig(".env.gateway")
var registerConfig = config.NewConfig(".env.register")
var containers tc.AllContainers

func TestMain(m *testing.M) {
	const testName = "provider-admin"
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
	gConfBuilder := fcrgatewayadmin.CreateSettings()
	gConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	gConfBuilder.SetRegisterURL(registerApiEndpoint)
	gConf := gConfBuilder.Build()
	gwAdmin := fcrgatewayadmin.NewFilecoinRetrievalGatewayAdmin(*gConf)

	// Configure provider admin
	pConfBuilder := fcrprovideradmin.CreateSettings()
	pConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	pConfBuilder.SetRegisterURL(registerApiEndpoint)
	conf := pConfBuilder.Build()
	pAdmin := fcrprovideradmin.NewFilecoinRetrievalProviderAdmin(*conf)

	// Configure client
	clientConfBuilder := fcrclient.CreateSettings()
	clientConfBuilder.SetEstablishmentTTL(101)
	clientConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	clientConfBuilder.SetRegisterURL(registerApiEndpoint)
	clientConf := clientConfBuilder.Build()

	// Create client
	client, err := fcrclient.NewFilecoinRetrievalClient(*clientConf, rm)
	if err != nil {
		logging.Error("error creating retrieval client: %s", err.Error())
		t.FailNow()
	}

	// Initialise gateway
	gatewayRootPubKey, gatewayRetrievalPubKey, gatewayRetrievalPrivateKey, err := cr.GenerateKeys()
	gatewayID, err := nodeid.NewNodeIDFromHexString("ebc134a429ba7dc4811bf64ccb67057f5bd57ca4676800e2f71731cbcc5eb518")
	if err != nil {
		logging.Error("error generating gateway id")
		os.Exit(1)
	}
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
		panic(err)
	}

	if err = rm.RegisterGateway(gatewayRegistrar); err != nil {
		t.Errorf("can't register gateway")
	}
	// Add the gateways to the passive list
	added := client.AddGatewaysToUse([]*nodeid.NodeID{gatewayID})
	if !assert.Equal(t, 1, added, "1 gateway should be added") {
		t.FailNow()
	}
	// Make the gateways active, this involves doing an establishment
	addedActive := client.AddActiveGateways(gatewayClientApiEndpoint, []*nodeid.NodeID{gatewayID})
	if !assert.Equal(t, 1, addedActive, "1 gateway should be activated") {
		t.FailNow()
	}

	// Initialise provider
	providerRootPubKey, providerRetrievalPubKey, providerRetrievalPrivateKey, err := cr.GenerateKeys()
	providerID, err := nodeid.NewNodeIDFromHexString("ebc134a429ba7dc4811bf64ccb67057f5bd57ca4676800e2f71731cbcc5eb518")
	if err != nil {
		logging.Error("can't generate provider ID: %s", err.Error())
		os.Exit(1)
	}
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
	// Initialise the provider using provider admin
	err = pAdmin.InitialiseProvider(providerAdminApiEndpoint, providerRegistrar, providerRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1))
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
	err = pAdmin.ForceUpdate(providerAdminApiEndpoint, providerID)
	if err != nil {
		panic(err)
	}
	err = gwAdmin.ForceUpdate(gatewayAdminApiEndpoint, gatewayID)
	if err != nil {
		panic(err)
	}

	// Publish Group CID
	err = pAdmin.PublishGroupCID(providerAdminApiEndpoint, providerID, pieceCIDs, 42, expiryDate, 42)
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}

	// Get all offers
	var gatewayIDs []nodeid.NodeID
	gatewayIDs = make([]nodeid.NodeID, 0)
	logging.Info("Get all offers")
	_, cidgroupInfo, err := pAdmin.GetGroupCIDOffer(providerAdminApiEndpoint, providerID, gatewayIDs)
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
	_, cidgroupInfo, err = pAdmin.GetGroupCIDOffer(providerAdminApiEndpoint, providerID, gatewayIDs)
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}
	logging.Info("Get offers by real gatewayID=%s: %d", realNodeID.ToString(), len(cidgroupInfo))
	assert.GreaterOrEqual(t, len(cidgroupInfo), 1, "Get offers by gatewayIDs real should be found")

	// Get offers by gatewayIDs fake
	fakeNodeID, _ := nodeid.NewNodeIDFromHexString("101112131415161718191A1B1C1D1E1F202122232425262728292A2B2C2DFA43")
	gatewayIDs[0] = *fakeNodeID
	logging.Info("Get offers by fake gatewayID=%s", fakeNodeID.ToString())
	_, cidgroupInfo, err = pAdmin.GetGroupCIDOffer(providerAdminApiEndpoint, providerID, gatewayIDs)
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
