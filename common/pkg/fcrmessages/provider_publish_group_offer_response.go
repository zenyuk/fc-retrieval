package fcrmessages

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
	"encoding/json"
	"errors"

	"github.com/ConsenSys/fc-retrieval/common/pkg/nodeid"
)

// providerPublishGroupOfferResponse is the response to providerPublishGroupOfferRequest
type providerPublishGroupOfferResponse struct {
	GatewaydID string `json:"gateway_id"`
	Digest     string `json:"digest"`
}

// EncodeProviderPublishGroupOfferResponse is used to get the FCRMessage of ProviderPublishGroupOfferResponse
func EncodeProviderPublishGroupOfferResponse(
	gatewayID nodeid.NodeID,
	digest string,
) (*FCRMessage, error) {
	body, err := json.Marshal(providerPublishGroupOfferResponse{
		GatewaydID: gatewayID.ToString(),
		Digest:     digest,
	})
	if err != nil {
		return nil, err
	}
	return CreateFCRMessage(ProviderPublishGroupOfferResponseType, body), nil
}

// DecodeProviderPublishGroupOfferResponse is used to get the fields from FCRMessage of ProviderPublishGroupOfferResponse
func DecodeProviderPublishGroupOfferResponse(fcrMsg *FCRMessage) (
	*nodeid.NodeID, // gatewayID
	string, // digest
	error, // error
) {
	if fcrMsg.GetMessageType() != ProviderPublishGroupOfferResponseType {
		return nil, "", errors.New("message type mismatch")
	}
	msg := providerPublishGroupOfferResponse{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return nil, "", err
	}
	nodeID, _ := nodeid.NewNodeIDFromHexString(msg.GatewaydID)
	return nodeID, msg.Digest, nil
}
