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
	"io/ioutil"
	"net/http"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages/fcrmsgbasic"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-provider/internal/core"
	"github.com/ConsenSys/fc-retrieval-provider/internal/util/settings"
	"github.com/ant0ine/go-json-rest/rest"
)

// StartAdminRestAPI starts the REST API as a separate go routine.
// Any start-up errors are returned.
func StartAdminRestAPI(settings settings.AppSettings) error {
	// Start the REST API and block until the error code is set.
	errChan := make(chan error, 1)
	go startRestAPI(settings, errChan)
	return <-errChan
}

func startRestAPI(settings settings.AppSettings, errChannel chan<- error) {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		// TODO: Remove these debug APIs prior to production release.
		// rest.Get("/time", getTime),     // Get system time.
		// rest.Get("/ip", getIP),         // Get IP address.
		// rest.Get("/host", getHostname), // Get host name.
		rest.Post("/v1", msgRouter),
	)
	if err != nil {
		logging.Error1(err)
		errChannel <- err
		return
	}
	logging.Info("Running Admin API on: %s", settings.BindAdminAPI)
	api.SetApp(router)
	errChannel <- nil
	logging.Error(http.ListenAndServe(":"+settings.BindAdminAPI, api.MakeHandler()).Error())
	panic("Error binding")
}

func msgRouter(w rest.ResponseWriter, r *rest.Request) {
	// Get core structure
	c := core.GetSingleInstance()

	logging.Trace("Received request via /v1 API")
	content, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		logging.Error("Error reading request: %s.", err.Error())
		rest.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}
	if len(content) == 0 {
		logging.Error("Error empty request")
		rest.Error(w, "Error empty request", http.StatusBadRequest)
		return
	}

	request, err := fcrmessages.FCRMsgFromBytes(content)
	if err != nil {
		logging.Error("Failed to decode payload: %s.", err.Error())
		rest.Error(w, "Failed to decode payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	checked := checkProtocol(w, request, c)
	if checked == false {
		return
	}

	switch request.GetMessageType() {
	case fcrmessages.ProviderAdminPublishGroupOfferRequestType:
		handleProviderPublishGroupCID(w, request)
	case fcrmessages.ProviderAdminGetPublishedOfferRequestType:
		handleProviderGetGroupCID(w, request)
	case fcrmessages.ProviderAdminInitialiseKeyRequestType:
		handleKeyManagement(w, request)
	default:
		logging.Warn("Client Request: Unknown message type: %d", request.GetMessageType())
		rest.Error(w, "Unknown message type", http.StatusBadRequest)
	}
}

func checkProtocol(w rest.ResponseWriter, request *fcrmessages.FCRMessage, c *core.Core) bool {
	protocolVersion := c.ProtocolVersion
	protocolSupported := c.ProtocolSupported
	// Only process the rest of the message if the protocol version is understood.
	if request.GetProtocolVersion() != protocolVersion {
		// Check to see if the client supports the gateway's preferred version
		for _, clientProvVer := range request.GetProtocolSupported() {
			if clientProvVer == protocolVersion {
				// Request the client switch to this protocol version
				// TODO what can we get from request object?
				logging.Info("Requesting client (TODO) switch protocol versions from %d to %d", request.GetProtocolVersion(), protocolVersion)
				response, _ := fcrmsgbasic.EncodeProtocolChangeRequest(protocolVersion)
				w.WriteJson(response)
				return false
			}
		}

		// Go through the protocol versions supported by the client and the
		// gateway to search for any common version, prioritising
		// the gateway preference over the client preference.
		for _, clientProvVer := range request.GetProtocolSupported() {
			for _, gatewayProtVer := range protocolSupported {
				if clientProvVer == gatewayProtVer {
					// When we support more than one version of the protocol, this code will change the gateway
					// to using the other (common version)
					logging.Error("Not implemented yet")
					panic("Multiple protocol versions not implemented yet")
				}
			}
		}
		// No common protocol versions supported.
		// TODO what can we get from request object?
		logging.Warn("Client Request: Unsupported protocol version(s): %d", request.GetProtocolVersion())
		response, _ := fcrmsgbasic.EncodeProtocolChangeResponse(false)
		w.WriteJson(response)
		return false
	}
	return true
}
