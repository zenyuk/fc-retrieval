package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

func handleGatewayCommunication(conn net.Conn, g *Gateway) {
	// Initialise a reader
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	defer conn.Close()
	for {
		// Read message.
		// There is only one possible message type:
		// Gateway DHT Discover Request.
		msgType, data, err := readTCPMessage(reader)
		if err != nil {
			log.Println(err.Error())
			return
		}
		if msgType == GatewayDHTDiscoverRequestType {
			request := GatewayDHTDiscoverRequest{}
			err = json.Unmarshal(data, &request)
			if err != nil {
				log.Println(err.Error())
				return
			}
			// Object is created, for now just print the object
			log.Println(request)
			// TODO: Dummy response
			response := ClientEstablishmentResponse{}
			response.CommonFields.ProtocolVersion = g.ProtocolVersion
			response.CommonFields.ProtocolSupported = g.ProtocolSupported
			// Reply
			data, err = json.Marshal(&response)
			if err != nil {
				log.Println(err.Error())
				return
			}
			err = sendTCPMessage(writer, GatewayDHTDiscoverResponseType, data)
			if err != nil {
				log.Println(err.Error())
				return
			}
		} else {
			log.Println("Message not supported.")
			return
		}
	}
}

func requestDHTDiscover(addrs []string, g *Gateway) error {
	for i, addr := range addrs {
		// TODO: Need to handle this request.
		fmt.Printf("%d: %s", i, addr)
	}
	return nil
}
