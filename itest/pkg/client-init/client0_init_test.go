/*
Package client_init - integration and end-to-end tests, specific to Retrieval Client
*/
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

	"github.com/testcontainers/testcontainers-go"

	"github.com/ConsenSys/fc-retrieval/itest/config"
	tc "github.com/ConsenSys/fc-retrieval/itest/pkg/util/test-containers"

	"github.com/stretchr/testify/assert"

	"github.com/ConsenSys/fc-retrieval/client/pkg/fcrclient"
	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrregistermgr"
	"github.com/ConsenSys/fc-retrieval/common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval/common/pkg/nodeid"
)

// Tests in this file use the Client API, but don't need the rest of the system to be
// configured. These tests need to be run prior to the other client tests.
var fClient *fcrclient.FilecoinRetrievalClient

var containers tc.AllContainers

func TestMain(m *testing.M) {
	const testName = "client-init"
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

func TestGetClientVersion(t *testing.T) {
	versionInfo := fcrclient.GetVersion()
	// Verify that the client version is an integer number.
	ver, err := strconv.Atoi(versionInfo.Version)
	if err != nil {
		panic(err)
	}

	// The version must be 1 or more.
	assert.LessOrEqual(t, 1, ver)
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

	var rm = fcrregistermgr.NewFCRRegisterMgr(conf.RegisterURL(), false, false, 10*time.Second)
	fClient, err = fcrclient.NewFilecoinRetrievalClient(*conf, rm)
	assert.Nil(t, err)
}

func TestNoConfiguredGateways(t *testing.T) {
	// The current configuration means that there should only be one connected gateway
	gateways := fClient.GetGatewaysToUse()
	assert.Equal(t, 0, len(gateways), "Unexpected number of gateways returned")
}

func TestUnknownGatewayAdded(t *testing.T) {
	var rm = fcrregistermgr.NewFCRRegisterMgr("http://localhost/fakeurl", false, true, 10*time.Second)
	if err := rm.Start(); err != nil {
		logging.Error("error starting Register Manager: %s", err.Error())
	}

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
	assert.Equal(t, 0, len(gateways), "Unexpected number of gateways returned")
}
