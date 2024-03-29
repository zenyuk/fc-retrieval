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

	"github.com/ConsenSys/fc-retrieval/common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval/common/pkg/cidoffer"
)

// clientStandardDiscoverOfferResponse is the response to clientStandardDiscoverRequest
type clientStandardDiscoverOfferResponse struct {
	PieceCID             string                 `json:"piece_cid"`
	Nonce                int64                  `json:"nonce"`
	Found                bool                   `json:"found"`
	SubCIDOffers         []cidoffer.SubCIDOffer `json:"sub_cid_offers"`
	FundedPaymentChannel []bool                 `json:"funded_payment_channel"`
	PaymentRequired      bool                   `json:"payment_required"` // when true means caller have to pay first, using the PaymentChannel field
	PaymentChannel       string                 `json:"payment_channel"`  // payment channel address used in conjunction with PaymentRequired field
}

// EncodeClientStandardDiscoverOfferResponse is used to get the FCRMessage of clientStandardDiscoverOfferResponse
func EncodeClientStandardDiscoverOfferResponse(
	pieceCID *cid.ContentID,
	nonce int64,
	found bool,
	offers []cidoffer.SubCIDOffer,
	fundedPaymentChannel []bool,
	paymentRequired bool,
	paymentChannel string,
) (*FCRMessage, error) {
	body, err := json.Marshal(clientStandardDiscoverOfferResponse{
		PieceCID:             pieceCID.ToString(),
		Nonce:                nonce,
		Found:                found,
		SubCIDOffers:         offers,
		FundedPaymentChannel: fundedPaymentChannel,
		PaymentRequired:      paymentRequired,
		PaymentChannel:       paymentChannel,
	})
	if err != nil {
		return nil, err
	}
	return CreateFCRMessage(ClientStandardDiscoverOfferResponseType, body), nil
}

// DecodeClientStandardDiscoverOfferResponse is used to get the fields from FCRMessage of clientStandardDiscoverOfferResponse
func DecodeClientStandardDiscoverOfferResponse(fcrMsg *FCRMessage) (
	*cid.ContentID, // piece cid
	int64, // nonce
	bool, // found
	[]cidoffer.SubCIDOffer, // sub cid offers
	[]bool, // fundedPaymentChannel
	bool, // paymentRequired
	string, // paymentChannel
	error, // error
) {
	if fcrMsg.GetMessageType() != ClientStandardDiscoverOfferResponseType {
		return nil, 0, false, nil, nil, false, "", errors.New("message type mismatch")
	}
	msg := clientStandardDiscoverOfferResponse{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return nil, 0, false, nil, nil, false, "", err
	}
	contentID, _ := cid.NewContentIDFromHexString(msg.PieceCID)
	return contentID, msg.Nonce, msg.Found, msg.SubCIDOffers, msg.FundedPaymentChannel, msg.PaymentRequired, msg.PaymentChannel, nil
}
