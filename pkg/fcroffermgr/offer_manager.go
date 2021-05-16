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
	"strconv"

	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
)

// FCROfferMgr manages offer storage
type FCROfferMgr struct {
	dhtOffers   *offerStorage
	groupOffers *offerStorage
}

// NewFCROfferMgr returns
func NewFCROfferMgr() *FCROfferMgr {
	return &FCROfferMgr{
		dhtOffers:   newOfferStorage(),
		groupOffers: newOfferStorage(),
	}
}

// AddGroupOffer stores a group offer
func (mgr *FCROfferMgr) AddGroupOffer(offer *cidoffer.CIDOffer) error {
	if len(offer.GetCIDs()) <= 1 {
		return errors.New("Not a group offer")
	}
	return mgr.groupOffers.add(offer)
}

// AddDHTOffer stores a dht offer
func (mgr *FCROfferMgr) AddDHTOffer(offer *cidoffer.CIDOffer) error {
	if len(offer.GetCIDs()) != 1 {
		return errors.New("Not a DHT offer")
	}
	return mgr.dhtOffers.add(offer)
}

// GetGroupOffers returns a list of group offers that contain the given cid
func (mgr *FCROfferMgr) GetGroupOffers(cid *cid.ContentID) ([]cidoffer.CIDOffer, bool) {
	res := mgr.groupOffers.get(cid)
	return res, len(res) > 0
}

// GetDHTOffers returns a list of dht offers that contain the given cid
func (mgr *FCROfferMgr) GetDHTOffers(cid *cid.ContentID) ([]cidoffer.CIDOffer, bool) {
	res := mgr.dhtOffers.get(cid)
	return res, len(res) > 0
}

// GetDHTOffersWithinRange returns a list of dht offers contains a cid within the given range
func (mgr *FCROfferMgr) GetDHTOffersWithinRange(cidMin, cidMax *cid.ContentID, maxOffers int) ([]cidoffer.CIDOffer, bool) {
	// TODO: Have a more efficient implementation, using Ring, but with the ability to remove expired entry
	offers := make([]cidoffer.CIDOffer, 0)

	min, err := strconv.ParseInt(cidMin.ToString(), 16, 32) // TODO, CHECK IF THIS IS CORRECT
	if err != nil {
		return offers, false
	}
	max, err := strconv.ParseInt(cidMax.ToString(), 16, 32) // TODO, CHECK IF THIS IS CORRECT
	if err != nil {
		return offers, false
	}
	if max < min {
		// TODO, Test boundary cases
		cidNewMax, err := cid.NewContentIDFromHexString("0xFFFFFFFF")
		if err != nil {
			logging.Error("Error in getting maximum cid")
			return offers, false
		}
		tempOffers, _ := mgr.GetDHTOffersWithinRange(cidMax, cidNewMax, maxOffers)
		max = min
		min = 0
		offers = append(offers, tempOffers...)
		maxOffers = maxOffers - len(tempOffers)
	}

	for i := min; i <= max; i++ {
		id, err := cid.NewContentID(big.NewInt(i))
		if err != nil {
			return offers, false
		}
		offers, exists := mgr.GetDHTOffers(id)
		if exists {
			for _, offer := range offers {
				offers = append(offers, offer)
				if len(offers) >= maxOffers {
					break
				}
			}
		}
		if len(offers) >= maxOffers {
			break
		}
	}

	return offers, len(offers) > 0
}

// GetOffers returns a list of all offers (group or dht) that contain the given cid
func (mgr *FCROfferMgr) GetOffers(cid *cid.ContentID) ([]cidoffer.CIDOffer, bool) {
	res := append(mgr.groupOffers.get(cid), mgr.dhtOffers.get(cid)...)
	return res, len(res) > 0
}
