package gateway

import (
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrtcpcomms"
	log "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)

// SendMessage to gateway
func SendMessage(message *fcrmessages.FCRMessage, nodeID *nodeid.NodeID, gCommPool *fcrtcpcomms.CommunicationPool) error {
	gComm, err := gCommPool.GetConnForRequestingNode(nodeID)
	if err != nil {
		log.Error("Connection issue: %v", err)
		if gComm != nil {
			log.Debug("Closing connection ...")
			gComm.Conn.Close()
		}
		log.Debug("Removing connection from pool ...")
		gCommPool.DeregisterNodeCommunication(nodeID)
		return err
	}
	gComm.CommsLock.Lock()
	defer gComm.CommsLock.Unlock()
	log.Info("Send message to: %v, message: %v", nodeID.ToString(), message)
	err = fcrtcpcomms.SendTCPMessage(
		gComm.Conn,
		message,
		30000)
	if err != nil {
		log.Error("Message not sent: %v", err)
		if gComm != nil {
			log.Debug("Closing connection ...")
			gComm.Conn.Close()
		}
		log.Debug("Removing connection from pool ...")
		gCommPool.DeregisterNodeCommunication(nodeID)
		return err
	}
	return nil
}
