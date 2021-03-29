package fcrmessagesprovideradmin

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

// providerAdminGetPublishedGroupOfferRequest is the requset from provideradmin to provider to ask for published group offers for given gateway ids
type providerAdminGetPublishedGroupOfferRequest struct {
	GatewayIDs []nodeid.NodeID `json:"gateway_id"`
}

// EncodeProviderAdminGetPublishedGroupOfferRequest is used to get the FCRMessage of providerAdminGetPublishedGroupOfferRequest
func EncodeProviderAdminGetPublishedGroupOfferRequest(
	gatewayIDs []nodeid.NodeID,
) (*fcrmessages.FCRMessage, error) {
	body, err := json.Marshal(providerAdminGetPublishedGroupOfferRequest{
		GatewayIDs: gatewayIDs,
	})
	if err != nil {
		return nil, err
	}
	return fcrmessages.CreateFCRMessage(fcrmessages.ProviderAdminGetPublishedGroupOfferRequestType, body), nil
}

// DecodeProviderAdminGetGroupCIDRequest is used to get the fields from FCRMessage of providerAdminGetGroupCIDRequest
func DecodeProviderAdminGetGroupCIDRequest(fcrMsg *fcrmessages.FCRMessage) (
	[]nodeid.NodeID, // piece cids
	error, // error
) {
	if fcrMsg.GetMessageType() != fcrmessages.ProviderAdminGetPublishedGroupOfferRequestType {
		return nil, errors.New("Message type mismatch")
	}
	msg := providerAdminGetPublishedGroupOfferRequest{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return nil, err
	}
	return msg.GatewayIDs, nil
}
