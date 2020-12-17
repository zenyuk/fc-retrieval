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
)



func TestKeyEncodeDecode(t *testing.T) {
    privateKey, err := GenKeyPair()
    if err != nil {
        panic(err)
    }

    pKeyStr := EncodePrivateKey(privateKey)
    pKey := DecodePrivateKey(pKeyStr)
    pKeyStr1 := EncodePrivateKey(pKey)
    assert.Equal(t, pKeyStr, pKeyStr1, "Private Key round trip not working")

    pubKeyStr := EncodePublicKey(&privateKey.PublicKey)
    pubKey := DecodePublicKey(pubKeyStr)
    pubKeyStr1 := EncodePublicKey(pubKey)
    assert.Equal(t, pubKeyStr, pubKeyStr1, "Public Key round trip not working")
}

