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
	"encoding/hex"
	"fmt"
	"math/big"
)

const wordSize = 32 // 32 bytes

// ContentID represents a CID
type ContentID struct {
	//id big.Int
	id []byte
}

// NewContentID creates a CID object
func NewContentID(id *big.Int) (*ContentID, error) {
	var n = ContentID{}
	//n.id = *id
	b := id.Bytes()
	l := len(b)
	if l > wordSize {
		return nil, fmt.Errorf("NodeID: Incorrect size1: %d", l)
	}
	idBytes := id.Bytes()
	n.id = make([]byte, len(idBytes))
	copy(n.id, idBytes)
	return &n, nil
}

// ToString returns a string for the CID.
func (n *ContentID) ToString() string {
	//return n.id.Text(16)
	str := hex.EncodeToString(n.id)
	if str == "" {
		str = "00"
	}
	return str
}

// ToBytes returns the byte array representation of the CID.
func (n *ContentID) ToBytes() []byte {
	return n.id
}

// MarshalJSON is used to marshal NodeID into bytes
func (n ContentID) MarshalJSON() ([]byte, error) {
	return n.id, nil
}

// UnmarshalJSON is used to unmarshal bytes into NodeID
func (n *ContentID) UnmarshalJSON(p []byte) error {
	// if string(p) == "null" {
	// 	return nil
	// }
	// var z big.Int
	// _, ok := z.SetString(string(p), 16)
	// if !ok {
	// 	return errors.New("Not a valid big integer: " + string(p))
	// }
	// n.id = z
	// return nil
	if len(p) != wordSize {
		return fmt.Errorf("NodeID: Incorrect size: %d", len(p))
	}
	copy(p, n.id)
	return nil
}
