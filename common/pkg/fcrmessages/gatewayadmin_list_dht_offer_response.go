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

// gatewayAdminInitialiseKeyResponse is the response to gatewayAdminInitialiseKeyRequest
type gatewayAdminListDHTOfferResponse struct {
	Refreshed bool `json:"refreshed"`
}

// EncodeGatewayAdminListDHTOfferResponse is used to get the FCRMessage of gatewayAdminListDHTOfferResponse
func EncodeGatewayAdminListDHTOfferResponse(refreshed bool) (*FCRMessage, error) {
	body, err := json.Marshal(gatewayAdminListDHTOfferResponse{
		Refreshed: refreshed,
	})
	if err != nil {
		return nil, err
	}
	return CreateFCRMessage(GatewayAdminListDHTOfferResponseType, body), nil
}

// DecodeGatewayAdminListDHTOfferResponse is used to get the fields from FCRMessage of gatewayAdminListDHTOfferResponse
func DecodeGatewayAdminListDHTOfferResponse(fcrMsg *FCRMessage) (
	bool, // refresh
	error, // error
) {
	if fcrMsg.GetMessageType() != GatewayAdminListDHTOfferResponseType {
		return false, errors.New("message type mismatch")
	}
	msg := gatewayAdminListDHTOfferResponse{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return false, err
	}
	return msg.Refreshed, nil
}
