package fcrmessages

import (
	"encoding/json"
	"testing"

	"github.com/ConsenSys/fc-retrieval/common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval/common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval/common/pkg/nodeid"
	"github.com/stretchr/testify/assert"
)

// // TestEncodeProviderPublishGroupOfferResponse success test
func TestEncodeProviderPublishGroupOfferResponse(t *testing.T) {
	mockProviderID, _ := nodeid.NewNodeIDFromHexString("42")
	contentID, _ := cid.NewContentIDFromBytes([]byte{1})
	mockCids := []cid.ContentID{*contentID}
	var mockPrice uint64 = 41
	var mockExpiry int64 = 42
	var mockQos uint64 = 43
	mockOffer, err := cidoffer.NewCIDOffer(mockProviderID, mockCids, mockPrice, mockExpiry, mockQos)
	mockMsgDigest := cidoffer.EncodeMessageDigest(mockOffer.GetMessageDigest())
	validMsg := &FCRMessage{
		messageType:       301,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       []byte(`{"gateway_id":"0000000000000000000000000000000000000000000000000000000000000042","digest":"syc/N3PHhXukj991WOragY0UbgulbK8EUw8AeMNK2Ms="}`),
		signature:         "",
	}

	msg, err := EncodeProviderPublishGroupOfferResponse(*mockProviderID, mockMsgDigest)
	assert.Empty(t, err)
	assert.Equal(t, msg, validMsg)
}

// TestDecodeProviderPublishGroupOfferResponse success test
func TestDecodeProviderPublishGroupOfferResponse(t *testing.T) {
	mockProviderID, _ := nodeid.NewNodeIDFromHexString("42")
	contentID, _ := cid.NewContentIDFromBytes([]byte{1})
	mockCids := []cid.ContentID{*contentID}
	var mockPrice uint64 = 41
	var mockExpiry int64 = 42
	var mockQos uint64 = 43
	mockOffer, err := cidoffer.NewCIDOffer(mockProviderID, mockCids, mockPrice, mockExpiry, mockQos)
	mockMsgDigest := cidoffer.EncodeMessageDigest(mockOffer.GetMessageDigest())

	mockOfferResponse, _ := json.Marshal(providerPublishGroupOfferResponse{
		GatewaydID: mockProviderID.ToString(),
		Digest:     mockMsgDigest,
	})

	validMsg := &FCRMessage{
		messageType:       301,
		protocolVersion:   1,
		protocolSupported: []int32{1, 1},
		messageBody:       mockOfferResponse,
		signature:         "",
	}
	nodeID, digest, err := DecodeProviderPublishGroupOfferResponse(validMsg)
	assert.Empty(t, err)
	assert.Equal(t, nodeID, mockProviderID)
	assert.Equal(t, digest, mockMsgDigest)
}
