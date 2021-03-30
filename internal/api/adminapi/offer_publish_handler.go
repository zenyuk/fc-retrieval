package adminapi

import (
	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages/fcrmsgpvdadmin"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-provider/internal/api/providerapi"
	"github.com/ConsenSys/fc-retrieval-provider/internal/core"
	"github.com/ant0ine/go-json-rest/rest"
)

func handleProviderPublishGroupCID(w rest.ResponseWriter, request *fcrmessages.FCRMessage) {
	// Get core structure
	c := core.GetSingleInstance()
	if c.ProviderPrivateKey == nil {
		logging.Error("This provider hasn't been initialised by the admin.")
		return
	}
	logging.Info("handleProviderPublishGroupCID: %+v", request)

	cids, price, expiry, qos, err := fcrmsgpvdadmin.DecodeProviderAdminPublishGroupOfferRequest(request)
	if err != nil {
		logging.Error("Error in decoding the incoming request")
		return
	}
	offer, err := cidoffer.NewCIDOffer(c.ProviderID, cids, price, expiry, qos)
	if err != nil {
		logging.Error("Error in creating offer")
		return
	}
	// Sign the offer
	if offer.Sign(c.ProviderPrivateKey, c.ProviderPrivateKeyVersion) != nil {
		logging.Error("Error in signing the offer.")
		return
	}

	// Add offer to storage
	c.GroupOffers.Add(offer)

	c.RegisteredGatewaysMapLock.RLock()
	defer c.RegisteredGatewaysMapLock.RUnlock()

	for _, gw := range c.RegisteredGatewaysMap {
		gatewayID, err := nodeid.NewNodeIDFromHexString(gw.GetNodeID())
		if err != nil {
			logging.Error("Error with nodeID %v: %v", gw.GetNodeID(), err)
			continue
		}

		err = providerapi.RequestProviderPublishGroupCID(offer, gatewayID)
		if err != nil {
			logging.Error("Error in publishing group offer :%v", err)
		}
	}

	// Respond to admin
	response, err := fcrmsgpvdadmin.EncodeProviderAdminPublishGroupOfferResponse(true)
	if err != nil {
		logging.Error("Error in encoding response.")
		return
	}
	// Sign the response
	if response.Sign(c.ProviderPrivateKey, c.ProviderPrivateKeyVersion) != nil {
		logging.Error("Error in signing the response.")
		return
	}
	w.WriteJson(response)
}
