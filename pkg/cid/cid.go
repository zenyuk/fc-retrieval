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
    "math/big"
)



// ContentID represents a CID
type ContentID struct {
    id big.Int
}


// NewContentID creates a CID object
func NewContentID(id *big.Int) (*ContentID) {
	var n = ContentID{}
    n.id = *id
	return &n
}

// ToString returns a string for the CID.
func (n *ContentID) ToString() (string) {
    return n.id.Text(16)
}

// ToBytes returns the byte array representation of the CID.
func (n *ContentID) ToBytes() ([]byte) {
    return n.id.Bytes()
}