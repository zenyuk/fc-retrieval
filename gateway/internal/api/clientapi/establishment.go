package clientapi

// Copyright (C) 2020 ConsenSys Software Inc
import (
	"encoding/json"
	"net/http"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/util"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/messages"
	"github.com/ant0ine/go-json-rest/rest"
)

// HandleClientNetworkEstablishment is used to handle initial establishment http request from client
func (g *ClientAPI) HandleClientNetworkEstablishment(w rest.ResponseWriter, content []byte) {
	payload := messages.ClientEstablishmentRequest{}
	err := json.Unmarshal(content, &payload)
	if err != nil {
		s := "Client Establishment: Failed to decode payload."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusBadRequest)
		return
	}

	logging.Trace("Client Establishment %+v", payload)

	now := util.GetTimeImpl().Now().Unix()
	if payload.TTL > now {
		// TODO how to just drop the connection?

	}

	response, err := g.gateway.GatewayClient.Establishment(&payload)
	if err != nil {
		s := "Client Establishment: Error decodeing payload."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusBadRequest)
		return
	}

	response.ProtocolVersion = clientAPIProtocolVersion
	response.GatewayID = g.gateway.GatewayID.ToString()

	sig, err1 := fcrcrypto.Sign(g.gateway.GatewayPrivateKey, *g.gateway.GatewayPrivateKeyVersion, *g.gateway.GatewayPrivateKeySigAlg, *response)
	if err1 != nil {
		// TODO for the moment just blow up!
		panic(err1)
		// s := "Client Establishment: Internal error signing."
		// logging.Error(s + err.Error())
		// rest.Error(w, s, http.StatusInternalServerError)

	}
	response.Signature = *sig
	w.WriteJson(response)
}
