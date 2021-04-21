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
	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
)

// TestGetClosestNodeIDs
func TestGetClosestNodeIDs(t *testing.T) {
	pieceCID, _ := cid.NewContentIDFromHexString("01")
	nodeID00, _ := nodeid.NewNodeIDFromHexString("00")
	nodeID01, _ := nodeid.NewNodeIDFromHexString("01")
	nodeID02, _ := nodeid.NewNodeIDFromHexString("02")
	nodeID5A, _ := nodeid.NewNodeIDFromHexString("5A")
	nodeIDFFFF, _ := nodeid.NewNodeIDFromHexString("FFFF")
	
	actual1, _ := GetClosestNodeIDs(
		*pieceCID,
		[]*nodeid.NodeID{nodeID5A, nodeID01, nodeID00, nodeIDFFFF, nodeID02},
		1,
	)
	actual2, _ := GetClosestNodeIDs(
		*pieceCID,
		[]*nodeid.NodeID{nodeID5A, nodeID01, nodeID00, nodeIDFFFF, nodeID02},
		2,
	)
	actual3, _ := GetClosestNodeIDs(
		*pieceCID,
		[]*nodeid.NodeID{nodeID5A, nodeID01, nodeID00, nodeIDFFFF, nodeID02},
		16,
	)

	assert.ElementsMatch(t, []*nodeid.NodeID{nodeID01}, actual1)
	assert.ElementsMatch(t, []*nodeid.NodeID{nodeID01, nodeID00}, actual2)
	assert.ElementsMatch(t, []*nodeid.NodeID{nodeID01, nodeID00, nodeID02, nodeID5A, nodeIDFFFF}, actual3)
}