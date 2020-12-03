package api

import (
	"log"
	"net"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/util"
)

// StartTCPAPI starts two go-routines.
// One for handling communication with provider.
// The other for handling communication with gateway.
func StartTCPAPI(settings util.AppSettings, g *Gateway) error {
	// For handling requests from provider.
	ln, err := net.Listen("tcp", settings.BindProviderAPI)
	if err != nil {
		return err
	}
	go func(ln net.Listener) {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err.Error())
				continue
			}
			log.Printf("Incoming connection from provider at: %s", conn.RemoteAddr())
			go handleProviderCommunication(conn, g)
		}
	}(ln)
	log.Println("Running provider communication API on: " + settings.BindProviderAPI)

	// For handling requests from gateway.
	ln2, err := net.Listen("tcp", settings.BindGatewayAPI)
	if err != nil {
		return err
	}
	go func(ln net.Listener) {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err.Error())
				continue
			}
			log.Printf("Incoming connection from gateway at: %s", conn.RemoteAddr())
			go handleGatewayCommunication(conn, g)
		}
	}(ln2)
	log.Println("Running gateway communication API on: " + settings.BindGatewayAPI)

	return nil
}
