package fcrmsggwadmin

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

// gatewayAdminSetReputationRequest is the request from an admin client to a gateway to set a client's reputation
type gatewayAdminSetReputationRequest struct {
	ClientID   nodeid.NodeID `json:"client_id"`
	Reputation int64         `json:"reputation"`
}

// EncodeGatewayAdminSetReputationRequest is used to get the FCRMessage of gatewayAdminSetReputationRequest
func EncodeGatewayAdminSetReputationRequest(
	clientID *nodeid.NodeID,
	reputation int64,
) (*fcrmessages.FCRMessage, error) {
	body, err := json.Marshal(gatewayAdminSetReputationRequest{
		ClientID:   *clientID,
		Reputation: reputation,
	})
	if err != nil {
		return nil, err
	}
	return fcrmessages.CreateFCRMessage(fcrmessages.GatewayAdminSetReputationRequestType, body), nil
}

// DecodeGatewayAdminSetReputationRequest is used to get the fields from FCRMessage of gatewayAdminSetReputationRequest
func DecodeGatewayAdminSetReputationRequest(fcrMsg *fcrmessages.FCRMessage) (
	*nodeid.NodeID, // client id
	int64, // reputation
	error, // error
) {
	if fcrMsg.GetMessageType() != fcrmessages.GatewayAdminSetReputationRequestType {
		return nil, 0, errors.New("Message type mismatch")
	}
	msg := gatewayAdminSetReputationRequest{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return nil, 0, err
	}
	return &msg.ClientID, msg.Reputation, nil
}
