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
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/register"
)

// RequestDHTDiscoverV2 requests a dht discover to a given gateway for a given contentID, nonce and ttl.
func (c *Client) RequestDHTDiscoverV2(
	gatewayRegistrar register.GatewayRegistrar,
	contentID *cid.ContentID,
	nonce int64,
	ttl int64,
	numDHT int64,
	incrementalResult bool,
	paychAddr string,
	voucher string,
) ([]nodeid.NodeID, []fcrmessages.FCRMessage, []nodeid.NodeID, error) {
	// Construct request
	request, err := fcrmessages.EncodeClientDHTDiscoverRequestV2(contentID, nonce, ttl, numDHT, incrementalResult, paychAddr, voucher)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error encoding Client DHT Discover Request: %s", err.Error())
	}

	// Send request and get response
	response, err := c.httpCommunicator.SendMessage(gatewayRegistrar.GetNetworkInfoClient(), request)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error sending DHT discover message to gateway ID: %s, error: %s", gatewayRegistrar.GetNodeID(), err.Error())
	}

	// Get the gateway's public key
	pubKey, err := gatewayRegistrar.GetSigningKey()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error getting signing key of gateway ID: %s, error: %s", gatewayRegistrar.GetNodeID(), err.Error())
	}

	// Verify the response
	if response.Verify(pubKey) != nil {
		return nil, nil, nil, fmt.Errorf("DHT discover response verification failed for gateway ID: %s, message type ID: %d", gatewayRegistrar.GetNodeID(), request.GetMessageType())
	}

	contacted, contactedResp, uncontactable, recvNonce, paymentRequired, paymentChannelAddrToTopup, err := fcrmessages.DecodeClientDHTDiscoverResponseV2(response)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error decoding DHT discover response: %s, gateway ID: %s", err.Error(), gatewayRegistrar.GetNodeID())
	}
	if recvNonce != nonce {
		return nil, nil, nil, fmt.Errorf("error validating nonce for DHT discover response for gateway ID: %s; expected nonce: %d, actual nonce: %d", gatewayRegistrar.GetNodeID(), nonce, recvNonce)
	}
	if len(contacted) != len(contactedResp) {
		return nil, nil, nil, fmt.Errorf("length mismatch error during DHT discover response validation for gateway ID: %s", gatewayRegistrar.GetNodeID())
	}
	if paymentRequired {
		return nil, nil, nil, fmt.Errorf("payment required, in order to proceed topup your balance for payment channel address: %d", paymentChannelAddrToTopup)
	}

	return contacted, contactedResp, uncontactable, nil
}
