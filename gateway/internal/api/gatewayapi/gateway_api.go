package gatewayapi

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

// StartGatewayAPI starts the TCP API as a separate go routine.
func StartGatewayAPI(settings settings.AppSettings, g *gateway.Gateway) error {
	// Start server
	ln, err := net.Listen("tcp", settings.BindGatewayAPI)
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
			logging.Info("Incoming connection from gateway at :%s", conn.RemoteAddr())
			go handleIncomingGatewayConnection(conn, g)
		}
	}(ln)
	logging.Info("Listening on %s for connections from Gateways", settings.BindGatewayAPI)
	return nil
}

func handleIncomingGatewayConnection(conn net.Conn, g *gateway.Gateway) {
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
		if msgType == messages.GatewayDHTDiscoverRequestType {
			request := messages.GatewayDHTDiscoverRequest{}
			if json.Unmarshal(data, &request) == nil {
				// Message is valid.
				err = handleGatewayDHTDiscoverRequest(conn, &request)
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

// GetConnForRequestingGateway returns the connection for sending request to a gateway with given id.
// It will reuse any active connection.
func GetConnForRequestingGateway(gatewayID *nodeid.NodeID, g *gateway.Gateway) (*gateway.CommunicationChannel, error) {
	// Check if there is an active connection.
	g.ActiveGatewaysLock.RLock()
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
		if gateway.RegisterGatewayCommunication(gatewayID, gComm) != nil {
			conn.Close()
			return nil, err
		}
	}
	return gComm, nil
}
