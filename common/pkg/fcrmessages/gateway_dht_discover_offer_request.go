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
	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
)

// gatewayDHTDiscoverOfferRequest is the request to implement the payment system
type gatewayDHTDiscoverOfferRequest struct {
	PieceCID    cid.ContentID                     `json:"piece_cid"`
	Nonce       int64                             `json:"nonce"`
	TTL         int64                             `json:"ttl"`
	OfferDigest [cidoffer.CIDOfferDigestSize]byte `json:"offer_digest"`
	PaychAddr   string                            `json:"payment_channel_address"`
	Voucher     string                            `json:"voucher"`
}

// EncodeGatewayDHTDiscoverOfferRequest is used to get the FCRMessage of gatewayDHTDiscoverOfferRequest
func EncodeGatewayDHTDiscoverOfferRequest(
	pieceCID *cid.ContentID,
	nonce int64,
	ttl int64,
	offerDigest [cidoffer.CIDOfferDigestSize]byte,
	paychAddr string,
	voucher string,
) (*FCRMessage, error) {
	body, err := json.Marshal(gatewayDHTDiscoverOfferRequest{
		PieceCID:    *pieceCID,
		Nonce:       nonce,
		TTL:         ttl,
		OfferDigest: offerDigest,
		PaychAddr:   paychAddr,
		Voucher:     voucher,
	})
	if err != nil {
		return nil, err
	}
	return CreateFCRMessage(GatewayDHTDiscoverOfferRequestType, body), nil
}

// DecodeGatewayDHTDiscoverOfferRequest is used to get the fields from FCRMessage of gatewayDHTDiscoverOfferRequest
func DecodeGatewayDHTDiscoverOfferRequest(fcrMsg *FCRMessage) (
	*cid.ContentID, // piece cid
	int64, // nonce
	int64, // ttl
	[cidoffer.CIDOfferDigestSize]byte, // offer_digest
	string, // payment channel address
	string, // voucher
	error, // error
) {
	if fcrMsg.GetMessageType() != GatewayDHTDiscoverOfferRequestType {
		return nil, 0, 0, [cidoffer.CIDOfferDigestSize]byte{}, "", "", errors.New("Message type mismatch")
	}
	msg := gatewayDHTDiscoverOfferRequest{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return nil, 0, 0, [cidoffer.CIDOfferDigestSize]byte{}, "", "", err
	}
	return &msg.PieceCID, msg.Nonce, msg.TTL, msg.OfferDigest, msg.PaychAddr, msg.Voucher, nil
}
