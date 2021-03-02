package fcrmessages

import (
	"encoding/json"
	"fmt"

	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmerkletree"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
)

// ProviderAdminPublishGroupCIDRequest is the request to publish group CID
type ProviderAdminPublishGroupCIDRequest struct {
	CIDs   []cid.ContentID `json:"cids"`
	Price  uint64          `json:"price"`
	Expiry int64           `json:"expiry"`
	QoS    uint64          `json:"qos"`
}

// EncodeProviderAdminPublishGroupCIDRequest is used to get the FCRMessage of ProviderAdminPublishGroupCIDRequest
func EncodeProviderAdminPublishGroupCIDRequest(
	cids []cid.ContentID,
	price uint64,
	expiry int64,
	qos uint64,
) (*FCRMessage, error) {
	body, err := json.Marshal(ProviderAdminPublishGroupCIDRequest{
		CIDs:   cids,
		Price:  price,
		Expiry: expiry,
		QoS:    qos,
	})
	if err != nil {
		return nil, err
	}
	return &FCRMessage{
		MessageType:       ProviderAdminPublishGroupCIDRequestType,
		ProtocolVersion:   protocolVersion,
		ProtocolSupported: protocolSupported,
		MessageBody:       body,
	}, nil
}

// DecodeProviderAdminPublishGroupCIDRequest is used to get the fields from FCRMessage of ProviderAdminPublishGroupCIDRequest
func DecodeProviderAdminPublishGroupCIDRequest(fcrMsg *FCRMessage) (
	[]cid.ContentID, // cids
	uint64, // price
	int64, // expity
	uint64, // qos
	error, // error
) {
	if fcrMsg.MessageType != ProviderAdminPublishGroupCIDRequestType {
		return nil, 0, 0, 0, fmt.Errorf("Message type mismatch")
	}
	msg := ProviderAdminPublishGroupCIDRequest{}
	err := json.Unmarshal(fcrMsg.MessageBody, &msg)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	return msg.CIDs, msg.Price, msg.Expiry, msg.QoS, nil
}

// ProviderAdminDHTPublishGroupCIDRequest is the request to publish single CID in DHT
type ProviderAdminDHTPublishGroupCIDRequest struct {
	CIDs   []cid.ContentID `json:"cids"`
	Price  []uint64        `json:"price"`
	Expiry []int64         `json:"expiry"`
	QoS    []uint64        `json:"qos"`
}

// EncodeProviderAdminDHTPublishGroupCIDRequest is used to get the FCRMessage of ProviderAdminDHTPublishGroupCIDRequest
func EncodeProviderAdminDHTPublishGroupCIDRequest(
	cids []cid.ContentID,
	price []uint64,
	expiry []int64,
	qos []uint64,
) (*FCRMessage, error) {
	body, err := json.Marshal(ProviderAdminDHTPublishGroupCIDRequest{
		CIDs:   cids,
		Price:  price,
		Expiry: expiry,
		QoS:    qos,
	})
	if err != nil {
		return nil, err
	}
	return &FCRMessage{
		MessageType:       ProviderAdminDHTPublishGroupCIDRequestType,
		ProtocolVersion:   protocolVersion,
		ProtocolSupported: protocolSupported,
		MessageBody:       body,
	}, nil
}

// DecodeProviderAdminDHTPublishGroupCIDRequest is used to get the fields from FCRMessage of ProviderAdminDHTPublishGroupCIDRequest
func DecodeProviderAdminDHTPublishGroupCIDRequest(fcrMsg *FCRMessage) (
	[]cid.ContentID, // cids
	[]uint64, // price
	[]int64, // expity
	[]uint64, // qos
	error, // error
) {
	if fcrMsg.MessageType != ProviderAdminDHTPublishGroupCIDRequestType {
		return nil, nil, nil, nil, fmt.Errorf("Message type mismatch")
	}
	msg := ProviderAdminDHTPublishGroupCIDRequest{}
	err := json.Unmarshal(fcrMsg.MessageBody, &msg)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return msg.CIDs, msg.Price, msg.Expiry, msg.QoS, nil
}

// ProviderAdminPublishOfferAck is the response to provider admin publish group cid or single cid offer
type ProviderAdminPublishOfferAck struct {
	Received bool `json:"received"`
}

// EncodeProviderAdminPublishOfferAck is used to get the FCRMessage of ProviderAdminPublishOfferAck
func EncodeProviderAdminPublishOfferAck(received bool) (*FCRMessage, error) {
	body, err := json.Marshal(ProviderAdminPublishOfferAck{
		Received: received,
	})
	if err != nil {
		return nil, err
	}
	return &FCRMessage{
		MessageType:       ProviderAdminPublishOfferAckType,
		ProtocolVersion:   protocolVersion,
		ProtocolSupported: protocolSupported,
		MessageBody:       body,
	}, nil
}

// DecodeProviderAdminPublishOfferAck is used to get the fields from FCRMessage of ProviderAdminPublishOfferAck
func DecodeProviderAdminPublishOfferAck(fcrMsg *FCRMessage) (
	bool, // received
	error, // error
) {
	if fcrMsg.MessageType != ProviderAdminPublishOfferAckType {
		return false, fmt.Errorf("Message type mismatch")
	}
	msg := ProviderAdminPublishOfferAck{}
	err := json.Unmarshal(fcrMsg.MessageBody, &msg)
	if err != nil {
		return false, err
	}
	return msg.Received, nil
}

// ProviderAdminGetGroupCIDRequest is the requset from client to gateway to ask for cid offer
type ProviderAdminGetGroupCIDRequest struct {
	GatewayIDs []nodeid.NodeID `json:"gateway_id"`
}

// EncodeProviderAdminGetGroupCIDRequest is used to get the FCRMessage of ProviderAdminGetGroupCIDRequest
func EncodeProviderAdminGetGroupCIDRequest(
	gatewayIDs []nodeid.NodeID,
) (*FCRMessage, error) {
	fmt.Printf("EncodeProviderAdminGetGroupCIDRequest\n")
	body, err := json.Marshal(ProviderAdminGetGroupCIDRequest{
		GatewayIDs: gatewayIDs,
	})
	fmt.Printf("Body: %+v\n", body)
	if err != nil {
		fmt.Printf("Error when marshalling request: %+v\n", err)
		return nil, err
	}
	return &FCRMessage{
		MessageType:       ProviderAdminGetGroupCIDRequestType,
		ProtocolVersion:   protocolVersion,
		ProtocolSupported: protocolSupported,
		MessageBody:       body,
	}, nil
}

// DecodeProviderAdminGetGroupCIDRequest is used to get the fields from FCRMessage of ProviderAdminGetGroupCIDRequest
func DecodeProviderAdminGetGroupCIDRequest(fcrMsg *FCRMessage) (
	[]nodeid.NodeID, // piece cids
	error, // error
) {
	if fcrMsg.MessageType != ProviderAdminGetGroupCIDRequestType {
		return nil, fmt.Errorf("Message type mismatch")
	}
	msg := ProviderAdminGetGroupCIDRequest{}
	err := json.Unmarshal(fcrMsg.MessageBody, &msg)
	if err != nil {
		return nil, err
	}
	return msg.GatewayIDs, nil
}

// ProviderAdminGetGroupCIDResponse is the response to ProviderAdminGetGroupCIDResponse
type ProviderAdminGetGroupCIDResponse struct {
	Found        bool                  `json:"found"`
	CIDGroupInfo []CIDGroupInformation `json:"cid_group_information"`
}

// EncodeProviderAdminGetGroupCIDResponse is used to get the FCRMessage of ProviderAdminGetGroupCIDResponse
func EncodeProviderAdminGetGroupCIDResponse(
	found bool,
	offers []*cidoffer.CidGroupOffer,
	roots []string,
	fundedPaymentChannel []bool,
) (*FCRMessage, error) {
	cidGroupInfo := make([]CIDGroupInformation, len(offers))
	if found {
		for i := 0; i < len(offers); i++ {
			offer := offers[i]
			cidGroupInfo[i] = CIDGroupInformation{
				ProviderID:           *offer.NodeID,
				Price:                offer.Price,
				Expiry:               offer.Expiry,
				QoS:                  offer.QoS,
				Signature:            offer.Signature,
				MerkleRoot:           roots[i],
				MerkleProof:					fcrmerkletree.FCRMerkleProof{},
				FundedPaymentChannel: fundedPaymentChannel[i],
			}
		}
	}
	body, err := json.Marshal(ProviderAdminGetGroupCIDResponse{
		Found:        found,
		CIDGroupInfo: cidGroupInfo,
	})
	if err != nil {
		return nil, err
	}
	return &FCRMessage{
		MessageType:       ProviderAdminGetGroupCIDResponseType,
		ProtocolVersion:   protocolVersion,
		ProtocolSupported: protocolSupported,
		MessageBody:       body,
	}, nil
}

// DecodeProviderAdminGetGroupCIDResponse is used to get the fields from FCRMessage of ProviderAdminGetGroupCIDResponse
func DecodeProviderAdminGetGroupCIDResponse(fcrMsg *FCRMessage) (
	bool, // found
	[]CIDGroupInformation, // offers
	error, // error
) {
	if fcrMsg.MessageType != ProviderAdminGetGroupCIDResponseType {
		return false, nil, fmt.Errorf("Message type mismatch")
	}
	msg := ProviderAdminGetGroupCIDResponse{}
	err := json.Unmarshal(fcrMsg.MessageBody, &msg)
	if err != nil {
		return false, nil, err
	}
	return msg.Found, msg.CIDGroupInfo, nil
}
