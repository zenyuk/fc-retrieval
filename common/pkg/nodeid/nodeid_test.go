package nodeid

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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/stretchr/testify/assert"
)

func TestNewNode(t *testing.T) {
	id, err := NewNodeID(big.NewInt(5))
	assert.Empty(t, err)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000005", id.ToString())
}

func TestNewNodeTooLong(t *testing.T) {
	idBytes := make([]byte, 33)
	for i := 0; i < 33; i++ {
		idBytes[i] = 0xff
	}
	_, err := NewNodeID(big.NewInt(0).SetBytes(idBytes))
	assert.NotEmpty(t, err)
}

func TestNewNodeEmpty(t *testing.T) {
	id, err := NewNodeIDFromBytes([]byte{})
	assert.Empty(t, err)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", id.ToString())
}

func TestNewNodeNil(t *testing.T) {
	id := &NodeID{}
	assert.Equal(t, "00", id.ToString())
}

func TestNewNodeIDTooLong(t *testing.T) {
	idBytes := make([]byte, 33)
	for i := 0; i < 33; i++ {
		idBytes[i] = 0xff
	}
	_, err := NewNodeIDFromBytes(idBytes)
	assert.NotEmpty(t, err)
}

func TestNewNodeIDFromBytes(t *testing.T) {
	id, err := NewNodeIDFromBytes([]byte{1})
	assert.Empty(t, err)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000001", id.ToString())
}

func TestNewNodeIDFromEmptyBytes(t *testing.T) {
	id, err := NewNodeIDFromBytes([]byte{})
	assert.Empty(t, err)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", id.ToString())
}

func TestNewNodeIDFromMaxBytes(t *testing.T) {
	idBytes := make([]byte, 32)
	for i := 0; i < 32; i++ {
		idBytes[i] = 0xff
	}
	id, err := NewNodeIDFromBytes(idBytes)
	assert.Empty(t, err)
	assert.Equal(t, "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", id.ToString())
}

func TestNewNodeIDFromBytesWithError(t *testing.T) {
	longBytes := make([]byte, 33)
	id, err := NewNodeIDFromBytes(longBytes)
	assert.NotEmpty(t, err)
	assert.Empty(t, id)
}

func TestNewNodeIDFromString(t *testing.T) {
	id, err := NewNodeIDFromHexString("10")
	assert.Empty(t, err)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000010", id.ToString())
}

func TestNewNodeIDFromEmptyString(t *testing.T) {
	id, err := NewNodeIDFromHexString("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	assert.Empty(t, err)
	assert.Equal(t, "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", id.ToString())
}

func TestNewNodeIDFromMaxString(t *testing.T) {
	id, err := NewNodeIDFromHexString("")
	assert.Empty(t, err)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", id.ToString())
}

func TestNewNodeIDFromStringExactLength(t *testing.T) {
	id, err := NewNodeIDFromHexString("1010101010101010101010101010101010101010101010101010101010101010")
	assert.Empty(t, err)
	assert.Equal(t, "1010101010101010101010101010101010101010101010101010101010101010", id.ToString())
}

func TestNewNodeIDFromStringLongLength(t *testing.T) {
	id, err := NewNodeIDFromHexString("101010101010101010101010101010101010101010101010101010101010101010")
	assert.NotEmpty(t, err)
	assert.Empty(t, id)
}

func TestNewNodeIDFromInvalidString(t *testing.T) {
	id, err := NewNodeIDFromHexString("abcdefghijkl")
	assert.NotEmpty(t, err)
	assert.Empty(t, id)
}

func TestNewNodeIDFromPublicKey(t *testing.T) {
	key, err := fcrcrypto.GenerateRetrievalV1KeyPair()
	assert.Empty(t, err)
	id, err := NewNodeIDFromPublicKey(key)
	assert.Empty(t, err)
	assert.NotEmpty(t, id)
}

func TestRandomNodeID(t *testing.T) {
	id := NewRandomNodeID()
	assert.NotEmpty(t, id)
}

func TestToBytes(t *testing.T) {
	id, err := NewNodeIDFromHexString("10")
	assert.Empty(t, err)
	res := make([]byte, 32)
	res[31] = 0x10
	assert.Equal(t, res, id.ToBytes())
}

func TestJSON(t *testing.T) {
	id1, err := NewNodeIDFromHexString("10")
	assert.Empty(t, err)
	p, err := id1.MarshalJSON()
	assert.Empty(t, err)
	id2 := NodeID{}
	err = id2.UnmarshalJSON(p)
	assert.Empty(t, err)
	assert.Equal(t, id1.ToBytes(), id2.ToBytes())
}

func TestJSONWithError(t *testing.T) {
	id := NodeID{}
	p := make([]byte, 32)
	err := id.UnmarshalJSON(p)
	assert.NotEmpty(t, err)
}

func TestJSONWithWrongLength(t *testing.T) {
	id := NodeID{}
	p0 := make([]byte, 30)
	p1, err := json.Marshal(p0)
	assert.Empty(t, err)
	err = id.UnmarshalJSON(p1)
	assert.NotEmpty(t, err)
}

func TestAsBytes32(t *testing.T) {
	id, err := NewNodeIDFromBytes([]byte{1})
	assert.Empty(t, err)
	var res [32]byte
	res[31] = 1
	assert.Equal(t, res, id.AsBytes32())
}

func TestRoundTripBigInt(t *testing.T) {
	testRoundTripBigInt(t, "00")
	testRoundTripBigInt(t, "01")
	testRoundTripBigInt(t, "fe")
	testRoundTripBigInt(t, "0100")
	testRoundTripBigInt(t, "30010203040506070809")
	testRoundTripBigInt(t, "80010203040506070809")
}

func TestRoundTripFromBytes(t *testing.T) {
	testRoundTripFromBytes(t, "00")
	testRoundTripFromBytes(t, "01")
	testRoundTripFromBytes(t, "fe")
	testRoundTripFromBytes(t, "0100")
	testRoundTripFromBytes(t, "30010203040506070809")
	testRoundTripFromBytes(t, "80010203040506070809")
}

func TestRoundTripFromString(t *testing.T) {
	testRoundTripFromString(t, "00")
	testRoundTripFromString(t, "01")
	testRoundTripFromString(t, "fe")
	testRoundTripFromString(t, "0100")
	testRoundTripFromString(t, "30010203040506070809")
	testRoundTripFromString(t, "80010203040506070809")
}

func testRoundTripBigInt(t *testing.T, value string) {
	id := new(big.Int)
	_, ok := id.SetString(value, 16)
	if !ok {
		panic("Number format issue: " + value)
	}

	nodeID, err := NewNodeID(id)
	if err != nil {
		panic(err)
	}
	idStr := nodeID.ToString()
	assert.Equal(t, fmt.Sprintf("%064s", value), idStr, "NewNodeID failed")

}

func testRoundTripFromBytes(t *testing.T, value string) {
	bytes, err := hex.DecodeString(value)
	if err != nil {
		panic(err)
	}

	nodeID, err := NewNodeIDFromBytes(bytes)
	if err != nil {
		panic(err)
	}
	idStr := nodeID.ToString()
	assert.Equal(t, fmt.Sprintf("%064s", value), idStr, "NewNodeIDFromString failed")

}

func testRoundTripFromString(t *testing.T, value string) {
	nodeID, err := NewNodeIDFromHexString(value)
	if err != nil {
		panic(err)
	}
	idStr := nodeID.ToString()
	assert.Equal(t, fmt.Sprintf("%064s", value), idStr, "NewNodeIDFromString failed")

}

// TestSortClockwiseNodeID
func TestSortClockwiseNodeID(t *testing.T) {
	nodeID, _ := NewNodeIDFromHexString("03")
	nodeID00, _ := NewNodeIDFromHexString("00")
	nodeID01, _ := NewNodeIDFromHexString("01")
	nodeID02, _ := NewNodeIDFromHexString("02")
	nodeID5A, _ := NewNodeIDFromHexString("5A")
	nodeIDFFFF, _ := NewNodeIDFromHexString("FFFF")

	ids := []*NodeID{nodeIDFFFF, nodeID5A, nodeID01, nodeID00, nodeID02}
	ids = SortClockwise(nodeID, ids)

	assert.Equal(t, nodeID5A, ids[0])
	assert.Equal(t, nodeIDFFFF, ids[1])
	assert.Equal(t, nodeID00, ids[2])
	assert.Equal(t, nodeID01, ids[3])
	assert.Equal(t, nodeID02, ids[4])
}

// TestSortClockwiseNodeID
func TestSortClockwiseNodeIDOneElement(t *testing.T) {
	nodeID, _ := NewNodeIDFromHexString("03")
	nodeID00, _ := NewNodeIDFromHexString("00")

	ids := []*NodeID{nodeID00}
	ids = SortClockwise(nodeID, ids)

	assert.Equal(t, nodeID00, ids[0])
}

func TestSortClockwiseNodeIDEmptyList(t *testing.T) {
	nodeID, _ := NewNodeIDFromHexString("03")

	ids := []*NodeID{}
	ids = SortClockwise(nodeID, ids)

	assert.Equal(t, 0, len(ids))
}

// TestSortNodeID
func TestSortNodeID(t *testing.T) {
	nodeID, _ := NewNodeIDFromHexString("01")
	nodeID00, _ := NewNodeIDFromHexString("00")
	nodeID01, _ := NewNodeIDFromHexString("01")
	nodeID02, _ := NewNodeIDFromHexString("02")
	nodeID5A, _ := NewNodeIDFromHexString("5A")
	nodeIDFFFF, _ := NewNodeIDFromHexString("FFFF")

	ids := []*NodeID{nodeIDFFFF, nodeID5A, nodeID01, nodeID00, nodeID02, nodeID}
	sortByteArrays(ids)

	assert.ElementsMatch(t, []*NodeID{nodeID00, nodeID, nodeID01, nodeID02, nodeID5A, nodeIDFFFF}, ids)
}
