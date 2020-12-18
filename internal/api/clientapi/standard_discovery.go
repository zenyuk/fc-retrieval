package clientapi

// Copyright (C) 2020 ConsenSys Software Inc
import (
	"net/http"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/messages"
	"github.com/ant0ine/go-json-rest/rest"
)

// HandleClientStandardCIDDiscover is used to handle client request for cid offer
func (g *ClientAPI) HandleClientStandardCIDDiscover(w rest.ResponseWriter, r *rest.Request) {
	payload := messages.ClientStandardDiscoverRequest{}
	err := r.DecodeJsonPayload(&payload)
	if err != nil {
		s := "Client Standard CID Discovery: Failed to decode payload."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusBadRequest)
		return
	}

	logging.Trace("Client Standard CID Discovery %+v", payload)

	// Dummy response
	response := messages.ClientStandardDiscoverResponse{MessageType: messages.ClientStandardDiscoverResponseType}
	w.WriteJson(response)
}
