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

// AppSettings defines the server configuraiton
type AppSettings struct {
	BindRestAPI     string `json:"bindrestapi"`     // Port number to bind to for client REST API.
	BindProviderAPI string `json:"bindproviderapi"` // Port number to bind to for provider TCP communication API.
	BindGatewayAPI  string `json:"bindgatewayapi"`  // Port number to bind to for gateway TCP communication API.
	LogLevel        string `json:"loglevel"`        // Log Level: NONE, ERROR, WARN, INFO, TRACE
	LogTarget       string `json:"logtarget"`       // Log Level: STDOUT
}

var defaults = AppSettings{
	settingsDefaultBindRestAPI,
	settingsDefaultBindProviderAPI,
	settingsDefaultBindGatewayAPI,
	settingsDefaultLogLevel,
	settingsDefaultLogTarget,
}

// TODO at present there is no way to get this global object. Do we need this?
var settings = defaults

// LoadSettings loads the app settings from the settings file.
func LoadSettings() (set AppSettings) {
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

// TODO if we want to have time and date in the output we should use this code.

// GetTimeNowString returns the time now in a standard format.
// func GetTimeNowString() string {
// 	t := time.Now()
// 	return t.Format("2006-01-02 15:04:05")
// }

