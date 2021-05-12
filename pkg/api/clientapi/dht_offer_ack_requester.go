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
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/register"
	req "github.com/ConsenSys/fc-retrieval-common/pkg/request"
)

// RequestDHTOfferAck requests a dht offer ack to a given provider for a pair of cid and gateway id
func RequestDHTOfferAck(
	providerInfo *register.ProviderRegister,
	contentID *cid.ContentID,
	gatewayID *nodeid.NodeID,
) ([]nodeid.NodeID, []cidoffer.SubCIDOffer, error) {
	// Construct request
	request, err := fcrmessages.EncodeClientDHTOfferAckRequest(contentID, gatewayID)
	if err != nil {
		logging.Error("Error encoding Client DHT Offer Ack Request: %+v", err)
		return nil, nil, err
	}

	// Send request and get response
	response, err := req.SendMessage(providerInfo.NetworkInfoClient, request)
	if err != nil {
		return nil, nil, err
	}

	// Get the gateway's public key
	pubKey, err := providerInfo.GetSigningKey()
	if err != nil {
		return nil, nil, err
	}

	// Verify the response
	if response.Verify(pubKey) != nil {
		return nil, nil, errors.New("Verification failed")
	}

	// TODO interpret the response.
	logging.Info("Response from server: %s", response.DumpMessage())
	return nil, nil, nil
}
