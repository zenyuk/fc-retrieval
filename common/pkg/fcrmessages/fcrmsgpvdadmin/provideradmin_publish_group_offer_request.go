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

	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
)

// providerAdminPublishGroupOfferRequest is the request to publish group CID
type providerAdminPublishGroupOfferRequest struct {
	CIDs   []cid.ContentID `json:"cids"`
	Price  uint64          `json:"price"`
	Expiry int64           `json:"expiry"`
	QoS    uint64          `json:"qos"`
}

// EncodeProviderAdminPublishGroupOfferRequest is used to get the FCRMessage of providerAdminPublishGroupOfferRequest
func EncodeProviderAdminPublishGroupOfferRequest(
	cids []cid.ContentID,
	price uint64,
	expiry int64,
	qos uint64,
) (*fcrmessages.FCRMessage, error) {
	body, err := json.Marshal(providerAdminPublishGroupOfferRequest{
		CIDs:   cids,
		Price:  price,
		Expiry: expiry,
		QoS:    qos,
	})
	if err != nil {
		return nil, err
	}
	return fcrmessages.CreateFCRMessage(fcrmessages.ProviderAdminPublishGroupOfferRequestType, body), nil
}

// DecodeProviderAdminPublishGroupOfferRequest is used to get the fields from FCRMessage of providerAdminPublishGroupOfferRequest
func DecodeProviderAdminPublishGroupOfferRequest(fcrMsg *fcrmessages.FCRMessage) (
	[]cid.ContentID, // cids
	uint64, // price
	int64, // expity
	uint64, // qos
	error, // error
) {
	if fcrMsg.GetMessageType() != fcrmessages.ProviderAdminPublishGroupOfferRequestType {
		return nil, 0, 0, 0, errors.New("Message type mismatch")
	}
	msg := providerAdminPublishGroupOfferRequest{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	return msg.CIDs, msg.Price, msg.Expiry, msg.QoS, nil
}
