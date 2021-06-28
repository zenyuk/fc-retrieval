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
	"fmt"

	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/register"
)

// RequestStandardDiscoverV2 requests a standard discover to a given gateway for a given contentID, nonce and ttl.
func (c *Client) RequestStandardDiscoverV2(
	gatewayRegistrar register.GatewayRegistrar,
	contentID *cid.ContentID,
	nonce int64,
	ttl int64,
	paychAddr string,
	voucher string,
) ([][cidoffer.CIDOfferDigestSize]byte, error) {

	// Construct request
	request, err := fcrmessages.EncodeClientStandardDiscoverRequestV2(contentID, nonce, ttl, paychAddr, voucher)
	if err != nil {
		return nil, fmt.Errorf("error encoding Client Standard Discover Request: %s", err.Error())
	}

	// Send request and get response
	response, err := c.httpCommunicator.SendMessage(gatewayRegistrar.GetNetworkInfoClient(), request)
	if err != nil {
		return nil, fmt.Errorf("error sending message to gateway ID: %s, error: %s", gatewayRegistrar.GetNodeID(), err.Error())
	}

	// Get the gateway's public key
	pubKey, err := gatewayRegistrar.GetSigningKey()
	if err != nil {
		return nil, fmt.Errorf("error getting signing key of gateway ID: %s, error: %s", gatewayRegistrar.GetNodeID(), err.Error())
	}

	// Verify the response
	if response.Verify(pubKey) != nil {
		return nil, fmt.Errorf("response verification failed for gateway ID: %s, message type ID: %d", gatewayRegistrar.GetNodeID(), request.GetMessageType())
	}

	// Decode the response, TODO deal with funded payment channels and found
	cID, nonceRecv, _, offerDigests, _, paymentRequired, paymentChannelAddrToTopup, err := fcrmessages.DecodeClientStandardDiscoverResponseV2(response)
	if err != nil {
		return nil, fmt.Errorf("error decoding Client Standard Discover Response: %s, gateway ID: %s", err.Error(), gatewayRegistrar.GetNodeID())
	}
	if cID.ToString() != contentID.ToString() {
		return nil, fmt.Errorf("error validating CID for Client Standard Discover Response for gateway ID: %s; expected CID: %s, actual CID: %s", gatewayRegistrar.GetNodeID(), contentID, cID)
	}
	if nonce != nonceRecv {
		return nil, fmt.Errorf("error validating nonce for Client Standard Discover Response for gateway ID: %s; expected nonce: %d, actual nonce: %d", gatewayRegistrar.GetNodeID(), nonce, nonceRecv)
	}
	if paymentRequired {
		return nil, fmt.Errorf("payment required, in order to proceed topup your balance for payment channel address: %d", paymentChannelAddrToTopup)
	}

	return offerDigests, nil
}
