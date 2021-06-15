package clientapi

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
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/register"
)

// RequestStandardDiscover requests a standard discover to a given gateway for a given contentID, nonce and ttl.
func (c *Client) RequestStandardDiscover(
	gatewayInfo *register.GatewayRegister,
	contentID *cid.ContentID,
	nonce int64,
	ttl int64,
	paychAddr string,
	voucher string,
) ([]cidoffer.SubCIDOffer, error) {
	// Construct request
	request, err := fcrmessages.EncodeClientStandardDiscoverRequest(contentID, nonce, ttl, paychAddr, voucher)
	if err != nil {
		logging.Error("Error encoding Client Standard Discover Request: %+v", err)
		return nil, err
	}

	// Send request and get response
	response, err := c.httpCommunicator.SendMessage(gatewayInfo.NetworkInfoClient, request)
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
		return nil, errors.New("verification failed")
	}

	// Decode the response, TODO deal with fundedpayment channels and found
	pieceCid, nonceRecv, _, offers, _, err := fcrmessages.DecodeClientStandardDiscoverResponse(response)
	if err != nil {
		return nil, err
	}
	if pieceCid.ToString() != contentID.ToString() {
		return nil, errors.New("CID Mismatch")
	}
	if nonce != nonceRecv {
		return nil, errors.New("nonce mismatch")
	}

	return offers, nil
}
