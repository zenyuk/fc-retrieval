package adminapi

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
  "net/http"

  "github.com/ant0ine/go-json-rest/rest"

  "github.com/ConsenSys/fc-retrieval/common/pkg/fcrmessages"
  "github.com/ConsenSys/fc-retrieval/common/pkg/logging"
  "github.com/ConsenSys/fc-retrieval/gateway/internal/core"
)

// HandleGatewayAdminGetReputationRequest handles admin get reputation request
func HandleGatewayAdminGetReputationRequest(w rest.ResponseWriter, request *fcrmessages.FCRMessage) {
	// Get core structure
	c := core.GetSingleInstance()

	if c.GatewayPrivateKey == nil {
		s := "This gateway hasn't been initialised by the admin"
		logging.Error(s)
		rest.Error(w, s, http.StatusBadRequest)
		return
	}

	clientID, err := fcrmessages.DecodeGatewayAdminGetReputationRequest(request)
	if err != nil {
		s := "Fail to decode message."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusBadRequest)
		return
	}

	// Get reputation db
	reputation, exists := c.ReputationMgr.GetClientReputation(clientID)

	// Construct message
	response, err := fcrmessages.EncodeGatewayAdminGetReputationResponse(clientID, reputation, exists)
	if err != nil {
		s := "Internal error: Fail to encode message."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusInternalServerError)
		return
	}
	// Sign message
	err = response.Sign(c.GatewayPrivateKey, c.GatewayPrivateKeyVersion)
	if err != nil {
		s := "Internal error: Fail to sign message."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusInternalServerError)
		return
	}
  if err := w.WriteJson(response); err != nil {
    logging.Error("can't write JSON during HandleGatewayAdminGetReputationRequest %s", err.Error())
  }
}
