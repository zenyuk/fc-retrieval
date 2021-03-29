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
	assert.Empty(t, err)
	cid2, err := cid.NewContentIDFromHexString("02")
	assert.Empty(t, err)
	cid3, err := cid.NewContentIDFromHexString("03")
	assert.Empty(t, err)
	cid4, err := cid.NewContentIDFromHexString("04")
	assert.Empty(t, err)
	cid5, err := cid.NewContentIDFromHexString("05")
	assert.Empty(t, err)
	tree, err := CreateMerkleTree([]merkletree.Content{cid1, cid2, cid3, cid4, cid5})
	assert.Empty(t, err)
	assert.NotEmpty(t, tree)
	_, err = CreateMerkleTree([]merkletree.Content{})
	assert.NotEmpty(t, err)
	assert.Equal(t, "b3704b0f54c4070b36ace72ae3f4879e91eda3f81cef4cb97653d2a0b8592ce1", tree.GetMerkleRoot())
}

func TestCreateTreeOneElement(t *testing.T) {
	cid1, err := cid.NewContentIDFromHexString("01")
	assert.Empty(t, err)
	tree, err := CreateMerkleTree([]merkletree.Content{cid1})
	assert.Empty(t, err)
	assert.NotEmpty(t, tree)
	assert.Equal(t, "c3c3a46684c07d12a9c238787df3049a6f258e7af203e5ddb66a8bd66637e108", tree.GetMerkleRoot())
}

func TestCreateTreeManyElements(t *testing.T) {
	elements := make([]merkletree.Content, 0)
	for i := 0; i < 100; i++ {
		cid, err := cid.NewContentID(big.NewInt(int64(i)))
		assert.Empty(t, err)
		elements = append(elements, cid)
	}
	tree, err := CreateMerkleTree(elements)
	assert.Empty(t, err)
	assert.NotEmpty(t, tree)
	assert.Equal(t, "eeb5943906dc9937a7b42ca9b57e2afe51d7551245a0ba00695b91cd261dafb0", tree.GetMerkleRoot())
}
