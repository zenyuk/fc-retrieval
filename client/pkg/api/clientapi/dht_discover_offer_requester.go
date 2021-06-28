/*
Package clientapi - a collection of Retrieval Client requests methods used to call FileCoin Secondary Retrieval network.
These functions are playing a role of intermediate facade to call the FileCoin Secondary Retrieval network and
hiding away network transport details
*/
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
	"fmt"

	"github.com/ConsenSys/fc-retrieval/common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval/common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval/common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval/common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval/common/pkg/register"
)

// GatewaySubOffers - relation between a Gateway and the Sub-Offers received through the Gateway
type GatewaySubOffers struct {
	GatewayID *nodeid.NodeID         `json:"gateway_id"`
	SubOffers []cidoffer.SubCIDOffer `json:"sub_cid_offers"`
}

func (c *Client) RequestDHTOfferDiscover(
	gatewayRegistrar register.GatewayRegistrar,
	gatewayIDs []nodeid.NodeID,
	contentID *cid.ContentID,
	nonce int64,
	offersDigests [][][cidoffer.CIDOfferDigestSize]byte,
	paymentChannelAddr string,
	voucher string,
) ([]GatewaySubOffers, error) {
	request, err := fcrmessages.EncodeClientDHTDiscoverOfferRequest(contentID, nonce, offersDigests, gatewayIDs, paymentChannelAddr, voucher)
	if err != nil {
		logging.Error("error encoding Client DHT Discover Request: %+v", err)
		return nil, err
	}

	// Send request and get response
	response, err := c.httpCommunicator.SendMessage(gatewayRegistrar.GetNetworkInfoClient(), request)
	if err != nil {
		return nil, err
	}

	// Get the gateway's public key
	pubKey, err := gatewayRegistrar.GetSigningKey()
	if err != nil {
		return nil, err
	}

	// Verify the response
	if response.Verify(pubKey) != nil {
		return nil, errors.New("verification failed")
	}

	_, _, gatewayIDs, fcrMessages, paymentRequiredCl, paymentChannelAddrToTopupCl, err := fcrmessages.DecodeClientDHTDiscoverOfferResponse(response)
	if err != nil {
		return nil, fmt.Errorf("error decoding client DHT discover offer response %s", err.Error())
	}
	if len(gatewayIDs) != len(fcrMessages) {
		return nil, fmt.Errorf("error decoding client DHT discover offer response, lengths of gateway IDs = %d and FCR messages = %d do not match", len(gatewayIDs), len(fcrMessages))
	}
	if paymentRequiredCl {
		return nil, fmt.Errorf("payment required, in order to proceed topup your balance for payment channel address: %d", paymentChannelAddrToTopupCl)
	}
	var result []GatewaySubOffers
	for idx, fcrMessage := range fcrMessages {
		_, _, found, subCIDOffers, _,  paymentRequired, paymentChannelAddrToTopup, decodeErr := fcrmessages.DecodeGatewayDHTDiscoverOfferResponse(&fcrMessage)
		if decodeErr != nil {
			logging.Error("error decoding gateway DHT discover offer response %s", decodeErr.Error())
		}
		if paymentRequired {
			return nil, fmt.Errorf("payment required, in order to proceed topup your balance for payment channel address: %d", paymentChannelAddrToTopup)
		}
		// return first good one
		if found && len(subCIDOffers) > 0 {
			gatewaysSubOffers := GatewaySubOffers{
				GatewayID: &gatewayIDs[idx],
				SubOffers: subCIDOffers,
			}
			result = append(result, gatewaysSubOffers)
		}
	}
	return result, nil
}
