package fcrmsgclient

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
)

// clientStandardDiscoverRequest is the requset from client to gateway to ask for cid offer
type clientStandardDiscoverRequest struct {
	PieceCID  cid.ContentID `json:"piece_cid"`
	Nonce     int64         `json:"nonce"`
	TTL       int64         `json:"ttl"`
	PaychAddr string        `json:"payment_channel_address"`
	Voucher   string        `json:"voucher"`
}

// EncodeClientStandardDiscoverRequest is used to get the FCRMessage of clientStandardDiscoverRequest
func EncodeClientStandardDiscoverRequest(
	pieceCID *cid.ContentID,
	nonce int64,
	ttl int64,
	paychAddr string,
	voucher string,
) (*fcrmessages.FCRMessage, error) {
	body, err := json.Marshal(clientStandardDiscoverRequest{
		PieceCID:  *pieceCID,
		Nonce:     nonce,
		TTL:       ttl,
		PaychAddr: paychAddr,
		Voucher:   voucher,
	})
	if err != nil {
		return nil, err
	}
	return fcrmessages.CreateFCRMessage(fcrmessages.ClientStandardDiscoverRequestType, body), nil
}

// DecodeClientStandardDiscoverRequest is used to get the fields from FCRMessage of clientStandardDiscoverRequest
func DecodeClientStandardDiscoverRequest(fcrMsg *fcrmessages.FCRMessage) (
	*cid.ContentID, // piece cid
	int64, // nonce
	int64, // ttl
	string, // payment channel address
	string, // voucher
	error, // error
) {
	if fcrMsg.GetMessageType() != fcrmessages.ClientStandardDiscoverRequestType {
		return nil, 0, 0, "", "", errors.New("Message type mismatch")
	}
	msg := clientStandardDiscoverRequest{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return nil, 0, 0, "", "", err
	}
	return &msg.PieceCID, msg.Nonce, msg.TTL, msg.PaychAddr, msg.Voucher, nil
}
