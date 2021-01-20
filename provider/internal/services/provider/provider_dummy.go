package provider

import (
	"math/rand"
	"time"
)

// CIDMessage data model
type CIDMessage struct {
	ProtocolVersion   string
	ProtocolSupported string
	MessageType       string
	Nonce             int
	ProviderID        string
	PricePerByte      int
	ExpiryDate        int64
	QosMetric         int
	Signature         string
	PieceCID          []string
}

func generateDummyMessage() CIDMessage {
	expiryDate := time.Now().Local().Add(time.Hour * time.Duration(24)).Unix()

	pieceCIDs := []string{"a", "b", "c", "d", "e"}
	dummyMessage := CIDMessage{
		ProtocolVersion:   "ProtocolVersion",
		ProtocolSupported: "ProtocolSupported",
		MessageType:       "MessageType",
		Nonce:             rand.Intn(100000),
		ProviderID:        "ProviderID",
		PricePerByte:      42,
		ExpiryDate:        expiryDate,
		QosMetric:         42,
		Signature:         "Signature",
		PieceCID:          pieceCIDs,
	}
	return dummyMessage
}
