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

	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
)

// FCROfferMgr manages offer storage
type FCROfferMgr struct {
	DHTOffers   *offerStorage
	GroupOffers *offerStorage
}

// NewFCROfferMgr returns
func NewFCROfferMgr() *FCROfferMgr {
	return &FCROfferMgr{
		DHTOffers:   newOfferStorage(),
		GroupOffers: newOfferStorage(),
	}
}

// AddGroupOffer stores a group offer
func (mgr *FCROfferMgr) AddGroupOffer(offer *cidoffer.CIDOffer) error {
	if len(offer.GetCIDs()) <= 1 {
		return errors.New("Not a group offer")
	}
	return mgr.GroupOffers.add(offer)
}

// AddDHTOffer stores a dht offer
func (mgr *FCROfferMgr) AddDHTOffer(offer *cidoffer.CIDOffer) error {
	if len(offer.GetCIDs()) != 1 {
		return errors.New("Not a DHT offer")
	}
	return mgr.DHTOffers.add(offer)
}

// GetGroupOffers returns a list of group offers that contain the given cid
func (mgr *FCROfferMgr) GetGroupOffers(cid *cid.ContentID) ([]cidoffer.CIDOffer, bool) {
	res := mgr.GroupOffers.get(cid)
	return res, len(res) > 0
}

// GetDHTOffers returns a list of dht offers that contain the given cid
func (mgr *FCROfferMgr) GetDHTOffers(cid *cid.ContentID) ([]cidoffer.CIDOffer, bool) {
	res := mgr.DHTOffers.get(cid)
	return res, len(res) > 0
}

// GetOffers returns a list of all offers (group or dht) that contain the given cid
func (mgr *FCROfferMgr) GetOffers(cid *cid.ContentID) ([]cidoffer.CIDOffer, bool) {
	res := append(mgr.GroupOffers.get(cid), mgr.DHTOffers.get(cid)...)
	return res, len(res) > 0
}
