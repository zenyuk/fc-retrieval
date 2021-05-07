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
	"encoding/json"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	tc "github.com/ConsenSys/fc-retrieval-itest/pkg/test-containers"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	composeID, err := tc.StartContainers()
	if err != nil {
		logging.Error("Can't start containers %s", err.Error())
		os.Exit(1)
	}
	defer tc.StopContainers(composeID)
	m.Run()
}

func TestLotusConnectivity(t *testing.T) {

	method := "Filecoin.ChainHead"
	id := 1
	lotusUrl := "http://127.0.0.1:1234/rpc/v0"
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
