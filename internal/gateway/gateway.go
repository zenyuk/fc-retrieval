package gateway

import (
	"net"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrtcpcomms"
	log "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
)

// SendMessage to gateway
func SendMessage(gwURL string, message *fcrmessages.FCRMessage) {
	log.Info("Send message to: %v, message: %v", gwURL, message)
	conn, err := net.Dial("tcp", gwURL)
	if err != nil {
		log.Panic("Fail to dial: %v", gwURL)
	}
	err = fcrtcpcomms.SendTCPMessage(
		conn,
		message,
		30000)
	if err != nil {
		log.Error("Message sent with error: %v", err)
	}
}
