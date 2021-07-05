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
	"github.com/ConsenSys/fc-retrieval/common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval/common/pkg/register"
)

// RequestInitialiseKey initialise a given gateway
func (a *Admin) RequestInitialiseKey(
	adminApiEndpoint string,
	gatewayRegistrar register.GatewayRegistrar,
	gatewayPrivKey *fcrcrypto.KeyPair,
	gatewayPrivKeyVer *fcrcrypto.KeyVersion,
	signingPrivkey *fcrcrypto.KeyPair,
	signingPrivKeyVer *fcrcrypto.KeyVersion) error {
	// First, Get pubkey
	pubKey, err := gatewayRegistrar.GetSigningKey()
	if err != nil {
		logging.Error("Error in obtaining signing key from register info.")
		return err
	}
	nodeID, err := nodeid.NewNodeIDFromHexString(gatewayRegistrar.GetNodeID())
	if err != nil {
		logging.Error("Error in generating nodeID.")
		return err
	}
	// Second, send key exchange to activate the given gateway
	request, err := fcrmessages.EncodeGatewayAdminInitialiseKeyRequest(nodeID, gatewayPrivKey, gatewayPrivKeyVer)
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
	ok, err := fcrmessages.DecodeGatewayAdminInitialiseKeyResponse(response)
	if err != nil {
		logging.Error("Error in decoding the message.")
		return err
	}
	if !ok {
		logging.Error("Initialise gateway failed.")
		return errors.New("fail to initialise gateway")
	}

	return nil
}
