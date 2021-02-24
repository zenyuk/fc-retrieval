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
	"encoding/binary"
	"fmt"
)


/**
 * Key Version: Keys used by the Gateway will be versioned to allow for key roll-over.
 *
 */
const (
	bitsInKeyVersion = 32
	lengthOfKeyVersionInBytes = 4

	initialKeyVersion = 1
	keyVersionIncrement = 1
)


// KeyVersion wraps a key version number.
type KeyVersion struct {
	ver uint32
}

// InitialKeyVersion returns the first key version
func InitialKeyVersion() *KeyVersion {
	return &KeyVersion{initialKeyVersion}
}

// NextKeyVersion creates a new key version
func (k *KeyVersion) NextKeyVersion() *KeyVersion {
	return &KeyVersion{k.ver + keyVersionIncrement}
}

// DecodeKeyVersion converts a number to an object.
func DecodeKeyVersion(ver uint32) *KeyVersion {
	return &KeyVersion{ver}
}

// DecodeKeyVersionFromBytes converts a byte array to ann object.
func DecodeKeyVersionFromBytes(version []byte) (*KeyVersion, error) {
	if len(version) < lengthOfKeyVersionInBytes {
		return nil, fmt.Errorf("Version bytes incorrect length: %d", len(version))
	}
	k := KeyVersion{}
	k.ver = binary.BigEndian.Uint32(version)
	return &k, nil
}

// EncodeKeyVersion converts an object to a number
func (k *KeyVersion) EncodeKeyVersion() (uint32) {
	return k.ver
}

// EncodeKeyVersionAsBytes converts an object to a byte array
func (k *KeyVersion) EncodeKeyVersionAsBytes() []byte {
	result := make([]byte, lengthOfKeyVersionInBytes)
	binary.BigEndian.PutUint32(result, k.ver)
	return result
}

// Equals returns true if the value passed in matches the version.
func (k *KeyVersion) Equals(other *KeyVersion) bool {
	return other.ver == k.ver
}

// EqualsRaw returns true if the value passed in matches the version.
func (k *KeyVersion) EqualsRaw(other uint32) bool {
	return other == k.ver
}

// NotEquals returns true if the value passed in does not match the version.
func (k *KeyVersion) NotEquals(other *KeyVersion) bool {
	return other.ver != k.ver
}

// NotEqualsRaw returns true if the value passed in does not match the version.
func (k *KeyVersion) NotEqualsRaw(other uint32) bool {
	return other != k.ver
}
