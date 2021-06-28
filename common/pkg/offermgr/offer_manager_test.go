// FCROfferMgr manages offer storage
package offermgr

import (
	//	"errors"
	// "fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ConsenSys/fc-retrieval/common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval/common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval/common/pkg/nodeid"
	"github.com/stretchr/testify/assert"
)

func init() {
	inUnitTest = true
}

func TestGetGroupOffers02(t *testing.T) {
	mgr := NewFCROfferMgr()

	offerGroup, err := getOfferGroup()
	assert.Equal(t, nil, err)

	err = mgr.AddGroupOffer(offerGroup)
	assert.Equal(t, nil, err)

	offers, _ := mgr.GetOffers(intToCid(6))
	assert.Equal(t, 0, len(offers))

	offers, _ = mgr.GetOffers(intToCid(7))
	assert.Equal(t, 1, len(offers))

	offerSingle, err := getOfferSingle(7)
	assert.Equal(t, nil, err)

	offers, _ = mgr.GetDHTOffersWithinRange(intToCid(6), intToCid(8), 3)
	assert.Equal(t, 1, len(offers))

	err = mgr.AddDHTOffer(offerSingle)
	assert.Equal(t, nil, err)

	offers, _ = mgr.GetDHTOffersWithinRange(intToCid(6), intToCid(8), 3)
	assert.Equal(t, 2, len(offers))

	offers, _ = mgr.GetOffers(intToCid(7))
	assert.Equal(t, 2, len(offers))

	err = mgr.AddDHTOffer(offerSingle)
	assert.Equal(t, nil, err)

	offers, _ = mgr.GetDHTOffersWithinRange(intToCid(6), intToCid(8), 3)
	assert.Equal(t, 2, len(offers))

	offerSingle, err = getOfferSingle(8)
	assert.Equal(t, nil, err)

	err = mgr.AddDHTOffer(offerSingle)
	offers, _ = mgr.GetDHTOffersWithinRange(intToCid(6), intToCid(9), 3)
	assert.Equal(t, 3, len(offers))

	offers, _ = mgr.GetDHTOffersWithinRange(intToCid(6), intToCid(9), 2)
	assert.Equal(t, 2, len(offers))

	_, find := mgr.GetOfferByDigest(offers[0].GetMessageDigest())
	assert.Equal(t, true, find)

	_, find = mgr.GetOfferByDigest([cidoffer.CIDOfferDigestSize]byte{})
	assert.Equal(t, false, find)
}

func TestGetDTHOffers01(t *testing.T) {
	offerSingle, err := getOfferSingle(7)
	assert.Equal(t, err, nil)
	mgr := NewFCROfferMgr()

	err = mgr.AddDHTOffer(offerSingle)
	assert.Equal(t, err, nil)

	err = mgr.AddDHTOffer(offerSingle)
	assert.Equal(t, err, nil)

	_, find := mgr.GetDHTOffers(intToCid(7))
	assert.Equal(t, true, find)

	_, find = mgr.GetDHTOffers(intToCid(8))
	assert.Equal(t, false, find)
}

func TestGetDTHOffers02(t *testing.T) {
	offerSingle, err := getOfferSingleExpired()
	assert.Equal(t, err, nil)
	mgr := NewFCROfferMgr()

	err = mgr.AddDHTOffer(offerSingle)
	assert.Equal(t, err, nil)

	_, find := mgr.GetDHTOffers(intToCid(7))
	assert.Equal(t, false, find)
}


// Helper functions

func intToCid(n int64) *cid.ContentID {
	aCid, _ := cid.NewContentID(big.NewInt(n))
	return aCid
}

func getOfferSingle(n int64) (*cidoffer.CIDOffer, error) {
	aNodeID, err := nodeid.NewNodeID(big.NewInt(n))
	aCid, err := cid.NewContentID(big.NewInt(n))
	cids := []cid.ContentID{*aCid}
	price := uint64(5)
	expiry := time.Now().Add(12 * time.Hour).Unix()
	qos := uint64(5)
	offer, err := cidoffer.NewCIDOffer(aNodeID, cids, price, expiry, qos)
	return offer, err
}

func getOfferSingleExpired() (*cidoffer.CIDOffer, error) {
	aNodeID, err := nodeid.NewNodeID(big.NewInt(7))
	aCid := intToCid(7)
	cids := []cid.ContentID{*aCid}
	price := uint64(5)
	expiry := int64(10)
	qos := uint64(5)
	offer, err := cidoffer.NewCIDOffer(aNodeID, cids, price, expiry, qos)
	return offer, err
}

func getOfferGroupExpired() (*cidoffer.CIDOffer, error) {
	aNodeID, err := nodeid.NewNodeID(big.NewInt(7))
	aCid1 := intToCid(7)
	aCid2 := intToCid(8)
	aCid3 := intToCid(9)
	cids := []cid.ContentID{*aCid1, *aCid2, *aCid3}
	price := uint64(5)
	expiry := int64(2)
	qos := uint64(5)
	offer, err := cidoffer.NewCIDOffer(aNodeID, cids, price, expiry, qos)
	return offer, err
}

func getOfferGroup() (*cidoffer.CIDOffer, error) {
	aNodeID, err := nodeid.NewNodeID(big.NewInt(7))
	aCid1 := intToCid(7)
	aCid2 := intToCid(8)
	aCid3 := intToCid(9)
	cids := []cid.ContentID{*aCid1, *aCid2, *aCid3}
	price := uint64(5)
	expiry := time.Now().Add(12 * time.Hour).Unix()
	qos := uint64(5)
	offer, err := cidoffer.NewCIDOffer(aNodeID, cids, price, expiry, qos)
	return offer, err
}
