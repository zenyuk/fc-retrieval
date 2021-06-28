package fcrmessages

import (
	"testing"

	"github.com/ConsenSys/fc-retrieval-common/pkg/challenge"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/stretchr/testify/assert"
)

// TestEncodeClientEstablishmentResponse success test
func TestEncodeClientEstablishmentResponse(t *testing.T) {
	mockNodeID, _ := nodeid.NewNodeIDFromHexString("42")
	mockChallenge := challenge.NewRandomChallenge()
	validMsg := &FCRMessage{
		messageType:       101,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"gateway_id":"0000000000000000000000000000000000000000000000000000000000000042","challenge":"` + mockChallenge + `"}`),
		signature:         "",
	}

	msg, err := EncodeClientEstablishmentResponse(
		mockNodeID,
		mockChallenge,
	)
	assert.Empty(t, err)
	assert.Equal(t, msg, validMsg)
}

// TestDecodeClientEstablishmentResponse success test
func TestDecodeClientEstablishmentResponse(t *testing.T) {
	mockNodeID, _ := nodeid.NewNodeIDFromHexString("42")
	mockChallenge := challenge.NewRandomChallenge()
	validMsg := &FCRMessage{
		messageType:       101,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"gateway_id":"0000000000000000000000000000000000000000000000000000000000000042","challenge":"` + mockChallenge + `"}`),
		signature:         "",
	}

	gatewayID, challenge, err := DecodeClientEstablishmentResponse(validMsg)
	assert.Empty(t, err)
	assert.Equal(t, gatewayID, mockNodeID)
	assert.Equal(t, challenge, mockChallenge)
}
