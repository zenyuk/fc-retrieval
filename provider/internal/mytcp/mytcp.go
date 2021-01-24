package mytcp

import (
	"fmt"
	"net"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cid"
	log "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/messages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/tcpcomms"
)

// SendMessage to gateway
func SendMessage() {
	providerID, _ := nodeid.NewRandomNodeID()
	contentID, _ := cid.NewRandomContentID()
	pieceCIDs := []cid.ContentID{*contentID}
	conn, err := net.Dial("tcp", "localhost:8090")
	if err != nil {
		log.Error("Fail to dial")
		os.Exit(1)
	}
	request := messages.ProviderPublishGroupCIDRequest{
		MessageType:       messages.ProviderPublishGroupCIDRequestType,
		ProtocolVersion:   1,
		ProtocolSupported: []int32{1},
		Nonce:             1,                                       // TODO, Add nonce
		ProviderID:        *providerID,
		Price:						 100,
		Expiry:						 time.Now().Add(10 * time.Second).Unix(),
		QoS:							 50,
		PieceCIDs:         pieceCIDs,
	}
	err = tcpcomms.SendMessageWithType(
		conn,
		messages.ProviderPublishGroupCIDRequestType,
		&request,
		30000)

	fmt.Println("MY MESSAGE", request)
	fmt.Println("-------------------------------------------")
	fmt.Println(err)

	// if err != nil && tcpcomms.IsTimeoutError(err) {
	// 	// Timeout can be ignored. Since this message can expire.
	// 	return nil, nil
	// } else if err != nil {
	// 	pComm.Conn.Close()
	// 	gateway.DeregisterGatewayCommunication(gatewayID)
	// 	return nil, err
	// }
	// if msgType == messages.GatewayDHTDiscoverResponseType {
	// 	response := messages.GatewayDHTDiscoverResponse{}
	// 	if json.Unmarshal(data, &response) == nil {
	// 		// Message is valid.
	// 		return &response, nil
	// 	}
	// }

	fmt.Println("-------------------------------------------")
}
