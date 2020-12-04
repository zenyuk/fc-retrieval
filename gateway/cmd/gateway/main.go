package main

import (
	"log"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/api"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util"
)

func main() {
	log.Println("Filecoin Gateway Start-up Start: " + util.GetTimeNowString())

	// Load settings.
	settings, err := util.LoadSettings()
	if err != nil {
		log.Println("Error starting server: Settings: " + err.Error())
		return
	}

	// Initialise a dummy gateway instance.
	g := api.Gateway{ProtocolVersion: 1, ProtocolSupported: []int{1, 2}}

	// Set-up the REST API
	err = api.StartRestAPI(settings, &g)
	if err != nil {
		log.Println("Error starting server: REST API: " + err.Error())
		return
	}
	err = api.StartTCPAPI(settings, &g)
	if err != nil {
		log.Println("Error starting tcp server: " + err.Error())
		return
	}

	// Configure what should be called if Control-C is hit.
	util.SetUpCtrlCExit(gracefulExit)
	log.Println("Filecoin Gateway Start-up Done: " + util.GetTimeNowString())

	// Wait forever.
	select {}
}

func gracefulExit() {
	log.Println("Filecoin Gateway Start: " + util.GetTimeNowString())

	// TODO

	log.Println("Filecoin Gateway Shutdown End: " + util.GetTimeNowString())
}
