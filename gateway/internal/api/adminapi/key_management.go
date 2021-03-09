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
	"net"
	"sync"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrtcpcomms"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
)

func handleAdminAcceptKeysChallenge(conn net.Conn, request *fcrmessages.FCRMessage, wg *sync.WaitGroup) error {
	logging.Info("In handleAdminAcceptKeysChallenge")

	logging.Info("Message received \n%s", request.DumpMessage())

	// Get the core structure
	g := gateway.GetSingleInstance()

	nodeID, encprivatekey, encprivatekeyversion, err := fcrmessages.DecodeAdminAcceptKeyChallenge(request)
	if err != nil {
		return err
	}

	// Decode private key from hex string to *fcrCrypto.KeyPair
	privatekey, err := fcrcrypto.DecodePrivateKey(encprivatekey)
	if err != nil {
		return err
	}
	// Decode from int32 to *fcrCrypto.KeyVersion
	privatekeyversion := fcrcrypto.DecodeKeyVersion(encprivatekeyversion)

	g.GatewayID = nodeID
	g.GatewayPrivateKey = privatekey
	g.GatewayPrivateKeyVersion = privatekeyversion
	wg.Done() // need mutex to protect g

	// Construct messaqe
	exists := true
	response, err := fcrmessages.EncodeAdminAcceptKeyResponse(exists)
	if err != nil {
		return err
	}

	response.SignMessage(func(msg interface{}) (string, error) {
		return fcrcrypto.SignMessage(g.GatewayPrivateKey, g.GatewayPrivateKeyVersion, msg)
	})

	logging.Info("Admin action: Key installation complete")
	// Send message
	return fcrtcpcomms.SendTCPMessage(conn, response, settings.DefaultTCPInactivityTimeout)
}
