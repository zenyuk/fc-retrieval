package adminapi

import (
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages/fcrmsgpvdadmin"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-provider/internal/core"
	"github.com/ant0ine/go-json-rest/rest"
)

func handleKeyManagement(w rest.ResponseWriter, request *fcrmessages.FCRMessage) {
	// Get core structure
	c := core.GetSingleInstance()
	logging.Info("handle key management.")

	nodeID, privKey, privKeyVer, err := fcrmsgpvdadmin.DecodeProviderAdminInitialiseKeyRequest(request)
	if err != nil {
		logging.Error("Error in decoding message.")
		return
	}

	// Set the node id
	logging.Info("Check if c is nil :%v", c == nil)
	logging.Info("Setting node id")
	c.ProviderID = nodeID
	c.ProviderPrivateKey = privKey
	c.ProviderPrivateKeyVersion = privKeyVer

	// Construct messaqe
	response, err := fcrmsgpvdadmin.EncodeProviderAdminInitialiseKeyResponse(true)
	if err != nil {
		logging.Error("Error in encoding message")
		return
	}

	logging.Info("Signing response.")
	// Sign the response
	if response.Sign(c.ProviderPrivateKey, c.ProviderPrivateKeyVersion) != nil {
		logging.Error("Error in signing message")
		return
	}
	w.WriteJson(response)
}
