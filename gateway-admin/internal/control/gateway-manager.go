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
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages/fcrmsggwadmin"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrtcpcomms"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-register/pkg/register"

	log "github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway-admin/internal/gatewayapi"
	"github.com/ConsenSys/fc-retrieval-gateway-admin/internal/settings"
)

// GatewayManager managers the pool of gateways and the connections to them.
type GatewayManager struct {
	settings      settings.ClientGatewayAdminSettings
	gateway       ActiveGateway
	gatewaysLock  sync.RWMutex
	registeredMap map[string]register.RegisteredNode
	conxPool      *fcrtcpcomms.CommunicationPool
}

// ActiveGateway contains information for a single gateway
type ActiveGateway struct {
	comms *gatewayapi.Comms
}

// NewGatewayManager creates a gateway manager.
func NewGatewayManager(conf settings.ClientGatewayAdminSettings) *GatewayManager {
	g := GatewayManager{}
	g.settings = conf
	g.registeredMap = make(map[string]register.RegisteredNode)
	g.conxPool = fcrtcpcomms.NewCommunicationPool(&g.registeredMap, &sync.RWMutex{})
	return &g
}

// InitializeGateway initialise a new gateway
func (g *GatewayManager) InitializeGateway(gatewayInfo *register.GatewayRegister, gatewayPrivKey *fcrcrypto.KeyPair, gatewayPrivKeyVer *fcrcrypto.KeyVersion) error {
	// TODO check whether gateway not initialized.
	// TODO check whether contract indicates initialised
	// TODO: Check given gatewayInfo is correct
	// First, Get pubkey
	pubKey, err := gatewayInfo.GetSigningKey()
	if err != nil {
		log.Error("Error in obtaining signing key from register info.")
		return err
	}

	nodeID, err := nodeid.NewNodeIDFromHexString(gatewayInfo.NodeID)
	if err != nil {
		log.Error("Error in generating nodeID.")
		return err
	}

	// Second, send key exchange to activate the given gateway
	request, err := fcrmsggwadmin.EncodeGatewayAdminInitialiseKeyRequest(nodeID, gatewayPrivKey, gatewayPrivKeyVer)
	if err != nil {
		log.Error("Error in encoding message.")
		return err
	}

	// Sign the request
	if request.Sign(g.settings.GatewayAdminPrivateKey(), g.settings.GatewayAdminPrivateKeyVer()) != nil {
		log.Error("Error signing message for sending private key to gateway: %+v", err)
		return err
	}

	log.Info("Sending message to gateway: %v, message: %s", nodeID.ToString(), request.DumpMessage())

	conn, err := g.getConnection(nodeID, gatewayInfo.NetworkInfoAdmin) //"gateway:9013"
	if err != nil {
		return err
	}
	err = fcrtcpcomms.SendTCPMessage(conn, request, settings.DefaultTCPInactivityTimeout)
	if err != nil {
		log.Error("Error sending private key to Gateway: %s", err)
		return err
	}

	// Process the response from the gateway.
	response, err := fcrtcpcomms.ReadTCPMessage(conn, time.Second*1)
	log.Info("Response message: %+v", response)
	if response.GetMessageType() != fcrmessages.GatewayAdminInitialiseKeyResponseType {
		// TODO other types of messages such as protocol version negotiation need to be handled.
		return fmt.Errorf("Unexpected message in response to set-up Gateway message: %d", response.GetMessageType())
	}

	// Verify the response
	if response.Verify(pubKey) != nil {
		return errors.New("Fail to verify the response")
	}

	keyAccepted, err := fcrmsggwadmin.DecodeGatewayAdminInitialiseKeyResponse(response)
	if err != nil {
		return err
	}
	if !keyAccepted {
		return fmt.Errorf("Key not accepted for unspecified reason")
	}

	return gatewayInfo.RegisterGateway(g.settings.RegisterURL())
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

func (g *GatewayManager) getConnection(gatewayNodeID *nodeid.NodeID, addr string) (net.Conn, error) {
	// Add new gateway to the connection pool.
	g.registeredMap[strings.ToLower(gatewayNodeID.ToString())] = &register.GatewayRegister{
		NodeID:             gatewayNodeID.ToString(),
		NetworkInfoGateway: addr,
	}

	// Get conn for the right gateway
	channel, err := g.conxPool.GetConnForRequestingNode(gatewayNodeID, fcrtcpcomms.AccessFromGateway)
	if err != nil {
		return nil, err
	}
	conn := channel.Conn
	if err != nil {
		log.Error("Error getting a connection to gateway %v: %s", gatewayNodeID.ToString(), err)
		return nil, err
	}
	return conn, nil

}
