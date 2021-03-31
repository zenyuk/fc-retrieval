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

	"github.com/ConsenSys/fc-retrieval-client/internal/network"
	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-register/pkg/register"
)

// GatewayStdCIDDiscovery sends a  and processes a response.
func (c *ClientManager) GatewayStdCIDDiscovery(gatewayInfo *register.GatewayRegister, contentID *cid.ContentID, nonce int64, ttl int64) ([]cidoffer.SubCIDOffer, error) {
	// Construct request
	request, err := fcrmessages.EncodeClientStandardDiscoverRequest(contentID, nonce, ttl, "", "")
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
	if response.Verify(pubKey) != nil {
		return nil, errors.New("Verification failed")
	}

	// Decode the response, TODO deal with fundedpayment channels and found
	cid, nonceRecv, _, offers, _, err := fcrmessages.DecodeClientStandardDiscoverResponse(response)
	if err != nil {
		return nil, err
	}
	if cid.ToString() != contentID.ToString() {
		return nil, errors.New("CID Mismatch")
	}
	if nonce != nonceRecv {
		return nil, errors.New("Nonce mismatch")
	}

	result := make([]cidoffer.SubCIDOffer, 0)
	// Verify each offer
	for _, offer := range offers {
		// Get the provider's pubkey
		providerInfo, err := register.GetProviderByID(c.Settings.RegisterURL(), offer.GetProviderID())
		if err != nil {
			logging.Error("Provider who created this offer does not exist in local storage.")
			continue
		}
		if !validateProviderInfo(&providerInfo) {
			logging.Error("Provider register info not valid.")
			continue
		}
		pubKey, _ := providerInfo.GetSigningKey()
		// Verify the offer
		if offer.Verify(pubKey) != nil {
			logging.Error("Offer signature fail to verify.")
			continue
		}
		// Now Verify the merkle proof
		if offer.VerifyMerkleProof() != nil {
			logging.Error("Merkle proof verification failed.")
			continue
		}
		// Offer pass verification
		logging.Info("Offer pass every verification, added to result")
		result = append(result, offer)
	}

	logging.Info("Total received offer: %d, total verified offer: %d", len(offers), len(result))
	return result, nil
}
