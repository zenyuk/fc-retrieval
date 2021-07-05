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

	"github.com/testcontainers/testcontainers-go"

	"github.com/ConsenSys/fc-retrieval/itest/config"
	tc "github.com/ConsenSys/fc-retrieval/itest/pkg/util/test-containers"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/api/apistruct"
	"github.com/stretchr/testify/assert"

	"github.com/ConsenSys/fc-retrieval/common/pkg/logging"
)

var containers tc.AllContainers

func TestMain(m *testing.M) {
	const testName = "lotus-connectivity"
	ctx := context.Background()
	var gatewayConfig = config.NewConfig(".env.gateway")
	var providerConfig = config.NewConfig(".env.provider")
	var registerConfig = config.NewConfig(".env.register")
	var network *testcontainers.Network
	var err error
	containers, network, err = tc.StartContainers(ctx, 1, 1, testName, true, gatewayConfig, providerConfig, registerConfig)
	if err != nil {
		logging.Error("%s test failed, container starting error: %s", testName, err.Error())
		tc.StopContainers(ctx, testName, containers, network)
		os.Exit(1)
	}
	defer tc.StopContainers(ctx, testName, containers, network)
	m.Run()
}

func TestLotusFullNodeConnectivityHttp(t *testing.T) {
	method := "Filecoin.ChainHead"
	id := 1
	lotusDaemonEndpoint, _ := containers.Lotus.GetLostHostApiEndpoints()
	lotusUrl := "http://" + lotusDaemonEndpoint + "/rpc/v0"
	requestBody := `{
		"jsonrpc": "2.0",
		"method": "` + method + `",
		"id": ` + strconv.Itoa(id) + `,
		"params": []
	}`

	resp, err := http.Post(lotusUrl, "application/json", bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		t.Errorf("Can't establish connection to Lotus %s", lotusUrl)
		t.FailNow()
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Can't read response from Lotus, method: %s", method)
		t.FailNow()
	}
	var fields map[string]interface{}
	err = json.Unmarshal(body, &fields)
	if err != nil {
		t.Errorf("Can't parse json response from Lotus, method: %s", method)
		t.FailNow()
	}

	assert.Equal(t, float64(id), fields["id"])
}

func TestLotusFullNodeConnectivityWs(t *testing.T) {
	var lotusApi apistruct.FullNodeStruct
	bgCtx := context.Background()
	ctx, _ := context.WithTimeout(bgCtx, time.Minute*3)

	lotusDaemonEndpoint, _ := containers.Lotus.GetLostHostApiEndpoints()
	clientClose, err := jsonrpc.NewMergeClient(
		ctx,
		"ws://"+lotusDaemonEndpoint+"/rpc/v0",
		"Filecoin",
		[]interface{}{
			&lotusApi.CommonStruct.Internal,
			&lotusApi.Internal,
		},
		http.Header{})
	if err != nil {
		t.Errorf("Can't construct a Lotus client, error: %s", err.Error())
		t.FailNow()
	}
	defer clientClose()

	head, err := lotusApi.ChainHead(context.Background())
	if err != nil {
		t.Errorf("Can't call method ChainHead of Lotus API, error: %s", err.Error())
		t.FailNow()
	}
	assert.Greater(t, len(head.Cids()), 0)
	assert.NotEqualf(t, "", head.Cids()[0].KeyString(), "Head CID Key is empty")
}

func TestLotusMinerConnectivityWs(t *testing.T) {
	var lotusApi apistruct.StorageMinerStruct
	bgCtx := context.Background()
	ctx, _ := context.WithTimeout(bgCtx, time.Minute*3)

	_, lotusMinerEndpoint := containers.Lotus.GetLostHostApiEndpoints()
	clientClose, err := jsonrpc.NewMergeClient(
		ctx,
		"ws://"+lotusMinerEndpoint+"/rpc/v0",
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
		t.FailNow()
	}
	assert.NotNil(t, id)
}
