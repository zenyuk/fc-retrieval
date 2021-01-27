package fcrgatewayadmin

/*
 * Copyright 2020 ConsenSys Software Inc.
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
	"github.com/ConsenSys/fc-retrieval-gateway-admin/internal/control"
	"github.com/ConsenSys/fc-retrieval-gateway-admin/internal/settings"
	log "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)

// FilecoinRetrievalGatewayAdminClient holds information about the interaction of
// the Filecoin Retrieval Gateway Admin Client with Filecoin Retrieval Gateways.
type FilecoinRetrievalGatewayAdminClient struct {
	gatewayManager *control.GatewayManager
	// TODO have a list of gateway objects of all the current gateways being interacted with
}

var singleInstance *FilecoinRetrievalGatewayAdminClient
var initialised = false

// InitFilecoinRetrievalGatewayAdminClient initialise the Filecoin Retreival Client library
func InitFilecoinRetrievalGatewayAdminClient(settings Settings) *FilecoinRetrievalGatewayAdminClient {
	if initialised {
		log.ErrorAndPanic("Attempt to init Filecoin Retrieval Gateway Admin Client a second time")
	}
	var c = FilecoinRetrievalGatewayAdminClient{}
	c.startUp(settings)
	singleInstance = &c
	initialised = true
	return singleInstance

}

// GetFilecoinRetrievalGatewayAdminClient creates a Filecoin Retrieval Gateway Admin Client
func GetFilecoinRetrievalGatewayAdminClient() *FilecoinRetrievalGatewayAdminClient {
	if !initialised {
		log.ErrorAndPanic("Filecoin Retrieval Gateway Admin Client not initialised")
	}

	return singleInstance
}

func (c *FilecoinRetrievalGatewayAdminClient) startUp(conf Settings) {
	log.Info("Filecoin Retrieval Gateway Admin Client started")
	clientSettings := conf.(*settings.ClientGatewayAdminSettings)
	c.gatewayManager = control.GetGatewayManager(*clientSettings)
}

// SetClientReputation requests a Gateway to set a client's reputation to a specified value.
func (c *FilecoinRetrievalGatewayAdminClient) SetClientReputation(clientID *nodeid.NodeID, rep int64) bool {
	log.Info("Filecoin Retrieval Gateway Admin Client: SetClientReputation(clientID: %s, reputation: %d", clientID, rep)
	// TODO
	log.Info("Filecoin Retrieval Gateway Admin Client: SetClientReputation(clientID: %s, reputation: %d) failed to set reputation.", clientID, rep)
	return false
}

// Shutdown releases all resources used by the library
func (c *FilecoinRetrievalGatewayAdminClient) Shutdown() {
	log.Info("Filecoin Retrieval Gateway Admin Client shutting down")
	c.gatewayManager.Shutdown()
}
