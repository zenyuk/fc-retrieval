package util

// Copyright (C) 2020 ConsenSys Software Inc

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

const settingsFileName = "settings.json"
const settingsLocDev = settingsFileName
const settingsLocContainer = "/etc/gateway/" + settingsFileName

const settingsDefaultBindRestAPI = "8080"
const settingsDefaultBindProviderAPI = "8090"
const settingsDefaultBindGatewayAPI = "8091"
const settingsDefaultVerbose = true

// AppSettings defines the server configuraiton
type AppSettings struct {
	BindRestAPI     string `json:"bindrestapi"`     // Port number to bind to for client REST API.
	BindProviderAPI string `json:"bindproviderapi"` // Port number to bind to for provider TCP communication API.
	BindGatewayAPI  string `json:"bindgatewayapi"`  // Port number to bind to for gateway TCP communication API.
	Verbose         bool   `json:"verbose"`         // If true, then more logging is shown.
}

var defaults = AppSettings{
	settingsDefaultBindRestAPI,
	settingsDefaultBindProviderAPI,
	settingsDefaultBindGatewayAPI,
	settingsDefaultVerbose}

var settings = defaults

// LoadSettings loads the app settings from the settings file.
func LoadSettings() (set AppSettings, err error) {
	// Load settings.
	settingsBytes, err := ioutil.ReadFile(settingsLocContainer)
	if err != nil {
		settingsBytes, err = ioutil.ReadFile(settingsLocDev)
		if err != nil {
			log.Println("Failed to read settings.json" + err.Error())
			return
		}
	}

	err = json.Unmarshal(settingsBytes, &settings)
	if err != nil {
		log.Println("Failed to parse settings.json: " + err.Error())
	}

	if settings.Verbose {
		log.Printf("Settings: (%+v)\n", settings)
	}

	set = settings
	return
}
