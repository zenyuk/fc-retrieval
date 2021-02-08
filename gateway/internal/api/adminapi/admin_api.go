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
	"net"
	"sync"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrtcpcomms"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/tcpcomms"
)

// StartAdminAPI starts the TCP API as a separate go routine.
func StartAdminAPI(settings settings.AppSettings, g *gateway.Gateway) error {
	// Start server
	ln, err := net.Listen("tcp", ":"+settings.BindAdminAPI)
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
		message, err := fcrtcpcomms.ReadTCPMessage(conn, settings.DefaultTCPInactivityTimeout)
		if err != nil && !fcrtcpcomms.IsTimeoutError(err) {
			// Error in tcp communication, drop the connection.
			logging.Error1(err)
			return
		}
		// Respond to requests for a client's reputation.
		if err == nil {
			if message.MessageType == fcrmessages.AdminGetReputationChallengeType {
				err = handleAdminGetReputationChallenge(conn, message)
				if err != nil && !tcpcomms.IsTimeoutError(err) {
					// Error in tcp communication, drop the connection.
					logging.Error1(err)
					return
				}
				continue
			} else if message.MessageType == fcrmessages.AdminSetReputationChallengeType {
				err = handleAdminSetReputationChallenge(conn, message)
				if err != nil && !tcpcomms.IsTimeoutError(err) {
					// Error in tcp communication, drop the connection.
					logging.Error1(err)
					return
				}
				continue
			} else if message.MessageType == fcrmessages.AdminAcceptKeyChallengeType {
				err = handleAdminAcceptKeysChallenge(conn, message)
				if err != nil && !tcpcomms.IsTimeoutError(err) {
					// Error in tcp communication, drop the connection.
					logging.Error1(err)
					return
				}
				continue
			}
		}

		/*
			   TODO: Add additional message types:
			   ✔︎ Get a client's reputation.
			   - Set reputation of client arbitrarily.
			   - Set reputation of client based on various actions (e.g. using existing functionality).
			   - Set reputation of other gateway.
			   - Set reputation of provider.
			   - Get id of random client (for testing purposes).
			   - Remove Piece CID offers from the standard cache.
			   - Remove Piece CID offers from the DHT cache.
			   - Remove all Piece CID offers from a certain provider from the standard or DHT cache.
			   - generate a key pair for the gateway. The API should have an optional parameter which
				is protocol version. The API should return a hex encoded private key that the user
				could put into the gateway settings file (for the moment this is where the private
				key will live, though before this goes into production we will need to use something
				like EthSigner so the private key can be in a HSM). There will probably also need to
				be an API to register the public key on the blockchain.
		*/

		// Message is invalid.
		fcrtcpcomms.SendInvalidMessage(conn, settings.DefaultTCPInactivityTimeout)
	}
}

// GetConnForRequestingAdminClient returns the connection for sending request to an admin client with given id.
// It will reuse any active connection.
func GetConnForRequestingAdminClient(gatewayID nodeid.NodeID, g *gateway.Gateway) (*gateway.CommunicationChannel, error) {
	// Check if there is an active connection.
	g.ActiveGatewaysLock.RLock() // TODO: Check this - Will need to add active admin connections to core structure
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
