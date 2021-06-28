package fcrmessages

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEncodeGatewayAdminForceRefreshResponse success test
func TestEncodeGatewayAdminForceRefreshResponse(t *testing.T) {

	validMsg := &FCRMessage{
		messageType:       407,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"refreshed":true}`),
	}

	msg, err := EncodeGatewayAdminForceRefreshResponse(true)
	assert.Empty(t, err)
	assert.Equal(t, validMsg, msg)
}

// TestDecodeGatewayAdminForceRefreshResponse success test
func TestDecodeGatewayAdminForceRefreshResponse(t *testing.T) {
	validMsg := &FCRMessage{
		messageType:       407,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"refreshed":true}`),
	}

	enrolled, err := DecodeGatewayAdminForceRefreshResponse(validMsg)

	assert.Empty(t, err)
	assert.Equal(t, true, enrolled)
}
