package fcrmessages

import (
	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestEncodeClientDHTDiscoverOfferResponse success test
func TestEncodeClientDHTDiscoverOfferResponse(t *testing.T) {
	mockContentID, _ := cid.NewContentIDFromBytes([]byte{1})
	mockNonce := int64(42)
	mockNodeID, _ := nodeid.NewNodeIDFromHexString("42")

	mockMessage := &FCRMessage{
		messageType:       115,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(``),
		signature:         "",
	}
	mockResponse := []FCRMessage{*mockMessage}

	validMsg := &FCRMessage{
		messageType:       115,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"piece_cid":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE=","nonce":42,"gateway_ids":["AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEI="],"response":[{"message_type":115,"protocol_version":1,"protocol_supported":[1,1],"message_body":"","message_signature":""}]}`),
		signature:         "",
	}

	msg, err := EncodeClientDHTDiscoverOfferResponse(
		mockContentID,
		mockNonce,
		[]nodeid.NodeID{*mockNodeID},
		mockResponse,
	)
	assert.Empty(t, err)
	assert.Equal(t, msg, validMsg)
}

// TestDecodeClientDHTDiscoverOfferResponse success test
func TestDecodeClientDHTDiscoverOfferResponse(t *testing.T) {
	mockContentID, _ := cid.NewContentIDFromBytes([]byte{1})
	mockNonce := int64(42)
	mockNodeID, _ := nodeid.NewNodeIDFromHexString("42")

	mockMessage := &FCRMessage{
		messageType:       115,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(``),
		signature:         "",
	}
	mockResponse := []FCRMessage{*mockMessage}

	validMsg := &FCRMessage{
		messageType:       115,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"piece_cid":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE=","nonce":42,"gateway_ids":["AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEI="],"response":[{"message_type":115,"protocol_version":1,"protocol_supported":[1,1],"message_body":"","message_signature":""}]}`),
		signature:         "",
	}

	contentID, nonce, gateway_ids, response, err := DecodeClientDHTDiscoverOfferResponse(validMsg)
	assert.Empty(t, err)
	assert.Equal(t, contentID, mockContentID)
	assert.Equal(t, nonce, mockNonce)
	assert.Equal(t, gateway_ids, []nodeid.NodeID{*mockNodeID})
	assert.Equal(t, response, mockResponse)
}
