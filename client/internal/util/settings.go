package util
// Copyright (C) 2020 ConsenSys Software Inc

import (
"log"
"io/ioutil"
"encoding/json"
)

const settingsFileName = "settings.json"
const settingsLocDev = settingsFileName
const settingsLocContainer = "/etc/client/" + settingsFileName


const settingsDefaultBindRestAPI = "8080"
const settingsDefaultVerbose = true

// AppSettings defines the server configuraiton
type AppSettings struct {
	BindRestAPI string // Port number to bind to for REST API.
	Verbose bool       // If true, then more logging is shown.
}

var defaults = AppSettings{settingsDefaultBindRestAPI, settingsDefaultVerbose}

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

