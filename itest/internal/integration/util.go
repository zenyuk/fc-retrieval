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
	"github.com/ConsenSys/fc-retrieval-client/pkg/fcrclient"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-provider/pkg/provider"

	"github.com/ConsenSys/fc-retrieval-itest/config"
)

// InitClient initialises a Filecoin Retrieval Client
func InitClient() *fcrclient.FilecoinRetrievalClient {
	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		panic(err)
	}

	confBuilder := fcrclient.CreateSettings()
	confBuilder.SetEstablishmentTTL(101)
	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	conf := confBuilder.Build()

	c, err := fcrclient.NewFilecoinRetrievalClient(*conf)
	if err != nil {
		panic(err)
	}

	return c
}

// CloseClient shuts down a Filecoin Retrieval Client
func CloseClient(client *fcrclient.FilecoinRetrievalClient) {
	client.Shutdown()
}

// InitProvider initialises a Filecoin Retrieval Provider
func InitProvider() *provider.Provider {
	conf := config.NewConfig(".env.provider")
	logging.Init(conf)

	p := provider.GetSingleInstance(conf)
	return p
}

// CloseProvider shuts down a Filecoin Retrieval Provider
func CloseProvider(provider *provider.Provider) {
	// provider.Shutdown()
}
