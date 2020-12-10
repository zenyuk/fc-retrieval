package providerapi

import (
	"bufio"
	"encoding/json"
	"log"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/messages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/tcpcomms"
)

func handleProviderDHTPublishGroupCIDRequest(reader *bufio.Reader, writer *bufio.Writer, request *messages.ProviderDHTPublishGroupCIDRequest) error {
	// Do something about the internal request
	// Will need gateway instance to read from db
	// gateway, err := gateway.GetSingleInstance()
	log.Printf("Provider request from: %s", request.ProviderID.ToString())
	// Respond to provider.
	response, _ := json.Marshal(messages.ProviderDHTPublishGroupCIDResponse{
		MessageType: messages.GatewayDHTDiscoverResponseType,
		// This is just a dummy response
	})
	return tcpcomms.SendTCPMessage(writer, messages.ProviderDHTPublishGroupCIDResponseType, response)
}
