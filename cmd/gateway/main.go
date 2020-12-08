package main

import (
	"log"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/api"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
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

	_, err = gateway.Create(settings)
	if err != nil {
		log.Println("Error starting server: Client REST API: " + err.Error())
		return
	}


	// Initialise a dummy gateway instance.
	g1 := api.Gateway{ProtocolVersion: 1, ProtocolSupported: []int{1, 2}}

	err = api.StartTCPAPI(settings, &g1)
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
