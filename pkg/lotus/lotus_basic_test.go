package client_gateway

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

	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-itest/pkg/util"
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
	util.CleanContainers()
	// We need a lotus
	tag := util.GetCurrentBranch()
	network := "itest-shared"

	// Create shared net
	ctx := context.Background()
	net := *util.CreateNetwork(ctx, network)
	defer net.Remove(ctx)

	// Start lotus
	lotus := *util.StartLotus(ctx, network, true)
	defer lotus.Terminate(ctx)
	defer lotus.StopLogProducer() // Call when verbose is set to true

	// Start itest
	done := make(chan bool)
	itest := *util.StartItest(ctx, tag, network, util.ColorGreen, "./pkg/lotus", done)
	defer itest.Terminate(ctx)
	defer itest.StopLogProducer()

	// Block until done.
	if <-done {
		logging.Info("Tests passed, shutdown...")
	} else {
		logging.Fatal("Tests failed, shutdown...")
	}
}

func TestLotusConnectivity(t *testing.T) {

	method := "Filecoin.ChainHead"
	id := 1
	lotusUrl := "http://lotus:1234/rpc/v0"
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
