package messages

// provider.go contains all messages originting from the provider

import (
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)

// ProviderPublishGroupCIDRequest is the request from provider to gateway to publish group cid offer
// It does not require a response.
type ProviderPublishGroupCIDRequest struct {
	MessageType       int32           `json:"message_type"`
	ProtocolVersion   int32           `json:"protocol_version"`
	ProtocolSupported []int32         `json:"protocol_supported"`
	Nonce             int64           `json:"nonce"`
	ProviderID        nodeid.NodeID   `json:"provider_id"`
	Price             int64           `json:"price_per_byte"`
	Expiry            int64           `json:"expiry_date"`
	QoS               int64           `json:"qos"`
	Signature         string          `json:"signature"`
	PieceCIDs         []cid.ContentID `json:"piece_cids"`
}

// ProviderDHTPublishGroupCIDRequest is the request from provider to gateway to publish group cid offer using DHT
type ProviderDHTPublishGroupCIDRequest struct {
	MessageType       int32         `json:"message_type"`
	ProtocolVersion   int32         `json:"protocol_version"`
	ProtocolSupported []int32       `json:"protocol_supported"`
	Nonce             int64         `json:"nonce"`
	ProviderID        nodeid.NodeID `json:"provider_id"`
	NumOffers         int64         `json:"num_of_offers"`
	CIDOffers         []struct {
		Price     int64         `json:"price_per_byte"`
		Expiry    int64         `json:"expiry_date"`
		QoS       int64         `json:"qos"`
		Signature string        `json:"signature"`
		PieceCID  cid.ContentID `json:"piece_cid"`
	} `json:"cid_offers"`
}

// ProviderDHTPublishGroupCIDResponse is the response to ProviderDHTPublishGroupCIDRequest
type ProviderDHTPublishGroupCIDResponse struct {
	MessageType     int32  `json:"message_type"`
	ProtocolVersion int32  `json:"protocol_version"`
	Nonce           int64  `json:"nonce"`
	Signature       string `json:"signature"`
}
