package api

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"net"
)

func handleProviderCommunication(conn net.Conn, g *Gateway) {
	// Initialise a reader
	reader := bufio.NewReader(conn)
	//writer := bufio.NewWriter(conn)
	defer conn.Close()
	for {
		// Read message.
		// There are two possible messages:
		// 1. Publish Group CID
		// 2. Publish Group CID DHT
		msgType, data, err := readTCPMessage(reader)
		if err != nil {
			log.Println(err.Error())
			return
		}
		if msgType == ProviderPublishGroupCIDRequestType {
			// TODO: Handle this.
			log.Println(data)
		} else if msgType == ProviderDHTPublishGroupCIDRequestType {
			// TODO: Handle this.
			log.Println(data)
		} else {
			log.Println("Message not supported.")
			return
		}
	}
}

func requestSingleCIDGroupPublish(addr string, g *Gateway) error {
	conn, err := net.Dial("tcp", addr)
	defer conn.Close()
	if err != nil {
		return err
	}
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	// TODO: Dummy request
	request := GatewaySingleCIDOfferPublishRequest{}
	request.CommonFields.ProtocolVersion = g.ProtocolVersion
	request.CommonFields.ProtocolSupported = g.ProtocolSupported
	// Send request
	data, err := json.Marshal(&request)
	err = sendTCPMessage(writer, GatewaySingleCIDOfferPublishRequestType, data)
	if err != nil {
		return err
	}
	// Waiting for response.
	msgType, data, err := readTCPMessage(reader)
	if err != nil {
		return nil
	}
	// There is only one message possible, that is GatewaySingleCIDOfferPublishResponse
	if msgType == GatewaySingleCIDOfferPublishResponseType {
		response := GatewaySingleCIDOfferPublishResponse{}
		err = json.Unmarshal(data, &response)
		if err != nil {
			return err
		}
		// For now just print the response.
		log.Println(response)
		// TODO: Dummy ACK
		ack := GatewaySingleCIDOfferPublishResponseAck{}
		ack.CommonFields.ProtocolVersion = g.ProtocolVersion
		ack.CommonFields.ProtocolSupported = g.ProtocolSupported
		// Send ack
		data, err = json.Marshal(&ack)
		if err != nil {
			return err
		}
		return sendTCPMessage(writer, GatewaySingleCIDOfferPublishResponseAckType, data)
	} else {
		return errors.New("Message not supported")
	}
}
