package control

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
	"strings"

	"github.com/ConsenSys/fc-retrieval-client/internal/network"
	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
)

// GatewayStdCIDDiscovery sends a  and processes a response.
func (c *ClientManager) GatewayStdCIDDiscovery(nodeID *nodeid.NodeID, contentID *cid.ContentID, nonce int64, ttl int64) ([]cidoffer.CidGroupOffer, error) {
	// Try to get gateway's ap
	c.GatewaysInUseLock.RLock()
	defer c.GatewaysInUseLock.RUnlock()
	gatewayInfo, exist := c.GatewaysInUse[strings.ToLower(nodeID.ToString())]
	if !exist {
		return nil, errors.New("GatewayID not found")
	}

	// Construct request
	request, err := fcrmessages.EncodeClientStandardDiscoverRequest(contentID, nonce, ttl)
	if err != nil {
		logging.Error("Error encoding Client Standard Discover Request: %+v", err)
		return nil, err
	}

	// Send request and get response
	response, err := network.SendMessage(gatewayInfo.NetworkInfoClient, request)
	if err != nil {
		return nil, err
	}

	// Get the gateway's public key
	pubKey, err := gatewayInfo.GetSigningKey()
	if err != nil {
		return nil, err
	}

	// Verify the response
	ok, err := response.VerifySignature(func(sig string, msg interface{}) (bool, error) {
		return fcrcrypto.VerifyMessage(pubKey, sig, msg)
	})
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("Verification failed")
	}

	// Decode the response, TODO, fundedpayment channels? and found?, how to deal with them
	cid, nonceRecv, _, offers, merkleRoots, merkleProofs, _, err := fcrmessages.DecodeClientStandardDiscoverResponse(response)
	if err != nil {
		return nil, err
	}
	if cid.ToString() != contentID.ToString() {
		return nil, errors.New("CID Mismatch")
	}
	if nonce != nonceRecv {
		return nil, errors.New("Nonce mismatch")
	}
	if len(offers) != len(merkleRoots) || len(offers) != len(merkleProofs) {
		return nil, errors.New("offer length and proof length not match")
	}

	result := make([]cidoffer.CidGroupOffer, 0)
	// Verify each offer
	for i, offer := range offers {
		// Get the provider's pubkey
		c.ProvidersLock.RLock()
		providerInfo, exist := c.Providers[strings.ToLower(offer.NodeID.ToString())]
		if !exist {
			logging.Error("Provider who created this offer does not exist in local storage.")
			continue
		}
		c.ProvidersLock.RUnlock()
		pubKey, _ := providerInfo.GetSigningKey()
		// Verify the offer
		ok, err := offer.VerifySignature(func(sig string, msg interface{}) (bool, error) {
			return fcrcrypto.VerifyMessage(pubKey, sig, msg)
		})
		if err != nil {
			logging.Error("Error in verifying the offer.")
			continue
		}
		if !ok {
			logging.Error("Offer signature fail to verify.")
			continue
		}
		// Now Verify the merkle root
		if offer.GetMerkleRoot() != merkleRoots[i] {
			logging.Error("Merkle root does not match.")
			continue
		}
		// Now verify the merkle proof
		ok = merkleProofs[i].VerifyContent(*cid, merkleRoots[i]) //TODO, Need to check this, if using pointer?
		if !ok {
			logging.Error("Merkle Proof verifcation failed.")
			continue
		}
		// Offer pass verification
		logging.Info("Offer pass every verification, added to result")
		result = append(result, offer)
	}

	logging.Info("Total received offer: %d, total verified offer: %d", len(offers), len(result))
	return result, nil
}
