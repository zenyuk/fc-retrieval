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
	"time"

	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-register/pkg/register"
)

// To implement register manager.
// 1. NewFCRRegisterMgr() initialises a register manager.
//		If given gateway param is true, initialise a map[nodeid](register.GatewayRegister) <- this can be a db
//		If given provider param is true, initialise a map[nodeid](register.ProviderRegister) <- this can be a db
// 2. Start() starts a thread to auto update the internal map every given duration.
//		Make sure Start() only works at first time, later attempt to call should fail.
// 3. Refresh() refreshs the internal map immediately.
// 4. GetGateway() returns a gateway register if found.
// 5. GetProvider() returns a provider register if found.
// Note to handle access from multiple threads (RWMutex is required to access internal storage for example).

// Register Manager manages the internal storage of registered nodes
type FCRRegisterMgr struct {
}

// NewFCRRegisterMgr creates a new register manager.
func NewFCRRegisterMgr(registerAPI string, provider bool, gateway bool, refreshDuration time.Duration) (*FCRRegisterMgr, error) {
	return nil, errors.New("Not yet implemented")
}

// Start starts a thread to auto update the internal map every given duration.
func (mgr *FCRRegisterMgr) Start() error {
	return errors.New("Not yet implemented")
}

// Refresh refreshs the internal map immediately.
func (mgr *FCRRegisterMgr) Refresh() error {
	return errors.New("Not yet implemented")
}

// GetGateway returns a gateway register if found.
func (mgr *FCRRegisterMgr) GetGateway(id *nodeid.NodeID) (*register.GatewayRegister, error) {
	return nil, errors.New("Not yet implemented")
}

// GetProvider returns a provider register if found.
func (mgr *FCRRegisterMgr) GetProvider(id *nodeid.NodeID) (*register.ProviderRegister, error) {
	return nil, errors.New("Not yet implemented")
}
