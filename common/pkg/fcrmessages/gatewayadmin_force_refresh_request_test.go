package fcrmessages

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEncodeGatewayAdminForceRefreshRequest success test
func TestEncodeGatewayAdminForceRefreshRequest(t *testing.T) {

	validMsg := &FCRMessage{
		messageType:       406,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"refresh":true}`),
	}

	msg, err := EncodeGatewayAdminForceRefreshRequest(true)
	assert.Empty(t, err)
	assert.Equal(t, validMsg, msg)
}

// TestDecodeGatewayAdminForceRefreshRequest success test
func TestDecodeGatewayAdminForceRefreshRequest(t *testing.T) {
	validMsg := &FCRMessage{
		messageType:       406,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"refresh":true}`),
	}

	enrolled, err := DecodeGatewayAdminForceRefreshRequest(validMsg)

	assert.Empty(t, err)
	assert.Equal(t, true, enrolled)
}
