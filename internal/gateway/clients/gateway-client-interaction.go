package clients

import (
	"github.com/ConsenSys/fc-retrieval-gateway/internal/reputation"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/messages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)

// GatewayClientInteraction holds data to manage the interaction between client and gateway
type GatewayClientInteraction struct {

}


// NewGatewayClientInteraction creates a new object for handling client - gateways
func NewGatewayClientInteraction() *GatewayClientInteraction {
	g := GatewayClientInteraction{}
	return &g
}

// Establishment completes all processing for an establishment message
func (c *GatewayClientInteraction) Establishment(req *messages.ClientEstablishmentRequest) (*messages.ClientEstablishmentResponse, error) {
	clientID, err := nodeid.NewNodeIDFromString(req.ClientID)
	if err != nil {
		return nil, err
	}

	repSystem := reputation.GetSingleInstance()
	repSystem.EstablishClientReputation(clientID)

	resp := &messages.ClientEstablishmentResponse{}
	resp.Challenge = req.Challenge
	return resp, nil
}