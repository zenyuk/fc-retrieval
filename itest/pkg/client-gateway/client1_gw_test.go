/*
Package client_gateway - integration tests, specific to Retrieval Gateways
*/
package client_gateway

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
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"

	tc "github.com/ConsenSys/fc-retrieval/itest/pkg/util/test-containers"

	"github.com/stretchr/testify/assert"

	"github.com/ConsenSys/fc-retrieval/client/pkg/fcrclient"
	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrregistermgr"
	"github.com/ConsenSys/fc-retrieval/common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval/common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval/common/pkg/register"
	"github.com/ConsenSys/fc-retrieval/gateway-admin/pkg/fcrgatewayadmin"
	"github.com/ConsenSys/fc-retrieval/itest/config"
)

var containers tc.AllContainers

func TestMain(m *testing.M) {
	const testName = "client-gateway"
	ctx := context.Background()
	var gatewayConfig = config.NewConfig(".env.gateway")
	var providerConfig = config.NewConfig(".env.provider")
	var registerConfig = config.NewConfig(".env.register")
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

func TestOneGateway(t *testing.T) {
	gatewayConfig := config.NewConfig(".env.gateway")

	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		logging.Panic(err.Error())
	}

	confBuilder := fcrgatewayadmin.CreateSettings()
	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	confBuilder.SetRegisterURL("http://" + containers.Register.GetRegisterHostApiEndpoint())
	conf := confBuilder.Build()
	gwAdmin := fcrgatewayadmin.NewFilecoinRetrievalGatewayAdmin(*conf)

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

	var rm = fcrregistermgr.NewFCRRegisterMgr(conf.RegisterURL(), false, true, time.Second)
	if err := rm.Start(); err != nil {
		logging.Error("error starting Register Manager: %s", err.Error())
	}

	gatewayID := nodeid.NewRandomNodeID()
	gatewayName := "gateway-0"
	_, _, gatewayClientApiEndpoint, gatewayAdminApiEndpoint := containers.Gateways[gatewayName].GetGatewayHostApiEndpoints()
	gatewayRegistrar := register.NewGatewayRegister(
		gatewayID.ToString(),
		gatewayConfig.GetString("GATEWAY_ADDRESS"),
		gatewayRootSigningKey,
		gatewayRetrievalSigningKey,
		gatewayConfig.GetString("GATEWAY_REGION_CODE"),
		gatewayName+":"+gatewayConfig.GetString("BIND_GATEWAY_API"),
		gatewayName+":"+gatewayConfig.GetString("BIND_PROVIDER_API"),
		gatewayName+":"+gatewayConfig.GetString("BIND_REST_API"),
		gatewayName+":"+gatewayConfig.GetString("BIND_ADMIN_API"),
	)

	if err = gwAdmin.InitialiseGateway(gatewayAdminApiEndpoint, gatewayRegistrar, gatewayRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1)); err != nil {
		logging.Error("gateway initialisation error: %s", err.Error())
		t.FailNow()
	}

	if err = rm.RegisterGateway(gatewayRegistrar); err != nil {
		logging.Error("gateway registration error: %s", err.Error())
		t.FailNow()
	}

	logging.Info("Adding to client config gateway: %s", gatewayID.ToString())

	blockchainPrivateKey, err = fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		panic(err)
	}

	clientConfBuilder := fcrclient.CreateSettings()
	clientConfBuilder.SetEstablishmentTTL(101)
	clientConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	clientConf := clientConfBuilder.Build()

	client, err := fcrclient.NewFilecoinRetrievalClient(*clientConf, rm)
	assert.Nil(t, err)
	newGatewaysToBeAdded := make([]*nodeid.NodeID, 0)
	newGatewaysToBeAdded = append(newGatewaysToBeAdded, gatewayID)
	numAdded := client.AddGatewaysToUse(newGatewaysToBeAdded)
	assert.Equal(t, 1, numAdded, "expecting the new Gateway be added to the list of gateways this client can potentially use")
	gws := client.GetGatewaysToUse()
	assert.Equal(t, 1, len(gws))

	numAdded = client.AddActiveGateways(gatewayClientApiEndpoint, newGatewaysToBeAdded)
	assert.Equal(t, 1, numAdded, "expecting the new Gateway be added to the list of gateways in use")
	ga := client.GetActiveGateways()
	assert.Equal(t, 1, len(ga))
}
