package fcrmessages

import (
	"testing"

	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	"github.com/stretchr/testify/assert"
)

// TestEncodeGatewayDHTDiscoverOfferRequest success test
func TestEncodeGatewayDHTDiscoverOfferRequest(t *testing.T) {
	mockContentID, _ := cid.NewContentIDFromBytes([]byte{1})
	mockNonce := int64(42)

	mockOffersDigest := [][cidoffer.CIDOfferDigestSize]byte{{1, 2, 4}, {5, 6, 7}}
	mockPayChAddr := "t2twbvr2oaxqzyktxqqjrv37bh7gzfhuqonfioayq"
	mockVoucher := "i1UCnYNY6cC8M4VO8IJjXfwn-bJT0g4AAED2AAFJAEVjkYJE9AAAAIBYYQK3pJLhIR8XTVSmQzsEiE7NIId2-2DPbWF396mBPBJCdoSQ_ctibPesW-YMnzKhGAEScF09H_sldF1nTfizTbsjWea9MN6R3T0Ew0Lb4znHtJnucGAkcbdlIyDAHCScOXE"
	validMsg := &FCRMessage{
		messageType:       209,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"piece_cid":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE=","nonce":42,"offers_digest":[[1,2,4,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],[5,6,7,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]],"payment_channel_address":"t2twbvr2oaxqzyktxqqjrv37bh7gzfhuqonfioayq","voucher":"i1UCnYNY6cC8M4VO8IJjXfwn-bJT0g4AAED2AAFJAEVjkYJE9AAAAIBYYQK3pJLhIR8XTVSmQzsEiE7NIId2-2DPbWF396mBPBJCdoSQ_ctibPesW-YMnzKhGAEScF09H_sldF1nTfizTbsjWea9MN6R3T0Ew0Lb4znHtJnucGAkcbdlIyDAHCScOXE"}`),
		signature:         "",
	}

	msg, err := EncodeGatewayDHTDiscoverOfferRequest(
		mockContentID,
		mockNonce,
		mockOffersDigest,
		mockPayChAddr,
		mockVoucher,
	)

	assert.Empty(t, err)
	assert.Equal(t, msg, validMsg)
}

// TestDecodeGatewayDHTDiscoverOfferRequest success test
func TestDecodeGatewayDHTDiscoverOfferRequest(t *testing.T) {
	mockContentID, _ := cid.NewContentIDFromBytes([]byte{1})
	mockNonce := int64(42)
	mockOfferDigest := [][cidoffer.CIDOfferDigestSize]byte{{1, 2, 4}, {5, 6, 7}}
	mockPayChAddr := "t2twbvr2oaxqzyktxqqjrv37bh7gzfhuqonfioayq"
	mockVoucher := "i1UCnYNY6cC8M4VO8IJjXfwn-bJT0g4AAED2AAFJAEVjkYJE9AAAAIBYYQK3pJLhIR8XTVSmQzsEiE7NIId2-2DPbWF396mBPBJCdoSQ_ctibPesW-YMnzKhGAEScF09H_sldF1nTfizTbsjWea9MN6R3T0Ew0Lb4znHtJnucGAkcbdlIyDAHCScOXE"
	validMsg := &FCRMessage{
		messageType:       209,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"piece_cid":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE=","nonce":42,"offers_digest":[[1,2,4,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],[5,6,7,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]],"payment_channel_address":"t2twbvr2oaxqzyktxqqjrv37bh7gzfhuqonfioayq","voucher":"i1UCnYNY6cC8M4VO8IJjXfwn-bJT0g4AAED2AAFJAEVjkYJE9AAAAIBYYQK3pJLhIR8XTVSmQzsEiE7NIId2-2DPbWF396mBPBJCdoSQ_ctibPesW-YMnzKhGAEScF09H_sldF1nTfizTbsjWea9MN6R3T0Ew0Lb4znHtJnucGAkcbdlIyDAHCScOXE"}`),
		signature:         "",
	}

	pieceCID, nonce, offerDigest, paychAddr, voucher, err := DecodeGatewayDHTDiscoverOfferRequest(validMsg)
	assert.Empty(t, err)
	assert.Equal(t, pieceCID, mockContentID)
	assert.Equal(t, nonce, mockNonce)
	assert.Equal(t, offerDigest, mockOfferDigest)
	assert.Equal(t, paychAddr, mockPayChAddr)
	assert.Equal(t, voucher, mockVoucher)
}
