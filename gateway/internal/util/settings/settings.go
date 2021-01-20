package settings

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
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
)

const settingsFileName = "settings.json"
const settingsLocDev = settingsFileName
const settingsLocContainer = "/etc/gateway/" + settingsFileName

const settingsDefaultBindRestAPI = "8080"
const settingsDefaultBindProviderAPI = "8090"
const settingsDefaultBindGatewayAPI = "8091"
const settingsDefaultBindAdminAPI = "8092"
const settingsDefaultLogLevel = "TRACE"
const settingsDefaultLogTarget = "STDOUT"

// TODO id doesn't make sense to have defaults for these values.
const settingsDefaultGatewayID = "01"
const settingsDefaultPrivateKey = "01"
const settingsDefaultPrivKeyVer = 0xff
const settingsDefaultPrivKeySigAlg = 0xff

// DefaultTCPInactivityTimeout is the default timeout for TCP inactivity
const DefaultTCPInactivityTimeout = 100 * time.Millisecond

// DefaultLongTCPInactivityTimeout is the default timeout for long TCP inactivity. This timeout should never be ignored.
const DefaultLongTCPInactivityTimeout = 5000 * time.Millisecond

// AppSettings defines the server configuraiton
type AppSettings struct {
	BindRestAPI     			string `mapstructure:"BIND_REST_API"`     		// Port number to bind to for client REST API.
	BindProviderAPI 			string `mapstructure:"BIND_PROVIDER_API"` 		// Port number to bind to for provider TCP communication API.
	BindGatewayAPI  			string `mapstructure:"BIND_GATEWAY_API"`  		// Port number to bind to for gateway TCP communication API.
	BindAdminAPI          string `mapstructure:"BIND_ADMIN_API"`    		// Port number to bind to for admin TCP communication API.
	LogLevel        			string `mapstructure:"LOG_LEVEL"`        			// Log Level: NONE, ERROR, WARN, INFO, TRACE
	LogTarget       			string `mapstructure:"LOG_TARGET"`       			// Log Level: STDOUT
	GatewayID       			string `mapstructure:"GATEWAY_ID"`       			// Node id of this gateway
	GatewayPrivKey  			string `mapstructure:"GATEWAY_PRIVATE_KEY"`		// Gateway private key
	GatewayKeyVersion 		uint32 `mapstructure:"GATEWAY_KEY_VERSION"`   // Key version of gateway private key
	GatewaySigAlg   			uint8  `mapstructure:"GATEWAY_SIG_ALG"`       // Signature algorithm to be used by private key.
}

var defaults = AppSettings{
	settingsDefaultBindRestAPI,
	settingsDefaultBindProviderAPI,
	settingsDefaultBindGatewayAPI,
	settingsDefaultBindAdminAPI,
	settingsDefaultLogLevel,
	settingsDefaultLogTarget,
	settingsDefaultGatewayID,
	settingsDefaultPrivateKey,
	settingsDefaultPrivKeyVer,
	settingsDefaultPrivKeySigAlg,
}

// TODO at present there is no way to get this global object. Do we need this?
var settings = defaults

// LoadSettings loads the app settings from the settings file.
func LoadSettings() AppSettings {
	// Load settings.
	settingsBytes, err := ioutil.ReadFile(settingsLocContainer)
	if err != nil {
		settingsBytes, err = ioutil.ReadFile(settingsLocDev)
		if err != nil {
			e := "Failed to read settings.json" + err.Error()
			logging.Error(e)
			panic(e)
		}
	}

	err = json.Unmarshal(settingsBytes, &settings)
	if err != nil {
		e := ("Failed to parse settings.json: " + err.Error())
		logging.Error(e)
		panic(e)
	}

	logging.SetLogLevel(settings.LogLevel)
	logging.SetLogTarget(settings.LogTarget)
	logging.Info("Settings: (%+v)\n", settings)

	return settings
}
