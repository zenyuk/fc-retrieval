package main

// Copyright (C) 2020 ConsenSys Software Inc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	// "net"
	// "sync"

	// "github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"

	"github.com/ConsenSys/fc-retrieval-provider/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-provider/pkg/nodeid"
)

// func GetConnForRequestingGateway(gatewayID *nodeid.NodeID, g *gateway.Gateway) (*gateway.CommunicationChannel, error) {
// 	// Check if there is an active connection.
// 	g.ActiveGatewaysLock.RLock()
// 	gComm := g.ActiveGateways[gatewayID.ToString()]
// 	g.ActiveGatewaysLock.RUnlock()
// 	if gComm == nil {
// 		// No active connection, connect to peer.
// 		g.GatewayAddressMapLock.RLock()
// 		conn, err := net.Dial("tcp", g.GatewayAddressMap[gatewayID.ToString()])
// 		g.GatewayAddressMapLock.RUnlock()
// 		if err != nil {
// 			return nil, err
// 		}
// 		gComm = &gateway.CommunicationChannel{
// 			CommsLock: sync.RWMutex{},
// 			Conn:      conn}
// 		if gateway.RegisterGatewayCommunication(gatewayID, gComm) != nil {
// 			conn.Close()
// 			return nil, err
// 		}
// 	}
// 	return gComm, nil
// 		}

// RequestGatewayDHTDiscover is used to request a DHT CID Discover
func RequestGatewayDHTDiscover(cid *cid.ContentID, gatewayID *nodeid.NodeID) (*messages.GatewayDHTDiscoverResponse, error) {
	// Get the connection to the gateway.
	pComm, err := GetConnForRequestingGateway(gatewayID, g)
	if err != nil {
		pComm.Conn.Close()
		gateway.DeregisterGatewayCommunication(gatewayID)
		return nil, err
	}
	pComm.CommsLock.Lock()
	defer pComm.CommsLock.Unlock()
	// Construct message
	request := messages.GatewayDHTDiscoverRequest{
		MessageType:       messages.GatewayDHTDiscoverRequestType,
		ProtocolVersion:   42,
		ProtocolSupported: 42,
		PieceCID:          *cid,
		Nonce:             1,                                       // TODO, Add nonce
		TTL:               time.Now().Add(10 * time.Second).Unix(), // TODO, ADD TTL, for now 10 seconds
	}
	err = tcpcomms.SendMessageWithType(pComm.Conn, messages.GatewayDHTDiscoverRequestType, &request, settings.DefaultTCPInactivityTimeout)
	if err != nil {
		pComm.Conn.Close()
		gateway.DeregisterGatewayCommunication(gatewayID)
		return nil, err
	}
	// Get a response
	msgType, data, err := tcpcomms.ReadTCPMessage(pComm.Conn, settings.DefaultLongTCPInactivityTimeout)
	if err != nil && tcpcomms.IsTimeoutError(err) {
		// Timeout can be ignored. Since this message can expire.
		return nil, nil
	} else if err != nil {
		pComm.Conn.Close()
		gateway.DeregisterGatewayCommunication(gatewayID)
		return nil, err
	}
	if msgType == messages.GatewayDHTDiscoverResponseType {
		response := messages.GatewayDHTDiscoverResponse{}
		if json.Unmarshal(data, &response) == nil {
			// Message is valid.
			return &response, nil
		}
	}
	// Message is invalid.
	return nil, errors.New("invalid message")
}

func main() {
	peerID, _ := nodeid.NewRandomNodeID()
	cid, _ := cid.NewRandomContentID()

	fmt.Println(peerID)
	fmt.Println(cid)
	fmt.Println("toot")

	response, _ := RequestGatewayDHTDiscover(cid, peerID)
	data, _ := json.Marshal(response)
	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, data, "", "\t")
	fmt.Println(string(prettyJSON.Bytes()))
}
