package api

import "math/big"

// Message types
// TODO: Message between client and provider will need to be added here.
// The enum should remains the same for client,provider and gateway.
// Maybe we should use defined number rather than iota.
// Maybe this file should be shared among the three repos.
const (
	clientEstablishmentRequestType = iota
	clientEstablishmentResponseType
	clientStandardDiscoverRequestType
	clientStandardDiscoverResponseType
	clientDHTDiscoverRequestType
	clientDHTDiscoverResponseType
	ProviderPublishGroupCIDRequestType
	ProviderDHTPublishGroupCIDRequestType
	ProviderDHTPublishGroupCIDResponseType
	GatewaySingleCIDOfferPublishRequestType
	GatewaySingleCIDOfferPublishResponseType
	GatewaySingleCIDOfferPublishResponseAckType
	GatewayDHTDiscoverRequestType
	GatewayDHTDiscoverResponseType
)

// TODO: Time, signature and proofs are represented using string here, will need to be changed.

// Gateway struct stores information of this gateway.
// TODO: Gateway will likely to be initialised in somewhere higher level.
type Gateway struct {
	ProtocolVersion   int
	ProtocolSupported []int
	// TODO: Add more fields (privkey, gateway id, etc.)
	// TODO: Add mutex for accessing gateway information.
}

// CommonMessageFields are shared fileds of all http requests
type CommonMessageFields struct {
	ProtocolVersion   int   `json:"protocol_version"`
	ProtocolSupported []int `json:"protocol_supported"`
	MessageType       int   `json:"message_type"`
}

// CIDGroupInformation represents a cid group information
type CIDGroupInformation struct {
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
	PieceCID     big.Int             `json:"piece_cid"`
	Nonce        int                 `json:"nonce"`
	TTL          string              `json:"ttl"`
}

// ClientStandardDiscoverResponse is the response to ClientStandardDiscoverResponse
type ClientStandardDiscoverResponse struct {
	CommonFields CommonMessageFields   `json:"common_fields"`
	PieceCID     big.Int               `json:"piece_cid"`
	Nonce        int                   `json:"nonce"`
	Found        bool                  `json:"found"`
	Signature    string                `json:"signature"`
	CIDGroupInfo []CIDGroupInformation `json:"cid_group_information"`
}

// ClientDHTDiscoverRequest is the request from client to gateway to ask for cid offer using DHT
type ClientDHTDiscoverRequest struct {
	CommonFields       CommonMessageFields `json:"common_fields"`
	PieceCID           big.Int             `json:"piece_cid"`
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

// ProviderPublishGroupCIDRequest is the request from provider to gateway to publish group cid offer
type ProviderPublishGroupCIDRequest struct {
	CommonFields CommonMessageFields `json:"common_fields"`
	Nonce        int                 `json:"nonce"`
	ProviderID   big.Int             `json:"provider_id"`
	Price        int                 `json:"price_per_byte"`
	Expiry       string              `json:"expiry_date"`
	QoS          int                 `json:"qos"`
	Signature    string              `json:"signature"`
	PieceCIDs    []big.Int           `json:"piece_cids"`
}

// ProviderDHTPublishGroupCIDRequest is the request from provider to gateway to publish group cid offer using DHT
type ProviderDHTPublishGroupCIDRequest struct {
	CommonFields CommonMessageFields `json:"common_fields"`
	Nonce        int                 `json:"nonce"`
	ProviderID   big.Int             `json:"provider_id"`
	NumOffers    int                 `json:"num_of_offers"`
	CIDOffers    []struct {
		Price     int     `price_per_byte`
		Expiry    string  `json:"expiry_date"`
		QoS       int     `json:"qos"`
		Signature string  `json:"signature"`
		PieceCID  big.Int `json:"piece_cid"`
	} `json:"cid_offers"`
}

// ProviderDHTPublishGroupCIDResponse is the response to ProviderDHTPublishGroupCIDRequest
type ProviderDHTPublishGroupCIDResponse struct {
	CommonFields CommonMessageFields `json:"common_fields"`
	Nonce        int                 `json:"nonce"`
	Signature    string              `json:"signature"`
}

// GatewaySingleCIDOfferPublishRequest is the request from gateway to provider during start-up asking for cid offers
type GatewaySingleCIDOfferPublishRequest struct {
	CommonFields CommonMessageFields `json:"common_fields"`
	Nonce        int                 `json:"nonce"`
	Signature    string              `json:"signature"`
}

// GatewaySingleCIDOfferPublishResponse is the repsonse to GatewaySingleCIDOfferPublishRequest
type GatewaySingleCIDOfferPublishResponse struct {
	CommonFields CommonMessageFields `json:"common_fields"`
	CIDOffers    []struct {
		Price     int     `json:"price_per_byte"`
		Expiry    string  `json:"expiry_date"`
		QoS       int     `json:"qos"`
		Signature string  `json:"signature"`
		PieceCID  big.Int `json:"piece_cid"`
	} `json:"cid_offers"` // TODO: Need to check if this is right.
}

// GatewaySingleCIDOfferPublishResponseAck is the acknowledgement to GatewaySingleCIDOfferPublishResponse
type GatewaySingleCIDOfferPublishResponseAck struct {
	CommonFields CommonMessageFields `json:"common_fields"`
	CIDOffersAck []struct {
		Nonce     int    `json:"nonce"`
		Signature string `json:"signature"`
	} `json:"cid_offers_ack"` // TODO: Need to check if this is right.
}

// GatewayDHTDiscoverRequest is the request from gateway to discover cid offer
type GatewayDHTDiscoverRequest struct {
	CommonFields CommonMessageFields `json:"common_fields"`
	PieceCID     big.Int             `json:"piece_cid"`
	Nonce        int                 `json:"nonce"`
	TTL          string              `json:"ttl"`
}

// GatewayDHTDiscoverResponse is the response to GatewayDHTDiscoverRequest
type GatewayDHTDiscoverResponse struct {
	CommonFields CommonMessageFields   `json:"common_fields"`
	PieceCID     big.Int               `json:"piece_cid"`
	Nonce        int                   `json:"nonce"`
	Found        bool                  `json:"found"`
	Signature    string                `json:"signature"`
	CIDGroupInfo []CIDGroupInformation `json:"cid_group_information"`
}
