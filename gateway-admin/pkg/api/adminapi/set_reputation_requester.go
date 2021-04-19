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

// RequestSetReputation sets the reputation for a given client id
func RequestSetReputation(
	gatewayInfo *register.GatewayRegister,
	clientID *nodeid.NodeID,
	reputation int64,
	signingPrivkey *fcrcrypto.KeyPair,
	signingPrivKeyVer *fcrcrypto.KeyVersion,
) (bool, error) {
	// First, Get pubkey
	pubKey, err := gatewayInfo.GetSigningKey()
	if err != nil {
		logging.Error("Error in obtaining signing key from register info.")
		return false, err
	}

	// Second, send get reputation to given gateway
	request, err := fcrmessages.EncodeGatewayAdminSetReputationRequest(clientID, reputation)
	if err != nil {
		logging.Error("Error in encoding message.")
		return false, err
	}
	// Sign the request
	if request.Sign(signingPrivkey, signingPrivKeyVer) != nil {
		return false, errors.New("Error in signing the request")
	}

	response, err := req.SendMessage(gatewayInfo.NetworkInfoAdmin, request)
	if err != nil {
		logging.Error("Error in sending the message.")
		return false, err
	}

	// Verify the response
	if response.Verify(pubKey) != nil {
		return false, errors.New("Fail to verify the response")
	}

	targetID, reputationNew, exists, err := fcrmessages.DecodeGatewayAdminSetReputationResponse(response)
	if err != nil {
		return false, err
	}
	if targetID.ToString() != clientID.ToString() {
		return false, errors.New("Wrong client id")
	}
	if !exists {
		return false, errors.New("Client id not existed")
	}
	if reputationNew != reputation {
		return false, errors.New("Reputation not correctly set")
	}

	return true, nil
}
