package fcrmessages

import (
	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/stretchr/testify/assert"
	"testing"
	// "github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	// "github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
)

// TestEncodeClientStandardDiscoverResponseV2 success test
func TestEncodeClientStandardDiscoverResponseV2(t *testing.T) {
	// mockNodeID, _ := nodeid.NewNodeIDFromHexString("42")
	mockNonce := int64(42)
	mockFound := true
	mockContentID, _ := cid.NewContentIDFromBytes([]byte{1})
	// mockCids := []cid.ContentID{*mockContentID}
	// var mockPrice uint64 = 41
	// var mockExpiry int64 = 42
	// var mockQos uint64 = 43
	// offer, _ := cidoffer.NewCIDOffer(mockNodeID, mockCids, mockPrice, mockExpiry, mockQos)
	// subOffer, _ := offer.GenerateSubCIDOffer(mockContentID)
	// mockSubCidOffers := []cidoffer.SubCIDOffer{*subOffer}
	mockSubCidOffers := [][32]byte{{1, 2}}
	mockFPCs := []bool{true}
	validMsg := &FCRMessage{
		messageType:       109,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"piece_cid":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE=","nonce":42,"found":true,"sub_cid_offer_digests":[[1,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]],"funded_payment_channel":[true]}`),
		signature:         "",
	}

	msg, err := EncodeClientStandardDiscoverResponseV2(
		mockContentID,
		mockNonce,
		mockFound,
		mockSubCidOffers,
		mockFPCs,
	)
	assert.Empty(t, err)
	assert.Equal(t, msg, validMsg)
}

// TestDecodeClientStandardDiscoverResponseV2 success test
func TestDecodeClientStandardDiscoverResponseV2(t *testing.T) {
	// mockNodeID, _ := nodeid.NewNodeIDFromHexString("42")
	mockNonce := int64(42)
	mockFound := true
	mockContentID, _ := cid.NewContentIDFromBytes([]byte{1})
	// mockCids := []cid.ContentID{*mockContentID}
	// var mockPrice uint64 = 41
	// var mockExpiry int64 = 42
	// var mockQos uint64 = 43
	// offer, _ := cidoffer.NewCIDOffer(mockNodeID, mockCids, mockPrice, mockExpiry, mockQos)
	// subOffer, _ := offer.GenerateSubCIDOffer(mockContentID)
	// mockSubCidOffers := []cidoffer.SubCIDOffer{*subOffer}
	mockSubCidOffers := [][32]byte{{1, 2}}
	mockFPCs := []bool{true}
	validMsg := &FCRMessage{
		messageType:       109,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"piece_cid":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE=","nonce":42,"found":true,"sub_cid_offer_digests":[[1,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]],"funded_payment_channel":[true]}`),
		signature:         "",
	}

	contentID, nonce, found, subOffers, FPCs, err := DecodeClientStandardDiscoverResponseV2(validMsg)
	assert.Empty(t, err)
	assert.Equal(t, contentID, mockContentID)
	assert.Equal(t, nonce, mockNonce)
	assert.Equal(t, found, mockFound)
	assert.Equal(t, subOffers, mockSubCidOffers)
	assert.Equal(t, FPCs, mockFPCs)
}
