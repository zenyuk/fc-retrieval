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
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/register"
	req "github.com/ConsenSys/fc-retrieval-common/pkg/request"
)

// RequestInitialiseKeyV2 initialise a given provider
func RequestInitialiseKeyV2(
	providerInfo *register.ProviderRegister,
	providerPrivKey *fcrcrypto.KeyPair,
	providerPrivKeyVer *fcrcrypto.KeyVersion,
	signingPrivkey *fcrcrypto.KeyPair,
	signingPrivKeyVer *fcrcrypto.KeyVersion,
	lotusWalletPrivateKey string,
	lotusAP string,
	lotusAuthToken string,
) error {
	// First, Get pubkey
	pubKey, err := providerInfo.GetSigningKey()
	if err != nil {
		logging.Error("Error in obtaining signing key from register info.")
		return err
	}
	nodeID, err := nodeid.NewNodeIDFromHexString(providerInfo.NodeID)
	if err != nil {
		logging.Error("Error in generating nodeID.")
		return err
	}
	// Second, send key exchange to activate the given provider
	request, err := fcrmessages.EncodeProviderAdminInitialiseKeyRequest(nodeID, providerPrivKey, providerPrivKeyVer)
	if err != nil {
		logging.Error("Error in encoding message.")
		return err
	}
	// Sign the request
	if request.Sign(signingPrivkey, signingPrivKeyVer) != nil {
		return errors.New("Error in signing the request")
	}

	response, err := req.SendMessage(providerInfo.NetworkInfoAdmin, request)
	if err != nil {
		logging.Error("Error in sending the message.")
		return err
	}

	// Verify the response
	if response.Verify(pubKey) != nil {
		return errors.New("Fail to verify the response")
	}

	ok, err := fcrmessages.DecodeProviderAdminInitialiseKeyResponse(response)
	if err != nil {
		logging.Error("Error in decoding the message.")
		return err
	}
	if !ok {
		logging.Error("Initialise provider failed.")
		return errors.New("Fail to initialise provider")
	}

	return nil
}
