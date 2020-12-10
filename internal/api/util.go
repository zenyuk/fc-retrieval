package api

import (
	"sync"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/api/clientapi"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)

// CommunicationThread stores channels of a gateway/provider communication thread for internal request
type CommunicationThread struct {
	NodeID            nodeid.NodeID
	CommsRequestChan  chan []byte
	CommsResponseChan chan []byte
}

// Gateway holds the main data structure for the whole gateway.
type Gateway struct {
	ProtocolVersion   int
	ProtocolSupported []int
	ClientAPI         *clientapi.ClientAPI

	// ActiveGateways store connected active gateways
	// A map from gateway id (big int in string repr)
	// to communication struct
	ActiveGateways     map[string](*CommunicationThread)
	ActiveGatewaysLock sync.RWMutex

	// ActiveProviders store connected active providers
	// A map from provider id (big in in string repr)
	// to a channel of byte slice (for data transfer)
	ActiveProviders     map[string](*CommunicationThread)
	ActiveProvidersLock sync.RWMutex
}
