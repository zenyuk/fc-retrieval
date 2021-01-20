package provider

import (
	"math/rand"
	"time"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)

// CIDMessage data model
type CIDMessage struct {
	MessageType       int32           `json:"message_type"`
	ProtocolVersion   int32           `json:"protocol_version"`
	ProtocolSupported []int32         `json:"protocol_supported"`
	Nonce             int64           `json:"nonce"`
	ProviderID        nodeid.NodeID   `json:"provider_id"`
	Price             uint64          `json:"price_per_byte"`
	Expiry            int64           `json:"expiry_date"`
	QoS               uint64          `json:"qos"`
	Signature         string          `json:"signature"`
	PieceCIDs         []cid.ContentID `json:"piece_cids"`
}

func generateDummyMessage() CIDMessage {
	expiryDate := time.Now().Local().Add(time.Hour * time.Duration(24)).Unix()

	protocolSupported := []int32{1, 2}
	providerID, _ := nodeid.NewRandomNodeID()

	contentID, _ := cid.NewRandomContentID()
	// pieceCIDs := []string{"a", "b", "c", "d", "e"}
	pieceCIDs := []cid.ContentID{*contentID}
	dummyMessage := CIDMessage{

		MessageType:       123,
		ProtocolVersion:   1,
		ProtocolSupported: protocolSupported,
		Nonce:             rand.Int63n(100000),
		ProviderID:        *providerID,
		Price:             42,
		Expiry:            expiryDate,
		QoS:               42,
		Signature:         "Signature",
		PieceCIDs:         pieceCIDs,
	}
	return dummyMessage
}
