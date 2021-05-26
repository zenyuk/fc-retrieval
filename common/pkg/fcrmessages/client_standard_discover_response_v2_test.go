package fcrmessages

import (
	"testing"

	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/stretchr/testify/assert"
)

// TestEncodeClientStandardDiscoverResponseV2 success test
func TestEncodeClientStandardDiscoverResponseV2(t *testing.T) {
	mockNonce := int64(42)
	mockFound := true
	mockContentID, _ := cid.NewContentIDFromBytes([]byte{1})
	mockSubCIDOfferDigests := [][32]byte{{1, 2}}
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
		mockSubCIDOfferDigests,
		mockFPCs,
	)
	assert.Empty(t, err)
	assert.Equal(t, msg, validMsg)
}

// TestDecodeClientStandardDiscoverResponseV2 success test
func TestDecodeClientStandardDiscoverResponseV2(t *testing.T) {
	mockNonce := int64(42)
	mockFound := true
	mockContentID, _ := cid.NewContentIDFromBytes([]byte{1})
	mockSubCIDOfferDigests := [][32]byte{{1, 2}}
	mockFPCs := []bool{true}
	validMsg := &FCRMessage{
		messageType:       109,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"piece_cid":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE=","nonce":42,"found":true,"sub_cid_offer_digests":[[1,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]],"funded_payment_channel":[true]}`),
		signature:         "",
	}

	contentID, nonce, found, subCIDOfferDigests, FPCs, err := DecodeClientStandardDiscoverResponseV2(validMsg)
	assert.Empty(t, err)
	assert.Equal(t, contentID, mockContentID)
	assert.Equal(t, nonce, mockNonce)
	assert.Equal(t, found, mockFound)
	assert.Equal(t, subCIDOfferDigests, mockSubCIDOfferDigests)
	assert.Equal(t, FPCs, mockFPCs)
}
