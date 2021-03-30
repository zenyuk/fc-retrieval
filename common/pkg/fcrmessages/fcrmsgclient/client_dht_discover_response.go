package fcrmsgclient

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

// clientDHTDiscoverResponse is the response to clientDHTDiscoverRequest
type clientDHTDiscoverResponse struct {
	Contacted     []fcrmessages.FCRMessage `json:"contacted_gateways"`
	UnContactable []nodeid.NodeID          `json:"uncontactable_gateways"`
	Nonce         int64                    `json:"nonce"`
}

// EncodeClientDHTDiscoverResponse is used to get the FCRMessage of ClientDHTDiscoverResponse
func EncodeClientDHTDiscoverResponse(
	contacted []fcrmessages.FCRMessage,
	unContactable []nodeid.NodeID,
	nonce int64,
) (*fcrmessages.FCRMessage, error) {
	body, err := json.Marshal(clientDHTDiscoverResponse{
		Contacted:     contacted,
		UnContactable: unContactable,
		Nonce:         nonce,
	})
	if err != nil {
		return nil, err
	}
	return fcrmessages.CreateFCRMessage(fcrmessages.ClientDHTDiscoverResponseType, body), nil
}

// DecodeClientDHTDiscoverResponse is used to get the fields from FCRMessage of ClientDHTDiscoverResponse
func DecodeClientDHTDiscoverResponse(fcrMsg *fcrmessages.FCRMessage) (
	[]fcrmessages.FCRMessage, // contacted
	[]nodeid.NodeID, // uncontactable
	int64, // nonce
	error, // error
) {
	if fcrMsg.GetMessageType() != fcrmessages.ClientDHTDiscoverResponseType {
		return nil, nil, 0, errors.New("Message type mismatch")
	}
	msg := clientDHTDiscoverResponse{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return nil, nil, 0, err
	}
	return msg.Contacted, msg.UnContactable, msg.Nonce, nil
}
