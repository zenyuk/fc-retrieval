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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyEncodeDecode(t *testing.T) {
	keyPair, err := GenerateRetrievalV1KeyPair()
	if err != nil {
		panic(err)
	}

	pKeyStr := keyPair.EncodePrivateKey()
	pKey, err := DecodePrivateKey(pKeyStr)
	if err != nil {
		panic(err)
	}
	pKeyStr1 := pKey.EncodePrivateKey()
	assert.Equal(t, pKeyStr, pKeyStr1, "Private Key round trip not working")

	pubKeyStr, err := keyPair.EncodePublicKey()
	if err != nil {
		panic(err)
	}
	pubKey, err := DecodePublicKey(pubKeyStr)
	if err != nil {
		panic(err)
	}
	pubKeyStr1, err := pubKey.EncodePublicKey()
	if err != nil {
		panic(err)
	}
	assert.Equal(t, pubKeyStr, pubKeyStr1, "Public Key round trip not working")
}

func TestKeyRawEncodeDecode(t *testing.T) {
	keyPair, err := GenerateRetrievalV1KeyPair()
	if err != nil {
		panic(err)
	}

	pKeyStr, err := keyPair.EncodeRawPrivateKey()
	if err != nil {
		panic(err)
	}
	pKey, err := DecodeRawPrivateKey(pKeyStr)
	if err != nil {
		panic(err)
	}
	pKeyStr1, err := pKey.EncodeRawPrivateKey()
	if err != nil {
		panic(err)
	}
	assert.Equal(t, pKeyStr, pKeyStr1, "Private Key round trip not working")
}

func TestSignVerify(t *testing.T) {
	keyPair, err := GenerateRetrievalV1KeyPair()
	if err != nil {
		panic(err)
	}

	tbs := make([]byte, 100)
	signature, err := keyPair.Sign(tbs)
	if err != nil {
		panic(err)
	}
	verified, err := keyPair.Verify(signature, tbs)
	if err != nil {
		panic(err)
	}
	assert.True(t, verified)

	tbs[0] = 1
	verified, err = keyPair.Verify(signature, tbs)
	if err != nil {
		panic(err)
	}
	assert.False(t, verified)
}

func TestExReceoverVerify(t *testing.T) {
	keyPair, err := GenerateRetrievalV1KeyPair()
	if err != nil {
		panic(err)
	}
	hashOfPublicKey, err := keyPair.HashPublicKey()
	if err != nil {
		panic(err)
	}

	tbs := make([]byte, 100)
	signature, err := keyPair.Sign(tbs)
	if err != nil {
		panic(err)
	}

	verified, err := RetrievalV1Verify(signature, tbs, hashOfPublicKey)
	if err != nil {
		panic(err)
	}
	assert.True(t, verified)

	tbs[0] = 1
	verified, err = RetrievalV1Verify(signature, tbs, hashOfPublicKey)
	if err != nil {
		panic(err)
	}
	assert.False(t, verified)
}

func TestJson(t *testing.T) {
	key, err := GenerateRetrievalV1KeyPair()
	if err != nil {
		panic(err)
	}
	data, err := json.Marshal(key)
	if err != nil {
		panic(err)
	}
	var pubKey KeyPair
	err = json.Unmarshal(data, &pubKey)
	if err != nil {
		panic(err)
	}
	key1, err := key.EncodePublicKey()
	if err != nil {
		panic(err)
	}
	key2, err := pubKey.EncodePublicKey()
	if err != nil {
		panic(err)
	}
	assert.Equal(t, key1, key2)
}

func TestEncodeRawPrivKeyWithErr(t *testing.T) {
	keyPair, err := GenerateRetrievalV1KeyPair()
	if err != nil {
		panic(err)
	}
	keyPair.alg.algorithm = SigAlgEcdsaSecP256K1Blake2b + 1

	str, err := keyPair.EncodeRawPrivateKey()
	assert.NotEmpty(t, err)
	assert.Empty(t, str)
}

func TestDecodePrivKeyWithError(t *testing.T) {
	res, err := DecodePrivateKey("abcdefghijklmn")
	assert.NotEmpty(t, err)
	assert.Empty(t, res)
}

func TestDecodePrivKeyWithUnknownVersion(t *testing.T) {
	keyPair, err := GenerateRetrievalV1KeyPair()
	if err != nil {
		panic(err)
	}
	keyPair.alg.algorithm = SigAlgEcdsaSecP256K1Blake2b + 1
	keyPairStr := keyPair.EncodePrivateKey()
	res, err := DecodePrivateKey(keyPairStr)
	assert.NotEmpty(t, err)
	assert.Empty(t, res)
}

func TestDecodeRawPrivKeyWithError(t *testing.T) {
	res, err := DecodeRawPrivateKey("abcdefghijklmn")
	assert.NotEmpty(t, err)
	assert.Empty(t, res)
}

func TestDecodeSECP256K1WithError(t *testing.T) {
	keyBytes := []byte{0x12, 0x13}
	res, err := decodeSecP256K1PrivateKey(keyBytes)
	assert.NotEmpty(t, err)
	assert.Empty(t, res)
}

func TestEncodePubKeyWithErr(t *testing.T) {
	keyPair, err := GenerateRetrievalV1KeyPair()
	if err != nil {
		panic(err)
	}
	keyPair.alg.algorithm = SigAlgEcdsaSecP256K1Blake2b + 1

	str, err := keyPair.EncodePublicKey()
	assert.NotEmpty(t, err)
	assert.Empty(t, str)
}

func TestDecodePubKeyWithError(t *testing.T) {
	res, err := DecodePublicKey("abcdefghijklmn")
	assert.NotEmpty(t, err)
	assert.Empty(t, res)
}

func TestDecodePubSECP256K1WithError(t *testing.T) {
	keyBytes := []byte{0x12, 0x13}
	res, err := decodeSecP256K1PublicKey(keyBytes)
	assert.NotEmpty(t, err)
	assert.Empty(t, res)
}
