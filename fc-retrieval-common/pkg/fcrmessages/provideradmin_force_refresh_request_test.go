package fcrmessages

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEncodeProviderAdminForceRefreshRequest success test
func TestEncodeProviderAdminForceRefreshRequest(t *testing.T) {

	validMsg := &FCRMessage{
		messageType:       508,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"refresh":true}`),
	}

	msg, err := EncodeProviderAdminForceRefreshRequest(true)
	assert.Empty(t, err)
	assert.Equal(t, validMsg, msg)
}

// TestDecodeProviderAdminForceRefreshRequest success test
func TestDecodeProviderAdminForceRefreshRequest(t *testing.T) {
	validMsg := &FCRMessage{
		messageType:       508,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"refresh":true}`),
	}

	enrolled, err := DecodeProviderAdminForceRefreshRequest(validMsg)

	assert.Empty(t, err)
	assert.Equal(t, true, enrolled)
}
