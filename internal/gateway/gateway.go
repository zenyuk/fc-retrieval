package gateway

import (
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrtcpcomms"
	log "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-provider/pkg/communication"
)

// SendMessage to gateway
func SendMessage(message *fcrmessages.FCRMessage, nodeID *nodeid.NodeID, gCommPool *communication.CommunicationPool) {
	gComm, err := gCommPool.GetConnForRequestingNode(nodeID)
	if err != nil {
		gComm.Conn.Close()
		gCommPool.DeregisterNodeCommunication(nodeID)
	}
	gComm.CommsLock.Lock()
	defer gComm.CommsLock.Unlock()
	log.Info("Send message to: %v, message: %v", nodeID.ToString(), message)
	err = fcrtcpcomms.SendTCPMessage(
		gComm.Conn,
		message,
		30000)
	if err != nil {
		log.Error("Message sent with error: %v", err)
		gComm.Conn.Close()
		gCommPool.DeregisterNodeCommunication(nodeID)
	}
}
