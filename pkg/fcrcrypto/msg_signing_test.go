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
	CopiedClientEstablishmentResponseType = 1
	PrivKey                               = "015ed053eab6fdf18c03954373ff7f89089992017d56beb8b05305b19800d6afe0"
	PubKey                                = "01047799f37b014564e23578447d718e5c70a786b0e4e58ca25cb2a086b822434594d910b9b8c0fcbfe9f4c2db321e874819e0614be5b57fbb5080accd69adb2eaad"
)

func TestSignMsgWithError(t *testing.T) {
	keyPair, err := DecodePrivateKey(PrivKey)
	if err != nil {
		panic(err)
	}
	// Change the key algorithm to unsupported
	keyPair.alg.algorithm = SigAlgEcdsaSecP256K1Blake2b + 1

	sig, err := SignMessage(keyPair, InitialKeyVersion(), CopiedClientEstablishmentResponse{})
	assert.NotEmpty(t, err)
	assert.Empty(t, sig)
}

func TestSignEmptyMsg(t *testing.T) {
	keyPair, err := DecodePrivateKey(PrivKey)
	if err != nil {
		panic(err)
	}

	sig, err := SignMessage(keyPair, InitialKeyVersion(), CopiedClientEstablishmentResponse{})
	assert.Empty(t, err)
	assert.Equal(t, "00000001e79fdaa275888de3b3171ddf219d61cd19df42f1fa942d88150722438efcdcaf55195e10fd8267536eb9e807061460bec41d8b04a88cf53a68a0c17f383c625801", sig)
}

func TestSignNonEmptyMsg(t *testing.T) {
	keyPair, err := DecodePrivateKey(PrivKey)
	if err != nil {
		panic(err)
	}

	sig, err := SignMessage(keyPair, InitialKeyVersion(), CopiedClientEstablishmentResponse{
		MessageType:     CopiedClientEstablishmentResponseType,
		ProtocolVersion: 1,
		GatewayID:       "1234567890abcdef01234567890abcdef01234567890abcdef01234567890abcdef0",
		Challenge:       "a4b2345654665646461234567890abcdef01234567890abcdef01234567890abcdef",
	})
	assert.Empty(t, err)
	assert.Equal(t, "00000001f6ef07bbb3522b0d087bdd5ff4aff73d5f55ed4febd45c02b8f375aa0082a9f85a2e152c03786842a9ee04d4b4d56d4d0a33ffac14f4601d65dd755dfa6ffafe01", sig)
}

func TestSignNonEmptyMsgPtr(t *testing.T) {
	keyPair, err := DecodePrivateKey(PrivKey)
	if err != nil {
		panic(err)
	}

	sig, err := SignMessage(keyPair, InitialKeyVersion(), &CopiedClientEstablishmentResponse{
		MessageType:     CopiedClientEstablishmentResponseType,
		ProtocolVersion: 1,
		GatewayID:       "1234567890abcdef01234567890abcdef01234567890abcdef01234567890abcdef0",
		Challenge:       "a4b2345654665646461234567890abcdef01234567890abcdef01234567890abcdef",
	})
	assert.Empty(t, err)
	assert.Equal(t, "00000001f6ef07bbb3522b0d087bdd5ff4aff73d5f55ed4febd45c02b8f375aa0082a9f85a2e152c03786842a9ee04d4b4d56d4d0a33ffac14f4601d65dd755dfa6ffafe01", sig)
}

func TestExtractKeyVersionFromMsgWithError(t *testing.T) {
	keyVer, err := ExtractKeyVersionFromMessage("abcdefghijklmn")
	assert.NotEmpty(t, err)
	assert.Empty(t, keyVer)
}

func TestExtractKeyVersionFromEmptyMsg(t *testing.T) {
	keyVer, err := ExtractKeyVersionFromMessage("")
	assert.Empty(t, err)
	assert.Equal(t, uint32(0), keyVer.ver)
}

func TestExtractKeyVersionFromMsgInitialVer(t *testing.T) {
	keyPair, err := DecodePrivateKey(PrivKey)
	if err != nil {
		panic(err)
	}

	sig, err := SignMessage(keyPair, InitialKeyVersion(), &CopiedClientEstablishmentResponse{
		MessageType:     CopiedClientEstablishmentResponseType,
		ProtocolVersion: 1,
		GatewayID:       "1234567890abcdef01234567890abcdef01234567890abcdef01234567890abcdef0",
		Challenge:       "a4b2345654665646461234567890abcdef01234567890abcdef01234567890abcdef",
	})
	if err != nil {
		panic(err)
	}

	keyVer, err := ExtractKeyVersionFromMessage(sig)
	assert.Empty(t, err)
	assert.Equal(t, InitialKeyVersion().ver, keyVer.ver)
}

func TestExtractKeyVersionFromMsgZeroVer(t *testing.T) {
	keyPair, err := DecodePrivateKey(PrivKey)
	if err != nil {
		panic(err)
	}

	sig, err := SignMessage(keyPair, DecodeKeyVersion(uint32(0)), &CopiedClientEstablishmentResponse{
		MessageType:     CopiedClientEstablishmentResponseType,
		ProtocolVersion: 1,
		GatewayID:       "1234567890abcdef01234567890abcdef01234567890abcdef01234567890abcdef0",
		Challenge:       "a4b2345654665646461234567890abcdef01234567890abcdef01234567890abcdef",
	})
	if err != nil {
		panic(err)
	}

	keyVer, err := ExtractKeyVersionFromMessage(sig)
	assert.Equal(t, uint32(0), keyVer.ver)
	assert.Empty(t, err)
}

func TestExtractKeyVersionFromMsgVer(t *testing.T) {
	keyPair, err := DecodePrivateKey(PrivKey)
	if err != nil {
		panic(err)
	}

	sig, err := SignMessage(keyPair, DecodeKeyVersion(uint32(0xff)), &CopiedClientEstablishmentResponse{
		MessageType:     CopiedClientEstablishmentResponseType,
		ProtocolVersion: 1,
		GatewayID:       "1234567890abcdef01234567890abcdef01234567890abcdef01234567890abcdef0",
		Challenge:       "a4b2345654665646461234567890abcdef01234567890abcdef01234567890abcdef",
	})
	if err != nil {
		panic(err)
	}

	keyVer, err := ExtractKeyVersionFromMessage(sig)
	assert.Empty(t, err)
	assert.Equal(t, uint32(0xff), keyVer.ver)
}

func TestExtractKeyVersionFromMsgMaxVer(t *testing.T) {
	keyPair, err := DecodePrivateKey(PrivKey)
	if err != nil {
		panic(err)
	}

	sig, err := SignMessage(keyPair, DecodeKeyVersion(uint32(0xffff)), &CopiedClientEstablishmentResponse{
		MessageType:     CopiedClientEstablishmentResponseType,
		ProtocolVersion: 1,
		GatewayID:       "1234567890abcdef01234567890abcdef01234567890abcdef01234567890abcdef0",
		Challenge:       "a4b2345654665646461234567890abcdef01234567890abcdef01234567890abcdef",
	})
	if err != nil {
		panic(err)
	}

	keyVer, err := ExtractKeyVersionFromMessage(sig)
	assert.Empty(t, err)
	assert.Equal(t, uint32(0xffff), keyVer.ver)
}

func TestVerifyMsgWithError(t *testing.T) {
	keyPair, err := DecodePublicKey(PubKey)
	if err != nil {
		panic(err)
	}
	res, err := VerifyMessage(keyPair, "abcdefghijklmn", CopiedClientEstablishmentResponse{})
	assert.NotEmpty(t, err)
	assert.False(t, res)
}

func TestVerifyMsgEmptySig(t *testing.T) {
	keyPair, err := DecodePublicKey(PubKey)
	if err != nil {
		panic(err)
	}
	res, err := VerifyMessage(keyPair, "", CopiedClientEstablishmentResponse{})
	assert.NotEmpty(t, err)
	assert.False(t, res)
}

func TestVerifyMsgShortSig(t *testing.T) {
	keyPair, err := DecodePublicKey(PubKey)
	if err != nil {
		panic(err)
	}
	res, err := VerifyMessage(keyPair, "0x12", CopiedClientEstablishmentResponse{})
	assert.NotEmpty(t, err)
	assert.False(t, res)
}

func TestVerifyMsg(t *testing.T) {
	keyPair, err := DecodePublicKey(PubKey)
	if err != nil {
		panic(err)
	}
	res, err := VerifyMessage(keyPair, "00000001f6ef07bbb3522b0d087bdd5aa4aff73d5f55ed4febd45c02b8f375aa0082a9f85a2e152c03786842a9ee04d4b4d56d4d0a33ffac14f4601d65dd755dfa6ffafe01", CopiedClientEstablishmentResponse{
		MessageType:     CopiedClientEstablishmentResponseType,
		ProtocolVersion: 1,
		GatewayID:       "1234567890abcdef01234567890abcdef01234567890abcdef01234567890abcdef0",
		Challenge:       "a4b2345654665646461234567890abcdef01234567890abcdef01234567890abcdef",
	})
	assert.Empty(t, err)
	assert.False(t, res)

	res, err = VerifyMessage(keyPair, "00000001f6ef07bbb3522b0d087bdd5ff4aff73d5f55ed4febd45c02b8f375aa0082a9f85a2e152c03786842a9ee04d4b4d56d4d0a33ffac14f4601d65dd755dfa6ffafe01", CopiedClientEstablishmentResponse{
		MessageType:     CopiedClientEstablishmentResponseType,
		ProtocolVersion: 1,
		GatewayID:       "1234567890abcdef01234567890abcdef01234567890abcdef01234567890abcdef0",
		Challenge:       "a4b2345654665646461234567890abcdef01234567890abcdef01234567890abcdef",
	})
	assert.Empty(t, err)
	assert.True(t, res)
}
