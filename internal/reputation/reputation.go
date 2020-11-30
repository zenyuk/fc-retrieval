package reputation
// Copyright (C) 2020 ConsenSys Software Inc

import (
	"sync"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)


// Single instance of the reputation system
var instance = newInstance()



// Reputation manages the reputation of all other nodes in the system from this 
// node's perspective.
type Reputation struct {
	clients          map[string]int
	clientsMapLock   sync.RWMutex
	gateways         map[string]int
	gatewaysMapLock  sync.RWMutex
	providers        map[string]int
	providersMapLock sync.RWMutex
}


// Create a new instance
func newInstance() *Reputation {
	var r = Reputation{}
	r.clients = make(map[string]int)
	r.clientsMapLock = sync.RWMutex{}
	r.gateways = make(map[string]int)
	r.gatewaysMapLock = sync.RWMutex{}
	r.providers = make(map[string]int)
	r.providersMapLock = sync.RWMutex{}
	return &r
}


// GetSingleInstance is a factory method to get the single instance of the reputation system
func GetSingleInstance() *Reputation {
	return instance
}

// EstablishClientReputation initialise the reputation of a Retrieval Client
func (r *Reputation) EstablishClientReputation(clientNodeID *nodeid.NodeID) {
	r.setClientReputation(clientNodeID, clientInitialReputation)
}

// ClientExists determines if a client has reputation
func (r *Reputation) ClientExists(clientNodeID *nodeid.NodeID) (exists bool) {
	_, exists = r.getClientReputation(clientNodeID)
	return
}

// GetClientReputation returns the client reputation, and creates a reputation for the Retrival Client 
// if the client doesn't have a reputation yet.
func (r *Reputation) GetClientReputation(clientNodeID *nodeid.NodeID) (val int) {
	var exists bool
	val, exists = r.getClientReputation(clientNodeID)
	if (!exists) {
		r.EstablishClientReputation(clientNodeID)
	}
	return
}

// ClientEstablishmentChallenge updates a Retrieval Client's reputation based on an
// Establishment Challenge being received.
func (r *Reputation) ClientEstablishmentChallenge(clientNodeID *nodeid.NodeID) {
	r.changeClientReputation(clientNodeID, clientEstablishmentChallenge);
}

// OnChainDeposit updates a Retrieval Client's reputation based on an
// a deposit of Establishment Challenge being received.
func (r *Reputation) OnChainDeposit(clientNodeID *nodeid.NodeID) {
	r.changeClientReputation(clientNodeID, clientOnChainDeposit)
}