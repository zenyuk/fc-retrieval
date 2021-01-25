package gateway

import (
	"net"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cid"
	log "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/messages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/tcpcomms"
)

// CIDMessage data model
type CIDMessage struct {
	MessageType       int32           `json:"message_type"`
	ProtocolVersion   int32           `json:"protocol_version"`
	ProtocolSupported []int32         `json:"protocol_supported"`
	Nonce             int64           `json:"nonce"`
	ProviderID        nodeid.NodeID   `json:"provider_id"`
	Price             uint64          `json:"price_per_byte"`
	Expiry            int64           `json:"expiry_date"`
	QoS               uint64          `json:"qos"`
	Signature         string          `json:"signature"`
	PieceCIDs         []cid.ContentID `json:"piece_cids"`
}

// SendMessage to gateway
func SendMessage(gwURL string, message messages.ProviderPublishGroupCIDRequest) {
	log.Info("Send message to: %v, message: %v", gwURL, message)
	conn, err := net.Dial("tcp", gwURL)
	if err != nil {
		log.Panic("Fail to dial: %v", gwURL)
	}

	cid, _ := cid.NewRandomContentID()
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

}
