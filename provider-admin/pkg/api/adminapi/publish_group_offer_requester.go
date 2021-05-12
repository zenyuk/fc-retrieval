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

	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	req "github.com/ConsenSys/fc-retrieval-common/pkg/request"
)

// RequestPublishGroupOffer publish a group cid offer to a given provider
func RequestPublishGroupOffer(
	adminAP string,
	providerPubKey *fcrcrypto.KeyPair,
	cids []cid.ContentID,
	price uint64,
	expiry int64,
	qos uint64,
	signingPrivkey *fcrcrypto.KeyPair,
	signingPrivKeyVer *fcrcrypto.KeyVersion,
) error {
	// Construct request
	request, err := fcrmessages.EncodeProviderAdminPublishGroupOfferRequest(cids, price, expiry, qos)
	// Sign the request
	if request.Sign(signingPrivkey, signingPrivKeyVer) != nil {
		return errors.New("Error in signing the request")
	}

	response, err := req.SendMessage(adminAP, request)
	if err != nil {
		logging.Error("Error in sending the message.")
		return err
	}

	// Verify the response
	if response.Verify(providerPubKey) != nil {
		return errors.New("Fail to verify the response")
	}

	received, err := fcrmessages.DecodeProviderAdminPublishGroupOfferResponse(response)
	if err != nil {
		logging.Error("Error in decoding the message.")
		return err
	}
	if !received {
		logging.Error("Publish offer failed.")
		return errors.New("Fail to publish offer")
	}
	return nil
}
