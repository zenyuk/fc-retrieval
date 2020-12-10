package gatewayapi

import (
	"bufio"
	"encoding/json"
	"log"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/messages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/tcpcomms"
)

func handleGatewayDHTDiscoverRequest(reader *bufio.Reader, writer *bufio.Writer, request *messages.GatewayDHTDiscoverRequest) error {
	// Do something about the internal request
	// Will need gateway instance to read from db
	// gateway, err := gateway.GetSingleInstance()
	log.Printf("Gateway request from: %s", request.GatewayID.ToString())
	// Respond to gateway.
	response, _ := json.Marshal(messages.GatewayDHTDiscoverResponse{
		MessageType: messages.GatewayDHTDiscoverResponseType,
		// This is just a dummy response
	})
	return tcpcomms.SendTCPMessage(writer, messages.GatewayDHTDiscoverResponseType, response)
}
