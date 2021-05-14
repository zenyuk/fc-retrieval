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
)

// gatewayPingResponse is the response to gatewayPingRequest
type gatewayPingResponse struct {
	Nonce   int64 `json:"nonce"`
	IsAlive bool  `json:"isAlive"`
}

// EncodeGatewayPingResponse is used to get the FCRMessage of gatewayPingResponse
func EncodeGatewayPingResponse(
	nonce int64,
	alive bool,
) (*FCRMessage, error) {
	body, err := json.Marshal(gatewayPingResponse{
		Nonce:   nonce,
		IsAlive: alive,
	})

	if err != nil {
		return nil, err
	}

	return CreateFCRMessage(GatewayPingResponseType, body), nil
}

// DecodeGatewayPingResponse is used to get the fields from FCRMessage of GatewayPingResponse
func DecodeGatewayPingResponse(fcrMsg *FCRMessage) (
	int64, // nonce
	bool, // found
	error, // error
) {
	if fcrMsg.GetMessageType() != GatewayPingResponseType {
		return 0, false, errors.New("message type mismatch")
	}

	msg := gatewayPingResponse{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)

	if err != nil {
		return 0, false, err
	}

	return msg.Nonce, true, nil
}
