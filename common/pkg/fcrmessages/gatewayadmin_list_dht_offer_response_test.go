package fcrmessages

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEncodeGatewayAdminListDHTOfferResponse success test
func TestEncodeGatewayAdminListDHTOfferResponse(t *testing.T) {

	validMsg := &FCRMessage{
		messageType:       411,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"refreshed":true}`),
	}

	msg, err := EncodeGatewayAdminListDHTOfferResponse(true)
	assert.Empty(t, err)
	assert.Equal(t, validMsg, msg)
}

// TestDecodeGatewayAdminListDHTOfferResponse success test
func TestDecodeGatewayAdminListDHTOfferResponse(t *testing.T) {
	validMsg := &FCRMessage{
		messageType:       411,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"refreshed":true}`),
	}

	enrolled, err := DecodeGatewayAdminListDHTOfferResponse(validMsg)

	assert.Empty(t, err)
	assert.Equal(t, true, enrolled)
}
