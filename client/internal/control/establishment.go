package control

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
	"encoding/base64"
	"errors"

	"github.com/ConsenSys/fc-retrieval-client/internal/network"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-register/pkg/register"
)

// GatewayClientEstablishment sends a GatewayClientEstablishmentRequest and processes a response.
// Return true if the establishment request was successful.
func (c *ClientManager) GatewayClientEstablishment(nodeInfo *register.GatewayRegister, challenge []byte, clientID *nodeid.NodeID, ttl int64) error {
	if len(challenge) != 32 {
		return errors.New("Challenge is not 32 bytes")
	}
	b := make([]byte, base64.StdEncoding.EncodedLen(len(challenge)))
	base64.StdEncoding.Encode(b, challenge[:])

	request, err := fcrmessages.EncodeClientEstablishmentRequest(clientID, string(b), ttl)
	if err != nil {
		logging.Error("Error encoding Client Establishment Request: %+v", err)
		return err
	}

	response, err := network.SendMessage(nodeInfo.NetworkInfoClient, request)
	if err != nil {
		return err
	}

	// Get the gateway's public key
	pubKey, err := nodeInfo.GetSigningKey()
	if err != nil {
		return err
	}

	// Verify the response
	if response.Verify(pubKey) != nil {
		return errors.New("Fail to verify response")
	}
	// Finally check if gatewayID and received challenge matches.
	gatewayID, recvChallenge, err := fcrmessages.DecodeClientEstablishmentResponse(response)
	if err != nil {
		return err
	}

	if nodeInfo.NodeID != gatewayID.ToString() {
		return errors.New("Gateway ID not match")
	}
	if recvChallenge != string(b) {
		return errors.New("Challenge mismatch")
	}

	return nil
}
