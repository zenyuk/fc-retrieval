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
	"crypto/sha512"
	"encoding/binary"
	"errors"
	"time"

	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmerkletree"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/cbergoon/merkletree"
)

// CidGroupOfferDigestSize is the size of message digest used for CidGroupOffers
const CidGroupOfferDigestSize = sha512.Size256

// CidGroupOffer represents a CID Group Offer. That is, an offer to deliver content
// for Piece CIDs at a certain price
type CidGroupOffer struct {
	NodeID     *nodeid.NodeID
	Cids       []cid.ContentID
	Price      uint64
	Expiry     int64
	QoS        uint64
	MerkleRoot string
	MerkleTrie *fcrmerkletree.FCRMerkleTrie
	Signature  string
}

// NewCidGroupOffer creates an unsigned CID Group Offer.
func NewCidGroupOffer(providerID *nodeid.NodeID, cids *[]cid.ContentID, price uint64, expiry int64, qos uint64) (*CidGroupOffer, error) {
	var c = CidGroupOffer{}
	c.NodeID = providerID
	if len(*cids) < 1 {
		return nil, errors.New("CID Group Offer: provide 1 or more CIDs")
	}
	c.Cids = *cids
	c.Price = price
	// TODO check that the expiry is in the future (are there scenarios where an expired offer should be loadable?)
	c.Expiry = expiry
	c.QoS = qos

	// Create merkle tree
	list := make([]merkletree.Content, len(*cids))
	for i := 0; i < len(*cids); i++ {
		list[i] = (*cids)[i]
	}
	var err error
	c.MerkleTrie, err = fcrmerkletree.CreateMerkleTrie(list)
	if err != nil {
		return nil, err
	}
	c.MerkleRoot = c.MerkleTrie.GetMerkleRoot()

	return &c, nil
}

// GetCIDs returns the CIDs this offer relates to.
func (c *CidGroupOffer) GetCIDs() *[]cid.ContentID {
	return &c.Cids
}

// GetPrice returns the price that the content for the CIDs will be supplied at.
func (c *CidGroupOffer) GetPrice() uint64 {
	return c.Price
}

// GetExpiry returns the expiry time of the offer
func (c *CidGroupOffer) GetExpiry() int64 {
	return c.Expiry
}

// GetQoS returns the qos of the offer
func (c *CidGroupOffer) GetQoS() uint64 {
	return c.QoS
}

// GetMerkleTrie returns the merkle tree of the cids
func (c *CidGroupOffer) GetMerkleTrie() *fcrmerkletree.FCRMerkleTrie {
	return c.MerkleTrie
}

// GetMerkleRoot returns the merkle root of the cids
func (c *CidGroupOffer) GetMerkleRoot() string {
	return c.MerkleRoot
}

// GetMessageDigest calculate the message digest of this CID Group Offer.
// Note that the methodology used here should not be externally visible. The
// message digest should only be used within the gateway.
func (c *CidGroupOffer) GetMessageDigest() (sum256 [CidGroupOfferDigestSize]byte) {
	b := c.NodeID.ToBytes()

	for _, aCid := range c.Cids {
		cidBytes := aCid.ToBytes()
		b = append(b[:], cidBytes[:]...)
	}

	bPrice := make([]byte, 8)
	binary.BigEndian.PutUint64(bPrice, uint64(c.Price))
	b = append(b[:], bPrice[:]...)

	bExpiry := make([]byte, 8)
	binary.BigEndian.PutUint64(bExpiry, uint64(c.Expiry))
	b = append(b[:], bExpiry[:]...)

	return sha512.Sum512_256(b)
}

// HasExpired returns true if the offer expiry date is in the past.
func (c *CidGroupOffer) HasExpired() bool {
	expiryTime := time.Unix(c.Expiry, 0)
	now := time.Now()
	return expiryTime.Before(now)
}

// VerifySignature is used to verify the signature
func (c *CidGroupOffer) VerifySignature(verify func(sig string, msg interface{}) (bool, error)) (bool, error) {
	// Clear signature, tree and cid array
	sig := c.Signature
	c.Signature = ""
	tree := c.MerkleTrie
	c.MerkleTrie = nil
	cids := c.Cids
	c.Cids = nil

	// Verify the offer
	res, err := verify(sig, c)
	if err != nil {
		return false, err
	}
	// Recover signature, tree and cid array
	c.Signature = sig
	c.MerkleTrie = tree
	c.Cids = cids
	return res, nil
}

// SignOffer is used to sign the offer
func (c *CidGroupOffer) SignOffer(sign func(msg interface{}) (string, error)) error {
	// Clear signature, tree and cid array
	c.Signature = ""
	tree := c.MerkleTrie
	c.MerkleTrie = nil
	cids := c.Cids
	c.Cids = nil

	// Sign the offer
	sig, err := sign(c)
	if err != nil {
		return err
	}
	c.Signature = sig

	// Recover the tree and cid array
	c.MerkleTrie = tree
	c.Cids = cids
	return nil
}
