package fcrmessages

import (
	"testing"
	// "errors"
	"github.com/stretchr/testify/assert"
)

func TestEncodeProtocolChangeRequest(t *testing.T) {
	var desiredVersion int32 = 42
	version, err := EncodeProtocolChangeRequest(desiredVersion)
	assert.Empty(t, err)
	assert.Equal(t, version, "temp")
}

func TestDecodeProtocolChangeRequest(t *testing.T) {
	var mockMessageBodyValid = []byte(`{"payment_channel_id":42}`)
	validMsg := &FCRMessage{
		messageType:900,
		protocolVersion:1,
		protocolSupported:[]int32{1, 1},
		messageBody:mockMessageBodyValid, 
		signature:"",
	}
	msg, err := DecodeProtocolChangeRequest(validMsg)
	assert.Empty(t, err)
	assert.Equal(t, msg, "temp")
}

