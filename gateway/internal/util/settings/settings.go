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
	"time"
)

const settingsFileName = "settings.json"
const settingsLocDev = settingsFileName
const settingsLocContainer = "/etc/gateway/" + settingsFileName

const settingsDefaultBindRestAPI = "9010"
const settingsDefaultBindProviderAPI = "9011"
const settingsDefaultBindGatewayAPI = "9012"
const settingsDefaultBindAdminAPI = "9013"
const settingsDefaultLogLevel = "TRACE"
const settingsDefaultLogTarget = "STDOUT"
const settingsDefaultLogDir = "/var/log/fc-retrieval/fc-retrieval-gateway"
const settingsDefaultLogFile = "gateway.log"
const settingsDefaultLogMaxBackups = 3
const settingsDefaultLogMaxAge = 28
const settingsDefaultLogMaxSize = 500
const settingsDefaultLogCompress = false

// TODO id doesn't make sense to have defaults for these values.
const settingsDefaultGatewayID = "01"

const settingsDefaultRegisterAPIURL = "http://localhost:9020"
const settingsDefaultGatewayAddress = "f0121345"
const settingsDefaultNetworkInfoClient = "127.0.0.1:9010"
const settingsDefaultNetworkInfoProvider = "127.0.0.1:9011"
const settingsDefaultNetworkInfoGateway = "127.0.0.1:9012"
const settingsDefaultNetworkInfoAdmin = "127.0.0.1:9013"
const settingsDefaultGatewayRegionCode = "US"
const settingsDefaultGatewayRootSigningKey = "0xABCDE123456789"
const settingsDefaultGatewaySigningKey = "0x987654321EDCBA"

// DefaultTCPInactivityTimeout is the default timeout for TCP inactivity
const DefaultTCPInactivityTimeout = 100 * time.Millisecond

// DefaultLongTCPInactivityTimeout is the default timeout for long TCP inactivity. This timeout should never be ignored.
const DefaultLongTCPInactivityTimeout = 5000 * time.Millisecond

// AppSettings defines the server configuraiton
type AppSettings struct {
	BindRestAPI     string `mapstructure:"BIND_REST_API"`     // Port number to bind to for client REST API.
	BindProviderAPI string `mapstructure:"BIND_PROVIDER_API"` // Port number to bind to for provider TCP communication API.
	BindGatewayAPI  string `mapstructure:"BIND_GATEWAY_API"`  // Port number to bind to for gateway TCP communication API.
	BindAdminAPI    string `mapstructure:"BIND_ADMIN_API"`    // Port number to bind to for admin TCP communication API.
	LogLevel        string `mapstructure:"LOG_LEVEL"`         // Log Level: NONE, ERROR, WARN, INFO, TRACE
	LogTarget       string `mapstructure:"LOG_TARGET"`        // Log Level: STDOUT
	LogDir          string `mapstructure:"LOG_DIR"`           // Log Dir: /var/log/fc-retrieval/fc-retrieval-gateway
	LogFile         string `mapstructure:"LOG_FILE"`          // Log File: gateway.log
	LogMaxBackups   int    `mapstructure:"LOG_MAX_BACKUPS"`   // Log max backups: 3
	LogMaxAge       int    `mapstructure:"LOG_MAX_AGE"`       // Log max age (days): 28
	LogMaxSize      int    `mapstructure:"LOG_MAX_SIZE"`      // Log max size (MB): 500
	LogCompress     bool   `mapstructure:"LOG_COMPRESS"`      // Log compress: false
	GatewayID       string `mapstructure:"GATEWAY_ID"`        // Node id of this gateway

	RegisterAPIURL        string `mapstructure:"REGISTER_API_URL"`         // Register service url
	GatewayAddress        string `mapstructure:"GATEWAY_ADDRESS"`          // Gateway address
	NetworkInfoGateway    string `mapstructure:"GATEWAY_NETWORK_INFO"`     // Gateway network info
	GatewayRegionCode     string `mapstructure:"GATEWAY_REGION_CODE"`      // Gateway region code
	GatewayRootSigningKey string `mapstructure:"GATEWAY_ROOT_SIGNING_KEY"` // Gateway root signing key
	GatewaySigningKey     string `mapstructure:"GATEWAY_SIGNING_KEY"`      // Gateway signing key

	NetworkInfoClient   string `mapstructure:"CLIENT_NETWORK_INFO"`   // Gateway client network info
	NetworkInfoProvider string `mapstructure:"PROVIDER_NETWORK_INFO"` // Gateway provider network info
	NetworkInfoAdmin    string `mapstructure:"ADMIN_NETWORK_INFO"`    // Gateway admin network info
}

var defaults = AppSettings{
	settingsDefaultBindRestAPI,
	settingsDefaultBindProviderAPI,
	settingsDefaultBindGatewayAPI,
	settingsDefaultBindAdminAPI,
	settingsDefaultLogLevel,
	settingsDefaultLogTarget,
	settingsDefaultLogDir,
	settingsDefaultLogFile,
	settingsDefaultLogMaxBackups,
	settingsDefaultLogMaxAge,
	settingsDefaultLogMaxSize,
	settingsDefaultLogCompress,
	settingsDefaultGatewayID,

	settingsDefaultRegisterAPIURL,
	settingsDefaultGatewayAddress,
	settingsDefaultNetworkInfoGateway,
	settingsDefaultGatewayRegionCode,
	settingsDefaultGatewayRootSigningKey,
	settingsDefaultGatewaySigningKey,

	settingsDefaultNetworkInfoClient,
	settingsDefaultNetworkInfoProvider,
	settingsDefaultNetworkInfoAdmin,
}
