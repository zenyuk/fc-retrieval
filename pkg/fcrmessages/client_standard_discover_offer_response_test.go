package fcrmessages

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
)

// TestEncodeClientStandardDiscoverOfferResponse success test
func TestEncodeClientStandardDiscoverOfferResponse(t *testing.T) {
	mockNodeID, _ := nodeid.NewNodeIDFromHexString("42")
	mockNonce := int64(42)
	mockFound := true
	mockContentID, _ := cid.NewContentIDFromBytes([]byte{1})
	mockCids := []cid.ContentID{*mockContentID}
	var mockPrice uint64 = 41
	var mockExpiry int64 = 42
	var mockQos uint64 = 43
	offer, _ := cidoffer.NewCIDOffer(mockNodeID, mockCids, mockPrice, mockExpiry, mockQos)
	subOffer, _ := offer.GenerateSubCIDOffer(mockContentID)
	mockSubCidOffers := []cidoffer.SubCIDOffer{*subOffer}
	mockFPCs := []bool{true}
	validMsg := &FCRMessage{
		messageType:       111,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"piece_cid":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE=","nonce":42,"found":true,"sub_cid_offers":[{"provider_id":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEI=","sub_cid":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE=","merkle_root":"c3c3a46684c07d12a9c238787df3049a6f258e7af203e5ddb66a8bd66637e108","merkle_proof":"AAAAMFsiQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBRT0iXQAAAANbMV0=","price":41,"expiry":42,"qos":43,"signature":""}],"funded_payment_channel":[true],"payment_required":true,"payment_channel":43}`),
		signature:         "",
	}
	fakePaymentRequired := true
	fakePaymentChannel := int64(43)

	msg, err := EncodeClientStandardDiscoverOfferResponse(
		mockContentID,
		mockNonce,
		mockFound,
		mockSubCidOffers,
		mockFPCs,
		fakePaymentRequired,
		fakePaymentChannel,
	)
	assert.Empty(t, err)
	assert.Equal(t, msg, validMsg)
}

// TestDecodeClientStandardDiscoverOfferResponse success test
func TestDecodeClientStandardDiscoverOfferResponse(t *testing.T) {
	mockNodeID, _ := nodeid.NewNodeIDFromHexString("42")
	mockNonce := int64(42)
	mockFound := true
	mockContentID, _ := cid.NewContentIDFromBytes([]byte{1})
	mockCids := []cid.ContentID{*mockContentID}
	var mockPrice uint64 = 41
	var mockExpiry int64 = 42
	var mockQos uint64 = 43
	offer, _ := cidoffer.NewCIDOffer(mockNodeID, mockCids, mockPrice, mockExpiry, mockQos)
	subOffer, _ := offer.GenerateSubCIDOffer(mockContentID)
	mockSubCidOffers := []cidoffer.SubCIDOffer{*subOffer}
	mockFPCs := []bool{true}
	validMsg := &FCRMessage{
		messageType:       111,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"piece_cid":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE=","nonce":42,"found":true,"sub_cid_offers":[{"provider_id":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEI=","sub_cid":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE=","merkle_root":"c3c3a46684c07d12a9c238787df3049a6f258e7af203e5ddb66a8bd66637e108","merkle_proof":"AAAAMFsiQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBRT0iXQAAAANbMV0=","price":41,"expiry":42,"qos":43,"signature":""}],"funded_payment_channel":[true],"payment_required":true,"payment_channel":43}`),
		signature:         "",
	}
	fakePaymentRequired := true
	fakePaymentChannel := int64(43)

	contentID, nonce, found, subOffers, FPCs, paymentRequired, paymentChannel, err := DecodeClientStandardDiscoverOfferResponse(validMsg)
	assert.Empty(t, err)
	assert.Equal(t, contentID, mockContentID)
	assert.Equal(t, nonce, mockNonce)
	assert.Equal(t, found, mockFound)
	assert.Equal(t, subOffers, mockSubCidOffers)
	assert.Equal(t, FPCs, mockFPCs)
	assert.Equal(t, fakePaymentRequired, paymentRequired)
	assert.Equal(t, fakePaymentChannel, paymentChannel)
}
