package api

import "math/big"

// NOTE: Time, signature and proofs are represented using string here.

// Gateway struct stores information of this gateway.
// TODO: Gateway will likely to be initialised in somewhere higher level.
type Gateway struct {
	ProtocolVersion   int
	ProtocolSupported []int
	// TODO: Add more fields (privkey, gateway id, etc.)
}

// CommonMessageFields are shared fileds of all http requests
type CommonMessageFields struct {
	ProtocolVersion   int   `json:"protocol_version"`
	ProtocolSupported []int `json:"protocol_supported"`
	MessageType       int   `json:"message_type"`
}

// CIDGroupOffer represents a cid group offer
type CIDGroupOffer struct {
	ProviderID           big.Int `json:"provider_id"`
	Price                int     `json:"price_per_byte"`
	Expiry               string  `json:"expiry_date"`
	QoS                  int     `json:"qos"`
	Signature            string  `json:"signature"`
	MerkleProof          string  `json:"merkle_proof"`
	FundedPaymentChannel bool    `json:"funded_payment_channel"` // Is this boolean?
}

// ClientEstablishmentRequest is the request from client to gateway to establish connection
type ClientEstablishmentRequest struct {
	CommonFields CommonMessageFields `json:"common_fields"`
	Challenge    string              `json:"challenge"`
	TTL          string              `json:"ttl"`
}

// ClientEstablishmentResponse is the response to ClientEstablishmentRequest
type ClientEstablishmentResponse struct {
	CommonFields CommonMessageFields `json:"common_fields"`
	Challenge    string              `json:"challenge"`
	Signature    string              `json:"signature"`
}

// ClientStandardDiscoverRequest is the requset from client to gateway to ask for cid offer
type ClientStandardDiscoverRequest struct {
	CommonFields CommonMessageFields `json:"common_fields"`
	PieceCID     big.Int             `json:"cid"`
	Nonce        int                 `json:"nonce"`
	TTL          string              `json:"ttl"`
}

// ClientStandardDiscoverResponse is the response to ClientStandardDiscoverResponse
type ClientStandardDiscoverResponse struct {
	CommonFields        CommonMessageFields `json:"common_fields"`
	PieceCID            big.Int             `json:"cid"`
	Nonce               int                 `json:"nonce"`
	Found               bool                `json:"found"`
	Signature           string              `json:"signature"`
	CIDGroupInformation []CIDGroupOffer     `json:cid_group_information`
}

// ClientDHTDiscoverRequest is the request from client to gateway to ask for cid offer using DHT
type ClientDHTDiscoverRequest struct {
	CommonFields       CommonMessageFields `json:"common_fields"`
	ieceCID            big.Int             `json:"cid"`
	Nonce              int                 `json:"nonce"`
	TTL                string              `json:"ttl"`
	NumDHT             int                 `json:"num_dht"`
	IncrementalResults bool                `json:"incremental_results"`
}

// ClientDHTDiscoverResponse is the response to ClientDHTDiscoverRequest
type ClientDHTDiscoverResponse struct {
	CommonFields  CommonMessageFields              `json:"common_fields"`
	Contacted     []ClientStandardDiscoverResponse `json:"contacted_gateways"`
	UnContactable []struct {
		GatewayID big.Int `json:"gateway_id"`
		Nonce     int     `json:"nonce"`
	} `json:"uncontactable_gateways"`
}
