package dhtring

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

func TestInsertOne(t *testing.T) {
	r := CreateRing()
	r.Insert("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, 1, r.Size())

	r = CreateRing()
	r.Insert("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	assert.Equal(t, 1, r.Size())

	r = CreateRing()
	r.Insert("101112131415161718191A1B1C1D1E1F202122232425262728292A2B2C2D2E2F")
	assert.Equal(t, 1, r.Size())

	r = CreateRing()
	r.Insert("101112131415161718191A1B1C1D1E1F202122232425262728292A2B2C2D2E2")
	assert.Equal(t, 0, r.Size())

	r = CreateRing()
	r.Insert("Invalid")
	assert.Equal(t, 0, r.Size())
}

func TestInsertTwo(t *testing.T) {
	r := CreateRing()
	r.Insert("1000000000000000000000000000000000000000000000000000000000000000")
	r.Insert("2000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, 2, r.Size())

	r = CreateRing()
	r.Insert("2000000000000000000000000000000000000000000000000000000000000000")
	r.Insert("1000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, 2, r.Size())
}

func TestInsertThree(t *testing.T) {
	r := CreateRing()
	r.Insert("1000000000000000000000000000000000000000000000000000000000000000")
	r.Insert("2000000000000000000000000000000000000000000000000000000000000000")
	r.Insert("3000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, 3, r.Size())

	r = CreateRing()
	r.Insert("1000000000000000000000000000000000000000000000000000000000000000")
	r.Insert("3000000000000000000000000000000000000000000000000000000000000000")
	r.Insert("2000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, 3, r.Size())

	r = CreateRing()
	r.Insert("2000000000000000000000000000000000000000000000000000000000000000")
	r.Insert("1000000000000000000000000000000000000000000000000000000000000000")
	r.Insert("3000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, 3, r.Size())

	r = CreateRing()
	r.Insert("2000000000000000000000000000000000000000000000000000000000000000")
	r.Insert("3000000000000000000000000000000000000000000000000000000000000000")
	r.Insert("3000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, 2, r.Size())
}

func TestRemove(t *testing.T) {
	r := CreateRing()
	r.Insert("F000000000000000000000000000000000000000000000000000000000000000")
	r.Insert("1000000000000000000000000000000000000000000000000000000000000000")
	r.Insert("3000000000000000000000000000000000000000000000000000000000000000")
	r.Insert("2000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, 4, r.Size())

	// Test remove not existed
	r.Remove("2000000000000000000000000000000000000000000000000000000000000001")
	assert.Equal(t, 4, r.Size())

	// Test remove head
	r.Remove("1000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, 3, r.Size())
	r.Insert("1000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, 4, r.Size())

	// Test remove tail
	r.Remove("F000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, 3, r.Size())
	r.Insert("F000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, 4, r.Size())

	// Test remove middle
	r.Remove("2000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, 3, r.Size())
	r.Insert("2000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, 4, r.Size())

	r.Remove("F000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, 3, r.Size())
	r.Remove("1000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, 2, r.Size())
	r.Remove("2000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, 1, r.Size())
	r.Remove("3000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, 0, r.Size())
}

func TestGetNodeIDsClosestToContentID(t *testing.T) {

	cid1 := "7080000000000000000000000000000000000000000000000000000000000000"
	cid2 := "1080000000000000000000000000000000000000000000000000000000000000"
	cid3 := "F080000000000000000000000000000000000000000000000000000000000000"
	cid4 := "5000000000000000000000000000000000000000000000000000000000000000"

	node0 := "0000000000000000000000000000000000000000000000000000000000000000"
	node1 := "1000000000000000000000000000000000000000000000000000000000000000"
	node2 := "2000000000000000000000000000000000000000000000000000000000000000"
	node3 := "3000000000000000000000000000000000000000000000000000000000000000"
	node4 := "4000000000000000000000000000000000000000000000000000000000000000"
	node5 := "5000000000000000000000000000000000000000000000000000000000000000"
	node6 := "6000000000000000000000000000000000000000000000000000000000000000"
	node7 := "6800000000000000000000000000000000000000000000000000000000000000"
	node8 := "7000000000000000000000000000000000000000000000000000000000000000"
	node9 := "8000000000000000000000000000000000000000000000000000000000000000"
	node10 := "9000000000000000000000000000000000000000000000000000000000000000"
	node11 := "A000000000000000000000000000000000000000000000000000000000000000"
	node12 := "B000000000000000000000000000000000000000000000000000000000000000"
	node13 := "C000000000000000000000000000000000000000000000000000000000000000"
	node14 := "D000000000000000000000000000000000000000000000000000000000000000"
	node15 := "E000000000000000000000000000000000000000000000000000000000000000"
	node16 := "F000000000000000000000000000000000000000000000000000000000000000"
	node17 := "F800000000000000000000000000000000000000000000000000000000000000"

	r := CreateRing()
	r.Insert(node0)
	r.Insert(node1)
	r.Insert(node2)
	r.Insert(node3)
	r.Insert(node4)
	r.Insert(node5)
	r.Insert(node6)
	r.Insert(node7)
	r.Insert(node8)
	r.Insert(node9)
	r.Insert(node10)
	r.Insert(node11)
	r.Insert(node12)
	r.Insert(node13)
	r.Insert(node14)
	r.Insert(node15)
	r.Insert(node16)
	r.Insert(node17)
	r.Dump()
	assert.Equal(t, 18, r.Size())

	res1, _ := r.GetClosest(cid1, 8, "")
	res2, _ := r.GetClosest(cid2, 8, "")
	res3, _ := r.GetClosest(cid3, 4, "")
	res4, _ := r.GetClosest(cid4, 3, "")
	res5, _ := r.GetClosest(node5, 3, node5)

	assert.Equal(t, 18, r.Size())

	assert.Equal(t, 8, len(res1))
	assert.Equal(t, node4, res1[0])
	assert.Equal(t, node5, res1[1])
	assert.Equal(t, node6, res1[2])
	assert.Equal(t, node7, res1[3])
	assert.Equal(t, node8, res1[4])
	assert.Equal(t, node9, res1[5])
	assert.Equal(t, node10, res1[6])
	assert.Equal(t, node11, res1[7])

	assert.Equal(t, 8, len(res2))
	assert.Equal(t, node15, res2[0])
	assert.Equal(t, node16, res2[1])
	assert.Equal(t, node17, res2[2])
	assert.Equal(t, node0, res2[3])
	assert.Equal(t, node1, res2[4])
	assert.Equal(t, node2, res2[5])
	assert.Equal(t, node3, res2[6])
	assert.Equal(t, node4, res2[7])

	assert.Equal(t, 4, len(res3))
	assert.Equal(t, node15, res3[0])
	assert.Equal(t, node16, res3[1])
	assert.Equal(t, node17, res3[2])
	assert.Equal(t, node0, res3[3])

	assert.Equal(t, 3, len(res4))
	assert.Equal(t, node4, res4[0])
	assert.Equal(t, node5, res4[1])
	assert.Equal(t, node6, res4[2])

	assert.Equal(t, 3, len(res5))
	assert.Equal(t, node4, res5[0])
	assert.Equal(t, node6, res5[1])
	assert.Equal(t, node7, res5[2])
}

func TestGetWithinRange(t *testing.T) {
	r := CreateRing()
	r.Insert("F000000000000000000000000000000000000000000000000000000000000000")
	r.Insert("1000000000000000000000000000000000000000000000000000000000000000")
	r.Insert("3000000000000000000000000000000000000000000000000000000000000000")
	r.Insert("2000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, 4, r.Size())

	res, err := r.GetWithinRange("1000000000000000000000000000000000000000000000000000000000000001", "2000000000000000000000000000000000000000000000000000000000000001")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, []string{
		"2000000000000000000000000000000000000000000000000000000000000000"}, res)

	res, err = r.GetWithinRange("F000000000000000000000000000000000000000000000000000000000000001", "2000000000000000000000000000000000000000000000000000000000000001")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, []string{
		"1000000000000000000000000000000000000000000000000000000000000000",
		"2000000000000000000000000000000000000000000000000000000000000000"}, res)

	res, err = r.GetWithinRange("F000000000000000000000000000000000000000000000000000000000000001", "3000000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, []string{
		"1000000000000000000000000000000000000000000000000000000000000000",
		"2000000000000000000000000000000000000000000000000000000000000000",
		"3000000000000000000000000000000000000000000000000000000000000000"}, res)

	res, err = r.GetWithinRange("F000000000000000000000000000000000000000000000000000000000000000", "2000000000000000000000000000000000000000000000000000000000000001")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, []string{
		"F000000000000000000000000000000000000000000000000000000000000000",
		"1000000000000000000000000000000000000000000000000000000000000000",
		"2000000000000000000000000000000000000000000000000000000000000000"}, res)

	res, err = r.GetWithinRange("F000000000000000000000000000000000000000000000000000000000000000", "3000000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, []string{
		"F000000000000000000000000000000000000000000000000000000000000000",
		"1000000000000000000000000000000000000000000000000000000000000000",
		"2000000000000000000000000000000000000000000000000000000000000000",
		"3000000000000000000000000000000000000000000000000000000000000000"}, res)
}

func TestGetRingNodeByAddress_ReturnsCorrectExisting(t *testing.T) {
	// arrange
	searchAddress := "3000000000000000000000000000000000000000000000000000000000000000"
	r := CreateRing()
	r.Insert("F000000000000000000000000000000000000000000000000000000000000000")
	r.Insert("1000000000000000000000000000000000000000000000000000000000000000")
	r.Insert("2000000000000000000000000000000000000000000000000000000000000000")
	r.Insert(searchAddress)

	// act
	ringNode := r.get(searchAddress)

	// assert
	assert.Equal(t, searchAddress, ringNode.val)
}

func TestGetRingNodeByAddress_ReturnsNilForNotExisting(t *testing.T) {
	// arrange
	notExistingAddress := "3000000000000000000000000000000000000000000000000000000000000099"
	r := CreateRing()
	r.Insert("F000000000000000000000000000000000000000000000000000000000000000")
	r.Insert("1000000000000000000000000000000000000000000000000000000000000000")
	r.Insert("2000000000000000000000000000000000000000000000000000000000000000")
	r.Insert("3000000000000000000000000000000000000000000000000000000000000000")

	// act
	ringNode := r.get(notExistingAddress)

	// assert
	assert.Nil(t, ringNode)
}
