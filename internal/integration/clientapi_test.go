package integration

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
	"strconv"
	"testing"

	"github.com/ConsenSys/fc-retrieval-client/pkg/fcrclient"
	// "github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	"github.com/stretchr/testify/assert"
)

// Test the Client API.

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

// func TestInitClientNoRetrievalKey(t *testing.T) {
// 	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
// 	if err != nil {
// 		panic(err)
// 	}

// 	confBuilder := fcrclient.CreateSettings()
// 	confBuilder.SetEstablishmentTTL(101)
// 	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
// 	conf := confBuilder.Build()

// 	client, err := fcrclient.NewFilecoinRetrievalClient(*conf)
// 	if err != nil {
// 		panic(err)
// 	}

// 	client.Shutdown()
// }

// func TestOneConnectedGateway(t *testing.T) {
// 	// The current configuration means that there should only be one connected gateway
// 	client := InitClient()
// 	gateways := client.ConnectedGateways()
// 	assert.Equal(t, 1, len(gateways), "Unexpected number of gateways returned")
// 	if len(gateways) > 0 {
// 		assert.Equal(t, "http://gateway:9011/", gateways[0])
// 	}
// 	CloseClient(client)

// }
