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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// CopiedClientEstablishmentResponse is a copy of a message used in the messages module. The
// struct has been copied here to remove the circular dependancy.
type CopiedClientEstablishmentResponse struct {
	messageType     int32  `json:"message_type"`
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

func (c CopiedClientEstablishmentResponse) MessageType() string {
	return fmt.Sprintf("  #messageType %v", c.messageType)
}

func TestGetToBeSigned(t *testing.T) {
	out := getToBeSigned(CopiedClientEstablishmentResponse{
		messageType:     CopiedClientEstablishmentResponseType,
		ProtocolVersion: 1,
		GatewayID:       "1234567890abcdef01234567890abcdef01234567890abcdef01234567890abcdef0",
		Challenge:       "a4b2345654665646461234567890abcdef01234567890abcdef01234567890abcdef",
	})
	assert.Equal(t, string(out), "11234567890abcdef01234567890abcdef01234567890abcdef01234567890abcdef0a4b2345654665646461234567890abcdef01234567890abcdef01234567890abcdef  #messageType 1")
}

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
	assert.Equal(t, "0000000149a4c1bab090ce563b003bc017cf255b89d26e1f7170da57e58e58afe51407aa51255338e77a828434cf97f71716b0e0c70a8f5c762fd64380916a8d98d0daf700", sig)
}

func TestSignNonEmptyMsg(t *testing.T) {
	keyPair, err := DecodePrivateKey(PrivKey)
	if err != nil {
		panic(err)
	}

	sig, err := SignMessage(keyPair, InitialKeyVersion(), CopiedClientEstablishmentResponse{
		messageType:     CopiedClientEstablishmentResponseType,
		ProtocolVersion: 1,
		GatewayID:       "1234567890abcdef01234567890abcdef01234567890abcdef01234567890abcdef0",
		Challenge:       "a4b2345654665646461234567890abcdef01234567890abcdef01234567890abcdef",
	})
	assert.Empty(t, err)
	assert.Equal(t, "00000001b29b643d232313afbbad00d6b10e23aa82e09b3183d619046de42cf56d9acc24411f8547fa761b416cc4804539ca859c3b4681b86cf0158a880668514855089000", sig)
}

func TestSignNonEmptyMsgPtr(t *testing.T) {
	keyPair, err := DecodePrivateKey(PrivKey)
	if err != nil {
		panic(err)
	}

	sig, err := SignMessage(keyPair, InitialKeyVersion(), &CopiedClientEstablishmentResponse{
		messageType:     CopiedClientEstablishmentResponseType,
		ProtocolVersion: 1,
		GatewayID:       "1234567890abcdef01234567890abcdef01234567890abcdef01234567890abcdef0",
		Challenge:       "a4b2345654665646461234567890abcdef01234567890abcdef01234567890abcdef",
	})
	assert.Empty(t, err)
	assert.Equal(t, "00000001b29b643d232313afbbad00d6b10e23aa82e09b3183d619046de42cf56d9acc24411f8547fa761b416cc4804539ca859c3b4681b86cf0158a880668514855089000", sig)
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
		messageType:     CopiedClientEstablishmentResponseType,
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
		messageType:     CopiedClientEstablishmentResponseType,
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
		messageType:     CopiedClientEstablishmentResponseType,
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
		messageType:     CopiedClientEstablishmentResponseType,
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
	res, err := VerifyMessage(keyPair, "00000001b29b643d232313afbbad00d6b10e73aa82e09b3183d619046de42cf56d9acc24411f8547fa761b416cc4804539ca859c3b4681b86cf0158a880668514855089000",
		CopiedClientEstablishmentResponse{
			messageType:     CopiedClientEstablishmentResponseType,
			ProtocolVersion: 1,
			GatewayID:       "1234567890abcdef01234567890abcdef01234567890abcdef01234567890abcdef0",
			Challenge:       "a4b2345654665646461234567890abcdef01234567890abcdef01234567890abcdef",
		})
	assert.Empty(t, err)
	assert.False(t, res)

	res, err = VerifyMessage(keyPair, "00000001b29b643d232313afbbad00d6b10e23aa82e09b3183d619046de42cf56d9acc24411f8547fa761b416cc4804539ca859c3b4681b86cf0158a880668514855089000",
		CopiedClientEstablishmentResponse{
			messageType:     CopiedClientEstablishmentResponseType,
			ProtocolVersion: 1,
			GatewayID:       "1234567890abcdef01234567890abcdef01234567890abcdef01234567890abcdef0",
			Challenge:       "a4b2345654665646461234567890abcdef01234567890abcdef01234567890abcdef",
		})
	assert.Empty(t, err)
	assert.True(t, res)
}

func TestSignMsgWithPtr(t *testing.T) {

	type (
		CopyNodeID struct {
			id []byte
		}
		CopyCidOfferSigning struct {
			NodeID     *CopyNodeID
			MerkleRoot string
			price      uint64
			expiry     int64
			qos        uint64
		}
	)

	keyPair, err := DecodePrivateKey(PrivKey)
	if err != nil {
		panic(err)
	}

	sig, err := SignMessage(keyPair, DecodeKeyVersion(uint32(0xffff)), CopyCidOfferSigning{
		NodeID: &CopyNodeID{
			id: []byte("dsafgasdgadfs"),
		},
		MerkleRoot: "asdfdasfas",
		expiry:     55,
		qos:        66,
	})

	assert.Empty(t, err)
	assert.Equal(t,
		"0000ffff26ed1e1a7cf86d7125983d1bed3a95ec23352f23e2cf8d07ae296e52255f3f1f76659c7a36b7e3e567592ebc7d731c66637a09d5d04ca13fcf22c4d5d18df37001",
		sig)
}
