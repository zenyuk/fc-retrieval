package fcrmessagesgatewayadmin

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

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
)

// gatewayAdminGetReputationResponse is the response to gatewayAdminGetReputationRequest
type gatewayAdminGetReputationResponse struct {
	ClientID   nodeid.NodeID `json:"client_id"`
	Reputation int64         `json:"reputation"`
	Exists     bool          `json:"exists"`
}

// EncodeGatewayAdminGetReputationResponse is used to get the FCRMessage of gatewayAdminGetReputationResponse
func EncodeGatewayAdminGetReputationResponse(
	clientID *nodeid.NodeID,
	reputation int64,
	exists bool,
) (*fcrmessages.FCRMessage, error) {
	body, err := json.Marshal(gatewayAdminGetReputationResponse{
		ClientID:   *clientID,
		Reputation: reputation,
		Exists:     exists,
	})
	if err != nil {
		return nil, err
	}
	return fcrmessages.CreateFCRMessage(fcrmessages.GatewayAdminGetReputationResponseType, body), nil
}

// DecodeGatewayAdminGetReputationResponse is used to get the fields from FCRMessage of gatewayAdminGetReputationResponse
func DecodeGatewayAdminGetReputationResponse(fcrMsg *fcrmessages.FCRMessage) (
	*nodeid.NodeID, // client id
	int64, // reputation
	bool, // exists
	error, // error
) {
	if fcrMsg.GetMessageType() != fcrmessages.GatewayAdminGetReputationResponseType {
		return nil, 0, false, errors.New("Message type mismatch")
	}
	msg := gatewayAdminGetReputationResponse{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return nil, 0, false, err
	}
	return &msg.ClientID, msg.Reputation, msg.Exists, nil
}
