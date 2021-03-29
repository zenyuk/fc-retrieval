package fcrmessagesgateway

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

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages2"
)

// gatewayListDHTOfferResponse is the repsonse to gatewayListDHTOfferRequest
type gatewayListDHTOfferResponse struct {
	PublishedDHTOffers []fcrmessages2.FCRMessage `json:"published_dht_offers"`
}

// EncodeGatewayListDHTOfferResponse is used to get the FCRMessage of GatewayListDHTOfferResponse
func EncodeGatewaySingleCIDOfferPublishResponse(
	publishedDHTOffers []fcrmessages2.FCRMessage,
) (*fcrmessages2.FCRMessage, error) {
	body, err := json.Marshal(gatewayListDHTOfferResponse{
		PublishedDHTOffers: publishedDHTOffers,
	})
	if err != nil {
		return nil, err
	}
	return fcrmessages2.CreateFCRMessage(fcrmessages2.GatewayListDHTOfferResponseType, body), nil
}

// DecodeGatewayListDHTOfferResponse is used to get the fields from FCRMessage of GatewayListDHTOfferResponse
func DecodeGatewayListDHTOfferResponse(fcrMsg *fcrmessages2.FCRMessage) (
	[]fcrmessages2.FCRMessage, // published dht cid offers
	error, // error
) {
	if fcrMsg.GetMessageType() != fcrmessages2.GatewayListDHTOfferResponseType {
		return nil, errors.New("Message type mismatch")
	}
	msg := gatewayListDHTOfferResponse{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return nil, err
	}
	return msg.PublishedDHTOffers, nil
}
