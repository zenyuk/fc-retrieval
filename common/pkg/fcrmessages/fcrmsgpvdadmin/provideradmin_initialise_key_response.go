package fcrmsgpvdadmin

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
)

// providerAdminInitialiseKeyResponse is the response to providerAdminInitialiseKeyRequest
type providerAdminInitialiseKeyResponse struct {
	Success bool `json:"success"`
}

// EncodeProviderAdminInitialiseKeyResponse is used to get the FCRMessage of providerAdminInitialiseKeyResponse
func EncodeProviderAdminInitialiseKeyResponse(
	success bool,
) (*fcrmessages.FCRMessage, error) {
	body, err := json.Marshal(providerAdminInitialiseKeyResponse{
		Success: success,
	})
	if err != nil {
		return nil, err
	}
	return fcrmessages.CreateFCRMessage(fcrmessages.ProviderAdminInitialiseKeyResponseType, body), nil
}

// DecodeProviderAdminInitialiseKeyResponse is used to get the fields from FCRMessage of ProviderAdminInitialiseKeyResponse
func DecodeProviderAdminInitialiseKeyResponse(fcrMsg *fcrmessages.FCRMessage) (
	bool, // success
	error, // error
) {
	if fcrMsg.GetMessageType() != fcrmessages.ProviderAdminInitialiseKeyResponseType {
		return false, errors.New("Message type mismatch")
	}
	msg := providerAdminInitialiseKeyResponse{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return false, err
	}
	return msg.Success, nil
}
