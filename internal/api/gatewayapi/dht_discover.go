package gatewayapi

import (
	"encoding/json"
	"net"
	"time"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/messages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/tcpcomms"
)

func handleGatewayDHTDiscoverRequest(conn net.Conn, request *messages.GatewayDHTDiscoverRequest) error {
	// First check if the message can be discarded.
	if time.Now().Unix() > request.TTL {
		// Message discarded.
		return nil
	}
	// Respond to DHT CID Discover Request
	// Get gateway core struct
	gateway := gateway.GetSingleInstance()
	offers, exists := gateway.Offers.GetOffers(&request.PieceCID)

	// Construct response
	response := messages.GatewayDHTDiscoverResponse{
		MessageType:     messages.GatewayDHTDiscoverResponseType,
		ProtocolVersion: gateway.ProtocolVersion,
		PieceCID:        request.PieceCID,
		Nonce:           request.Nonce}

	if exists {
		response.Found = true
		response.CIDGroupInfo = make([]messages.CIDGroupInformation, len(offers))
		for i, offer := range offers {
			response.CIDGroupInfo[i].ProviderID = *offer.NodeID
			response.CIDGroupInfo[i].Price = offer.Price
			response.CIDGroupInfo[i].Expiry = offer.Expiry
			response.CIDGroupInfo[i].QoS = offer.QoS
			// List of Todos
			response.CIDGroupInfo[i].Signature = "TODO"
			response.CIDGroupInfo[i].MerkleProof = "TODO"
			response.CIDGroupInfo[i].FundedPaymentChannel = false
		}
	} else {
		response.Found = false
		response.CIDGroupInfo = make([]messages.CIDGroupInformation, 0)
	}
	response.Signature = "TODO" // TODO, Sign the fields
	// Send message
	data, _ := json.Marshal(response)
	return tcpcomms.SendTCPMessage(conn, messages.GatewayDHTDiscoverResponseType, data, settings.DefaultTCPInactivityTimeoutMs*time.Millisecond)
}
