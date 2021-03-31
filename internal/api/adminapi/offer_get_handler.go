package adminapi

import (
	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-provider/internal/core"
	"github.com/ant0ine/go-json-rest/rest"
)

func handleProviderGetGroupCID(w rest.ResponseWriter, request *fcrmessages.FCRMessage) {
	// Get core structure
	c := core.GetSingleInstance()
	if c.ProviderPrivateKey == nil {
		logging.Error("This provider hasn't been initialised by the admin.")
		return
	}

	logging.Info("handleProviderGetGroupCID: %+v", request)
	gatewayIDs, err := fcrmessages.DecodeProviderAdminGetPublishedOfferRequest(request)
	if err != nil {
		logging.Info("Provider get group cid request fail to decode request.")
		return
	}
	logging.Info("Find offers: gatewayIDs=%+v", gatewayIDs)
	offers := make([]cidoffer.CIDOffer, 0)

	c.NodeOfferMapLock.Lock()
	defer c.NodeOfferMapLock.Unlock()
	if len(gatewayIDs) > 0 {
		for _, gatewayID := range gatewayIDs {
			offs := c.NodeOfferMap[gatewayID.ToString()]
			for _, off := range offs {
				offers = append(offers, off)
			}
		}
	} else {
		for _, values := range c.NodeOfferMap {
			for _, value := range values {
				offers = append(offers, value)
			}
		}
	}
	logging.Info("Found offers: %+v", len(offers))

	response, err := fcrmessages.EncodeProviderAdminGetPublishedOfferResponse(
		len(offers) > 0,
		offers,
	)
	if err != nil {
		logging.Info("Provider get group cid request fail to encode response.")
		panic(err)
	}
	// Sign the response
	if response.Sign(c.ProviderPrivateKey, c.ProviderPrivateKeyVersion) != nil {
		logging.Error("Error in signing the message")
		return
	}
	w.WriteJson(response)
}
