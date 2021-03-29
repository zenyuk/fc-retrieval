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

// gatewayListDHTOfferAck is the acknowledgement to gatewayListDHTOfferResponse
type gatewayListDHTOfferAck struct {
	PublishedDHTOffersAck []fcrmessages2.FCRMessage `json:"published_dht_offers_ack"`
}

// EncodeGatewayListDHTOfferAck is used to get the FCRMessage of gatewayListDHTOfferAck
func EncodeGatewayListDHTOfferAck(
	publishedDHTOffersAck []fcrmessages2.FCRMessage,
) (*fcrmessages2.FCRMessage, error) {
	body, err := json.Marshal(gatewayListDHTOfferAck{
		PublishedDHTOffersAck: publishedDHTOffersAck,
	})
	if err != nil {
		return nil, err
	}
	return fcrmessages2.CreateFCRMessage(fcrmessages2.GatewayListDHTOfferAckType, body), nil
}

// DecodeGatewayListDHTOfferAck is used to get the fields from FCRMessage of gatewayListDHTOfferAck
func DecodeGatewayListDHTOfferAck(fcrMsg *fcrmessages2.FCRMessage) (
	[]fcrmessages2.FCRMessage, // published dht cid offers ack
	error, // error
) {
	if fcrMsg.GetMessageType() != fcrmessages2.GatewayListDHTOfferAckType {
		return nil, errors.New("Message type mismatch")
	}
	msg := gatewayListDHTOfferAck{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return nil, err
	}
	return msg.PublishedDHTOffersAck, nil
}
