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

/**
 * Key and signature algorithm.
 *
 */

const (
	lengthOfSigAlgInBytes = 1

	// SigAlgEcdsaSecP256K1Blake2b indicates the signature algorithm ECDSA
	// using curve SecP256K1 with Blake2b message digest algorithm.
	SigAlgEcdsaSecP256K1Blake2b = uint8(1)
)

// KeySigAlg is holds the signature algorithm
type KeySigAlg struct {
	algorithm uint8
}

// DecodeSigAlg converts a number to an object.
func DecodeSigAlg(alg uint8) KeySigAlg {
	return KeySigAlg{algorithm: alg}
}

// DecodeSigAlgFromBytes converts bytes to an object
func DecodeSigAlgFromBytes(data []byte) *KeySigAlg {
	k := KeySigAlg{}
	k.algorithm = data[0]
	return &k
}

// EncodeSigAlg converts a number to an object.
func (k *KeySigAlg) EncodeSigAlg() uint8 {
	return k.algorithm
}

// EncodeSigAlgAsBytes converts an object to bytes.
func (k *KeySigAlg) EncodeSigAlgAsBytes() []byte {
	algBytes := make([]byte, lengthOfSigAlgInBytes)
	algBytes[0] = k.algorithm
	return algBytes
}

// Is returns true if the value passed in matches the algorithm.
func (k *KeySigAlg) Is(other uint8) bool {
	return other == k.algorithm
}

// Equals returns true if the value passed in matches the algorithm.
func (k *KeySigAlg) Equals(other *KeySigAlg) bool {
	return other.algorithm == k.algorithm
}

// IsNot returns true if the value passed in does not match the algorithm.
func (k *KeySigAlg) IsNot(other uint8) bool {
	return other != k.algorithm
}
