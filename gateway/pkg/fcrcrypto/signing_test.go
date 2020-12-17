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
    "log"
	"testing"

    "github.com/stretchr/testify/assert"
    "github.com/ConsenSys/fc-retrieval-gateway/pkg/messages"
)



func TestEstablishMessage(t *testing.T) {
    privateKey, err := GenKeyPair()
    if err != nil {
        panic(err)
    }

    log.Printf("Private Key: %s", EncodePrivateKey(privateKey))

    resp := messages.ClientEstablishmentResponse{}
    resp.MessageType = messages.ClientEstablishmentResponseType
    resp.ProtocolVersion = int32(1)
	resp.GatewayID = "1234567890abcdef01234567890abcdef01234567890abcdef01234567890abcdef0"
    resp.Challenge     = "a4b2345654665646461234567890abcdef01234567890abcdef01234567890abcdef"
    resp.Signature = ""

    publicKey := privateKey.PublicKey

    sigAlgorithm := SigAlg{Alg: SigAlgEcdsaP256Sha512_256}
    keyVersionInt := uint8(0x97)
    keyVersion := KeyVersion{Ver: keyVersionInt}

    signature, err2 := Sign(privateKey, keyVersion, sigAlgorithm, resp)
    if err2 != nil {
        panic(err2)
    }
    assert.NotEqual(t, "", *signature)

    resp.Signature = *signature

    // In the system, the message would be communicated to the entity receiving the data.


    foundKeyV, err1 := ExtractKeyVersion(&resp.Signature)
    if err1 != nil {
        panic(err1)
    }
    assert.Equal(t, keyVersionInt, foundKeyV.Ver)

    // In the system, the public key and signature algorithm would be fetched based on the key version.
    sigToBeVerified := resp.Signature
    resp.Signature = ""

    verified, err1 := Verify(&publicKey, sigAlgorithm, &sigToBeVerified, resp)
    if err1 != nil {
        panic(err1)
    }
    assert.True(t, verified, "Signature failed to verify")
}

