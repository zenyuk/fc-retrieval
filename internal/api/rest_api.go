package api

// Copyright (C) 2020 ConsenSys Software Inc

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/util"
	"github.com/ant0ine/go-json-rest/rest"
)

// StartRestAPI starts the REST API as a separate go routine.
// Any start-up errors are returned.
func StartRestAPI(settings util.AppSettings) error {
	// Start the REST API and block until the error code is set.
	errChan := make(chan error, 1)
	go startRestAPI(settings, errChan)
	return <-errChan
}

func startRestAPI(settings util.AppSettings, errChannel chan<- error) {

	// Initialise a dummy gateway instance.
	g := Gateway{ProtocolVersion: 1, ProtocolSupported: []int{1, 2}}

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/time", getTime),     // Get system time.
		rest.Get("/ip", getIP),         // Get IP address.
		rest.Get("/host", getHostname), // Get host name.
		rest.Post("/client/establishment", g.HandleClientNetworkEstablishment),       // Handle network establishment.
		rest.Post("/client/standard_request_cid", g.HandleClientStandardCIDDiscover), // Handle client standard cid request.
		rest.Post("/client/dht_request_cid", g.HandleClientDHTCIDDiscover),           // Handle DHT client cid request.
	)
	if err != nil {
		log.Fatal(err)
		errChannel <- err
		return
	}

	log.Println("Running REST API on: " + settings.BindRestAPI)
	api.SetApp(router)
	errChannel <- nil
	log.Fatal(http.ListenAndServe(":"+settings.BindRestAPI, api.MakeHandler()))
}

func getTime(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(time.Now())
}

func getHostname(w rest.ResponseWriter, r *rest.Request) {
	name, err := os.Hostname()
	if err != nil {
		log.Printf("Get host name1: %v\n", err)
		return
	}

	w.WriteJson(name)
}

func getIP(w rest.ResponseWriter, r *rest.Request) {
	name, err := os.Hostname()
	if err != nil {
		log.Printf("Get host name2: %v\n", err)
		return
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		log.Printf("Lookup host: %v\n", err)
		return
	}

	w.WriteJson(addrs)
}
