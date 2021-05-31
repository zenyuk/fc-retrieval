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

	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/register"
	req "github.com/ConsenSys/fc-retrieval-common/pkg/request"
)

func RequestDHTOfferDiscover(
	gatewayInfo *register.GatewayRegister,
	gatewayIDs []nodeid.NodeID,
	contentID *cid.ContentID,
	nonce int64,
	offersDigests [][][cidoffer.CIDOfferDigestSize]byte,
	paymentChannelAddr string,
	voucher string,
) (*cidoffer.SubCIDOffer, error) {
	request, err := fcrmessages.EncodeClientDHTDiscoverOfferRequest(contentID, nonce, offersDigests, gatewayIDs, paymentChannelAddr, voucher)
	if err != nil {
		logging.Error("error encoding Client DHT Discover Request: %+v", err)
		return nil, err
	}

	// Send request and get response
	response, err := req.SendMessage(gatewayInfo.NetworkInfoClient, request)
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

	// TODO: currently getting the first subCIDOffer
	_, _, _, fcrMessages, err := fcrmessages.DecodeClientDHTDiscoverOfferResponse(response)
	if err != nil {
		return nil, fmt.Errorf("error decoding client DHT discover offer response %s", err.Error())
	}
	for _, fcrMessage := range fcrMessages {
		_, _, found, subCIDOffers, _, decodeErr := fcrmessages.DecodeGatewayDHTDiscoverOfferResponse(&fcrMessage)
		if decodeErr != nil {
			logging.Error("error decoding gateway DHT discover offer response %s", decodeErr.Error())
		}
		// return first good one
		if found && len(subCIDOffers) > 0 {
			return &subCIDOffers[0], nil
		}
	}

	return nil, nil
}
