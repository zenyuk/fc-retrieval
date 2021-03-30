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

// gatewayAdminGetReputationRequest is the request from an admin client to a gateway to discover a client's reputation
type gatewayAdminGetReputationRequest struct {
	ClientID nodeid.NodeID `json:"client_id"`
}

// EncodeGatewayAdminGetReputationRequest is used to get the FCRMessage of gatewayAdminGetReputationRequest
func EncodeGatewayAdminGetReputationRequest(
	clientID *nodeid.NodeID,
) (*fcrmessages.FCRMessage, error) {
	body, err := json.Marshal(gatewayAdminGetReputationRequest{
		ClientID: *clientID,
	})
	if err != nil {
		return nil, err
	}
	return fcrmessages.CreateFCRMessage(fcrmessages.GatewayAdminGetReputationRequestType, body), nil
}

// DecodeGatewayAdminGetReputationRequest is used to get the fields from FCRMessage of gatewayAdminGetReputationRequest
func DecodeGatewayAdminGetReputationRequest(fcrMsg *fcrmessages.FCRMessage) (
	*nodeid.NodeID, // client id
	error, // error
) {
	if fcrMsg.GetMessageType() != fcrmessages.GatewayAdminGetReputationRequestType {
		return nil, errors.New("Message type mismatch")
	}
	msg := gatewayAdminGetReputationRequest{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return nil, err
	}
	return &msg.ClientID, nil
}
