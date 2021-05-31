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
	req "github.com/ConsenSys/fc-retrieval-common/pkg/request"
)

func RequestDHTOfferDiscover(
	gatewayInfo *register.GatewayRegister,
	contentID *cid.ContentID,
	nonce int64,
	ttl int64,
	numDHT int64,
	offerDigest [cidoffer.CIDOfferDigestSize]byte,
	paymentChannelAddr string,
	voucher string,
) (*cidoffer.SubCIDOffer, error) {
	//todo: ? what to do with offerDigest

	request, err := fcrmessages.EncodeClientDHTDiscoverRequestV2(contentID, nonce, ttl, numDHT, false, paymentChannelAddr, voucher)
	if err != nil {
		logging.Error("Error encoding Client DHT Discover Request: %+v", err)
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

	contacted, contactedResp, _, recvNonce, err := fcrmessages.DecodeClientDHTDiscoverResponseV2(response)
	var result *cidoffer.SubCIDOffer
	//todo: ? how to convert response to SubCIDOffer

	if err != nil {
		return nil, err
	}
	if recvNonce != nonce {
		return nil, errors.New("Nonce not matching")
	}
	if len(contacted) != len(contactedResp) {
		return nil, errors.New("Length mismtach")
	}

	return result, nil
}
