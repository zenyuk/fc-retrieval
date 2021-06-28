/*
Package cid - provides methods for ContentIDAdapter struct.

ContentIDAdapter is 32 bytes is a unique identifier of a file stored in a Filecoin blockchain network.
*/
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
	"crypto/sha256"
	"fmt"

	"github.com/cbergoon/merkletree"
)

// ContentIDAdapter represents a CID.
type ContentIDAdapter struct {
	Id string
}

// CalculateHash hashes the values of a ContentIDAdapter.
func (n ContentIDAdapter) CalculateHash() ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write([]byte(n.Id)); err != nil {
		return nil, err
	}

	fmt.Printf("h.Sum(nil) %+v\n", h.Sum(nil))

	return h.Sum(nil), nil
}

//Equals tests for equality of two Contents
func (n ContentIDAdapter) Equals(other merkletree.Content) (bool, error) {
	return n.ToString() == other.(*ContentIDAdapter).ToString(), nil
}

// ToString returns a string for the ContentID.
func (n *ContentIDAdapter) ToString() string {
	return n.Id
}
