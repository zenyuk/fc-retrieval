package fcroffermgr

/*
 * Copyright 2020 ConsenSys Software Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

import (
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ConsenSys/fc-retrieval/common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval/common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval/common/pkg/nodeid"
	"github.com/stretchr/testify/assert"
)

func TestAddGroupOffer01(t *testing.T) {
	mgr := NewFCROfferMgr()
	offerSingle, err := getOfferSingleExpired()
	assert.Equal(t, nil, err)
	err = mgr.AddGroupOffer(offerSingle)
	assert.Equal(t, errors.New("not a group offer"), err)
}

func TestAddGroupOffer02(t *testing.T) {
	offerGroupExpired, err := getOfferGroupExpired()
	assert.Equal(t, nil, err)
	mgr := NewFCROfferMgr()
	err = mgr.AddGroupOffer(offerGroupExpired)
	assert.Equal(t, errors.New("offers: Attempt to add an expired offer"), err)
}

func TestAddGroupOffer03(t *testing.T) {
	offerGroup, err := getOfferGroup()
	assert.Equal(t, nil, err)
	mgr := NewFCROfferMgr()

	err = mgr.AddGroupOffer(offerGroup)
	assert.Equal(t, nil, err)

	err = mgr.AddGroupOffer(offerGroup)
	assert.Equal(t, nil, err)
}

func TestAddDHTOffer01(t *testing.T) {
	mgr := NewFCROfferMgr()
	offerSingle, err := getOfferSingleExpired()
	assert.Equal(t, nil, err)
	err = mgr.AddDHTOffer(offerSingle)
	assert.Equal(t, errors.New("offers: Attempt to add an expired offer"), err)
}

func TestAddDHTOffer02(t *testing.T) {
	offerGroupExpired, err := getOfferGroupExpired()
	assert.Equal(t, nil, err)
	mgr := NewFCROfferMgr()
	err = mgr.AddDHTOffer(offerGroupExpired)
	assert.Equal(t, errors.New("not a DHT offer"), err)
}

func TestAddDHTOffer03(t *testing.T) {
	mgr := NewFCROfferMgr()

	offerSingle, err := getOfferSingle(7)
	assert.Equal(t, nil, err)

	err = mgr.AddDHTOffer(offerSingle)
	assert.Equal(t, nil, err)

	offerSingle, err = getOfferSingle(8)
	assert.Equal(t, nil, err)

	err = mgr.AddDHTOffer(offerSingle)
	assert.Equal(t, nil, err)
}

func TestGetGroupOffers01(t *testing.T) {
	offerGroup, err := getOfferGroup()
	assert.Equal(t, nil, err)
	mgr := NewFCROfferMgr()

	err = mgr.AddGroupOffer(offerGroup)
	assert.Equal(t, nil, err)

	err = mgr.AddGroupOffer(offerGroup)
	assert.Equal(t, nil, err)

	_, find := mgr.GetGroupOffers(intToCid(7))
	assert.Equal(t, true, find)

	_, find = mgr.GetGroupOffers(intToCid(8))
	assert.Equal(t, true, find)

	_, find = mgr.GetGroupOffers(intToCid(6))
	assert.Equal(t, false, find)
}

func TestGetGroupOffers02(t *testing.T) {
	mgr := NewFCROfferMgr()

	offerGroup, err := getOfferGroup()
	assert.Equal(t, nil, err)

	err = mgr.AddGroupOffer(offerGroup)
	assert.Equal(t, nil, err)

	offers, _ := mgr.GetDHTOffersWithinRange(intToCid(6), intToCid(8), 3)
	assert.Equal(t, 0, len(offers))

	offers, _ = mgr.GetOffers(intToCid(6))
	assert.Equal(t, 0, len(offers))

	offers, _ = mgr.GetOffers(intToCid(7))
	assert.Equal(t, 1, len(offers))

	offerSingle, err := getOfferSingle(7)
	assert.Equal(t, nil, err)

	err = mgr.AddDHTOffer(offerSingle)
	assert.Equal(t, nil, err)

	offers, _ = mgr.GetDHTOffersWithinRange(intToCid(6), intToCid(8), 3)
	assert.Equal(t, 1, len(offers))

	offers, _ = mgr.GetOffers(intToCid(7))
	assert.Equal(t, 2, len(offers))

	err = mgr.AddDHTOffer(offerSingle)
	assert.Equal(t, nil, err)

	offers, _ = mgr.GetDHTOffersWithinRange(intToCid(6), intToCid(8), 3)
	assert.Equal(t, 1, len(offers))

	offerSingle, err = getOfferSingle(8)
	assert.Equal(t, nil, err)

	err = mgr.AddDHTOffer(offerSingle)
	offers, _ = mgr.GetDHTOffersWithinRange(intToCid(6), intToCid(9), 3)
	assert.Equal(t, 2, len(offers))

	// offers, _ = mgr.GetDHTOffersWithinRange(intToCid(6), intToCid(9), 3)
	// assert.Equal(t, 1, len(offers))

	// _, find := mgr.GetOfferByDigest(offers[0].GetMessageDigest())
	// assert.Equal(t, true, find)

	// _, find = mgr.GetOfferByDigest([cidoffer.CIDOfferDigestSize]byte{})
	// assert.Equal(t, false, find)
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
