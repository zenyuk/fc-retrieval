package clientapi

// Copyright (C) 2020 ConsenSys Software Inc
import (
	"net/http"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/util"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/messages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ant0ine/go-json-rest/rest"
)

// HandleClientNetworkEstablishment is used to handle initial establishment http request from client
func (g *ClientAPI) HandleClientNetworkEstablishment(w rest.ResponseWriter, r *rest.Request) {
	payload := messages.ClientEstablishmentRequest{}
	err := r.DecodeJsonPayload(&payload)
	if err != nil {
		s := "Client Establishment: Failed to decode payload." + err.Error()
		logging.Error(s)
		rest.Error(w, s, http.StatusBadRequest)
		return
	}
	// TODO: For now just print the payload
	logging.Info("Payload %+v", payload)

	now := util.GetTimeImpl().Now().Unix()
	if payload.TTL > now {
		// TODO how to just drop the connection?

	}

	response := messages.ClientEstablishmentResponse{}
	response.ProtocolVersion = clientAPIProtocolVersion
	response.Challenge = payload.Challenge
	response.Signature = "TODO: NONE YET!!!"
	w.WriteJson(response)
}
