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
	"fmt"
	"sync"

	"github.com/ConsenSys/fc-retrieval-gateway-admin/internal/contracts"
	"github.com/ConsenSys/fc-retrieval-gateway-admin/internal/gatewayapi"
	"github.com/ConsenSys/fc-retrieval-gateway-admin/internal/settings"
	log "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
)

// GatewayManager managers the pool of gateways and the connections to them.
type GatewayManager struct {
	settings                    settings.ClientGatewayAdminSettings
	gatewayRegistrationContract *contracts.GatewayRegistrationContract
	gateways                    []ActiveGateway
	gatewaysLock                sync.RWMutex
}

// ActiveGateway contains information for a single gateway
type ActiveGateway struct {
	info  contracts.GatewayInformation
	comms *gatewayapi.Comms
}

// // GatewayManagerSettings is used to communicate the settings to be used by the
// // Gateway Manager.
// type GatewayManagerSettings struct {
// 	MaxEstablishmentTTL int64
// 	NodeID *nodeid.NodeID
// }

var doOnce sync.Once
var singleInstance *GatewayManager

// GetGatewayManager returns the single instance of the gateway manager.
// The settings parameter must be used with the first call to this function.
// After that, the settings parameter is ignored.
func GetGatewayManager(settings ...settings.ClientGatewayAdminSettings) *GatewayManager {
	doOnce.Do(func() {
		if len(settings) != 1 {
			log.ErrorAndPanic("Unexpected number of parameter passed to first call of GetGatewayManager")
		}
		startGatewayManager(settings[0])
	})
	return singleInstance
}

func startGatewayManager(settings settings.ClientGatewayAdminSettings) {
	g := GatewayManager{}
	g.settings = settings
	g.gatewayRegistrationContract = contracts.GetGatewayRegistrationContract()

	singleInstance = &g

	//	errChan := make(chan error, 1)
	//	go g.gatewayManagerRunner()
	g.gatewayManagerRunner()

	// TODO what should be done with error that is returned possibly in the future?
	// TODO would it be better just to have gatewayManagerRunner panic after emitting a log?
}

func (g *GatewayManager) gatewayManagerRunner() {
	log.Info("Gateway Manager: Management thread started")

	// Call this once each hour or maybe day.
	g.gatewayRegistrationContract.FetchUpdatedInformationFromContract()

	// TODO this loop is where the managing of gateways that the client is using
	// happens.

	// TODO given we are using dummy data, just grab the gateway information once.
	gatewayInfo := g.gatewayRegistrationContract.GetGateways(10)
	numGateways := len(gatewayInfo)
	numGatewaysStr := fmt.Sprint(numGateways)
	log.Info("Gateway Manager: GetGateways returned " + numGatewaysStr + "gateways")
	for _, info := range gatewayInfo {
		host := info.Hostname
		comms, err := gatewayapi.NewGatewayAPIComms(host, g.settings.ClientID())
		if err != nil {
			log.ErrorAndPanic("Could not establish communication with the gateway")
		}

		activeGateway := ActiveGateway{info, comms}
		g.gateways = append(g.gateways, activeGateway)
	}

	log.Info("Gateway Manager using %d gateways", len(g.gateways))
}

// BlockGateway adds a host to disallowed list of gateways
func (g *GatewayManager) BlockGateway(hostName string) {
	// TODO
}

// UnblockGateway add a host to allowed list of gateways
func (g *GatewayManager) UnblockGateway(hostName string) {
	// TODO

}

// Shutdown stops go routines and closes sockets. This should be called as part
// of the graceful library shutdown
func (g *GatewayManager) Shutdown() {
	// TODO
}
