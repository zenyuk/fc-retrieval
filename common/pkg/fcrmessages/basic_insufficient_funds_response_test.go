package fcrmessages

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestEncodeInsufficientFundsResponseValidMessage(t *testing.T) {
	var mockPaymentChannelIDValid int64 = 42
	var mockMessageBodyValid = []byte(`{"payment_channel_id":42}`)
	validMsg := &FCRMessage{
		messageType:903,
		protocolVersion:1, 
		protocolSupported:[]int32{1, 1}, 
		messageBody:mockMessageBodyValid, 
		signature:"",
	}

	msg, err := EncodeInsufficientFundsResponse(mockPaymentChannelIDValid)
	assert.Empty(t, err)
	assert.Equal(t, msg, validMsg)
}

func TestDecodeInsufficientFundsResponse(t *testing.T) {
	var mockPaymentChannelIDValid int64 = 42
	var mockMessageBodyValid = []byte(`{"payment_channel_id":42}`)
	validMsg := &FCRMessage{
		messageType:903,
		protocolVersion:1, 
		protocolSupported:[]int32{1, 1}, 
		messageBody:mockMessageBodyValid, 
		signature:"",
	}

	msg, err := DecodeInsufficientFundsResponse(validMsg)
	assert.Empty(t, err)
	assert.Equal(t, msg, mockPaymentChannelIDValid)
}