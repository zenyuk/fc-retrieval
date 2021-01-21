package fcrcrypto

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
	"testing"

    "github.com/stretchr/testify/assert"
    "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
)


// CopiedClientEstablishmentResponse is a copy of a message used in the messages module. The
// struct has been copied here to remove the circular dependancy.
type CopiedClientEstablishmentResponse struct {
	MessageType     int32  `json:"message_type"`
	ProtocolVersion int32  `json:"protocol_version"`
	GatewayID       string `json:"gateway_id"`
	Challenge       string `json:"challenge"`
	Signature       string `json:"signature"`
}

const (
	CopiedClientEstablishmentResponseType             = 1
)



func TestEstablishMessage(t *testing.T) {
    keyPair, err := GenerateRetrievalV1KeyPair()
    if err != nil {
        panic(err)
    }

    logging.Debug("Private Key: %s", keyPair.EncodePrivateKey())

    resp := CopiedClientEstablishmentResponse{}
    resp.MessageType = CopiedClientEstablishmentResponseType
    resp.ProtocolVersion = int32(1)
	resp.GatewayID = "1234567890abcdef01234567890abcdef01234567890abcdef01234567890abcdef0"
    resp.Challenge     = "a4b2345654665646461234567890abcdef01234567890abcdef01234567890abcdef"
    resp.Signature = ""


    keyVersionInt := uint32(0x97)
    keyVersion := DecodeKeyVersion(keyVersionInt)

    signature, err := SignMessage(keyPair, keyVersion, resp)
    if err != nil {
        panic(err)
    }
    assert.NotEqual(t, "", signature)

    resp.Signature = signature

    // In the system, the message would be communicated to the entity receiving the data.


    foundKeyV, err := ExtractKeyVersionFromMessage(resp.Signature)
    if err != nil {
        panic(err)
    }
    assert.True(t, foundKeyV.Equals(keyVersion))

    // In the system, the public key and signature algorithm would be fetched based on the key version.
    sigToBeVerified := resp.Signature
    resp.Signature = ""

    verified, err := VerifyMessage(keyPair, sigToBeVerified, resp)
    if err != nil {
        panic(err)
    }
    assert.True(t, verified, "Signature failed to verify")
}

