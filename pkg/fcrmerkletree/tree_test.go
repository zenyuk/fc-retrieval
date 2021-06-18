package fcrmerkletree

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
	"testing"

	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/cbergoon/merkletree"
	"github.com/stretchr/testify/assert"
)

func TestCreateTree(t *testing.T) {
	cid1, err := cid.NewContentIDFromHexString("01")
	cida1 := cid.ContentIDAdapter{Id: cid1.ToString()}
	assert.Empty(t, err)
	cid2, err := cid.NewContentIDFromHexString("02")
	cida2 := cid.ContentIDAdapter{Id: cid2.ToString()}
	assert.Empty(t, err)
	cid3, err := cid.NewContentIDFromHexString("03")
	cida3 := cid.ContentIDAdapter{Id: cid3.ToString()}
	assert.Empty(t, err)
	cid4, err := cid.NewContentIDFromHexString("04")
	cida4 := cid.ContentIDAdapter{Id: cid4.ToString()}
	assert.Empty(t, err)
	cid5, err := cid.NewContentIDFromHexString("05")
	cida5 := cid.ContentIDAdapter{Id: cid5.ToString()}
	assert.Empty(t, err)
	tree, err := CreateMerkleTree([]merkletree.Content{cida1, cida2, cida3, cida4, cida5})
	assert.Empty(t, err)
	assert.NotEmpty(t, tree)
	_, err = CreateMerkleTree([]merkletree.Content{})
	assert.NotEmpty(t, err)
	assert.Equal(t, "93fe3ef34b47ac151048bb83ddce0656f794ad70adde0f563e1e6fb129d5a15f", tree.GetMerkleRoot())
}

func TestCreateTreeOneElement(t *testing.T) {
	cid1, err := cid.NewContentIDFromHexString("01")

	cida1 := cid.ContentIDAdapter{Id: cid1.ToString()}

	assert.Empty(t, err)
	tree, err := CreateMerkleTree([]merkletree.Content{cida1})
	assert.Empty(t, err)
	assert.NotEmpty(t, tree)

	assert.Equal(t, "c8bc8219fe198ce0a3b7fa9d14c7f88c16fc9bae96622761643eaafbc7eb7303", tree.GetMerkleRoot())
}

func TestCreateTreeManyElements(t *testing.T) {
	elements := make([]merkletree.Content, 0)
	for i := 0; i < 100; i++ {
		cid1, err := cid.NewContentID(big.NewInt(int64(i)))
		cida := cid.ContentIDAdapter{Id: cid1.ToString()}
		assert.Empty(t, err)
		elements = append(elements, cida)
	}
	tree, err := CreateMerkleTree(elements)
	assert.Empty(t, err)
	assert.NotEmpty(t, tree)
	assert.Equal(t, "f217029c872d9d64f6f5a860fcd027dd8348229f3fb8eb0c7a9ad534ecfa3eae", tree.GetMerkleRoot())
}
