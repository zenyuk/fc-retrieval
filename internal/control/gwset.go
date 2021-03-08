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
	"crypto/rand"
	"strings"
	"time"

	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-register/pkg/register"
)

// RefreshLatestProviderInfo refreshs the registered provider information from the registry.
func (c *ClientManager) RefreshLatestProviderInfo() {
	providers, err := register.GetRegisteredProviders(c.Settings.RegisterURL())
	if err != nil {
		logging.Error("Error in getting registered providers: %v", err.Error())
		return
	}
	c.ProvidersLock.RLock()
	update := false
	if len(providers) != len(c.Providers) {
		update = true
	} else {
		for _, provider := range providers {
			if !validateProviderInfo(&provider) {
				logging.Error("Register info not valid.")
				continue
			}
			storedInfo, exist := c.Providers[strings.ToLower(provider.NodeID)]
			if !exist {
				update = true
				break
			} else {
				// TODO, There is a potential bug here, if error not checked right after it appears, panic occurs.
				// Needs to fix also in gateway & provider go routine.
				key, err := storedInfo.GetRootSigningKey()
				rootSigningKey, err2 := key.EncodePublicKey()
				key, err3 := storedInfo.GetSigningKey()
				signingKey, err4 := key.EncodePublicKey()
				if err != nil || err2 != nil || err3 != nil || err4 != nil {
					logging.Error("Error in generating key string")
					break
				}
				if provider.Address != storedInfo.GetAddress() ||
					provider.NetworkInfoAdmin != storedInfo.GetNetworkInfoAdmin() ||
					provider.NetworkInfoClient != storedInfo.GetNetworkInfoClient() ||
					provider.NetworkInfoGateway != storedInfo.GetNetworkInfoGateway() ||
					provider.RegionCode != storedInfo.GetRegionCode() ||
					provider.RootSigningKey != rootSigningKey ||
					provider.SigningKey != signingKey {
					update = true
					break
				}
			}
		}
	}
	c.ProvidersLock.RUnlock()
	if update {
		c.ProvidersLock.Lock()
		c.Providers = make(map[string]register.ProviderRegister)
		for _, provider := range providers {
			logging.Info("Add to registered providers map: nodeID=%+v", provider.NodeID)
			c.Providers[strings.ToLower(provider.NodeID)] = provider
		}
		c.ProvidersLock.Unlock()
	}

}

// RefreshLatestGatewayInfo refreshs the registered gateway information from the registry.
func (c *ClientManager) RefreshLatestGatewayInfo() {
	gateways, err := register.GetRegisteredGateways(c.Settings.RegisterURL())
	if err != nil {
		logging.Error("Error in getting registered gateways: %v", err.Error())
		return
	}
	c.GatewaysLock.RLock()
	update := false
	if len(gateways) != len(c.Gateways) {
		update = true
	} else {
		for _, gateway := range gateways {
			if !validateGatewayInfo(&gateway) {
				logging.Error("Register info not valid.")
				continue
			}
			storedInfo, exist := c.Gateways[strings.ToLower(gateway.NodeID)]
			if !exist {
				update = true
				break
			} else {
				// TODO, There is a potential bug here, if error not checked right after it appears, panic occurs.
				// Needs to fix also in gateway & provider go routine.
				key, err := storedInfo.GetRootSigningKey()
				rootSigningKey, err2 := key.EncodePublicKey()
				key, err3 := storedInfo.GetSigningKey()
				signingKey, err4 := key.EncodePublicKey()
				if err != nil || err2 != nil || err3 != nil || err4 != nil {
					logging.Error("Error in generating key string")
					break
				}
				if gateway.Address != storedInfo.GetAddress() ||
					gateway.NetworkInfoAdmin != storedInfo.GetNetworkInfoAdmin() ||
					gateway.NetworkInfoClient != storedInfo.GetNetworkInfoClient() ||
					gateway.NetworkInfoProvider != storedInfo.GetNetworkInfoProvider() ||
					gateway.NetworkInfoGateway != storedInfo.GetNetworkInfoGateway() ||
					gateway.RegionCode != storedInfo.GetRegionCode() ||
					gateway.RootSigningKey != rootSigningKey ||
					gateway.SigningKey != signingKey {
					update = true
					break
				}
			}
		}
	}
	c.GatewaysLock.RUnlock()
	if update {
		c.GatewaysLock.Lock()
		c.Gateways = make(map[string]register.GatewayRegister)
		for _, gateway := range gateways {
			logging.Info("Add to registered gateways map: nodeID=%+v", gateway.NodeID)
			c.Gateways[strings.ToLower(gateway.NodeID)] = gateway
		}
		// Need to check if any gateway in use should be removed.
		c.GatewaysInUseLock.Lock()
		for gw := range c.GatewaysInUse {
			if _, ok := c.Gateways[gw]; !ok {
				// Need to remove this gateway in use.
				delete(c.GatewaysInUse, gw)
			}
		}
		c.GatewaysInUseLock.Unlock()
		c.GatewaysLock.Unlock()
	}
}

// FindGateways find gateways located near to the specified location. Use AddGateways
// to use these gateways.
func (c *ClientManager) FindGateways(location string, maxNumToLocate int) ([]*nodeid.NodeID, error) {
	// Determine gateways to use. For the moment, this is just "use all of them"
	// TODO: This will have to become, use gateways that this client has FIL registered with.
	c.GatewaysLock.RLock()
	defer c.GatewaysLock.RUnlock()

	res := make([]*nodeid.NodeID, 0)
	count := 0
	for _, info := range c.Gateways {
		if info.GetRegionCode() == location {
			nodeID, err := nodeid.NewNodeIDFromString(info.GetNodeID())
			if err != nil {
				logging.Error("Error in generating node id, skipping: %v", info.GetNodeID())
				continue
			}
			res = append(res, nodeID)
			count++
			if count >= maxNumToLocate {
				break
			}
		}
	}
	return res, nil
}

// AddGateways adds one or more gateways to use.
// Returns the number of gateways added.
func (c *ClientManager) AddGateways(gwNodeIDs []*nodeid.NodeID) int {
	numAdded := 0
	for _, gwToAddID := range gwNodeIDs {
		c.GatewaysInUseLock.RLock()
		_, exist := c.GatewaysInUse[strings.ToLower(gwToAddID.ToString())] //TODO, Need to check if nodeID.tostring() returns lower case result
		c.GatewaysInUseLock.RUnlock()
		if exist {
			continue
		}
		c.GatewaysLock.RLock()
		info, exist := c.Gateways[strings.ToLower(gwToAddID.ToString())]
		c.GatewaysLock.RUnlock()
		if !exist {
			logging.Error("Given node id: %v does not exist in local cached registered map, consider refresh the map.", gwToAddID.ToString())
			continue
		}
		// Attempt an establishment
		challenge := make([]byte, 32)
		rand.Read(challenge)
		ttl := time.Now().Unix() + c.Settings.EstablishmentTTL()
		err := c.GatewayClientEstablishment(&info, challenge, c.Settings.ClientID(), ttl)
		if err != nil {
			logging.Error("Error in initial establishment: %v", err.Error())
			continue
		}
		// It is success
		c.GatewaysInUseLock.Lock()
		c.GatewaysInUse[strings.ToLower(gwToAddID.ToString())] = info
		c.GatewaysInUseLock.Unlock()
		numAdded++
	}
	return numAdded
}

// RemoveGateways removes one or more gateways from the list of Gateways to use.
func (c *ClientManager) RemoveGateways(gwNodeIDs []*nodeid.NodeID) int {
	c.GatewaysInUseLock.Lock()
	defer c.GatewaysInUseLock.Unlock()

	numRemoved := 0
	for _, gwToRemoveID := range gwNodeIDs {
		_, exist := c.GatewaysInUse[strings.ToLower(gwToRemoveID.ToString())]
		if exist {
			delete(c.GatewaysInUse, strings.ToLower(gwToRemoveID.ToString()))
			numRemoved++
		}
	}

	return numRemoved
}

// RemoveAllGateways removes all gateways from the list of Gateways to use.
func (c *ClientManager) RemoveAllGateways() int {
	c.GatewaysInUseLock.Lock()
	defer c.GatewaysInUseLock.Unlock()

	numRemoved := len(c.GatewaysInUse)
	c.GatewaysInUse = make(map[string]register.GatewayRegister)

	return numRemoved
}

// GetGateways returns the list of gateways that are being used.
func (c *ClientManager) GetGateways() []*nodeid.NodeID {
	c.GatewaysInUseLock.RLock()
	defer c.GatewaysInUseLock.RUnlock()

	res := make([]*nodeid.NodeID, 0)
	for key := range c.GatewaysInUse {
		nodeID, err := nodeid.NewNodeIDFromString(key)
		if err != nil {
			logging.Error("Error in generating node id.")
			continue
		}
		res = append(res, nodeID)
	}

	return res
}
