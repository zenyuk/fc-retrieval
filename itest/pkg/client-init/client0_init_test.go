package client_init

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
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-itest/pkg/util"
)

// Tests in this file use the Client API, but don't need the rest of the system to be
// configured. These tests need to be run prior to the other client tests.
var fClient *fcrclient.FilecoinRetrievalClient

func TestMain(m *testing.M) {
	// Need to make sure this env is not set in host machine
	itestEnv := os.Getenv("ITEST_CALLING_FROM_CONTAINER")

	if itestEnv != "" {
		// Env is set, we are calling from docker container
		m.Run()
		return
	}
	// Env is not set, we are calling from host
	// We don't need any running instance
	tag := util.GetCurrentBranch()
	network := "itest-shared"
	util.CleanContainers(network)

	// Create shared net
	ctx := context.Background()
	net := *util.CreateNetwork(ctx, network)
	defer net.Remove(ctx)

	// Start itest
	done := make(chan bool)
	itest := *util.StartItest(ctx, tag, network, util.ColorGreen, "./pkg/client-init", done, true)
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

func TestGetClientVersion(t *testing.T) {
	versionInfo := fcrclient.GetVersion()
	// Verify that the client version is an integer number.
	ver, err := strconv.Atoi(versionInfo.Version)
	if err != nil {
		panic(err)
	}

	// The version must be 1 or more.
	if !assert.LessOrEqual(t, 1, ver) {
		return
	}
}

// Test that the client can be initialized without causing an error.
func TestInitClientNoRetrievalKey(t *testing.T) {
	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		panic(err)
	}

	confBuilder := fcrclient.CreateSettings()
	confBuilder.SetEstablishmentTTL(101)
	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	conf := confBuilder.Build()

	fClient = fcrclient.NewFilecoinRetrievalClient(*conf)
}

func TestNoConfiguredGateways(t *testing.T) {
	// The current configuration means that there should only be one connected gateway
	gateways := fClient.GetGatewaysToUse()
	assert.Equal(t, 0, len(gateways), "Unexpected number of gateways returned")
}

func TestUnknownGatewayAdded(t *testing.T) {
	randomGatewayID := nodeid.NewRandomNodeID()
	newGatwaysToBeAdded := make([]*nodeid.NodeID, 0)
	newGatwaysToBeAdded = append(newGatwaysToBeAdded, randomGatewayID)
	numAdded := fClient.AddGatewaysToUse(newGatwaysToBeAdded)
	assert.Equal(t, 0, numAdded)
	gws := fClient.GetGatewaysToUse()
	assert.Equal(t, 0, len(gws))

	// Give the client time to fail the look-up for the random gateway id
	time.Sleep(500 * time.Millisecond)

	gateways := fClient.GetActiveGateways()
	if !assert.Equal(t, 0, len(gateways), "Unexpected number of gateways returned") {
		return
	}
}
