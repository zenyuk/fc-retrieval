package fcrmessages

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ConsenSys/fc-retrieval/common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval/common/pkg/cidoffer"
)

// TestEncodeGatewayDHTDiscoverResponseV2 success test
func TestEncodeGatewayDHTDiscoverResponseV2(t *testing.T) {
	mockNonce := int64(42)
	mockContentID, _ := cid.NewContentIDFromBytes([]byte{1})
	mockFound := true
	mockFPCs := []bool{true}
	mockSubCIDOfferDigest := cidoffer.EncodeMessageDigest([cidoffer.CIDOfferDigestSize]byte{1, 2, 4})
	mockSubCIDOfferDigests := []string{mockSubCIDOfferDigest}
	fakePaymentRequired := true
	fakePaymentChannel := "43"

	validMsg := &FCRMessage{
		messageType:       208,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"piece_cid":"0000000000000000000000000000000000000000000000000000000000000001","nonce":42,"found":true,"sub_cid_offer_digest":["AQIEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="],"funded_payment_channel":[true],"payment_required":true,"payment_channel":"43"}`),
		signature:         "",
	}

	msg, err := EncodeGatewayDHTDiscoverResponseV2(mockContentID, mockNonce, mockFound, mockSubCIDOfferDigests, mockFPCs, fakePaymentRequired, fakePaymentChannel)
	assert.Empty(t, err)
	assert.Equal(t, msg, validMsg)
}

// TestDecodeGatewayDHTDiscoverResponseV2 success test
func TestDecodeGatewayDHTDiscoverResponseV2(t *testing.T) {
	mockNonce := int64(42)
	mockContentID, _ := cid.NewContentIDFromBytes([]byte{1})
	mockFound := true
	mockFPCs := []bool{true}
	mockSubCIDOfferDigest := cidoffer.EncodeMessageDigest([cidoffer.CIDOfferDigestSize]byte{1, 2, 4})
	mockSubCIDOfferDigests := []string{mockSubCIDOfferDigest}

	validMsg := &FCRMessage{
		messageType:       208,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"piece_cid":"0000000000000000000000000000000000000000000000000000000000000001","nonce":42,"found":true,"sub_cid_offer_digest":["AQIEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="],"funded_payment_channel":[true],"payment_required":true,"payment_channel":"43"}`),
		signature:         "",
	}
	fakePaymentRequired := true
	fakePaymentChannel := "43"

	contentID, nonce, found, subCIDOfferDigest, FPCs, paymentRequired, paymentChannel, err := DecodeGatewayDHTDiscoverResponseV2(validMsg)
	assert.Empty(t, err)
	assert.Equal(t, contentID, mockContentID)
	assert.Equal(t, nonce, mockNonce)
	assert.Equal(t, found, mockFound)
	assert.Equal(t, subCIDOfferDigest, mockSubCIDOfferDigests)
	assert.Equal(t, FPCs, mockFPCs)
	assert.Equal(t, fakePaymentRequired, paymentRequired)
	assert.Equal(t, fakePaymentChannel, paymentChannel)
}
