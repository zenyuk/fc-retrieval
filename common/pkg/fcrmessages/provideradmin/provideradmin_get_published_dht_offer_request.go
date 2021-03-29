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

// providerAdminGetPublishedDHTOfferRequest is the requset from provideradmin to provider to ask for published dht offers for given gateway ids
type providerAdminGetPublishedDHTOfferRequest struct {
	GatewayIDs []nodeid.NodeID `json:"gateway_id"`
}

// EncodeProviderAdminGetPublishedDHTOfferRequest is used to get the FCRMessage of providerAdminGetPublishedDHTOfferRequest
func EncodeProviderAdminGetPublishedDHTOfferRequest(
	gatewayIDs []nodeid.NodeID,
) (*fcrmessages.FCRMessage, error) {
	body, err := json.Marshal(providerAdminGetPublishedDHTOfferRequest{
		GatewayIDs: gatewayIDs,
	})
	if err != nil {
		return nil, err
	}
	return fcrmessages.CreateFCRMessage(fcrmessages.ProviderAdminGetPublishedDHTOfferRequestType, body), nil
}

// DecodeProviderAdminGetDHTCIDRequest is used to get the fields from FCRMessage of providerAdminGetDHTCIDRequest
func DecodeProviderAdminGetDHTCIDRequest(fcrMsg *fcrmessages.FCRMessage) (
	[]nodeid.NodeID, // piece cids
	error, // error
) {
	if fcrMsg.GetMessageType() != fcrmessages.ProviderAdminGetPublishedDHTOfferRequestType {
		return nil, errors.New("Message type mismatch")
	}
	msg := providerAdminGetPublishedDHTOfferRequest{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return nil, err
	}
	return msg.GatewayIDs, nil
}
