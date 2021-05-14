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
	"math/big"
	"sort"

	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/math"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
)

func getClosestNodeIDs(landmark []byte, nodeIDs []*nodeid.NodeID, maxResults int) ([]*nodeid.NodeID, error) {
	var m = make(map[*big.Int][]*nodeid.NodeID)
	cardinality := big.NewInt(0).Exp(big.NewInt(2), big.NewInt(cid.WordSize*8), nil) // 2**(32*8)
	for _, nodeID := range nodeIDs {
		dist, _ := math.GetDistance(nodeID.ToBytes(), landmark, cardinality)
		a, ok := m[dist]
		if !ok {
			a = make([]*nodeid.NodeID, 0)
		}
		a = append(a, nodeID)
		m[dist] = a
	}
	keys := make([]*big.Int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i].Cmp(keys[j]) < 0
	})
	var r = make([]*nodeid.NodeID, 0)
	for _, k := range keys {
		a := m[k]
		r = append(r, a...)
	}
	if len(r) < maxResults {
		return r, nil
	} else {
		return r[0:maxResults], nil
	}
}

// GetNodeIDsClosestToContentID gets nodeIDs that are close to a contentID
func GetNodeIDsClosestToContentID(landmark []byte, nodeIDs []*nodeid.NodeID, maxResults int) ([]*nodeid.NodeID, error) {
	return getClosestNodeIDs(landmark, nodeIDs, maxResults)
}

// SortClosestNodesIDs sort nodeIDs that are close to a nodeID
func SortClosestNodesIDs(landmark []byte, nodeIDs []*nodeid.NodeID) ([]*nodeid.NodeID, error) {
	return getClosestNodeIDs(landmark, nodeIDs, len(nodeIDs))
}
