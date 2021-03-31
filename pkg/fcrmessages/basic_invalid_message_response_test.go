package fcrmessages

import (
	"testing"
	"errors"
	"github.com/stretchr/testify/assert"
)

func TestEncodeInvalidMessageResponse(t *testing.T) {
	var mockMessageBodyValid = []byte(`{}`)
	emptyMsg := &FCRMessage{
		messageType:902,
		protocolVersion:1,
		protocolSupported:[]int32{1, 1},
		messageBody:mockMessageBodyValid,
		signature:"",
	}

	msg, err := EncodeInvalidMessageResponse()
	assert.Empty(t, err)
	assert.Equal(t, msg, emptyMsg)
}

func TestDecodeInvalidMessageResponseValid(t *testing.T) {
	var mockMessageBodyValid = []byte(`{}`)
	validMsg := &FCRMessage{
		messageType:902,
		protocolVersion:1,
		protocolSupported:[]int32{1, 1},
		messageBody:mockMessageBodyValid, 
		signature:"",
	}

	err := DecodeInvalidMessageResponse(validMsg)
	assert.Equal(t, err, nil)
}

func TestDecodeInvalidMessageResponseInvalid(t *testing.T) {
	var mockMessageBodyValid = []byte(`{}`)
	invalidTypeMsg := &FCRMessage{
		messageType:-1,
		protocolVersion:1,
		protocolSupported:[]int32{1, 1},
		messageBody:mockMessageBodyValid, 
		signature:"",
	}

	err := DecodeInvalidMessageResponse(invalidTypeMsg)
	assert.Equal(t, err, errors.New("Message type mismatch"))
}

