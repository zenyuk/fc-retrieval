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

// CopiedClientEstablishmentResponse is a copy of a message used in the messages module. The
// struct has been copied here to remove the circular dependancy.
type CopiedClientEstablishmentResponse struct {
	messageType       int32
	ProtocolVersion   int32
	GatewayID         string
	Challenge         string
	Signature         string
	SliceIDs          []string
	ProtocolSupported []int32
	Msg               []byte
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

	raw, _ := json.Marshal(CopiedClientEstablishmentResponse{})
	sig, err := SignMessage(keyPair, InitialKeyVersion(), raw)
	assert.NotEmpty(t, err)
	assert.Empty(t, sig)
}

func TestSignEmptyMsg(t *testing.T) {
	keyPair, err := DecodePrivateKey(PrivKey)
	if err != nil {
		panic(err)
	}

	raw, _ := json.Marshal(CopiedClientEstablishmentResponse{})
	sig, err := SignMessage(keyPair, InitialKeyVersion(), raw)
	assert.Empty(t, err)
	assert.Equal(t, "00000001cb6b97a086eb94fccb5a3c7b37da13f78dd153e20bf10ebc74d0c3340b8e4ef03450fa264669c65f032b535295c9fccbb37334d3f7777b7977ee480bc99a3b7e00", sig)
}

func TestSignNonEmptyMsg(t *testing.T) {
	keyPair, err := DecodePrivateKey(PrivKey)
	if err != nil {
		panic(err)
	}

	raw, _ := json.Marshal(CopiedClientEstablishmentResponse{
		messageType:     CopiedClientEstablishmentResponseType,
		ProtocolVersion: 1,
		GatewayID:       "1234567890abcdef01234567890abcdef01234567890abcdef01234567890abcdef0",
		Challenge:       "a4b2345654665646461234567890abcdef01234567890abcdef01234567890abcdef",
	})
	sig, err := SignMessage(keyPair, InitialKeyVersion(), raw)
	assert.Empty(t, err)
	assert.Equal(t, "000000017989a43a3545120d9e9134b592653b582b230f80c6eb18aa18847e0ba47f7518261cac7232388ea18a977718e5f0f81e78c89fc60f7132bed878ac3fccf7063d00", sig)
}

func TestSignNonEmptyMsgPtr(t *testing.T) {
	keyPair, err := DecodePrivateKey(PrivKey)
	if err != nil {
		panic(err)
	}

	raw, _ := json.Marshal(&CopiedClientEstablishmentResponse{
		messageType:     CopiedClientEstablishmentResponseType,
		ProtocolVersion: 1,
		GatewayID:       "1234567890abcdef01234567890abcdef01234567890abcdef01234567890abcdef0",
		Challenge:       "a4b2345654665646461234567890abcdef01234567890abcdef01234567890abcdef",
	})
	sig, err := SignMessage(keyPair, InitialKeyVersion(), raw)
	assert.Empty(t, err)
	assert.Equal(t, "000000017989a43a3545120d9e9134b592653b582b230f80c6eb18aa18847e0ba47f7518261cac7232388ea18a977718e5f0f81e78c89fc60f7132bed878ac3fccf7063d00", sig)
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

	raw, _ := json.Marshal(&CopiedClientEstablishmentResponse{
		messageType:     CopiedClientEstablishmentResponseType,
		ProtocolVersion: 1,
		GatewayID:       "1234567890abcdef01234567890abcdef01234567890abcdef01234567890abcdef0",
		Challenge:       "a4b2345654665646461234567890abcdef01234567890abcdef01234567890abcdef",
	})
	sig, err := SignMessage(keyPair, InitialKeyVersion(), raw)
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

	raw, _ := json.Marshal(&CopiedClientEstablishmentResponse{
		messageType:     CopiedClientEstablishmentResponseType,
		ProtocolVersion: 1,
		GatewayID:       "1234567890abcdef01234567890abcdef01234567890abcdef01234567890abcdef0",
		Challenge:       "a4b2345654665646461234567890abcdef01234567890abcdef01234567890abcdef",
	})

	sig, err := SignMessage(keyPair, DecodeKeyVersion(uint32(0)), raw)
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
	raw, _ := json.Marshal(&CopiedClientEstablishmentResponse{
		messageType:     CopiedClientEstablishmentResponseType,
		ProtocolVersion: 1,
		GatewayID:       "1234567890abcdef01234567890abcdef01234567890abcdef01234567890abcdef0",
		Challenge:       "a4b2345654665646461234567890abcdef01234567890abcdef01234567890abcdef",
	})
	sig, err := SignMessage(keyPair, DecodeKeyVersion(uint32(0xff)), raw)
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

	raw, _ := json.Marshal(&CopiedClientEstablishmentResponse{
		messageType:     CopiedClientEstablishmentResponseType,
		ProtocolVersion: 1,
		GatewayID:       "1234567890abcdef01234567890abcdef01234567890abcdef01234567890abcdef0",
		Challenge:       "a4b2345654665646461234567890abcdef01234567890abcdef01234567890abcdef",
	})

	sig, err := SignMessage(keyPair, DecodeKeyVersion(uint32(0xffff)), raw)
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

	raw, _ := json.Marshal(CopiedClientEstablishmentResponse{})
	res, err := VerifyMessage(keyPair, "abcdefghijklmn", raw)
	assert.NotEmpty(t, err)
	assert.False(t, res)
}

func TestVerifyMsgEmptySig(t *testing.T) {
	keyPair, err := DecodePublicKey(PubKey)
	if err != nil {
		panic(err)
	}
	raw, _ := json.Marshal(CopiedClientEstablishmentResponse{})
	res, err := VerifyMessage(keyPair, "", raw)
	assert.NotEmpty(t, err)
	assert.False(t, res)
}

func TestVerifyMsgShortSig(t *testing.T) {
	keyPair, err := DecodePublicKey(PubKey)
	if err != nil {
		panic(err)
	}
	raw, _ := json.Marshal(CopiedClientEstablishmentResponse{})
	res, err := VerifyMessage(keyPair, "0x12", raw)
	assert.NotEmpty(t, err)
	assert.False(t, res)
}

func TestVerifyMsg(t *testing.T) {
	keyPair, err := DecodePublicKey(PubKey)
	if err != nil {
		panic(err)
	}
	raw, _ := json.Marshal(CopiedClientEstablishmentResponse{
		messageType:     CopiedClientEstablishmentResponseType,
		ProtocolVersion: 1,
		GatewayID:       "1234567890abcdef01234567890abcdef01234567890abcdef01234567890abcdef0",
		Challenge:       "a4b2345654665646461234567890abcdef01234567890abcdef01234567890abcdef",
	})
	res, err := VerifyMessage(keyPair, "00000001b29b643d232313afbbad00d6b10e73aa82e09b3183d619046de42cf56d9acc24411f8547fa761b416cc4804539ca859c3b4681b86cf0158a880668514855089000",
		raw)
	assert.Empty(t, err)
	assert.False(t, res)

	raw, _ = json.Marshal(CopiedClientEstablishmentResponse{
		messageType:     CopiedClientEstablishmentResponseType,
		ProtocolVersion: 1,
		GatewayID:       "1234567890abcdef01234567890abcdef01234567890abcdef01234567890abcdef0",
		Challenge:       "a4b2345654665646461234567890abcdef01234567890abcdef01234567890abcdef",
	})
	res, err = VerifyMessage(keyPair, "000000017989a43a3545120d9e9134b592653b582b230f80c6eb18aa18847e0ba47f7518261cac7232388ea18a977718e5f0f81e78c89fc60f7132bed878ac3fccf7063d00", raw)
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

	raw, _ := json.Marshal(CopyCidOfferSigning{
		NodeID: &CopyNodeID{
			id: []byte("dsafgasdgadfs"),
		},
		MerkleRoot: "asdfdasfas",
		expiry:     55,
		qos:        66,
	})

	sig, err := SignMessage(keyPair, DecodeKeyVersion(uint32(0xffff)), raw)

	assert.Empty(t, err)
	assert.Equal(t,
		"0000ffff6529cc4a868697c4d6f1332e40ed317d8351167c6d401f63d75d444e3dfb35c47403d069868501dfc629da671aa38d34976413334d2879b4c3bffebd32e60ea900",
		sig)
}
