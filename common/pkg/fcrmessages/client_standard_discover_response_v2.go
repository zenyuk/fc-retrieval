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
	"fmt"

	"github.com/ConsenSys/fc-retrieval/common/pkg/cid"
)

// clientStandardDiscoverResponseV2 is the response to clientStandardDiscoverRequest with sub cid offer digests
type clientStandardDiscoverResponseV2 struct {
	PieceCID             string   `json:"piece_cid"`
	Nonce                int64    `json:"nonce"`
	Found                bool     `json:"found"`
	SubCIDOfferDigests   []string `json:"sub_cid_offer_digests"`
	FundedPaymentChannel []bool   `json:"funded_payment_channel"`
	PaymentRequired      bool     `json:"payment_required"` // when true means caller have to pay first, using the PaymentChannel field
	PaymentChannel       string   `json:"payment_channel"`  // payment channel address used in conjunction with PaymentRequired field
}

// EncodeClientStandardDiscoverResponseV2 is used to get the FCRMessage of clientStandardDiscoverResponseV2
func EncodeClientStandardDiscoverResponseV2(
	pieceCID *cid.ContentID,
	nonce int64,
	found bool,
	offersDigests []string,
	fundedPaymentChannel []bool,
	paymentRequired bool,
	paymentChannel string,
) (*FCRMessage, error) {
	body, err := json.Marshal(clientStandardDiscoverResponseV2{
		PieceCID:             pieceCID.ToString(),
		Nonce:                nonce,
		Found:                found,
		SubCIDOfferDigests:   offersDigests,
		FundedPaymentChannel: fundedPaymentChannel,
		PaymentRequired:      paymentRequired,
		PaymentChannel:       paymentChannel,
	})
	if err != nil {
		return nil, err
	}
	return CreateFCRMessage(ClientStandardDiscoverResponseV2Type, body), nil
}

// DecodeClientStandardDiscoverResponseV2 is used to get the fields from FCRMessage of clientStandardDiscoverResponseV2
func DecodeClientStandardDiscoverResponseV2(fcrMsg *FCRMessage) (
	*cid.ContentID, // piece cid
	int64, // nonce
	bool, // found
	[]string, // subCIDOfferDigests
	[]bool, // fundedPaymentChannel
	bool, // paymentRequired
	string, // paymentChannel
	error, // error
) {
	if fcrMsg.GetMessageType() != ClientStandardDiscoverResponseV2Type {
		return nil, 0, false, nil, nil, false, "", fmt.Errorf("message type mismatch, expected type ID: %d, actual type ID: %d", ClientStandardDiscoverResponseV2Type, fcrMsg.GetMessageType())
	}
	msg := clientStandardDiscoverResponseV2{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return nil, 0, false, nil, nil, false, "", err
	}
	contentID, _ := cid.NewContentIDFromHexString(msg.PieceCID)
	return contentID, msg.Nonce, msg.Found, msg.SubCIDOfferDigests, msg.FundedPaymentChannel, msg.PaymentRequired, msg.PaymentChannel, nil
}
