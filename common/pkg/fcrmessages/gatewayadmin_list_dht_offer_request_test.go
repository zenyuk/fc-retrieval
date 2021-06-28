package fcrmessages

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEncodeGatewayAdminListDHTOfferRequest success test
func TestEncodeGatewayAdminListDHTOfferRequest(t *testing.T) {

	validMsg := &FCRMessage{
		messageType:       410,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"refresh":true}`),
	}

	msg, err := EncodeGatewayAdminListDHTOfferRequest(true)
	assert.Empty(t, err)
	assert.Equal(t, validMsg, msg)
}

// TestDecodeGatewayAdminListDHTOfferRequest success test
func TestDecodeGatewayAdminListDHTOfferRequest(t *testing.T) {
	validMsg := &FCRMessage{
		messageType:       410,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"refresh":true}`),
	}

	enrolled, err := DecodeGatewayAdminListDHTOfferRequest(validMsg)

	assert.Empty(t, err)
	assert.Equal(t, true, enrolled)
}
