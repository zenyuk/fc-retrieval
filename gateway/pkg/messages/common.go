package messages

import (
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)

// CommonRequestMessageFields are shared fileds of all messages // TODO: Do we really need this?
type CommonRequestMessageFields struct {
	MessageType       int32   `json:"message_type"`
	ProtocolVersion   int32   `json:"protocol_version"`
	ProtocolSupported []int32 `json:"protocol_supported"`
}

// ProtocolChangeResponse message is sent to indicate that the gateway is requesting the
// other entity change protocol version
type ProtocolChangeResponse struct {
	MessageType    int32 `json:"message_type"`
	DesiredVersion int32 `json:"desired_version"`
}

// ProtocolMismatchResponse message is sent to indicate that there are no common protocol
// versions between the gateway and the requesting entity.
type ProtocolMismatchResponse struct {
	MessageType int32 `json:"message_type"`
}

// InvalidMessageResponse message is sent to indicate that the message is invalid
type InvalidMessageResponse struct {
	MessageType int32 `json:"message_type"`
}

// InsufficientFundsResponse message is sent to indicate the there are not sufficient funds to finish the request
type InsufficientFundsResponse struct {
	PaymentChannelID int64 `json:"payment_channel_id"` // TODO: Is this int64?
}

// CIDGroupInformation represents a cid group information
// TODO: Maybe use class in cidgroupoffer
type CIDGroupInformation struct {
	ProviderID           nodeid.NodeID `json:"provider_id"`
	Price                int64         `json:"price_per_byte"`
	Expiry               int64         `json:"expiry_date"`
	QoS                  int64         `json:"qos"`
	Signature            string        `json:"signature"`
	MerkleProof          string        `json:"merkle_proof"`
	FundedPaymentChannel bool          `json:"funded_payment_channel"` // TODO: Is this boolean?
}
