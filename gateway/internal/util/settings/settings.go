package settings

// Copyright (C) 2020 ConsenSys Software Inc

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
)

const settingsFileName = "settings.json"
const settingsLocDev = settingsFileName
const settingsLocContainer = "/etc/gateway/" + settingsFileName

const settingsDefaultBindRestAPI = "8080"
const settingsDefaultBindProviderAPI = "8090"
const settingsDefaultBindGatewayAPI = "8091"
const settingsDefaultLogLevel = "TRACE"
const settingsDefaultLogTarget = "STDOUT"

// TODO id doesn't make sense to have defaults for these values.
const settingsDefaultGatewayID = "01"
const settingsDefaultPrivateKey = "01"
const settingsDefaultPrivKeyVer = 0xff
const settingsDefaultPrivKeySigAlg = 0xff

// DefaultTCPInactivityTimeoutMs is the default timeout for TCP inactivity
const DefaultTCPInactivityTimeoutMs = 100

// AppSettings defines the server configuraiton
type AppSettings struct {
	BindRestAPI     string `json:"bindrestapi"`     // Port number to bind to for client REST API.
	BindProviderAPI string `json:"bindproviderapi"` // Port number to bind to for provider TCP communication API.
	BindGatewayAPI  string `json:"bindgatewayapi"`  // Port number to bind to for gateway TCP communication API.
	LogLevel        string `json:"loglevel"`        // Log Level: NONE, ERROR, WARN, INFO, TRACE
	LogTarget       string `json:"logtarget"`       // Log Level: STDOUT
	GatewayID       string `json:"gatewayid"`       // Node id of this gateway
	GatewayPrivKey  string `json:"privatekey"`      // Gateway private key
	GatewayPrivKeyVersion uint8  `json:"keyversion"`      // Key version of gateway private key
	GatewaySigAlg   uint8  `json:"sigalg"`          // Signature algorithm to be used by private key.
}

var defaults = AppSettings{
	settingsDefaultBindRestAPI,
	settingsDefaultBindProviderAPI,
	settingsDefaultBindGatewayAPI,
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
func LoadSettings() (AppSettings) {
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


