package messages

// Message types
// The enum should remains the same for client,provider and gateway.
const (
	ClientEstablishmentRequestType              = 0
	ClientEstablishmentResponseType             = 1
	ClientStandardDiscoverRequestType           = 2
	ClientStandardDiscoverResponseType          = 3
	ClientDHTDiscoverRequestType                = 4
	ClientDHTDiscoverResponseType               = 5
	ProviderPublishGroupCIDRequestType          = 6
	ProviderDHTPublishGroupCIDRequestType       = 7
	ProviderDHTPublishGroupCIDResponseType      = 8
	GatewaySingleCIDOfferPublishRequestType     = 9
	GatewaySingleCIDOfferPublishResponseType    = 10
	GatewaySingleCIDOfferPublishResponseAckType = 11
	GatewayDHTDiscoverRequestType               = 12
	GatewayDHTDiscoverResponseType              = 13
	ProtocolChange                              = 100
	ProtocolMismatch                            = 101
	InvalidMessage                              = 102
	InsufficientFunds                           = 103
)
