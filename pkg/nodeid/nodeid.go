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

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
)

const wordSize = 32 // the NodeID length is 32 bytes.

// NodeID represents a NodeID.
type NodeID struct {
	id []byte
}

// NewNodeID creates a NodeID object.
func NewNodeID(id *big.Int) (*NodeID, error) {
	var n = NodeID{}
	b := id.Bytes()
	l := len(b)
	if l > wordSize {
		return nil, fmt.Errorf("NodeID: Incorrect size: %d, should be fewer than %d", l, wordSize)
	}
	n.id = make([]byte, wordSize)
	copy(n.id[wordSize-l:], b)
	return &n, nil
}

// NewNodeIDFromBytes creates a NodeID object.
func NewNodeIDFromBytes(id []byte) (*NodeID, error) {
	var n = NodeID{}
	lenID := len(id)
	if lenID > wordSize {
		return nil, fmt.Errorf("NodeID: Incorrect size: %d, should be fewer than %d", lenID, wordSize)
	}
	n.id = make([]byte, wordSize)
	copy(n.id[wordSize-len(id):], id)
	return &n, nil
}

// NewNodeIDFromHexString creates a NodeID from a string.
func NewNodeIDFromHexString(id string) (*NodeID, error) {
	var n = NodeID{}
	bytes, err := hex.DecodeString(id)
	if err != nil {
		return nil, err
	}

	if len(bytes) > wordSize {
		return nil, fmt.Errorf("NodeID: Incorrect size: %d, should be fewer than %d", len(id), wordSize)
	}
	n.id = make([]byte, wordSize)
	copy(n.id[wordSize-len(bytes):], bytes)
	return &n, nil
}

// NewNodeIDFromPublicKey create a NodeID based on a public key.
func NewNodeIDFromPublicKey(pubKey *fcrcrypto.KeyPair) (*NodeID, error) {
	hashedPubKey, err := pubKey.HashPublicKey()
	if err != nil {
		return nil, err
	}
	return NewNodeIDFromBytes(hashedPubKey)
}

// NewRandomNodeID creates a random NodeID object.
func NewRandomNodeID() *NodeID {
	var n = NodeID{}
	n.id = make([]byte, wordSize)
	fcrcrypto.GeneratePublicRandomBytes(n.id)
	return &n
}

// ToString returns a string for the NodeID.
func (n *NodeID) ToString() string {
	str := hex.EncodeToString(n.id)
	if str == "" {
		str = "00"
	}
	return str
}

// ToBytes returns the byte array representation of the NodeID.
func (n *NodeID) ToBytes() []byte {
	return n.id
}

// AsBytes32 returns the NodeID as a [32]byte.
func (n *NodeID) AsBytes32() (result [wordSize]byte) {
	copy(result[:], n.id)
	return
}

// MarshalJSON is used to marshal NodeID into bytes.
func (n NodeID) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.id)
}

// UnmarshalJSON is used to unmarshal bytes into NodeID.
func (n *NodeID) UnmarshalJSON(p []byte) error {
	var id []byte
	err := json.Unmarshal(p, &id)
	if err != nil {
		return err
	}

	if len(id) != wordSize {
		return fmt.Errorf("NodeID: Incorrect size: %d, should be %d", len(id), wordSize)
	}
	n.id = make([]byte, wordSize)
	copy(n.id, id)
	return nil
}
