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

	"github.com/stretchr/testify/assert"

	"github.com/ConsenSys/fc-retrieval-client/pkg/fcrclient"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/register"
	"github.com/ConsenSys/fc-retrieval-gateway-admin/pkg/fcrgatewayadmin"

	"github.com/ConsenSys/fc-retrieval-itest/config"
	"github.com/ConsenSys/fc-retrieval-itest/pkg/util"
)

func TestMain(m *testing.M) {
	// Need to make sure this env is not set in host machine
	itestEnv := os.Getenv("ITEST_CALLING_FROM_CONTAINER")

	if itestEnv != "" {
		// Env is set, we are calling from docker container
		m.Run()
		return
	}
	// Env is not set, we are calling from host
	// We need a redis, a register and a gateway
	tag := util.GetCurrentBranch()
	network := "itest-shared"
	util.CleanContainers(network)

	// Get env
	rgEnv := util.GetEnvMap("../../.env.register")
	gwEnv := util.GetEnvMap("../../.env.gateway")

	// Create shared net
	ctx := context.Background()
	net := *util.CreateNetwork(ctx, network)
	defer net.Remove(ctx)

	// Start redis
	redis := *util.StartRedis(ctx, network, true)
	defer redis.Terminate(ctx)
	defer redis.StopLogProducer()

	// Start register
	register := *util.StartRegister(ctx, tag, network, util.ColorYellow, rgEnv, true)
	defer register.Terminate(ctx)
	defer register.StopLogProducer()

	// Start gateway
	gateway := *util.StartGateway(ctx, "gateway", tag, network, util.ColorBlue, gwEnv, true)
	defer gateway.Terminate(ctx)
	defer gateway.StopLogProducer()

	// Start itest
	done := make(chan bool)
	itest := *util.StartItest(ctx, tag, network, util.ColorGreen, "./pkg/client-gateway", done, true)
	defer itest.Terminate(ctx)
	defer itest.StopLogProducer()

	// Block until done.
	if <-done {
		logging.Info("Tests passed, shutdown...")
	} else {
		logging.Fatal("Tests failed, shutdown...")
	}
	// Clean containers to shutdown
	util.CleanContainers(network)
}

func TestOneGateway(t *testing.T) {
	gatewayConfig := config.NewConfig(".env.gateway")

	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		logging.Panic(err.Error())
	}

	confBuilder := fcrgatewayadmin.CreateSettings()
	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	confBuilder.SetRegisterURL("http://register:9020")
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

	gatewayID := nodeid.NewRandomNodeID()

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

	logging.Info("Adding to client config gateway: %s", gatewayID.ToString())

	blockchainPrivateKey, err = fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		panic(err)
	}

	clientConfBuilder := fcrclient.CreateSettings()
	clientConfBuilder.SetEstablishmentTTL(101)
	clientConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	clientConf := clientConfBuilder.Build()

	client := fcrclient.NewFilecoinRetrievalClient(*clientConf)
	newGatewaysToBeAdded := make([]*nodeid.NodeID, 0)
	newGatewaysToBeAdded = append(newGatewaysToBeAdded, gatewayID)
	numAdded := client.AddGatewaysToUse(newGatewaysToBeAdded)
	assert.Equal(t, 1, numAdded)
	gws := client.GetGatewaysToUse()
	assert.Equal(t, 1, len(gws))

	numAdded = client.AddActiveGateways(newGatewaysToBeAdded)
	assert.Equal(t, 1, numAdded)
	ga := client.GetActiveGateways()
	assert.Equal(t, 1, len(ga))
}
