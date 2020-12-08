package clientapi
// Copyright (C) 2020 ConsenSys Software Inc
import (
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/api"
)

// HandleClientDHTCIDDiscover is used to handle client request for cid offer
func (c *ClientAPI) HandleClientDHTCIDDiscover(w rest.ResponseWriter, r *rest.Request) {
	payload := api.ClientDHTDiscoverRequest{}
	err := r.DecodeJsonPayload(&payload)
	if err != nil {
		log.Println(err.Error())
		rest.Error(w, "Fail to decode payload.", http.StatusBadRequest)
		return
	}
	// TODO: For now just print the payload
	log.Println(payload)

	// Dummy response
	response := api.ClientDHTDiscoverResponse{}
	response.CommonFields.ProtocolVersion = clientAPIProtocolVersion
	response.CommonFields.ProtocolSupported = clientAPIProtocolSupported
	w.WriteJson(response)
}
