package fcrmessagesclient

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

	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
)

// clientDHTOfferAckResponse is the response to clientDHTOfferAckRequest
type clientDHTOfferAckResponse struct {
	PieceCID                cid.ContentID          `json:"piece_cid"`
	GatewayID               nodeid.NodeID          `json:"gateway_id"`
	Found                   bool                   `json:"found"`
	PublishDHTOfferRequest  fcrmessages.FCRMessage `json:"publish_dht_offer_request"`
	PublishDHTOfferResponse fcrmessages.FCRMessage `json:"publish_dht_offer_response"`
}

// EncodeClientDHTOfferAckResponse is used to get the FCRMessage of clientDHTOfferAckResponse
func EncodeClientDHTOfferAckResponse(
	pieceCID *cid.ContentID,
	gatewayID *nodeid.NodeID,
	found bool,
	publishDHTOfferRequest *fcrmessages.FCRMessage,
	publishDHTOfferResponse *fcrmessages.FCRMessage,
) (*fcrmessages.FCRMessage, error) {
	body, err := json.Marshal(clientDHTOfferAckResponse{
		PieceCID:                *pieceCID,
		GatewayID:               *gatewayID,
		Found:                   found,
		PublishDHTOfferRequest:  *publishDHTOfferRequest,
		PublishDHTOfferResponse: *publishDHTOfferResponse,
	})
	if err != nil {
		return nil, err
	}
	return fcrmessages.CreateFCRMessage(fcrmessages.ClientDHTOfferAckResponseType, body), nil
}

// DecodeClientDHTOfferAckResponse is used to get the fields from FCRMessage of clientDHTOfferAckResponse
func DecodeClientDHTOfferAckResponse(fcrMsg *fcrmessages.FCRMessage) (
	*cid.ContentID, // piece cid
	*nodeid.NodeID, // gateway id
	bool, // found
	*fcrmessages.FCRMessage, // publish dht offer request
	*fcrmessages.FCRMessage, // publish dht offer resposne
	error, // error
) {
	if fcrMsg.GetMessageType() != fcrmessages.ClientDHTOfferAckResponseType {
		return nil, nil, false, nil, nil, errors.New("Message type mismatch")
	}
	msg := clientDHTOfferAckResponse{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return nil, nil, false, nil, nil, err
	}
	return &msg.PieceCID, &msg.GatewayID, msg.Found, &msg.PublishDHTOfferRequest, &msg.PublishDHTOfferResponse, nil
}
