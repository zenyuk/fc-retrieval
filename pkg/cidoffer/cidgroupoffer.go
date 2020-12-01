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
    "errors"
    "time"
    "github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
    "github.com/ConsenSys/fc-retrieval-gateway/pkg/cid"
    "crypto/sha512"
    "encoding/binary"
)

// CidGroupOfferDigestSize is the size of message digest used for CidGroupOffers
const CidGroupOfferDigestSize = sha512.Size256

// CidGroupOffer represents a CID Group Offer. That is, an offer to deliver content 
// for Piece CIDs at a certain price
type CidGroupOffer struct {
    nodeID *nodeid.NodeID
    cids []cid.ContentID
    price uint64
    expiry int64
    // TODO signature
}


// NewCidGroupOffer creates an unsigned CID Group Offer.
func NewCidGroupOffer(providerID *nodeid.NodeID, cids *[]cid.ContentID, price uint64, expiry int64) (*CidGroupOffer, error) {
    var c = CidGroupOffer{}
    c.nodeID = providerID
    if (len(*cids) < 1) {
        return nil, errors.New("CID Group Offer: provide 1 or more CIDs")
    }
    c.cids = *cids
    if (price < 0) {
        return nil, errors.New("CID Group Offer: price must be greater than or equal to zero")
    }
    c.price = price
    // TODO check that the expiry is in the future (are there scenarios where an expired offer should be loadable?)
    c.expiry = expiry
	return &c, nil
}

// GetCIDs returns the CIDs this offer relates to.
func (c *CidGroupOffer) GetCIDs() *[]cid.ContentID {
    return &c.cids
}

// GetPrice returns the price that the content for the CIDs will be supplied at.
func (c *CidGroupOffer) GetPrice() (uint64) {
    return c.price
}

// GetExpiry returns the expiry time of the offer
func (c *CidGroupOffer) GetExpiry() (int64) {
    return c.expiry
}


// GetMessageDigest calculate the message digest of this CID Group Offer.
// Note that the methodology used here should not be externally visible. The 
// message digest should only be used within the gateway.
func (c *CidGroupOffer) GetMessageDigest() (sum256 [CidGroupOfferDigestSize]byte) {
    b := c.nodeID.ToBytes()

    for _, aCid := range c.cids  {
        cidBytes := aCid.ToBytes()
        b = append(b[:], cidBytes[:]...)
    }

    bPrice := make([]byte, 8)
	binary.BigEndian.PutUint64(bPrice, uint64(c.price))
    b = append(b[:], bPrice[:]...)

    bExpiry := make([]byte, 8)
	binary.BigEndian.PutUint64(bExpiry, uint64(c.expiry))
    b = append(b[:], bExpiry[:]...)


    return sha512.Sum512_256(b)
}

// HasExpired returns true if the offer expiry date is in the past.
func (c *CidGroupOffer) HasExpired() bool {
    expiryTime := time.Unix(c.expiry, 0)
    now := time.Now()
    return expiryTime.Before(now)
}