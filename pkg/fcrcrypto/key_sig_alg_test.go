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



func TestKeySigAlgRoundTrip(t *testing.T) {
    algObj := DecodeSigAlg(SigAlgEcdsaSecP256K1Blake2b)
    algVal := algObj.EncodeSigAlg()
    assert.Equal(t, algVal, SigAlgEcdsaSecP256K1Blake2b)
}

func TestKeySigAlgRoundTripBytes(t *testing.T) {
    algObj := DecodeSigAlg(SigAlgEcdsaSecP256K1Blake2b)
    algBytes := algObj.EncodeSigAlgAsBytes()
    algObj2 := DecodeSigAlgFromBytes(algBytes)
    assert.True(t, algObj.Equals(algObj2))
}


func TestIs(t *testing.T) {
    algObj1 := DecodeSigAlg(SigAlgEcdsaSecP256K1Blake2b)
    algObj3 := DecodeSigAlg(101)
    assert.True(t, algObj1.Is(SigAlgEcdsaSecP256K1Blake2b))
    assert.False(t, algObj1.Is(101))
    assert.False(t, algObj3.Is(SigAlgEcdsaSecP256K1Blake2b))
}

func TestIsNot(t *testing.T) {
    algObj1 := DecodeSigAlg(SigAlgEcdsaSecP256K1Blake2b)
    algObj3 := DecodeSigAlg(101)
    assert.False(t, algObj1.IsNot(SigAlgEcdsaSecP256K1Blake2b))
    assert.True(t, algObj1.IsNot(101))
    assert.True(t, algObj3.IsNot(SigAlgEcdsaSecP256K1Blake2b))
}
