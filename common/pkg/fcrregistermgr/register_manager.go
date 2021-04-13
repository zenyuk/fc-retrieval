package fcrregistermgr

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
	"errors"
	"sync"
	"time"

	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-register/pkg/register"
)

// Register Manager manages the internal storage of registered nodes.
type FCRRegisterMgr struct {
	// Boolean indicates if the manager has started
	start bool

	// API for the register
	registerAPI string

	// Duration to wait between two updates
	refreshDuration time.Duration

	// Boolean indicates if to discover gateway/provider
	gatewayDiscv  bool
	providerDiscv bool

	// Channels to control the threads
	gatewayShutdownCh  chan bool
	providerShutdownCh chan bool
	gatewayRefreshCh   chan bool
	providerRefreshCh  chan bool

	// registeredGatewaysMap stores mapping from gateway id (big int in string repr) to its registration info
	registeredGatewaysMap     map[string]register.GatewayRegister
	registeredGatewaysMapLock sync.RWMutex

	// registeredProvidersMap stores mapping from provider id (big int in string repr) to its registration info
	registeredProvidersMap     map[string]register.ProviderRegister
	registeredProvidersMapLock sync.RWMutex
}

// NewFCRRegisterMgr creates a new register manager.
func NewFCRRegisterMgr(registerAPI string, providerDiscv bool, gatewayDiscv bool, refreshDuration time.Duration) *FCRRegisterMgr {
	res := &FCRRegisterMgr{
		start:           false,
		registerAPI:     registerAPI,
		refreshDuration: refreshDuration,
		gatewayDiscv:    gatewayDiscv,
		providerDiscv:   providerDiscv,
	}
	if gatewayDiscv {
		res.registeredGatewaysMap = make(map[string]register.GatewayRegister)
		res.registeredGatewaysMapLock = sync.RWMutex{}
		res.gatewayShutdownCh = make(chan bool)
		res.gatewayRefreshCh = make(chan bool)
	}
	if providerDiscv {
		res.registeredProvidersMap = make(map[string]register.ProviderRegister)
		res.registeredProvidersMapLock = sync.RWMutex{}
		res.providerShutdownCh = make(chan bool)
		res.providerRefreshCh = make(chan bool)
	}
	return res
}

// Start starts a thread to auto update the internal map every given duration.
func (mgr *FCRRegisterMgr) Start() error {
	if mgr.start {
		return errors.New("Manager has already started.")
	}
	mgr.start = true
	if mgr.gatewayDiscv {
		go mgr.updateGateways()
	}
	if mgr.providerDiscv {
		go mgr.updateProviders()
	}
	return nil
}

// Shutdown will shutdown the register manager.
func (mgr *FCRRegisterMgr) Shutdown() {
	if !mgr.start {
		return
	}
	if mgr.gatewayDiscv {
		mgr.gatewayShutdownCh <- true
	}
	if mgr.providerDiscv {
		mgr.providerShutdownCh <- true
	}
	mgr.start = false
}

// Refresh refreshs the internal map immediately.
func (mgr *FCRRegisterMgr) Refresh() {
	if !mgr.start {
		return
	}
	if mgr.gatewayDiscv {
		mgr.gatewayRefreshCh <- true
	}
	if mgr.providerDiscv {
		mgr.providerRefreshCh <- true
	}
}

// GetGateway returns a gateway register if found.
func (mgr *FCRRegisterMgr) GetGateway(id *nodeid.NodeID) *register.GatewayRegister {
	if !mgr.start || !mgr.gatewayDiscv {
		return nil
	}
	mgr.registeredGatewaysMapLock.RLock()
	gateway, ok := mgr.registeredGatewaysMap[id.ToString()]
	if !ok {
		mgr.registeredGatewaysMapLock.RUnlock()
		// TODO: Do we call refresh here, if can't find a gateway?
		// mgr.Refresh()
		mgr.registeredGatewaysMapLock.RLock()
		gateway = mgr.registeredGatewaysMap[id.ToString()]
	}
	defer mgr.registeredGatewaysMapLock.RUnlock()
	// Return the pointer of a copy of the register
	res := register.GatewayRegister{
		NodeID:              gateway.NodeID,
		Address:             gateway.Address,
		RootSigningKey:      gateway.RootSigningKey,
		SigningKey:          gateway.SigningKey,
		RegionCode:          gateway.RegionCode,
		NetworkInfoGateway:  gateway.NetworkInfoGateway,
		NetworkInfoProvider: gateway.NetworkInfoProvider,
		NetworkInfoClient:   gateway.NetworkInfoClient,
		NetworkInfoAdmin:    gateway.NetworkInfoAdmin,
	}
	return &res
}

// GetProvider returns a provider register if found.
func (mgr *FCRRegisterMgr) GetProvider(id *nodeid.NodeID) *register.ProviderRegister {
	if !mgr.start || !mgr.providerDiscv {
		return nil
	}
	mgr.registeredProvidersMapLock.RLock()
	provider, ok := mgr.registeredProvidersMap[id.ToString()]
	if !ok {
		mgr.registeredProvidersMapLock.RUnlock()
		// TODO: Do we call refresh here, if can't find a provider?
		// mgr.Refresh()
		mgr.registeredProvidersMapLock.RLock()
		provider = mgr.registeredProvidersMap[id.ToString()]
	}
	defer mgr.registeredProvidersMapLock.RUnlock()
	// Return the pointer of a copy of the register
	res := register.ProviderRegister{
		NodeID:             provider.NodeID,
		Address:            provider.Address,
		RootSigningKey:     provider.RootSigningKey,
		SigningKey:         provider.SigningKey,
		RegionCode:         provider.RegionCode,
		NetworkInfoGateway: provider.NetworkInfoGateway,
		NetworkInfoClient:  provider.NetworkInfoClient,
		NetworkInfoAdmin:   provider.NetworkInfoAdmin,
	}
	return &res
}

// GetAllGateways returns  all discovered gateways.
func (mgr *FCRRegisterMgr) GetAllGateways() []register.GatewayRegister {
	if !mgr.start || !mgr.gatewayDiscv {
		return nil
	}
	res := make([]register.GatewayRegister, 0)
	mgr.registeredGatewaysMapLock.RLock()
	defer mgr.registeredGatewaysMapLock.RUnlock()
	for _, gateway := range mgr.registeredGatewaysMap {
		res = append(res, register.GatewayRegister{
			NodeID:              gateway.NodeID,
			Address:             gateway.Address,
			RootSigningKey:      gateway.RootSigningKey,
			SigningKey:          gateway.SigningKey,
			RegionCode:          gateway.RegionCode,
			NetworkInfoGateway:  gateway.NetworkInfoGateway,
			NetworkInfoProvider: gateway.NetworkInfoProvider,
			NetworkInfoClient:   gateway.NetworkInfoClient,
			NetworkInfoAdmin:    gateway.NetworkInfoAdmin,
		})
	}
	return res
}

// GetAllProviders returns  all discovered providers.
func (mgr *FCRRegisterMgr) GetAllProviders() []register.ProviderRegister {
	if !mgr.start || !mgr.providerDiscv {
		return nil
	}
	res := make([]register.ProviderRegister, 0)
	mgr.registeredProvidersMapLock.RLock()
	defer mgr.registeredProvidersMapLock.RUnlock()
	for _, provider := range mgr.registeredProvidersMap {
		res = append(res, register.ProviderRegister{
			NodeID:             provider.NodeID,
			Address:            provider.Address,
			RootSigningKey:     provider.RootSigningKey,
			SigningKey:         provider.SigningKey,
			RegionCode:         provider.RegionCode,
			NetworkInfoGateway: provider.NetworkInfoGateway,
			NetworkInfoClient:  provider.NetworkInfoClient,
			NetworkInfoAdmin:   provider.NetworkInfoAdmin,
		})
	}
	return res
}

// GetGatewaysNearCID returns a list of gatewayRegisters whose id is close to the given cid.
func (mgr *FCRRegisterMgr) GetGatewaysNearCID(cid *cid.ContentID, numDHT int) ([]register.GatewayRegister, error) {
	// TODO: To implement.
	// Note: needs to return a list of copies.
	return nil, errors.New("Not yet implemented")
}

// updateGateways updates gateways.
func (mgr *FCRRegisterMgr) updateGateways() {
	for {
		gateways, err := register.GetRegisteredGateways(mgr.registerAPI)
		if err != nil {
			logging.Error("Register manager has error in getting registered gateways: %s", err.Error())
		} else {
			// Check for update
			for _, gateway := range gateways {
				mgr.registeredGatewaysMapLock.RLock()
				storedInfo, ok := mgr.registeredGatewaysMap[gateway.NodeID]
				mgr.registeredGatewaysMapLock.RUnlock()
				if !ok {
					// Not exist, we need to add a new entry
					mgr.registeredGatewaysMapLock.Lock()
					mgr.registeredGatewaysMap[gateway.NodeID] = gateway
					mgr.registeredGatewaysMapLock.Unlock()
				} else {
					// Exist, check if need update
					if gateway.Address != storedInfo.Address ||
						gateway.RootSigningKey != storedInfo.RootSigningKey ||
						gateway.SigningKey != storedInfo.SigningKey ||
						gateway.RegionCode != storedInfo.RegionCode ||
						gateway.NetworkInfoGateway != storedInfo.NetworkInfoGateway ||
						gateway.NetworkInfoProvider != storedInfo.NetworkInfoProvider ||
						gateway.NetworkInfoClient != storedInfo.NetworkInfoClient ||
						gateway.NetworkInfoAdmin != storedInfo.NetworkInfoAdmin {
						// Need update
						mgr.registeredGatewaysMapLock.Lock()
						mgr.registeredGatewaysMap[gateway.NodeID] = gateway
						mgr.registeredGatewaysMapLock.Unlock()
					}
				}
			}
		}
		afterChan := time.After(mgr.refreshDuration)
		select {
		case <-mgr.gatewayRefreshCh:
			// Need to refresh
			logging.Error("Register manager force update internal gateway map.")
		case <-afterChan:
			// Need to refresh
		case <-mgr.gatewayShutdownCh:
			// Need to shutdown
			logging.Error("Register manager shutdown gateway routine.")
			return
		}
	}
}

// updateProviders updates providers.
func (mgr *FCRRegisterMgr) updateProviders() {
	for {
		providers, err := register.GetRegisteredProviders(mgr.registerAPI)
		if err != nil {
			logging.Error("Register manager has error in getting registered providers: %s", err.Error())
		} else {
			// Check for update
			for _, provider := range providers {
				mgr.registeredProvidersMapLock.RLock()
				storedInfo, ok := mgr.registeredProvidersMap[provider.NodeID]
				mgr.registeredProvidersMapLock.RUnlock()
				if !ok {
					// Not exist, we need to add a new entry
					mgr.registeredProvidersMapLock.Lock()
					mgr.registeredProvidersMap[provider.NodeID] = provider
					mgr.registeredProvidersMapLock.Unlock()
				} else {
					// Exist, check if need update
					if provider.Address != storedInfo.Address ||
						provider.RootSigningKey != storedInfo.RootSigningKey ||
						provider.SigningKey != storedInfo.SigningKey ||
						provider.RegionCode != storedInfo.RegionCode ||
						provider.NetworkInfoGateway != storedInfo.NetworkInfoGateway ||
						provider.NetworkInfoClient != storedInfo.NetworkInfoClient ||
						provider.NetworkInfoAdmin != storedInfo.NetworkInfoAdmin {
						// Need update
						mgr.registeredProvidersMapLock.Lock()
						mgr.registeredProvidersMap[provider.NodeID] = provider
						mgr.registeredProvidersMapLock.Unlock()
					}
				}
			}
		}
		afterChan := time.After(mgr.refreshDuration)
		select {
		case <-mgr.providerRefreshCh:
			// Need to refresh
			logging.Error("Register manager force update internal provider map.")
		case <-afterChan:
			// Need to refresh
		case <-mgr.providerShutdownCh:
			// Need to shutdown
			logging.Error("Register manager shutdown provider routine.")
			return
		}
	}
}
