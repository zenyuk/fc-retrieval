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

	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
)

// clientDHTOfferAckRequest is the request from client to provider to request the signed ack of a dht offer publish
type clientDHTOfferAckRequest struct {
	PieceCID  cid.ContentID `json:"piece_cid"`
	GatewayID nodeid.NodeID `json:"gateway_id"`
}

// EncodeClientDHTOfferAckRequest is used to get the FCRMessage of clientDHTOfferAckRequest
func EncodeClientDHTOfferAckRequest(
	pieceCID *cid.ContentID,
	gatewayID *nodeid.NodeID,
) (*FCRMessage, error) {
	body, err := json.Marshal(clientDHTOfferAckRequest{
		PieceCID:  *pieceCID,
		GatewayID: *gatewayID,
	})
	if err != nil {
		return nil, err
	}
	return CreateFCRMessage(ClientDHTOfferAckRequestType, body), nil
}

// DeodeClientDHTOfferAckRequest is used to get the fields from FCRMessage of clientDHTOfferAckRequest
func DeodeClientDHTOfferAckRequest(fcrMsg *FCRMessage) (
	*cid.ContentID, // piece cid
	*nodeid.NodeID, // gateway id
	error, // error
) {
	if fcrMsg.GetMessageType() != ClientDHTDiscoverRequestType {
		return nil, nil, errors.New("Message type mismatch")
	}
	msg := clientDHTOfferAckRequest{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return nil, nil, err
	}
	return &msg.PieceCID, &msg.GatewayID, nil
}
