package fcrmessages

import (
	"testing"

	"github.com/ConsenSys/fc-retrieval/common/pkg/nodeid"
	"github.com/stretchr/testify/assert"
)

// TestEncodeProviderAdminGetPublishedOfferRequest success test
func TestEncodeProviderAdminGetPublishedOfferRequest(t *testing.T) {
	mockNodeID, _ := nodeid.NewNodeIDFromHexString("42")
	mockNodeIDs := []nodeid.NodeID{*mockNodeID}

	validMsg := &FCRMessage{
		messageType:       506,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"gateway_id":["0000000000000000000000000000000000000000000000000000000000000042"]}`),
		signature:         "",
	}

	msg, err := EncodeProviderAdminGetPublishedOfferRequest(mockNodeIDs)
	assert.Empty(t, err)
	assert.Equal(t, msg, validMsg)
}

// TestDecodeProviderAdminGetPublishedOfferRequest success test
func TestDecodeProviderAdminGetPublishedOfferRequest(t *testing.T) {

	mockNodeID, _ := nodeid.NewNodeIDFromHexString("42")
	mockNodeIDs := []nodeid.NodeID{*mockNodeID}
	validMsg := &FCRMessage{
		messageType:       506,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"gateway_id":["0000000000000000000000000000000000000000000000000000000000000042"]}`),
		signature:         "",
	}

	nodeIDs, err := DecodeProviderAdminGetPublishedOfferRequest(validMsg)
	assert.Empty(t, err)
	assert.Equal(t, nodeIDs, mockNodeIDs)
}
