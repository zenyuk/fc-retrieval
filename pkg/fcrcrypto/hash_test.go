package fcrcrypto

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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlockchainHasher(t *testing.T) {
	someBytes := make([]byte, 100)
	hasher := GetBlockchainHasher()
	hasher.Write(someBytes)
	digest1 := hasher.Sum(nil)

	hasher.Reset()
	hasher.Write(someBytes)
	digest2 := hasher.Sum(nil)

	digest4 := BlockchainHash(someBytes)

	assert.Equal(t, digest1, digest2)
	assert.Equal(t, digest1, digest4)
}

func TestRetrievalHasher(t *testing.T) {
	someBytes := make([]byte, 100)
	hasher := GetRetrievalV1Hasher()
	hasher.Write(someBytes)
	digest1 := hasher.Sum(nil)

	hasher.Reset()
	hasher.Write(someBytes)
	digest2 := hasher.Sum(nil)

	digest4 := RetrievalV1Hash(someBytes)

	assert.Equal(t, digest1, digest2)
	assert.Equal(t, digest1, digest4)
}
