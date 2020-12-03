package offers

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
	"math/big"
	"testing"
	"time"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
	"github.com/stretchr/testify/assert"
)


func TestOffersInitial(t *testing.T) {
	GetSingleInstance()
}

func TestEmpty(t *testing.T) {
	o := GetSingleInstance()

	_, exists := o.GetOffers(cidOne())
	assert.False(t, exists, "Found CID offer despite empty offer system")
}


func TestAddOld(t *testing.T) {
	o := GetSingleInstance()
	err := o.Add(createOldSingleCidGroupOfferCidOne(t))
	if err == nil {
        t.Errorf("Didn't return error due to CID having expired.")
	}

}

func TestAdd(t *testing.T) {
	o := newInstance()
	err := o.Add(createNewSingleCidGroupOfferCidOne(t))
	if err != nil {
		t.Errorf("Error returned by Add: %e", err)
	}

	cidOffers, exists := o.GetOffers(cidOne())
	assert.True(t, exists, "Can't find any offers for CID 1")
	assert.Equal(t, 1, len(cidOffers), "Should be only one CID offer")

	o.ExpireOffers()
	cidOffers, exists = o.GetOffers(cidOne())
	assert.True(t, exists, "Can't find any offers for CID 1")
	assert.Equal(t, 1, len(cidOffers), "Should be only one CID offer")
}


func TestAddTwo(t *testing.T) {
	o := newInstance()
	o.Add(createNewSingleCidGroupOfferCidOne(t))
	o.Add(createNewSingleCidGroupOfferCidOne(t))

	cidOffers, exists := o.GetOffers(cidOne())
	assert.True(t, exists, "Can't find any offers for CID 1")
	assert.Equal(t, 2, len(cidOffers), "Should be only two CID offers")
}


func createOldSingleCidGroupOfferCidOne(t *testing.T) (*cidoffer.CidGroupOffer) {
	return createSingleCidGroupOffer(t, cidOne(), true)
}

func createNewSingleCidGroupOfferCidOne(t *testing.T) (*cidoffer.CidGroupOffer) {
	return createSingleCidGroupOffer(t, cidOne(), false)
}


func createSingleCidGroupOffer(t *testing.T, theCid *cid.ContentID, old bool) (*cidoffer.CidGroupOffer) {
    aNodeID := nodeid.NewNodeID(nodeid.CreateRandomIdentifier())
    cids := make([]cid.ContentID, 0)
    cids = append(cids, *theCid)
	price := uint64(5)
	now := time.Now()
	nowSeconds := now.Unix()
	var expiry int64
	if (old) {
		expiry = nowSeconds - 1
	} else {
		expiry = nowSeconds + 1000
	}
    c, err := cidoffer.NewCidGroupOffer(aNodeID, &cids, price, expiry)
    if err != nil {
        t.Errorf("Error returned by NewCidGroupOffer: %e", err)
	}
	return c
}

func cidOne() *cid.ContentID {
	return cid.NewContentID(big.NewInt(1))
}

func cidTwo() *cid.ContentID {
	return cid.NewContentID(big.NewInt(2))
}