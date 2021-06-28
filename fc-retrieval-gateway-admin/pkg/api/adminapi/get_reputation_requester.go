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
)

// RequestGetReputation gets the reputation for a given client id
func (a *Admin) RequestGetReputation(
  gatewayRegistrar register.GatewayRegistrar,
	clientID *nodeid.NodeID,
	signingPrivkey *fcrcrypto.KeyPair,
	signingPrivKeyVer *fcrcrypto.KeyVersion,
) (int64, error) {
	// First, Get pubkey
	pubKey, err := gatewayRegistrar.GetSigningKey()
	if err != nil {
		logging.Error("Error in obtaining signing key from register info.")
		return 0, err
	}

	// Second, send get reputation to given gateway
	request, err := fcrmessages.EncodeGatewayAdminGetReputationRequest(clientID)
	if err != nil {
		logging.Error("Error in encoding message.")
		return 0, err
	}
	// Sign the request
	if request.Sign(signingPrivkey, signingPrivKeyVer) != nil {
		return 0, errors.New("error in signing the request")
	}

	response, err := a.httpCommunicator.SendMessage(gatewayRegistrar.GetNetworkInfoAdmin(), request)
	if err != nil {
		logging.Error("Error in sending the message.")
		return 0, err
	}

	// Verify the response
	if response.Verify(pubKey) != nil {
		return 0, errors.New("fail to verify the response")
	}

	targetID, reputation, exists, err := fcrmessages.DecodeGatewayAdminGetReputationResponse(response)
	if err != nil {
		return 0, err
	}
	if targetID.ToString() != clientID.ToString() {
		return 0, errors.New("wrong client id")
	}
	if !exists {
		return 0, errors.New("client id not existed")
	}

	return reputation, nil
}
