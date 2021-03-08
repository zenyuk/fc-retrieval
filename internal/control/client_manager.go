package control

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
	"sync"

	"github.com/ConsenSys/fc-retrieval-register/pkg/register"

	"github.com/ConsenSys/fc-retrieval-client/internal/settings"
)

// ClientManager manages the pool of gateways and the connections to them.
type ClientManager struct {
	Settings settings.ClientSettings

	Gateways     map[string]register.GatewayRegister
	GatewaysLock sync.RWMutex

	Providers     map[string]register.ProviderRegister
	ProvidersLock sync.RWMutex

	// List of gateway to use. A client may request a node be added to this list e
	GatewaysInUse     map[string]register.GatewayRegister
	GatewaysInUseLock sync.RWMutex
}

// NewClientManager returns an initialised instance of the client manager.
func NewClientManager(settings settings.ClientSettings) *ClientManager {
	return &ClientManager{
		Settings:          settings,
		Gateways:          make(map[string]register.GatewayRegister),
		GatewaysLock:      sync.RWMutex{},
		Providers:         make(map[string]register.ProviderRegister),
		ProvidersLock:     sync.RWMutex{},
		GatewaysInUse:     make(map[string]register.GatewayRegister),
		GatewaysInUseLock: sync.RWMutex{},
	}
}

// FindOffersStandardDiscovery finds offers using the standard discovery mechanism.
// func (c *ClientManager) FindOffersStandardDiscovery(contentID *cid.ContentID) ([]cidoffer.CidGroupOffer, error) {

// 	g.gatewaysLock.RLock()
// 	gatewaysSnapshot := g.gateways
// 	g.gatewaysLock.RUnlock()

// 	if len(gatewaysSnapshot) == 0 {
// 		return nil, fmt.Errorf("No gateways available")
// 	}

// 	var aggregateOffers []cidoffer.CidGroupOffer
// 	for _, gw := range gatewaysSnapshot {
// 		// TODO need to do nonce management
// 		// TODO need to do requests to all gateways in parallel, rather than serially
// 		offers, err := gw.comms.GatewayStdCIDDiscovery(contentID, 1)
// 		if err != nil {
// 			logging.Warn("GatewayStdDiscovery error. Gateway: %s, Error: %s", gw.info.NodeID, err)
// 		}
// 		// TODO: probably should remove duplicate offers at this point
// 		aggregateOffers = append(aggregateOffers, offers...)
// 	}
// 	return aggregateOffers, nil
// }

// GetConnectedGateways returns the list of domain names of gateways that the client
// is currently connected to.
// func (c *GatewayManager) GetConnectedGateways() []string {
// 	urls := make([]string, len(g.gateways))
// 	for i, gateway := range g.gateways {
// 		urls[i] = gateway.comms.ApiURL
// 	}
// 	return urls
// }

// func (g *GatewayManager) addGateway(nodeID *nodeid.NodeID) {
// 	// TODO add gateway by ID
// 	gws, err := register.GetRegisteredGateways(g.settings.RegisterURL())
// 	if err != nil {
// 		logging.Error("Unable to get registered gateways: %v", err)
// 		return
// 	}
// 	logging.Info("Register returned %d gateways", len(gws))
// 	if len(gws) == 0 {
// 		logging.Warn("Unable to get registered gateways: %v", err)
// 		return
// 	}
// 	if len(gws) != 1 {
// 		logging.Warn("Unexpectedly, multiple gateways returned: %d", len(gws))
// 		return
// 	}
// 	gw := gws[0]
// 	gatewayID, err := nodeid.NewNodeIDFromString(gw.NodeID)

// 	if !g.validateGatewayInfo(&gw) {
// 		logging.Warn("Gateway registration information for gateway (%s) is invalid. Ignoring.", gatewayID)
// 		return
// 	}

// 	logging.Info("Setting-up comms with: %+v", gw)
// 	comms, err := gatewayapi.NewGatewayAPIComms(&gw, &g.settings)
// 	if err != nil {
// 		logging.Error("Error encountered which contacting gateway (%s): %+v", gatewayID, err)
// 		return
// 	}

// 	// TODO this should only be done for new gateways.
// 	// Try to do the establishment with the new gateway
// 	var challenge [32]byte
// 	fcrcrypto.GeneratePublicRandomBytes(challenge[:])
// 	comms.GatewayClientEstablishment(challenge)
// 	if err != nil {
// 		logging.Warn("Error processing node id: %+v", err)
// 		return
// 	}

// 	activeGateway := ActiveGateway{gw, comms, gatewayID}
// 	g.gatewaysLock.RLock()
// 	g.gateways = append(g.gateways, activeGateway)
// 	g.gatewaysLock.RUnlock()

// 	if len(g.gateways) > 0 {
// 		logging.Info("Gateway Manager using %d gateways", len(g.gateways))
// 	} else {
// 		logging.Warn("No gateways available")
// 	}
// }
