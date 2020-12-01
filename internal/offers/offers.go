package offers
// Copyright (C) 2020 ConsenSys Software Inc

import (
    "errors"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cidoffer"
	"sync"
	"container/list"
)


// Single instance of the CID Offer system
var instance = newInstance()




// Offers manages all CID Offers on this gateway.
type Offers struct {
	// Map of CIDs to array of CID Group Offer message digests
	cidMap         map[string][][cidoffer.CidGroupOfferDigestSize]byte
	cidMapLock   sync.RWMutex
	// Map of CID Group Offer message digests to CID Group Offer
	cidOffers          map[[cidoffer.CidGroupOfferDigestSize]byte]*cidoffer.CidGroupOffer
	cidOffersLock   sync.RWMutex
	// Linked list of CID Group Offer message digests, in order of expiry time.
	offerExpiry		*list.List
	offerExpiryLock  sync.RWMutex
}

type expiringOffers struct {
	expiry int64
	offerDigest *[cidoffer.CidGroupOfferDigestSize]byte
}

// Create a new instance
func newInstance() *Offers {
	var o = Offers{}
	o.cidMap = make(map[string][][cidoffer.CidGroupOfferDigestSize]byte)
	o.cidMapLock = sync.RWMutex{}
	o.cidOffers = make(map[[cidoffer.CidGroupOfferDigestSize]byte]*cidoffer.CidGroupOffer)
	o.cidOffersLock = sync.RWMutex{}
	o.offerExpiry = list.New()
	o.offerExpiryLock = sync.RWMutex{}
	return &o
}


// GetSingleInstance is a factory method to get the single instance of the offers
func GetSingleInstance() *Offers {
	return instance
}

// Add adds a new CID Offer to the offers system
// Ignore the new offer if it already exists.
func (o *Offers) Add(newOffer *cidoffer.CidGroupOffer) error {
	digest := newOffer.GetMessageDigest()
	o.cidOffersLock.Lock()
	_, exists := o.cidOffers[digest]
	o.cidOffersLock.Unlock()
	if (exists) {
		// Ignore if offer already in the system.
		return nil
	}

	if newOffer.HasExpired() {
		return errors.New("Offers: Add: Attempt to add expired offer")
	}

	// First: Add to the map of Offer digests to Offers
	o.cidOffersLock.Lock()
	o.cidOffers[digest] = newOffer
	o.cidOffersLock.Unlock()

	// Second: Add the CIDs to the map of CIDs to Offer digests
	for _, aCid := range *newOffer.GetCIDs() {
		cidStr := aCid.ToString()
		// TODO: Would it be more efficient to do the locking outside of the loop?
		// TODO: however, could this lock the map for too long?
		o.cidMapLock.Lock()
		o.cidMap[cidStr] = append(o.cidMap[cidStr][:], digest)
		o.cidMapLock.Unlock()
	}

	// Now add the offer to the linked list of expiring offers
	expiry := newOffer.GetExpiry()
	newExpOff := expiringOffers{expiry: expiry, offerDigest: &digest}
	added := false
	o.offerExpiryLock.Lock()
	for e := o.offerExpiry.Front(); e != nil; e = e.Next() {
		var anOffer *expiringOffers
		anOffer = e.Value.(*expiringOffers)
		if expiry < anOffer.expiry {
			o.offerExpiry.InsertBefore(newExpOff, e)
			added = true
		}
	}
	if (!added) {
		o.offerExpiry.PushBack(newExpOff)
	}
	o.offerExpiryLock.Unlock()
	return nil
}

