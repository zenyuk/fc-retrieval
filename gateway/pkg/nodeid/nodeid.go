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
	"errors"
	"math/big"
)

const wordSize = 32 // 32 bytes

// NodeID represents a Gateway id
type NodeID struct {
	id big.Int
}

// NewNodeID creates a node id object
func NewNodeID(id *big.Int) *NodeID {
	var n = NodeID{}
	n.id = *id
	return &n
}

// ToString returns a string for the node id.
func (n *NodeID) ToString() string {
	return n.id.Text(16)
}

// ToBytes returns the byte array representation of the node id.
func (n *NodeID) ToBytes() []byte {
	return n.id.Bytes()
}

// MarshalJSON is used to marshal NodeID into bytes
func (n NodeID) MarshalJSON() ([]byte, error) {
	return []byte(n.ToString()), nil
}

// UnmarshalJSON is used to unmarshal bytes into NodeID
func (n *NodeID) UnmarshalJSON(p []byte) error {
	if string(p) == "null" {
		return nil
	}
	var z big.Int
	_, ok := z.SetString(string(p), 16)
	if !ok {
		return errors.New("Not a valid big integer: " + string(p))
	}
	n.id = z
	return nil
}
