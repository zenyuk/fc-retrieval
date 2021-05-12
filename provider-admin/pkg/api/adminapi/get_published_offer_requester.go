package adminapi

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
	"errors"

	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	req "github.com/ConsenSys/fc-retrieval-common/pkg/request"
)

// RequestGetPublishedOffer checks the group offer stored in the provider for a given list of gateways.
func RequestGetPublishedOffer(
	adminAP string,
	providerPubKey *fcrcrypto.KeyPair,
	gatewayIDs []nodeid.NodeID,
	signingPrivkey *fcrcrypto.KeyPair,
	signingPrivKeyVer *fcrcrypto.KeyVersion,
) (
	bool, // found
	[]cidoffer.CIDOffer, // offers
	error, // error
) {
	// Construct request
	request, err := fcrmessages.EncodeProviderAdminGetPublishedOfferRequest(gatewayIDs)
	// Sign the request
	if request.Sign(signingPrivkey, signingPrivKeyVer) != nil {
		return false, nil, errors.New("Error in signing the request")
	}

	response, err := req.SendMessage(adminAP, request)
	if err != nil {
		logging.Error("Error in sending the message.")
		return false, nil, err
	}

	// Verify the response
	if response.Verify(providerPubKey) != nil {
		return false, nil, errors.New("Fail to verify the response")
	}

	found, offers, err := fcrmessages.DecodeProviderAdminGetPublishedOfferResponse(response)
	if err != nil {
		logging.Error("Error in decoding the message")
		return false, nil, err
	}
	return found, offers, nil
}
