package messages

// gateway.go contains all messages originting from the gateway

import (
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)

// GatewaySingleCIDOfferPublishRequest is the request from gateway to provider during start-up asking for cid offers
type GatewaySingleCIDOfferPublishRequest struct {
	MessageType       int32         `json:"message_type"`
	ProtocolVersion   int32         `json:"protocol_version"`
	ProtocolSupported []int32       `json:"protocol_supported"`
	GatewayID         nodeid.NodeID `json:"gateway_id"`
	Nonce             int64         `json:"nonce"`
	Signature         string        `json:"signature"`
}

// GatewaySingleCIDOfferPublishResponse is the repsonse to GatewaySingleCIDOfferPublishRequest
type GatewaySingleCIDOfferPublishResponse struct {
	MessageType       int32   `json:"message_type"`
	ProtocolVersion   int32   `json:"protocol_version"`
	ProtocolSupported []int32 `json:"protocol_supported"`
	CIDOffers         []struct {
		Price     int64         `json:"price_per_byte"`
		Expiry    int64         `json:"expiry_date"`
		QoS       int64         `json:"qos"`
		Signature string        `json:"signature"`
		PieceCID  cid.ContentID `json:"piece_cid"`
	} `json:"cid_offers"` // TODO: Need to check if this is right.
}

// GatewaySingleCIDOfferPublishResponseAck is the acknowledgement to GatewaySingleCIDOfferPublishResponse
type GatewaySingleCIDOfferPublishResponseAck struct {
	MessageType     int32 `json:"message_type"`
	ProtocolVersion int32 `json:"protocol_version"`
	CIDOffersAck    []struct {
		Nonce     int64  `json:"nonce"`
		Signature string `json:"signature"`
	} `json:"cid_offers_ack"` // TODO: Need to check if this is right.
}

// GatewayDHTDiscoverRequest is the request from gateway to gateway to discover cid offer
type GatewayDHTDiscoverRequest struct {
	MessageType       int32         `json:"message_type"`
	ProtocolVersion   int32         `json:"protocol_version"`
	ProtocolSupported []int32       `json:"protocol_supported"`
	GatewayID         nodeid.NodeID `json:"gateway_id"`
	PieceCID          cid.ContentID `json:"piece_cid"`
	Nonce             int64         `json:"nonce"`
	TTL               int64         `json:"ttl"`
}

// GatewayDHTDiscoverResponse is the response to GatewayDHTDiscoverRequest
type GatewayDHTDiscoverResponse struct {
	MessageType     int32                 `json:"message_type"`
	ProtocolVersion int32                 `json:"protocol_version"`
	PieceCID        int64                 `json:"piece_cid"`
	Nonce           int64                 `json:"nonce"`
	Found           bool                  `json:"found"`
	Signature       string                `json:"signature"`
	CIDGroupInfo    []CIDGroupInformation `json:"cid_group_information"`
}
