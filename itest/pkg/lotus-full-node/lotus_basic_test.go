/*
Package lotus_full_node - end-to-end tests, calling private, locally hosted Lotus network running Lotus Full Node
(with Daemon and Miner).
*/
package lotus_full_node

/*
 * Copyright 2021 ConsenSys Software Inc.
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
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/ConsenSys/fc-retrieval/common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval/itest/pkg/util"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/api/apistruct"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Need to make sure this env is not set in host machine
	itestEnv := os.Getenv("ITEST_CALLING_FROM_CONTAINER")

	if itestEnv != "" {
		// Env is set, we are calling from docker container
		m.Run()
		return
	}
	// Env is not set, we are calling from host
	// We need a lotus
	tag := util.GetCurrentBranch()

	// Create shared net
	bgCtx := context.Background()
	ctx, _ := context.WithTimeout(bgCtx, time.Minute*2)
	network, networkName := util.CreateNetwork(ctx)

	// Start lotus
	lotusContainer := util.StartLotusFullNode(ctx, networkName, true)

	// Start itest
	done := make(chan bool)
	itestContainer := util.StartItest(ctx, tag, networkName, util.ColorGreen, "", "", done, true, "")
	// Block until done.
	if <-done {
		logging.Info("Tests passed, shutdown...")
	} else {
		logging.Error("Tests failed, shutdown...")
	}

	if err := itestContainer.Terminate(ctx); err != nil {
		logging.Error("error while terminating test container: %s", err.Error())
	}
	if err := lotusContainer.Terminate(ctx); err != nil {
		logging.Error("error while terminating test container: %s", err.Error())
	}
	if err :=  (*network).Remove(ctx); err != nil {
		logging.Error("error while terminating test container network: %s", err.Error())
	}
}

func TestLotusFullNodeConnectivityHttp(t *testing.T) {

	method := "Filecoin.ChainHead"
	id := 1
	lotusUrl := "http://lotus-full-node:1234/rpc/v0"
	requestBody := `{
		"jsonrpc": "2.0",
		"method": "` + method + `",
		"id": ` + strconv.Itoa(id) + `,
		"params": []
	}`

	resp, err := http.Post(lotusUrl, "application/json", bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		t.Errorf("Can't establish connection to Lotus %s", lotusUrl)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Can't read response from Lotus, method: %s", method)
	}
	var fields map[string]interface{}
	err = json.Unmarshal(body, &fields)
	if err != nil {
		t.Errorf("Can't parse json response from Lotus, method: %s", method)
	}

	assert.Equal(t, float64(id), fields["id"])
}

func TestLotusFullNodeConnectivityWs(t *testing.T) {
	var lotusApi apistruct.FullNodeStruct
	bgCtx := context.Background()
	ctx, _ := context.WithTimeout(bgCtx, time.Minute*3)

	clientClose, err := jsonrpc.NewMergeClient(
		ctx,
		"ws://lotus-full-node:1234/rpc/v0",
		"Filecoin",
		[]interface{}{
			&lotusApi.CommonStruct.Internal,
			&lotusApi.Internal,
		},
		http.Header{})
	if err != nil {
		t.Errorf("Can't construct a Lotus client, error: %s", err.Error())
	}
	defer clientClose()

	head, err := lotusApi.ChainHead(context.Background())
	if err != nil {
		t.Errorf("Can't call method ChainHead of Lotus API, error: %s", err.Error())
	}
	assert.Greater(t, len(head.Cids()), 0)
	assert.NotEqualf(t, "", head.Cids()[0].KeyString(), "Head CID Key is empty")
}

func TestLotusMinerConnectivityWs(t *testing.T) {
	var lotusApi apistruct.StorageMinerStruct
	bgCtx := context.Background()
	ctx, _ := context.WithTimeout(bgCtx, time.Minute*3)

	clientClose, err := jsonrpc.NewMergeClient(
		ctx,
		"ws://lotus-full-node:2345/rpc/v0",
		"Filecoin",
		[]interface{}{
			&lotusApi.CommonStruct.Internal,
			&lotusApi.Internal,
		},
		http.Header{})
	if err != nil {
		t.Errorf("Can't construct a Lotus client, error: %s", err.Error())
	}
	defer clientClose()

	id, err := lotusApi.ID(context.Background())
	if err != nil {
		t.Errorf("Can't call method ChainHead of Lotus API, error: %s", err.Error())
	}
	assert.NotNil(t, id)
}
