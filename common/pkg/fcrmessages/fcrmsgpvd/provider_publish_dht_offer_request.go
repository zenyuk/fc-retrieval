package fcrmsgpvd

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
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
)

// providerPublishDHTOfferRequest is the request from provider to gateway to publish dht offer
type providerPublishDHTOfferRequest struct {
	Nonce      int64               `json:"nonce"`
	ProviderID nodeid.NodeID       `json:"provider_id"`
	NumOffers  int64               `json:"num_of_offers"`
	Offers     []cidoffer.CIDOffer `json:"single_offers"`
}

// EncodeProviderPublishDHTOfferRequest is used to get the FCRMessage of providerPublishDHTOfferRequest
func EncodeProviderPublishDHTOfferRequest(
	nonce int64,
	providerID *nodeid.NodeID,
	offers []cidoffer.CIDOffer,
) (*fcrmessages.FCRMessage, error) {
	body, err := json.Marshal(providerPublishDHTOfferRequest{
		Nonce:      nonce,
		ProviderID: *providerID,
		NumOffers:  int64(len(offers)),
		Offers:     offers,
	})
	if err != nil {
		return nil, err
	}
	return fcrmessages.CreateFCRMessage(fcrmessages.ProviderPublishDHTOfferRequestType, body), nil
}

// DecodeProviderPublishDHTOfferRequest is used to get the fields from FCRMessage of providerPublishDHTOfferRequest
func DecodeProviderDHTPublishGroupCIDRequest(fcrMsg *fcrmessages.FCRMessage) (
	int64, // nonce
	*nodeid.NodeID, // provider id
	[]cidoffer.CIDOffer, // offers
	error, // error
) {
	if fcrMsg.GetMessageType() != fcrmessages.ProviderPublishDHTOfferRequestType {
		return 0, nil, nil, errors.New("Message type mismatch")
	}
	msg := providerPublishDHTOfferRequest{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return 0, nil, nil, err
	}
	// Check every offer is single offer
	for _, offer := range msg.Offers {
		if len(offer.GetCIDs()) != 1 {
			return 0, nil, nil, errors.New("Offers contain group offer")
		}
	}

	return msg.Nonce, &msg.ProviderID, msg.Offers, nil
}
