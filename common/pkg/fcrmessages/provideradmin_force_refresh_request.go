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

// providerAdminForceRefreshRequest is the request from an admin client to a provider to refresh internal register status
type providerAdminForceRefreshRequest struct {
	Refresh bool `json:"refresh"`
}

// EncodeProviderAdminForceRefreshRequest is used to get the FCRMessage of providerAdminForceRefreshRequest
func EncodeProviderAdminForceRefreshRequest(refresh bool) (*FCRMessage, error) {
	body, err := json.Marshal(providerAdminForceRefreshRequest{
		Refresh: refresh,
	})
	if err != nil {
		return nil, err
	}
	return CreateFCRMessage(ProviderAdminForceRefreshRequestType, body), nil
}

// DecodeProviderAdminForceRefreshRequest is used to get the fields from FCRMessage of providerAdminForceRefreshRequest
func DecodeProviderAdminForceRefreshRequest(fcrMsg *FCRMessage) (
	bool, // refresh
	error, // error
) {
	if fcrMsg.GetMessageType() != ProviderAdminForceRefreshRequestType {
		return false, errors.New("message type mismatch")
	}
	msg := providerAdminForceRefreshRequest{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return false, err
	}
	return msg.Refresh, nil
}
