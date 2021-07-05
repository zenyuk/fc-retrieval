package fcrmessages

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ConsenSys/fc-retrieval/common/pkg/cid"
)

// TestEncodeClientStandardDiscoverResponseV2 success test
func TestEncodeClientStandardDiscoverResponseV2(t *testing.T) {
	mockNonce := int64(42)
	mockFound := true
	mockContentID, _ := cid.NewContentIDFromBytes([]byte{1})
	mockSubCIDOfferDigests := []string{"AQIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="}
	mockFPCs := []bool{true}
	validMsg := &FCRMessage{
		messageType:       109,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"piece_cid":"0000000000000000000000000000000000000000000000000000000000000001","nonce":42,"found":true,"sub_cid_offer_digests":["AQIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="],"funded_payment_channel":[true],"payment_required":true,"payment_channel":"43"}`),
		signature:         "",
	}
	fakePaymentRequired := true
	fakePaymentChannel := "43"

	msg, err := EncodeClientStandardDiscoverResponseV2(
		mockContentID,
		mockNonce,
		mockFound,
		mockSubCIDOfferDigests,
		mockFPCs,
		fakePaymentRequired,
		fakePaymentChannel,
	)
	assert.Empty(t, err)
	assert.Equal(t, msg, validMsg)
}

// TestDecodeClientStandardDiscoverResponseV2 success test
func TestDecodeClientStandardDiscoverResponseV2(t *testing.T) {
	mockNonce := int64(42)
	mockFound := true
	mockContentID, _ := cid.NewContentIDFromBytes([]byte{1})
	mockSubCIDOfferDigests := []string{"AQIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="}
	mockFPCs := []bool{true}
	validMsg := &FCRMessage{
		messageType:       109,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"piece_cid":"0000000000000000000000000000000000000000000000000000000000000001","nonce":42,"found":true,"sub_cid_offer_digests":["AQIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="],"funded_payment_channel":[true],"payment_required":true,"payment_channel":"43"}`),
		signature:         "",
	}
	fakePaymentRequired := true
	fakePaymentChannel := "43"

	contentID, nonce, found, subCIDOfferDigests, FPCs, paymentRequired, paymentChannel, err := DecodeClientStandardDiscoverResponseV2(validMsg)
	assert.Empty(t, err)
	assert.Equal(t, contentID, mockContentID)
	assert.Equal(t, nonce, mockNonce)
	assert.Equal(t, found, mockFound)
	assert.Equal(t, subCIDOfferDigests, mockSubCIDOfferDigests)
	assert.Equal(t, FPCs, mockFPCs)
	assert.Equal(t, fakePaymentRequired, paymentRequired)
	assert.Equal(t, fakePaymentChannel, paymentChannel)
}
