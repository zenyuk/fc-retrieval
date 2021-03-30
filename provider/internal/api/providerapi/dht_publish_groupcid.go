package providerapi

import (
	"errors"

	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages/fcrmsgpvd"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrtcpcomms"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-provider/internal/core"
	"github.com/ConsenSys/fc-retrieval-provider/internal/util/settings"
)

// RequestDHTProviderPublishGroupCID is used to publish a dht group CID offer to a given gateway
func RequestDHTProviderPublishGroupCID(offers []cidoffer.CIDOffer, gatewayID *nodeid.NodeID) error {
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
	pubKey, err := c.RegisteredGatewaysMap[gatewayID.ToString()].GetSigningKey()
	if err != nil {
		return err
	}

	// Construct message, TODO: Add nonce
	request, err := fcrmsgpvd.EncodeProviderPublishDHTOfferRequest(c.ProviderID, 1, offers)
	if err != nil {
		return err
	}
	// Sign the request
	if request.Sign(c.ProviderPrivateKey, c.ProviderPrivateKeyVersion) != nil {
		return errors.New("Error in signing the message")
	}
	// Send request
	err = fcrtcpcomms.SendTCPMessage(gComm.Conn, request, settings.DefaultTCPInactivityTimeout)
	if err != nil {
		c.GatewayCommPool.DeregisterNodeCommunication(gatewayID)
		return err
	}
	// Get a response
	response, err := fcrtcpcomms.ReadTCPMessage(gComm.Conn, settings.DefaultLongTCPInactivityTimeout)
	if err != nil {
		c.GatewayCommPool.DeregisterNodeCommunication(gatewayID)
		return err
	}
	// Verify the response
	if response.Verify(pubKey) != nil {
		return errors.New("Fail to verify the response")
	}

	// Verify the acks
	// TODO: Check nonce
	_, sig, err := fcrmsgpvd.DecodeProviderPublishDHTOfferResponse(response)
	if err != nil {
		return err
	}
	ok, err := fcrcrypto.VerifyMessage(pubKey, sig, request)
	if err != nil {
		return errors.New("Internal error in verifying ack")
	}
	if !ok {
		return errors.New("Fail to verify the ack")
	}

	// Add offer to ack map and storage
	for _, offer := range offers {
		// Add offer to storage
		c.NodeOfferMapLock.Lock()
		sentOffers, ok := c.NodeOfferMap[gatewayID.ToString()]
		if !ok {
			sentOffers = make([]cidoffer.CIDOffer, 0)
		}
		sentOffers = append(sentOffers, offer)
		c.NodeOfferMap[gatewayID.ToString()] = sentOffers
		c.NodeOfferMapLock.Unlock()
		// Add offer to ack map
		c.AcknowledgementMapLock.Lock()
		cidMap, ok := c.AcknowledgementMap[(offer.GetCIDs())[0].ToString()]
		if !ok {
			cidMap = make(map[string]core.DHTAcknowledgement)
			c.AcknowledgementMap[(offer.GetCIDs())[0].ToString()] = cidMap
		}
		cidMap[gatewayID.ToString()] = core.DHTAcknowledgement{
			Msg:    *request,
			MsgAck: *response,
		}
		c.AcknowledgementMapLock.Unlock()
	}
	return nil
}
