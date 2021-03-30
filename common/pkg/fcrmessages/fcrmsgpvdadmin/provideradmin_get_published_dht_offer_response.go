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

	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
)

// providerAdminGetPublishedDHTOfferResponse is the response to providerAdminGetPublishedDHTOfferRequest
type providerAdminGetPublishedDHTOfferResponse struct {
	Exists bool                `json:"exists"`
	Offers []cidoffer.CIDOffer `json:"cid_offers"`
}

// EncodeProviderAdminGetPublishedDHTOfferResponse is used to get the FCRMessage of providerAdminGetPublishedDHTOfferResponse
func EncodeProviderAdminGetPublishedDHTOfferResponse(
	exists bool,
	offers []cidoffer.CIDOffer,
) (*fcrmessages.FCRMessage, error) {
	body, err := json.Marshal(providerAdminGetPublishedDHTOfferResponse{
		Exists: exists,
		Offers: offers,
	})
	if err != nil {
		return nil, err
	}
	return fcrmessages.CreateFCRMessage(fcrmessages.ProviderAdminGetPublishedDHTOfferResponseType, body), nil
}

// DecodeProviderAdminGetDHTCIDResponse is used to get the fields from FCRMessage of providerAdminGetDHTCIDResponse
func DecodeProviderAdminGetDHTCIDResponse(fcrMsg *fcrmessages.FCRMessage) (
	bool, // exists
	[]cidoffer.CIDOffer, // cid offers
	error, // error
) {
	if fcrMsg.GetMessageType() != fcrmessages.ProviderAdminGetPublishedDHTOfferResponseType {
		return false, nil, errors.New("Message type mismatch")
	}
	msg := providerAdminGetPublishedDHTOfferResponse{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return false, nil, err
	}
	// Check every offer is single offer
	for _, offer := range msg.Offers {
		if len(offer.GetCIDs()) != 1 {
			return false, nil, errors.New("Offers contain group offer")
		}
	}

	return msg.Exists, msg.Offers, nil
}
