package adminapi

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
	"encoding/json"
	"net"
	"sync"
	"time"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/messages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/tcpcomms"
)

// StartAdminAPI starts the TCP API as a separate go routine.
func StartAdminAPI(settings settings.AppSettings, g *gateway.Gateway) error {
	// Start server
	ln, err := net.Listen("tcp", settings.BindAdminAPI)
	if err != nil {
		return err
	}
	go func(ln net.Listener) {
		for {
			conn, err := ln.Accept()
			if err != nil {
				logging.Error1(err)
				continue
			}
			logging.Info("Incoming connection from admin client at :%s", conn.RemoteAddr())
			go handleIncomingAdminConnection(conn, g)
		}
	}(ln)
	logging.Info("Listening on %s for connections from admin clients", settings.BindAdminAPI)
	return nil
}

func handleIncomingAdminConnection(conn net.Conn, g *gateway.Gateway) {
	// Close connection on exit.
	defer conn.Close()

	// Loop until error occurs and connection is dropped.
	for {
		msgType, data, err := tcpcomms.ReadTCPMessage(conn, settings.DefaultTCPInactivityTimeoutMs*time.Millisecond)
		if err != nil && !tcpcomms.IsTimeoutError(err) {
			// Error in tcp communication, drop the connection.
			logging.Error1(err)
			return
		}
		//	TODO: Rework for new admin message types. Start with getReputationChallenge(clientId)
		if msgType == messages.AdminGetReputationChallengeType {
			request := messages.AdminGetReputationChallenge{}
			if json.Unmarshal(data, &request) == nil {
				// Message is valid.
				err = handleAdminGetReputationChallenge(conn, &request)
				if err != nil && !tcpcomms.IsTimeoutError(err) {
					// Error in tcp communication, drop the connection.
					logging.Error1(err)
					return
				}
				continue
			}
		}
		// Message is invalid.
		err = tcpcomms.SendInvalidMessage(conn, settings.DefaultTCPInactivityTimeoutMs*time.Millisecond)
		if err != nil && !tcpcomms.IsTimeoutError(err) {
			// Error in tcp communication, drop the connection.
			logging.Error1(err)
			return
		}
	}
}

// GetConnForRequestingAdminClient returns the connection for sending request to an admin client with given id.
// It will reuse any active connection.
func GetConnForRequestingAdminClient(gatewayID nodeid.NodeID, g *gateway.Gateway) (*gateway.CommunicationChannel, error) {
	// Check if there is an active connection.
	g.ActiveGatewaysLock.RLock() // TODO: Check this
	gComm := g.ActiveGateways[gatewayID.ToString()]
	g.ActiveGatewaysLock.RUnlock()
	if gComm == nil {
		// No active connection, connect to peer.
		g.GatewayAddressMapLock.RLock()
		conn, err := net.Dial("tcp", g.GatewayAddressMap[gatewayID.ToString()])
		g.GatewayAddressMapLock.RUnlock()
		if err != nil {
			return nil, err
		}
		gComm = &gateway.CommunicationChannel{
			CommsLock: sync.RWMutex{},
			Conn:      conn}
		if gateway.RegisterGatewayCommunication(&gatewayID, gComm) != nil {
			conn.Close()
			return nil, err
		}
	}
	return gComm, nil
}
