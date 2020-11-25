package cidoffer
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
)



// CidGroupOffer represents a CID Group Offer. That is, an offer to deliver content 
// for Piece CIDs at a certain price
type CidGroupOffer struct {
    cids []big.Int
    price uint64
    expiry uint64
}


// NewCidGroupOffer creates an unsigned CID Group Offer.
func NewCidGroupOffer(cids *[]big.Int, price, expiry uint64) (*CidGroupOffer) {
	var c = CidGroupOffer{}
    c.cids = *cids
    c.price = price
    c.expiry = expiry
	return &c
}

// GetPrice returns the price that the content for the CIDs will be supplied at.
func (c *CidGroupOffer) GetPrice() (uint64) {
    return c.price
}
