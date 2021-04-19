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

// gatewayAdminForceRefreshResponse is the response to gatewayAdminForceRefreshRequest
type gatewayAdminForceRefreshResponse struct {
	Refreshed bool `json:"refreshed"`
}

// EncodeGatewayAdminForceRefreshResponse is used to get the FCRMessage of gatewayAdminForceRefreshResponse
func EncodeGatewayAdminForceRefreshResponse(
	refreshed bool,
) (*FCRMessage, error) {
	body, err := json.Marshal(gatewayAdminForceRefreshResponse{
		Refreshed: refreshed,
	})
	if err != nil {
		return nil, err
	}
	return CreateFCRMessage(GatewayAdminForceRefreshResponseType, body), nil
}

// DecodeGatewayAdminForceRefreshResponse is used to get the fields from FCRMessage of gatewayAdminForceRefreshResponse
func DecodeGatewayAdminForceRefreshResponse(fcrMsg *FCRMessage) (
	bool, // refreshed
	error, // error
) {
	if fcrMsg.GetMessageType() != GatewayAdminForceRefreshResponseType {
		return false, errors.New("message type mismatch")
	}
	msg := gatewayAdminForceRefreshResponse{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return false, err
	}
	return msg.Refreshed, nil
}
