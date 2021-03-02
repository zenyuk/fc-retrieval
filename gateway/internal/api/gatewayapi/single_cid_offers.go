package gatewayapi

import (
	"errors"

	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrtcpcomms"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
)

// RequestSingleCIDOffers is used at start-up to request a set of single CID Offers
// from a provider with a given provider id.
func RequestSingleCIDOffers(cidMin, cidMax *cid.ContentID, providerID *nodeid.NodeID) (*fcrmessages.FCRMessage, error) {
	// Get the core structure
	g := gateway.GetSingleInstance()

	// Get the connection to provider.
	pComm, err := g.ProviderCommPool.GetConnForRequestingNode(providerID, fcrtcpcomms.AccessFromGateway)
	if err != nil {
		return nil, err
	}
	pComm.CommsLock.Lock()
	defer pComm.CommsLock.Unlock()
	// Construct message
	request, err := fcrmessages.EncodeGatewaySingleCIDOfferPublishRequest(
		g.GatewayID,
		cidMin,
		cidMax,
		g.RegistrationBlockHash,
		g.RegistrationTransactionReceipt,
		g.RegistrationMerkleRoot,
		g.RegistrationMerkleProof,
	)
	if err != nil {
		return nil, err
	}
	// Sign the request
	request.SignMessage(func(msg interface{}) (string, error) {
		return fcrcrypto.SignMessage(g.GatewayPrivateKey, g.GatewayPrivateKeyVersion, msg)
	})
	err = fcrtcpcomms.SendTCPMessage(pComm.Conn, request, settings.DefaultTCPInactivityTimeout)
	if err != nil {
		g.ProviderCommPool.DeregisterNodeCommunication(providerID)
		return nil, err
	}
	// Get a response.
	response, err := fcrtcpcomms.ReadTCPMessage(pComm.Conn, settings.DefaultLongTCPInactivityTimeout)
	if err != nil {
		g.ProviderCommPool.DeregisterNodeCommunication(providerID)
		return nil, err
	}
	// Verify the response
	// Get the provider's signing key
	g.RegisteredProvidersMapLock.RLock()
	defer g.RegisteredProvidersMapLock.RUnlock()
	_, ok := g.RegisteredProvidersMap[providerID.ToString()]
	if !ok {
		return nil, errors.New("Provider public key not found")
	}
	pubKey, err := g.RegisteredProvidersMap[providerID.ToString()].GetSigningKey()
	if err != nil {
		return nil, err
	}
	ok, err = response.VerifySignature(func(sig string, msg interface{}) (bool, error) {
		return fcrcrypto.VerifyMessage(pubKey, sig, msg)
	})
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("Fail to verify the request")
	}
	return response, nil
}

// AcknowledgeSingleCIDOffers is used to acknowledge a response
func AcknowledgeSingleCIDOffers(response *fcrmessages.FCRMessage, providerID *nodeid.NodeID) ([]fcrmessages.FCRMessage, error) {
	// Get the core structure
	g := gateway.GetSingleInstance()

	// Get the connection to provider.
	pComm, err := g.ProviderCommPool.GetConnForRequestingNode(providerID, fcrtcpcomms.AccessFromGateway)
	if err != nil {
		g.ProviderCommPool.DeregisterNodeCommunication(providerID)
		return nil, err
	}
	pComm.CommsLock.Lock()
	defer pComm.CommsLock.Unlock()

	// Verify the response
	// Get the provider's signing key
	g.RegisteredProvidersMapLock.RLock()
	defer g.RegisteredProvidersMapLock.RUnlock()
	_, ok := g.RegisteredProvidersMap[providerID.ToString()]
	if !ok {
		return nil, errors.New("Provider public key not found")
	}
	pubKey, err := g.RegisteredProvidersMap[providerID.ToString()].GetSigningKey()
	if err != nil {
		return nil, err
	}
	ok, err = response.VerifySignature(func(sig string, msg interface{}) (bool, error) {
		return fcrcrypto.VerifyMessage(pubKey, sig, msg)
	})
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("Fail to verify the request")
	}

	// Decode the response
	cidOffers, err := fcrmessages.DecodeGatewaySingleCIDOfferPublishResponse(response)
	if err != nil {
		return nil, err
	}
	// Construct the message
	cidOfferAcks := make([]fcrmessages.FCRMessage, 0)
	for _, cidOffer := range cidOffers {
		nonce, _, _, err := fcrmessages.DecodeProviderDHTPublishGroupCIDRequest(&cidOffer)
		if err != nil {
			return nil, err
		}
		// Sign the offer
		sig, err := fcrcrypto.SignMessage(g.GatewayPrivateKey, g.GatewayPrivateKeyVersion, cidOffer)
		if err != nil {
			return nil, err
		}
		cidOfferAck, err := fcrmessages.EncodeProviderDHTPublishGroupCIDAck(nonce, sig)
		if err != nil {
			return nil, err
		}
		cidOfferAcks = append(cidOfferAcks, *cidOfferAck)
	}
	ack, err := fcrmessages.EncodeGatewaySingleCIDOfferPublishResponseAck(cidOfferAcks)
	if err != nil {
		return nil, err
	}
	// Sign the ack
	ack.SignMessage(func(msg interface{}) (string, error) {
		return fcrcrypto.SignMessage(g.GatewayPrivateKey, g.GatewayPrivateKeyVersion, msg)
	})
	// Send ack
	err = fcrtcpcomms.SendTCPMessage(pComm.Conn, ack, settings.DefaultTCPInactivityTimeout)
	if err != nil {
		g.ProviderCommPool.DeregisterNodeCommunication(providerID)
		return nil, err
	}
	// Ack success
	return cidOffers, nil
}
