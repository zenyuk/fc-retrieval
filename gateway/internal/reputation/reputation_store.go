package reputation
// Copyright (C) 2020 ConsenSys Software Inc

import (
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)

// NOTE: At present all reputation is stored in memory. This will need to change as we move beyond a PoC.


func (r *Reputation) getClientReputation(clientNodeID *nodeid.NodeID) (val int, exists bool) {
	val, exists = r.clients[clientNodeID.ToString()]
	return
}

func (r *Reputation) setClientReputation(clientNodeID *nodeid.NodeID, val int) {
	r.clientsMapLock.Lock()
	r.clients[clientNodeID.ToString()] = val
	r.clientsMapLock.Unlock()
}

func (r *Reputation) changeClientReputation(clientNodeID *nodeid.NodeID, amount int) {
	clientNodeIDStr := clientNodeID.ToString()
	var val int
	var exists bool
	val, exists = r.clients[clientNodeIDStr]
	if (!exists) {
		panic("changeClientReputation for non-existant client: " + clientNodeIDStr)
	}
	r.clientsMapLock.Lock()
	r.clients[clientNodeIDStr] = val + amount
	r.clientsMapLock.Unlock()
}
