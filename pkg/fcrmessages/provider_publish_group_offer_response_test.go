package fcrmessages

import (
	"testing"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
)

// TestEncodeProviderPublishGroupOfferRequest success test
func TestEncodeProviderPublishGroupOfferRequest(t *testing.T) {
	mockProviderID, _ := nodeid.NewNodeIDFromHexString("42")
	mockNodeID, _ := nodeid.NewNodeIDFromHexString("42")
	mockNonce := int64(42)
	contentID, _ := cid.NewContentIDFromBytes([]byte{1})
	mockCids := []cid.ContentID{*contentID}
	var mockPrice uint64 = 41
	var mockExpiry int64 = 42
	var mockQos uint64 = 43
	mockOffer, err := cidoffer.NewCIDOffer(mockProviderID, mockCids, mockPrice, mockExpiry, mockQos)
		
	validMsg := &FCRMessage{
		messageType:300,
		protocolVersion:1,
		protocolSupported:[]int32{1, 1},
		messageBody:[]byte(`{"provider_id":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEI=","nonce":42,"offer":{"provider_id":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEI=","cids":["AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE="],"price":41,"expiry":42,"qos":43,"signature":""}}`), 
		signature:"",
	}

	msg, err := EncodeProviderPublishGroupOfferRequest(mockNodeID, mockNonce, mockOffer)
	assert.Empty(t, err)
	assert.Equal(t, msg, validMsg)
}

// TestDecodeProviderPublishGroupOfferResponse success test
func TestDecodeProviderPublishGroupOfferResponse(t *testing.T) {
	mockProviderID, _ := nodeid.NewNodeIDFromHexString("42")
	mockNodeID, _ := nodeid.NewNodeIDFromHexString("42")
	contentID, _ := cid.NewContentIDFromBytes([]byte{1})
	mockCids := []cid.ContentID{*contentID}
	var mockPrice uint64 = 41
	var mockExpiry int64 = 42
	var mockQos uint64 = 43
	mockOffer, err := cidoffer.NewCIDOffer(mockProviderID, mockCids, mockPrice, mockExpiry, mockQos)
	mockMsgDigest := mockOffer.GetMessageDigest()

	mockOfferREsponse, _ := json.Marshal(providerPublishGroupOfferResponse{
		GatewaydID: *mockNodeID,
		Digest:     mockMsgDigest,
	})

	validMsg := &FCRMessage{
		messageType:301,
		protocolVersion:1,
		protocolSupported:[]int32{1, 1},
		messageBody:mockOfferREsponse, 
		signature:"",
	}
	nodeID, digest, err := DecodeProviderPublishGroupOfferResponse(validMsg)
	assert.Empty(t, err)
	assert.Equal(t, nodeID, mockNodeID)
	assert.Equal(t, digest, mockMsgDigest)
}