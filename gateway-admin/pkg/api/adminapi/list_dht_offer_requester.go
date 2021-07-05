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

	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval/common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval/common/pkg/register"
)

// RequestListDHTOffer asks a given gateway to list dht offer
func (a *Admin) RequestListDHTOffer(
	adminApiEndpoint string,
	gatewayRegistrar register.GatewayRegistrar,
	signingPrivkey *fcrcrypto.KeyPair,
	signingPrivKeyVer *fcrcrypto.KeyVersion) error {
	// First, Get pubkey
	pubKey, err := gatewayRegistrar.GetSigningKey()
	if err != nil {
		logging.Error("Error in obtaining signing key from register info.")
		return err
	}

	// Second, send key exchange to activate the given gateway
	request, err := fcrmessages.EncodeGatewayAdminListDHTOfferRequest(true)
	if err != nil {
		logging.Error("Error in encoding message.")
		return err
	}
	// Sign the request
	if request.Sign(signingPrivkey, signingPrivKeyVer) != nil {
		return errors.New("error in signing the request")
	}

	response, err := a.httpCommunicator.SendMessage(adminApiEndpoint, request)
	if err != nil {
		logging.Error("Error in sending the message.")
		return err
	}

	// Verify the response
	if response.Verify(pubKey) != nil {
		return errors.New("fail to verify the response")
	}

	refreshed, err := fcrmessages.DecodeGatewayAdminListDHTOfferResponse(response)
	if err != nil {
		logging.Error("Error in decoding the message.")
		return err
	}
	if !refreshed {
		logging.Error("Force refresh failed.")
		return errors.New("fail to force the gateway to refresh")
	}

	return nil
}
