package math

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
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetDistance
func TestGetDistance(t *testing.T) {
	str00, _ := hex.DecodeString("00") // bigInt 0
	str01, _ := hex.DecodeString("01") // bigInt 1
	str02, _ := hex.DecodeString("02") // bigInt 2
	str5A, _ := hex.DecodeString("5A") // bigInt 90
	strFFFF, _ := hex.DecodeString("FFFF") // bigInt 65535
	mockCardinality := big.NewInt(65536)

	actual1, _ := GetDistance(str00, str00, mockCardinality)
	actual2, _ := GetDistance(str01, str00, mockCardinality)
	actual3, _ := GetDistance(str01, str02, mockCardinality)
	actual4, _ := GetDistance(str01, str5A, mockCardinality)
	actual5, _ := GetDistance(str01, strFFFF, mockCardinality)

	assert.Equal(t, big.NewInt(0), actual1)
	assert.Equal(t, big.NewInt(1), actual2)
	assert.Equal(t, big.NewInt(1), actual3)
	assert.Equal(t, big.NewInt(89), actual4)
	assert.Equal(t, big.NewInt(2), actual5)
}
