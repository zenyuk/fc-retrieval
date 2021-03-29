package fcrmessagesbasic

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

// invalidMessageResponse message is sent to indicate that the message is invalid
type invalidMessageResponse struct {
}

// EncodeInvalidMessageResponse is used to get the FCRMessage of invalidMessageResponse
func EncodeInvalidMessageResponse() (*fcrmessages.FCRMessage, error) {
	body, err := json.Marshal(invalidMessageResponse{})
	if err != nil {
		return nil, err
	}
	return fcrmessages.CreateFCRMessage(fcrmessages.InvalidMessageResponseType, body), nil
}

// DecodeInvalidMessageResponse is used to get the fields from FCRMessage of invalidMessageResponse
func DecodeInvalidMessageResponse(fcrMsg *fcrmessages.FCRMessage) error {
	if fcrMsg.GetMessageType() != fcrmessages.InvalidMessageResponseType {
		return errors.New("Message type mismatch")
	}
	msg := invalidMessageResponse{}
	return json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
}
