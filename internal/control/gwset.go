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

// FindGateways find gateways located near to the specified location. Use AddGateways
// to use these gateways.
func (c *ClientManager) FindGateways(location string, maxNumToLocate int) ([]*nodeid.NodeID, error) {
	// Determine gateways to use. For the moment, this is just "use all of them"
	// TODO: This will have to become, use gateways that this client has FIL registered with.
	gateways, err := register.GetRegisteredGateways(c.Settings.RegisterURL())
	if err != nil {
		logging.Error("Error in getting all registered gateways.")
		return nil, err
	}

	res := make([]*nodeid.NodeID, 0)
	count := 0
	for _, info := range gateways {
		if info.GetRegionCode() == location {
			nodeID, err := nodeid.NewNodeIDFromHexString(info.GetNodeID())
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

// AddGatewaysToUse adds one or more gateways to use.
func (c *ClientManager) AddGatewaysToUse(gwNodeIDs []*nodeid.NodeID) int {
	numAdded := 0
	for _, gwToAddID := range gwNodeIDs {
		c.GatewaysToUseLock.RLock()
		_, exist := c.GatewaysToUse[gwToAddID.ToString()]
		c.GatewaysToUseLock.RUnlock()
		if exist {
			continue
		}
		gateway, err := register.GetGatewayByID(c.Settings.RegisterURL(), gwToAddID)
		if err != nil {
			logging.Error("Error getting registered gateway %v: %v", gwToAddID, err.Error())
			continue
		}
		if !validateGatewayInfo(&gateway) {
			logging.Error("Register info not valid.")
			continue
		}
		// Success
		c.GatewaysToUseLock.Lock()
		c.GatewaysToUse[gwToAddID.ToString()] = gateway
		c.GatewaysToUseLock.Unlock()
		numAdded++
	}
	return numAdded
}

// RemoveGatewaysToUse removes one or more gateways from the list of Gateways to use.
// This also removes the gateway from gateways in active map.
func (c *ClientManager) RemoveGatewaysToUse(gwNodeIDs []*nodeid.NodeID) int {
	c.GatewaysToUseLock.Lock()
	defer c.GatewaysToUseLock.Unlock()

	numRemoved := 0
	for _, gwToRemoveID := range gwNodeIDs {
		_, exist := c.GatewaysToUse[gwToRemoveID.ToString()]
		if exist {
			delete(c.GatewaysToUse, gwToRemoveID.ToString())
			numRemoved++
			c.ActiveGatewaysLock.Lock()
			delete(c.ActiveGateways, gwToRemoveID.ToString())
			c.ActiveGatewaysLock.Unlock()
		}
	}

	return numRemoved
}

// RemoveAllGatewaysToUse removes all gateways from the list of Gateways.
// This also cleared all gateways in active
func (c *ClientManager) RemoveAllGatewaysToUse() int {
	c.GatewaysToUseLock.Lock()
	defer c.GatewaysToUseLock.Unlock()
	c.ActiveGatewaysLock.Lock()
	defer c.ActiveGatewaysLock.Unlock()

	numRemoved := len(c.GatewaysToUse)
	c.GatewaysToUse = make(map[string]register.GatewayRegister)
	c.ActiveGateways = make(map[string]register.GatewayRegister)

	return numRemoved
}

// GetGatewaysToUse returns the list of gateways to use.
func (c *ClientManager) GetGatewaysToUse() []*nodeid.NodeID {
	c.GatewaysToUseLock.RLock()
	defer c.GatewaysToUseLock.RUnlock()

	res := make([]*nodeid.NodeID, 0)
	for key := range c.GatewaysToUse {
		nodeID, err := nodeid.NewNodeIDFromHexString(key)
		if err != nil {
			logging.Error("Error in generating node id.")
			continue
		}
		res = append(res, nodeID)
	}

	return res
}

// AddActiveGateways adds one or more gateways to active gateway map.
// Returns the number of gateways added.
func (c *ClientManager) AddActiveGateways(gwNodeIDs []*nodeid.NodeID) int {
	numAdded := 0
	for _, gwToAddID := range gwNodeIDs {
		c.ActiveGatewaysLock.RLock()
		_, exist := c.ActiveGateways[gwToAddID.ToString()]
		c.ActiveGatewaysLock.RUnlock()
		if exist {
			continue
		}
		c.GatewaysToUseLock.RLock()
		info, exist := c.GatewaysToUse[strings.ToLower(gwToAddID.ToString())]
		c.GatewaysToUseLock.RUnlock()
		if !exist {
			logging.Error("Given node id: %v does not exist in gateways to use map, consider add the gateway first.", gwToAddID.ToString())
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
		c.ActiveGatewaysLock.Lock()
		c.ActiveGateways[gwToAddID.ToString()] = info
		c.ActiveGatewaysLock.Unlock()
		numAdded++
	}
	return numAdded
}

// RemoveActiveGateways removes one or more gateways from the list of Gateways in active.
func (c *ClientManager) RemoveActiveGateways(gwNodeIDs []*nodeid.NodeID) int {
	c.ActiveGatewaysLock.Lock()
	defer c.ActiveGatewaysLock.Unlock()

	numRemoved := 0
	for _, gwToRemoveID := range gwNodeIDs {
		_, exist := c.ActiveGateways[gwToRemoveID.ToString()]
		if exist {
			delete(c.ActiveGateways, gwToRemoveID.ToString())
			numRemoved++
		}
	}

	return numRemoved
}

// RemoveAllActiveGateways removes all gateways from the list of Gateways in active.
func (c *ClientManager) RemoveAllActiveGateways() int {
	c.ActiveGatewaysLock.Lock()
	defer c.ActiveGatewaysLock.Unlock()

	numRemoved := len(c.ActiveGateways)
	c.ActiveGateways = make(map[string]register.GatewayRegister)

	return numRemoved
}

// GetActiveGateways returns the list of gateways that are active.
func (c *ClientManager) GetActiveGateways() []*nodeid.NodeID {
	c.ActiveGatewaysLock.RLock()
	defer c.ActiveGatewaysLock.RUnlock()

	res := make([]*nodeid.NodeID, 0)
	for key := range c.ActiveGateways {
		nodeID, err := nodeid.NewNodeIDFromHexString(key)
		if err != nil {
			logging.Error("Error in generating node id.")
			continue
		}
		res = append(res, nodeID)
	}

	return res
}
