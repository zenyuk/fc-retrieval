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



func TestKeyVersionRoundTrip(t *testing.T) {
    kv := InitialKeyVersion()
    rawVer := kv.EncodeKeyVersion()
    kv1 := DecodeKeyVersion(rawVer)
    assert.True(t, kv.Equals(kv1))
    assert.False(t, kv.NotEquals(kv1))
}


func TestKeyVersionNext(t *testing.T) {
    kv := InitialKeyVersion()
    kv2 := kv.NextKeyVersion()
    assert.False(t, kv.Equals(kv2))
    assert.True(t, kv.NotEquals(kv2))
}

func TestKeyVersionRaw(t *testing.T) {
    kv := InitialKeyVersion()
    ver := kv.EncodeKeyVersion()
    ver2 := ver+1

    assert.True(t, kv.EqualsRaw(ver))
    assert.True(t, kv.NotEqualsRaw(ver2))
}



func TestKeyVersionBytesRoundTrip(t *testing.T) {
    kv := InitialKeyVersion()
    verBytes := kv.EncodeKeyVersionAsBytes()
    kv2, err := DecodeKeyVersionFromBytes(verBytes)
    if err != nil {
        panic(err)
    }
    assert.True(t, kv.Equals(kv2))
}

func TestKeyVersionBytesBadLen(t *testing.T) {
    verBytes := make([]byte, 1)
    _, err := DecodeKeyVersionFromBytes(verBytes)
    assert.Error(t, err)
}


