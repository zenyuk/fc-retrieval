package clientapi

// Copyright (C) 2020 ConsenSys Software Inc

// Contains debug APIs

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
)

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
