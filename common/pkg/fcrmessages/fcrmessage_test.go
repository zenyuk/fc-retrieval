package fcrmessages

import (
	"testing"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/stretchr/testify/assert"
)

const (
	PrivKey     = "015ed053eab6fdf18c03954373ff7f89089992017d56beb8b05305b19800d6afe0"
	PubKey      = "01047799f37b014564e23578447d718e5c70a786b0e4e58ca25cb2a086b822434594d910b9b8c0fcbfe9f4c2db321e874819e0614be5b57fbb5080accd69adb2eaad"
	PubKeyWrong = "01047799f37b014564e23578447d718e5c70a786b0e4e58ca25cb2a086b822434594d910b9b8c0fcbfe9f4c2db321e874819e0614be5b57fbb5080accd69adb2eaaa"
)

// TestFCRMessage success test
func TestFCRMessage(t *testing.T) {
	mockMsgType := int32(203)
	mockMsgBody := []byte(`{"gateway_id":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEI=","piece_cid":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE=","nonce":42,"ttl":43,"payment_channel_address":"t2twbvr2oaxqzyktxqqjrv37bh7gzfhuqonfioayq","voucher":"i1UCnYNY6cC8M4VO8IJjXfwn-"}`)
	validMsg := &FCRMessage{
		messageType:       203,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"gateway_id":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEI=","piece_cid":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE=","nonce":42,"ttl":43,"payment_channel_address":"t2twbvr2oaxqzyktxqqjrv37bh7gzfhuqonfioayq","voucher":"i1UCnYNY6cC8M4VO8IJjXfwn-"}`),
		signature:         "",
	}

	msg := CreateFCRMessage(mockMsgType, mockMsgBody)

	assert.Equal(t, msg, validMsg)
	assert.Equal(t, msg.GetMessageType(), validMsg.GetMessageType())
	assert.Equal(t, msg.GetProtocolVersion(), validMsg.GetProtocolVersion())
	assert.Equal(t, msg.GetProtocolSupported(), validMsg.GetProtocolSupported())
	assert.Equal(t, msg.GetMessageBody(), validMsg.GetMessageBody())
	assert.Equal(t, msg.GetSignature(), validMsg.GetSignature())
}

// TestSign success test
func TestSign(t *testing.T) {
	mockMsgType := int32(203)
	mockMsgBody := []byte(`{"gateway_id":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEI=","piece_cid":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE=","nonce":42,"ttl":43,"payment_channel_address":"t2twbvr2oaxqzyktxqqjrv37bh7gzfhuqonfioayq","voucher":"i1UCnYNY6cC8M4VO8IJjXfwn-"}`)
	mockPrivKey, _ := fcrcrypto.DecodePrivateKey(PrivKey)
	mockKeyVer := fcrcrypto.InitialKeyVersion()
	validMsg := &FCRMessage{
		messageType:       203,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"gateway_id":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEI=","piece_cid":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE=","nonce":42,"ttl":43,"payment_channel_address":"t2twbvr2oaxqzyktxqqjrv37bh7gzfhuqonfioayq","voucher":"i1UCnYNY6cC8M4VO8IJjXfwn-"}`),
		signature:         "00000001278d023bbdab1fa9ebd2ce74c18cda2eb77ad0990fce3595be5dc0b36da0a3341350683a5cf9e5de925d03180ea8a00ce01128d35a807fe91946ed7e84c0d0b900",
	}

	msg := CreateFCRMessage(mockMsgType, mockMsgBody)
	msg.Sign(mockPrivKey, mockKeyVer)

	assert.Equal(t, msg, validMsg)
	assert.Equal(t, msg.GetSignature(), validMsg.GetSignature())
}

// TestVerify success test
func TestVerify(t *testing.T) {
	mockMsgType := int32(203)
	mockMsgBody := []byte(`{"gateway_id":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEI=","piece_cid":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE=","nonce":42,"ttl":43,"payment_channel_address":"t2twbvr2oaxqzyktxqqjrv37bh7gzfhuqonfioayq","voucher":"i1UCnYNY6cC8M4VO8IJjXfwn-"}`)
	mockPrivKey, _ := fcrcrypto.DecodePrivateKey(PrivKey)
	mockPubKey, _ := fcrcrypto.DecodePublicKey(PubKey)
	mockKeyVer := fcrcrypto.InitialKeyVersion()
	validMsg := &FCRMessage{
		messageType:       203,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"gateway_id":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEI=","piece_cid":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE=","nonce":42,"ttl":43,"payment_channel_address":"t2twbvr2oaxqzyktxqqjrv37bh7gzfhuqonfioayq","voucher":"i1UCnYNY6cC8M4VO8IJjXfwn-"}`),
		signature:         "00000001278d023bbdab1fa9ebd2ce74c18cda2eb77ad0990fce3595be5dc0b36da0a3341350683a5cf9e5de925d03180ea8a00ce01128d35a807fe91946ed7e84c0d0b900",
	}

	msg := CreateFCRMessage(mockMsgType, mockMsgBody)
	msg.Sign(mockPrivKey, mockKeyVer)
	err := msg.Verify(mockPubKey)

	assert.Equal(t, msg, validMsg)
	assert.Equal(t, msg.GetSignature(), validMsg.GetSignature())
	assert.Empty(t, err)
}

// TestFCRMsgDump success test
func TestFCRMsgDump(t *testing.T) {
	mockMsgType := int32(203)
	mockMsgBody := []byte(`{"gateway_id":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEI=","piece_cid":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE=","nonce":42,"ttl":43,"payment_channel_address":"t2twbvr2oaxqzyktxqqjrv37bh7gzfhuqonfioayq","voucher":"i1UCnYNY6cC8M4VO8IJjXfwn-"}`)
	validMsg := &FCRMessage{
		messageType:       203,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"gateway_id":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEI=","piece_cid":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE=","nonce":42,"ttl":43,"payment_channel_address":"t2twbvr2oaxqzyktxqqjrv37bh7gzfhuqonfioayq","voucher":"i1UCnYNY6cC8M4VO8IJjXfwn-"}`),
		signature:         "",
	}
	validDump := validMsg.DumpMessage()

	msg := CreateFCRMessage(mockMsgType, mockMsgBody)
	dump := msg.DumpMessage()

	assert.Equal(t, msg, validMsg)
	assert.Equal(t, dump, validDump)
}

// TestFCRMsgToBytes success test
func TestFCRMsgToBytes(t *testing.T) {
	mockMsgType := int32(203)
	mockMsgBody := []byte(`{"gateway_id":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEI=","piece_cid":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE=","nonce":42,"ttl":43,"payment_channel_address":"t2twbvr2oaxqzyktxqqjrv37bh7gzfhuqonfioayq","voucher":"i1UCnYNY6cC8M4VO8IJjXfwn-"}`)
	validMsg := &FCRMessage{
		messageType:       203,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"gateway_id":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEI=","piece_cid":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE=","nonce":42,"ttl":43,"payment_channel_address":"t2twbvr2oaxqzyktxqqjrv37bh7gzfhuqonfioayq","voucher":"i1UCnYNY6cC8M4VO8IJjXfwn-"}`),
		signature:         "",
	}
	validFCRMsgToBytes, _ := validMsg.FCRMsgToBytes()

	msg := CreateFCRMessage(mockMsgType, mockMsgBody)
	FCRMsgToBytes, _ := msg.FCRMsgToBytes()

	assert.Equal(t, msg, validMsg)
	assert.Equal(t, FCRMsgToBytes, validFCRMsgToBytes)
}

// TestFCRMsgFromBytes success test
func TestFCRMsgFromBytes(t *testing.T) {
	mockMsgType := int32(203)
	mockMsgBody := []byte(`{"gateway_id":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEI=","piece_cid":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE=","nonce":42,"ttl":43,"payment_channel_address":"t2twbvr2oaxqzyktxqqjrv37bh7gzfhuqonfioayq","voucher":"i1UCnYNY6cC8M4VO8IJjXfwn-"}`)
	validMsg := &FCRMessage{
		messageType:       203,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"gateway_id":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEI=","piece_cid":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE=","nonce":42,"ttl":43,"payment_channel_address":"t2twbvr2oaxqzyktxqqjrv37bh7gzfhuqonfioayq","voucher":"i1UCnYNY6cC8M4VO8IJjXfwn-"}`),
		signature:         "",
	}

	msg := CreateFCRMessage(mockMsgType, mockMsgBody)
	FCRMsgToBytes, _ := msg.FCRMsgToBytes()
	FCRMsgFromBytes, _ := FCRMsgFromBytes(FCRMsgToBytes)

	assert.Equal(t, msg, validMsg)
	assert.Equal(t, FCRMsgFromBytes, validMsg)
}

// TestMarshalJSON success test
func TestMarshalJSON(t *testing.T) {
	mockMsgType := int32(203)
	mockMsgBody := []byte(`{"gateway_id":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEI=","piece_cid":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE=","nonce":42,"ttl":43,"payment_channel_address":"t2twbvr2oaxqzyktxqqjrv37bh7gzfhuqonfioayq","voucher":"i1UCnYNY6cC8M4VO8IJjXfwn-"}`)
	validMsg := &FCRMessage{
		messageType:       203,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"gateway_id":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEI=","piece_cid":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE=","nonce":42,"ttl":43,"payment_channel_address":"t2twbvr2oaxqzyktxqqjrv37bh7gzfhuqonfioayq","voucher":"i1UCnYNY6cC8M4VO8IJjXfwn-"}`),
		signature:         "",
	}
	validJson, _ := validMsg.MarshalJSON()

	msg := CreateFCRMessage(mockMsgType, mockMsgBody)
	json, _ := msg.MarshalJSON()

	assert.Equal(t, msg, validMsg)
	assert.Equal(t, json, validJson)
}

// TestUnmarshalJSON success test
func TestUnmarshalJSON(t *testing.T) {
	mockMsgType := int32(203)
	mockMsgBody := []byte(`{"gateway_id":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEI=","piece_cid":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE=","nonce":42,"ttl":43,"payment_channel_address":"t2twbvr2oaxqzyktxqqjrv37bh7gzfhuqonfioayq","voucher":"i1UCnYNY6cC8M4VO8IJjXfwn-"}`)
	validMsg := &FCRMessage{
		messageType:       203,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"gateway_id":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEI=","piece_cid":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE=","nonce":42,"ttl":43,"payment_channel_address":"t2twbvr2oaxqzyktxqqjrv37bh7gzfhuqonfioayq","voucher":"i1UCnYNY6cC8M4VO8IJjXfwn-"}`),
		signature:         "",
	}

	msg := CreateFCRMessage(mockMsgType, mockMsgBody)
	json, _ := msg.MarshalJSON()
	FCRMsg := &FCRMessage{}
	err := FCRMsg.UnmarshalJSON(json)

	assert.Equal(t, msg, validMsg)
	assert.Empty(t, err)
	assert.Equal(t, FCRMsg, validMsg)
	assert.Equal(t, FCRMsg.GetMessageType(), validMsg.GetMessageType())
	assert.Equal(t, FCRMsg.GetProtocolVersion(), validMsg.GetProtocolVersion())
	assert.Equal(t, FCRMsg.GetProtocolSupported(), validMsg.GetProtocolSupported())
	assert.Equal(t, FCRMsg.GetMessageBody(), validMsg.GetMessageBody())
	assert.Equal(t, FCRMsg.GetSignature(), validMsg.GetSignature())
}
