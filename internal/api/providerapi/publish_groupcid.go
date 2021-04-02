package providerapi

import (
	"bytes"
	"errors"
	"strings"

	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrtcpcomms"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-provider/internal/core"
	"github.com/ConsenSys/fc-retrieval-provider/internal/util/settings"
)

// RequestProviderPublishGroupCID is used to publish a group CID offer to a given gateway
func RequestProviderPublishGroupCID(offer *cidoffer.CIDOffer, gatewayID *nodeid.NodeID, settings settings.AppSettings) error {
	// Get the core structure
	c := core.GetSingleInstance()

	// Get the connection to the given gateway
	gComm, err := c.GatewayCommPool.GetConnForRequestingNode(gatewayID, fcrtcpcomms.AccessFromProvider)
	if err != nil {
		return err
	}
	gComm.CommsLock.Lock()
	defer gComm.CommsLock.Unlock()

	// Get the gateways's signing key
	c.RegisteredGatewaysMapLock.RLock()
	defer c.RegisteredGatewaysMapLock.RUnlock()
	pubKey, err := c.RegisteredGatewaysMap[strings.ToLower(gatewayID.ToString())].GetSigningKey()
	if err != nil {
		return err
	}

	// Construct message, TODO: Add nonce
	request, err := fcrmessages.EncodeProviderPublishGroupOfferRequest(c.ProviderID, 1, offer)
	if err != nil {
		return err
	}
	// Sign the request
	if request.Sign(c.ProviderPrivateKey, c.ProviderPrivateKeyVersion) != nil {
		return errors.New("Error in signing request")
	}
	// Send request
	err = fcrtcpcomms.SendTCPMessage(gComm.Conn, request, settings.TCPInactivityTimeout)
	if err != nil {
		c.GatewayCommPool.DeregisterNodeCommunication(gatewayID)
		return err
	}
	// Get a response
	response, err := fcrtcpcomms.ReadTCPMessage(gComm.Conn, settings.TCPInactivityTimeout)
	if err != nil {
		c.GatewayCommPool.DeregisterNodeCommunication(gatewayID)
		return err
	}

	if response.Verify(pubKey) != nil {
		logging.Error("Verify not ok")
		return errors.New("Fail to verify the response")
	}
	logging.Info("Got reponse from gateway=%v: %+v", gatewayID.ToString(), response)
	// TODO: Check nonce
	_, candidate, err := fcrmessages.DecodeProviderPublishGroupOfferResponse(response)
	if err != nil {
		logging.Error("Error with decode response: %v", err)
		return err
	}
	logging.Info("Received digest: %v", candidate)

	digest := offer.GetMessageDigest()

	// Add offer to storage
	if bytes.Equal(candidate[:], digest[:]) {
		logging.Info("Digest is OK! Add offer to storage")
		c.NodeOfferMapLock.Lock()
		defer c.NodeOfferMapLock.Unlock()
		sentOffers, ok := c.NodeOfferMap[gatewayID.ToString()]
		if !ok {
			sentOffers = make([]cidoffer.CIDOffer, 0)
		}
		sentOffers = append(sentOffers, *offer)
		c.NodeOfferMap[gatewayID.ToString()] = sentOffers
	} else {
		return errors.New("Digest not match")
	}
	return nil
}
