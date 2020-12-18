package clientapi

// Copyright (C) 2020 ConsenSys Software Inc
import (
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/messages"
)

// HandleClientDHTCIDDiscover is used to handle client request for cid offer
func (c *ClientAPI) HandleClientDHTCIDDiscover(w rest.ResponseWriter, r *rest.Request) {
	payload := messages.ClientDHTDiscoverRequest{}
	err := r.DecodeJsonPayload(&payload)
	if err != nil {
		s := "Client DHT CID Discovery: Failed to decode payload."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusBadRequest)
		return
	}

	logging.Trace("Client DHT CID Discovery %+v", payload)

	// Dummy response
	response := messages.ClientDHTDiscoverResponse{MessageType: messages.ClientDHTDiscoverRequestType}
	w.WriteJson(response)
}
