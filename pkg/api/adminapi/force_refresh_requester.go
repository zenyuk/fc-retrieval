/*
Package adminapi - remote API for FileCoin Secondary Retrieval - Provider Admin library
*/
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

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/register"
	req "github.com/ConsenSys/fc-retrieval-common/pkg/request"
)

// RequestForceRefresh forces a given provider to refresh its internal register
func RequestForceRefresh(
	providerInfo *register.ProviderRegister,
	signingPrivkey *fcrcrypto.KeyPair,
	signingPrivKeyVer *fcrcrypto.KeyVersion) error {
	// First, Get pubkey
	pubKey, err := providerInfo.GetSigningKey()
	if err != nil {
		logging.Error("Error in obtaining signing key from register info.")
		return err
	}

	// Second, send key exchange to activate the given provider
	request, err := fcrmessages.EncodeProviderAdminForceRefreshRequest(true)
	if err != nil {
		logging.Error("Error in encoding message.")
		return err
	}
	// Sign the request
	if request.Sign(signingPrivkey, signingPrivKeyVer) != nil {
		return errors.New("error in signing the request")
	}

	response, err := req.SendMessage(providerInfo.NetworkInfoAdmin, request)
	if err != nil {
		logging.Error("Error in sending the message.")
		return err
	}

	// Verify the response
	if response.Verify(pubKey) != nil {
		return errors.New("fail to verify the response")
	}

	refreshed, err := fcrmessages.DecodeProviderAdminForceRefreshResponse(response)
	if err != nil {
		logging.Error("Error in decoding the message.")
		return err
	}
	if !refreshed {
		logging.Error("Force refresh failed.")
		return errors.New("fail to force the provider to refresh")
	}

	return nil
}
