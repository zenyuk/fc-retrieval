package fcrmessages

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEncodeProviderAdminForceRefreshResponse success test
func TestEncodeProviderAdminForceRefreshResponse(t *testing.T) {

	validMsg := &FCRMessage{
		messageType:       509,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"refreshed":true}`),
	}

	msg, err := EncodeProviderAdminForceRefreshResponse(true)
	assert.Empty(t, err)
	assert.Equal(t, validMsg, msg)
}

// TestDecodeProviderAdminForceRefreshResponse success test
func TestDecodeProviderAdminForceRefreshResponse(t *testing.T) {
	validMsg := &FCRMessage{
		messageType:       509,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"refreshed":true}`),
	}

	enrolled, err := DecodeProviderAdminForceRefreshResponse(validMsg)

	assert.Empty(t, err)
	assert.Equal(t, true, enrolled)
}
