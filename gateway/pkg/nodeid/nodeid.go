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
	"fmt"
	"math/big"
	"encoding/hex"
)

const wordSize = 32 // 32 bytes

// NodeID represents a Gateway id
type NodeID struct {
	id []byte
}

// NewNodeID creates a node id object
func NewNodeID(id *big.Int) (*NodeID, error) {
	var n = NodeID{}
	b := id.Bytes()
	l := len(b)
	if l > wordSize {
		return nil, fmt.Errorf("NodeID: Incorrect size1: %d", l)
	}
	copy(n.id, id.Bytes())
	return &n, nil
}

// NewNodeIDFromBytes creates a node id object
func NewNodeIDFromBytes(id []byte) (*NodeID, error) {
	var n = NodeID{}
	if len(id) > wordSize {
		return nil, fmt.Errorf("NodeID: Incorrect size2: %d", len(id))
	}
	copy(n.id, id)
	return &n, nil
}

// NewNodeIDFromString creates a NodeID from a string
func NewNodeIDFromString(id string) (*NodeID, error) {
	bytes, err := hex.DecodeString(id)
	if err != nil {
		return nil, err
	}

	if len(bytes) > wordSize {
		return nil, fmt.Errorf("NodeID: Incorrect size2: %d", len(id))
	}
	var n = NodeID{}
	n.id = bytes
	return &n, nil

}

// ToString returns a string for the node id.
func (n *NodeID) ToString() string {
	return hex.EncodeToString(n.id)
}

// ToBytes returns the byte array representation of the node id.
func (n *NodeID) ToBytes() []byte {
	return n.id
}

// AsBytes32 returns the node id as a [32]byte
func (n *NodeID) AsBytes32() (result [wordSize]byte) {
	copy(result[:], n.id)
	return
}

// MarshalJSON is used to marshal NodeID into bytes
func (n NodeID) MarshalJSON() ([]byte, error) {
	return n.id, nil
}

// UnmarshalJSON is used to unmarshal bytes into NodeID
func (n *NodeID) UnmarshalJSON(p []byte) error {
	if len(p) != wordSize {
		return fmt.Errorf("NodeID: Incorrect size: %d", len(p))
	}
	copy(p, n.id)
	return nil
}
