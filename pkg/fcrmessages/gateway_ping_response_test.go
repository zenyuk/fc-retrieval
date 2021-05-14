package fcrmessages

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEncodeGatewayPingResponse success test
func TestEncodeGatewayPingResponse(t *testing.T) {
	mockIsAlive := true
	mockNonce := int64(42)

	validMsg := &FCRMessage{
		messageType:       206,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"nonce":42,"isAlive":true}`),
		signature:         "",
	}

	msg, err := EncodeGatewayPingResponse(mockNonce, mockIsAlive)
	assert.Empty(t, err)
	assert.Equal(t, msg, validMsg)
}

// TestDecodeGatewayPingResponse success test
func TestDecodeGatewayPingResponse(t *testing.T) {
	mockIsAlive := true
	mockNonce := int64(42)

	validMsg := &FCRMessage{
		messageType:       206,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"nonce":42,"isAlive":true}`),
		signature:         "",
	}

	nonce, isAlive, err := DecodeGatewayPingResponse(validMsg)
	assert.Empty(t, err)
	assert.Equal(t, nonce, mockNonce)
	assert.Equal(t, isAlive, mockIsAlive)
}

// TestDecodeGatewayPingRequest failure test
func TestDecodeGatewayPingResponseWrongMessageType(t *testing.T) {
	validMsg := &FCRMessage{
		messageType:       206000,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"dummy":43}`),
		signature:         "",
	}

	nonce, isAlive, err := DecodeGatewayPingResponse(validMsg)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "message type mismatch")
	assert.False(t, isAlive)
	assert.Empty(t, nonce)
}

// TestDecodeGatewayPingRequest failure test
func TestDecodeGatewayPingResponseUnmarshalError(t *testing.T) {
	validMsg := &FCRMessage{
		messageType:       206,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`invalid_message`),
		signature:         "",
	}

	nonce, isAlive, err := DecodeGatewayPingResponse(validMsg)
	assert.NotNil(t, err)
	assert.NotEmpty(t, err.Error())
	assert.Empty(t, nonce)
	assert.False(t, isAlive)
}
