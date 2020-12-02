package api

import (
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
)

// HandleClientNetworkEstablishment is used to handle initial establishment http request from client
func (g *Gateway) HandleClientNetworkEstablishment(w rest.ResponseWriter, r *rest.Request) {
	payload := ClientEstablishmentRequest{}
	err := r.DecodeJsonPayload(&payload)
	if err != nil {
		log.Println(err.Error())
		rest.Error(w, "Fail to decode payload.", http.StatusBadRequest)
		return
	}
	// TODO: For now just print the payload
	log.Println(payload)

	// Dummy response
	response := ClientEstablishmentResponse{}
	response.CommonFields.ProtocolVersion = g.ProtocolVersion
	response.CommonFields.ProtocolSupported = g.ProtocolSupported
	w.WriteJson(response)
}

// HandleClientStandardCIDDiscover is used to handle client request for cid offer
func (g *Gateway) HandleClientStandardCIDDiscover(w rest.ResponseWriter, r *rest.Request) {
	payload := ClientStandardDiscoverRequest{}
	err := r.DecodeJsonPayload(&payload)
	if err != nil {
		log.Println(err.Error())
		rest.Error(w, "Fail to decode payload.", http.StatusBadRequest)
		return
	}
	// TODO: For now just print the payload
	log.Println(payload)

	// Dummy response
	response := ClientStandardDiscoverResponse{}
	response.CommonFields.ProtocolVersion = g.ProtocolVersion
	response.CommonFields.ProtocolSupported = g.ProtocolSupported
	w.WriteJson(response)
}

// HandleClientDHTCIDDiscover is used to handle client request for cid offer
func (g *Gateway) HandleClientDHTCIDDiscover(w rest.ResponseWriter, r *rest.Request) {
	payload := ClientDHTDiscoverRequest{}
	err := r.DecodeJsonPayload(&payload)
	if err != nil {
		log.Println(err.Error())
		rest.Error(w, "Fail to decode payload.", http.StatusBadRequest)
		return
	}
	// TODO: For now just print the payload
	log.Println(payload)

	// Dummy response
	response := ClientDHTDiscoverResponse{}
	response.CommonFields.ProtocolVersion = g.ProtocolVersion
	response.CommonFields.ProtocolSupported = g.ProtocolSupported
	w.WriteJson(response)
}
