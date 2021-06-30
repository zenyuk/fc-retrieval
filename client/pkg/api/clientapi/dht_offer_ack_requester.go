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

	"github.com/ConsenSys/fc-retrieval/common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval/common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval/common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval/common/pkg/register"
)

// RequestDHTOfferAck requests a dht offer ack to a given provider for a pair of cid and gateway id
func (c *Client) RequestDHTOfferAck(
	gatewayRegistrar register.ProviderRegistrar,
	contentID *cid.ContentID,
	gatewayID *nodeid.NodeID,
) (bool, *fcrmessages.FCRMessage, *fcrmessages.FCRMessage, error) {
	// Construct request
	request, err := fcrmessages.EncodeClientDHTOfferAckRequest(contentID, gatewayID)
	if err != nil {
		logging.Error("Error encoding Client DHT Offer Ack Request: %+v", err)
		return false, nil, nil, err
	}

	// Send request and get response
	response, err := c.httpCommunicator.SendMessage(gatewayRegistrar.GetNetworkInfoClient(), request)
	if err != nil {
		return false, nil, nil, err
	}

	// Get the gateway's public key
	pubKey, err := gatewayRegistrar.GetSigningKey()
	if err != nil {
		return false, nil, nil, err
	}

	// Verify the response
	if response.Verify(pubKey) != nil {
		return false, nil, nil, errors.New("verification failed")
	}

	_, _, found, publishDhtOfferReq, publishDhtOfferRes, err := fcrmessages.DecodeClientDHTOfferAckResponse(response)
	if err != nil {
		return false, nil, nil, err
	}

	return found, publishDhtOfferReq, publishDhtOfferRes, nil
}
