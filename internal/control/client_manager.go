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
	"math/rand"
	"sync"
	"time"

	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-register/pkg/register"

	"github.com/ConsenSys/fc-retrieval-client/internal/settings"
)

// ClientManager manages the pool of gateways and the connections to them.
type ClientManager struct {
	Settings settings.ClientSettings

	// List of gateways this client can potentially use
	GatewaysToUse     map[string]register.GatewayRegister
	GatewaysToUseLock sync.RWMutex

	// List of gateway in use. A client may request a node be added to this list
	ActiveGateways     map[string]register.GatewayRegister
	ActiveGatewaysLock sync.RWMutex
}

// NewClientManager returns an initialised instance of the client manager.
func NewClientManager(settings settings.ClientSettings) *ClientManager {
	return &ClientManager{
		Settings:           settings,
		GatewaysToUse:      make(map[string]register.GatewayRegister),
		GatewaysToUseLock:  sync.RWMutex{},
		ActiveGateways:     make(map[string]register.GatewayRegister),
		ActiveGatewaysLock: sync.RWMutex{},
	}
}

// FindOffersStandardDiscovery finds offers using the standard discovery mechanism.
func (c *ClientManager) FindOffersStandardDiscovery(contentID *cid.ContentID) ([]cidoffer.SubCIDOffer, error) {
	c.ActiveGatewaysLock.RLock()
	defer c.ActiveGatewaysLock.RUnlock()

	aggregateOffers := make([]cidoffer.SubCIDOffer, 0)
	for _, gw := range c.ActiveGateways {
		// TODO need to do nonce management
		// TODO need to do requests to all gateways in parallel, rather than serially
		offers, err := c.GatewayStdCIDDiscovery(&gw, contentID, rand.Int63(), time.Now().Unix()+c.Settings.EstablishmentTTL())
		if err != nil {
			logging.Warn("GatewayStdDiscovery error. Gateway: %s, Error: %s", gw.NodeID, err)
			continue
		}
		// TODO: probably should remove duplicate offers at this point
		aggregateOffers = append(aggregateOffers, offers...)
	}
	return aggregateOffers, nil
}
