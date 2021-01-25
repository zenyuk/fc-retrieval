package gateway

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

	peerID, _ := nodeid.NewRandomNodeID()
	cid, _ := cid.NewRandomContentID()

	fmt.Println(peerID)
	fmt.Println(cid)

	// Initialise a dummy gateway instance.
	// g := gateway.GetSingleInstance(&settings)

	// response, _ := RequestGatewayDHTDiscover(cid, peerID, g)
	// data, _ := json.Marshal(response)
	// var prettyJSON bytes.Buffer
	// json.Indent(&prettyJSON, data, "", "\t")
	// fmt.Println(string(prettyJSON.Bytes()))

	conn, err := net.Dial("tcp", "localhost:8090")
	if err != nil {
		log.Error("Fail to dial")
		os.Exit(1)
	}

	// request := messages.GatewayDHTDiscoverRequest{
	// 	MessageType:       messages.GatewayDHTDiscoverRequestType,
	// 	ProtocolVersion:   1,
	// 	ProtocolSupported: []int32{1},
	// 	Nonce:             1,                                       // TODO, Add nonce
	// 	TTL:               time.Now().Add(10 * time.Second).Unix(), // TODO, ADD TTL, for now 10 seconds
	// }

	// err = tcpcomms.SendMessageWithType(
	// 	conn,
	// 	messages.GatewayDHTDiscoverRequestType,
	// 	&request,
	// 30000)

	request := messages.GatewayDHTDiscoverRequest{
		MessageType:       messages.GatewayDHTDiscoverRequestType,
		ProtocolVersion:   1,
		ProtocolSupported: []int32{1},
		PieceCID:          *cid,
		Nonce:             1,                                       // TODO, Add nonce
		TTL:               time.Now().Add(10 * time.Second).Unix(), // TODO, ADD TTL, for now 10 seconds
	}

	err = tcpcomms.SendMessageWithType(
		conn,
		messages.GatewayDHTDiscoverRequestType,
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
