package cid

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
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContent(t *testing.T) {
	cid, err := NewContentID(big.NewInt(5))
	assert.Empty(t, err)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000005", cid.ToString())
}

func TestNewContentIDFromBytes(t *testing.T) {
	cid, err := NewContentIDFromBytes([]byte{1})
	assert.Empty(t, err)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000001", cid.ToString())
}

func TestNewContentIDFromEmptyBytes(t *testing.T) {
	cid, err := NewContentIDFromBytes([]byte{})
	assert.Empty(t, err)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", cid.ToString())
}

func TestNewContentIDFromMaxBytes(t *testing.T) {
	cidBytes := make([]byte, 32)
	for i := 0; i < 32; i++ {
		cidBytes[i] = 0xff
	}
	cid, err := NewContentIDFromBytes(cidBytes)
	assert.Empty(t, err)
	assert.Equal(t, "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", cid.ToString())
}

func TestNewContentIDFromBytesWithError(t *testing.T) {
	longBytes := make([]byte, 33)
	cid, err := NewContentIDFromBytes(longBytes)
	assert.NotEmpty(t, err)
	assert.Empty(t, cid)
}

func TestNewContentIDFromString(t *testing.T) {
	cid, err := NewContentIDFromHexString("10")
	assert.Empty(t, err)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000010", cid.ToString())
}

func TestNewContentIDFromEmptyString(t *testing.T) {
	cid, err := NewContentIDFromHexString("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	assert.Empty(t, err)
	assert.Equal(t, "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", cid.ToString())
}

func TestNewContentIDFromMaxString(t *testing.T) {
	cid, err := NewContentIDFromHexString("")
	assert.Empty(t, err)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", cid.ToString())
}

func TestNewContentIDFromStringExactLength(t *testing.T) {
	cid, err := NewContentIDFromHexString("1010101010101010101010101010101010101010101010101010101010101010")
	assert.Empty(t, err)
	assert.Equal(t, "1010101010101010101010101010101010101010101010101010101010101010", cid.ToString())
}

func TestNewContentIDFromStringLongLength(t *testing.T) {
	cid, err := NewContentIDFromHexString("101010101010101010101010101010101010101010101010101010101010101010")
	assert.NotEmpty(t, err)
	assert.Empty(t, cid)
}

func TestNewContentIDFromInvalidString(t *testing.T) {
	cid, err := NewContentIDFromHexString("abcdefghijkl")
	assert.NotEmpty(t, err)
	assert.Empty(t, cid)
}

func TestRandomCID(t *testing.T) {
	cid := NewRandomContentID()
	assert.NotEmpty(t, cid)
}

func TestToString(t *testing.T) {
	cid, err := NewContentIDFromHexString("10")
	assert.Empty(t, err)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000010", cid.ToString())
}

func TestToStringEmpty(t *testing.T) {
	cid, err := NewContentIDFromBytes([]byte{})
	assert.Empty(t, err)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", cid.ToString())
}

func TestToBytes(t *testing.T) {
	cid, err := NewContentIDFromHexString("10")
	assert.Empty(t, err)
	res := make([]byte, 32)
	res[31] = 0x10
	assert.Equal(t, res, cid.ToBytes())
}

func TestJSON(t *testing.T) {
	cid1, err := NewContentIDFromHexString("10")
	assert.Empty(t, err)
	p, err := cid1.MarshalJSON()
	assert.Empty(t, err)
	cid2 := ContentID{}
	err = cid2.UnmarshalJSON(p)
	assert.Empty(t, err)
	assert.Equal(t, cid1.ToBytes(), cid2.ToBytes())
}

func TestJSONWithError(t *testing.T) {
	cid := ContentID{}
	p := make([]byte, 32)
	err := cid.UnmarshalJSON(p)
	assert.NotEmpty(t, err)
}

func TestJSONWithWrongLength(t *testing.T) {
	cid := ContentID{}
	p0 := make([]byte, 30)
	p1, err := json.Marshal(p0)
	assert.Empty(t, err)
	err = cid.UnmarshalJSON(p1)
	assert.NotEmpty(t, err)
}

func TestCalculateHash(t *testing.T) {
	cid, err := NewContentIDFromHexString("01")
	assert.Empty(t, err)
	hash, err := cid.CalculateHash()
	assert.Empty(t, err)
	assert.NotEmpty(t, hash)
}

func TestEqual(t *testing.T) {
	cid1, err := NewContentIDFromHexString("01")
	assert.Empty(t, err)
	cid2, err := NewContentIDFromBytes([]byte{0x01})
	assert.Empty(t, err)
	cid3, err := NewContentIDFromBytes([]byte{0x02})
	assert.Empty(t, err)
	res1, err := cid1.Equals(cid2)
	assert.Empty(t, err)
	assert.True(t, res1)
	res2, err := cid1.Equals(cid3)
	assert.Empty(t, err)
	assert.False(t, res2)
}
