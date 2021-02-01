package fcrtcpcomms

import (
	"errors"
	"net"
	"sync"

	log "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)

// CommunicationChannel holds the connection for sending outgoing TCP requests.
// CommsLock is used to ensure only one thread can access the tcp connection at any time.
// Conn is the net connection for sending outgoing TCP requests.
type CommunicationChannel struct {
	CommsLock sync.RWMutex
	Conn      net.Conn
}

// CommunicationPool holds the node address map and active node connections.
type CommunicationPool struct {
	// AddressMap stores mapping from node id (big int in string repr) to its address.
	NodeAddressMap     	map[string](string)
	NodeAddressMapLock 	sync.RWMutex

	// ActiveNodes store connected active nodes for outgoing request:
	// A map from node id (big int in string repr) to a CommunicationChannel.
	ActiveNodes     		map[string](*CommunicationChannel)
	ActiveNodesLock 		sync.RWMutex
}

// Create a new communication commPool.
func NewCommunicationPool() CommunicationPool {
	return CommunicationPool{
		NodeAddressMap:        make(map[string](string)),
		NodeAddressMapLock:    sync.RWMutex{},
		ActiveNodes:           make(map[string](*CommunicationChannel)),
		ActiveNodesLock:       sync.RWMutex{},
	}
}

// GetConnForRequestingNode returns the connection for sending request to a node with given id.
// It will reuse any active connection.
func (commPool *CommunicationPool) GetConnForRequestingNode(nodeID *nodeid.NodeID) (*CommunicationChannel, error) {
	log.Info("Get active connection, nodeID: %v", nodeID.ToString())
	commPool.ActiveNodesLock.RLock()
	comm := commPool.ActiveNodes[nodeID.ToString()]
	commPool.ActiveNodesLock.RUnlock()
	if comm == nil {
		log.Info("No active connection, connect to peer")
		commPool.NodeAddressMapLock.RLock()
		address := commPool.NodeAddressMap[nodeID.ToString()]
		log.Debug("Got address: %v", address)
		conn, err := net.Dial("tcp", address)
		commPool.NodeAddressMapLock.RUnlock()
		if err != nil {
			log.Error("Unable to get connection: %v", err)
			return nil, err
		}
		comm = &CommunicationChannel{
			CommsLock: sync.RWMutex{},
			Conn:      conn}
		err = commPool.RegisterNodeCommunication(nodeID, comm)
		if err != nil {
			log.Error("Unable to register node communication: %v", err)
			conn.Close()
			return nil, err
		}
	}
	return comm, nil
}

// RegisterNodeAddress registers a node address
func (commPool *CommunicationPool) RegisterNodeAddress(nodeID *nodeid.NodeID, address string) {
	commPool.NodeAddressMapLock.Lock()
	defer commPool.NodeAddressMapLock.Unlock()
	commPool.NodeAddressMap[nodeID.ToString()] = address
}

// DeregisterNodeAddress deregisters a node address
// Fail silently
func (commPool *CommunicationPool) DeregisterNodeAddress(nodeID *nodeid.NodeID) {
	commPool.NodeAddressMapLock.Lock()
	defer commPool.NodeAddressMapLock.Unlock()
	_, exist := commPool.NodeAddressMap[nodeID.ToString()]
	if exist {
		delete(commPool.NodeAddressMap, nodeID.ToString())
	}
}

// RegisterNodeCommunication registers a node communication
func (commPool *CommunicationPool) RegisterNodeCommunication(nodeID *nodeid.NodeID, comm *CommunicationChannel) error {
	commPool.ActiveNodesLock.Lock()
	defer commPool.ActiveNodesLock.Unlock()
	_, exist := commPool.ActiveNodes[nodeID.ToString()]
	if exist {
		return errors.New("Error: connection existed")
	}
	commPool.ActiveNodes[nodeID.ToString()] = comm
	return nil
}

// DeregisterNodeCommunication deregisters a node communication
// Fail silently
func (commPool *CommunicationPool) DeregisterNodeCommunication(nodeID *nodeid.NodeID) {
	commPool.ActiveNodesLock.Lock()
	defer commPool.ActiveNodesLock.Unlock()
	_, exist := commPool.ActiveNodes[nodeID.ToString()]
	if exist {
		delete(commPool.ActiveNodes, nodeID.ToString())
	}
}