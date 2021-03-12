package providerapi

import (
	"errors"
	"net"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrtcpcomms"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
)

func handleProviderPublishGroupCIDRequest(conn net.Conn, request *fcrmessages.FCRMessage) error {
	// Get the core structure
	g := gateway.GetSingleInstance()
	logging.Info("GatewayPrivateKey: %s", g.GatewayPrivateKey.EncodePrivateKey())

	logging.Info("handleProviderPublishGroupCIDRequest: %+v", request)
	// TODO Add nonce, it looks like nonce is not needed
	logging.Info("Decode provider publish group CID request: %+v", request)
	_, offer, err := fcrmessages.DecodeProviderPublishGroupCIDRequest(request)
	if err != nil {
		logging.Info("Provider publish group cid request fail to decode.")
		return err
	}
	// Get the provider's signing key
	g.RegisteredProvidersMapLock.RLock()
	defer g.RegisteredProvidersMapLock.RUnlock()
	_, ok := g.RegisteredProvidersMap[offer.NodeID.ToString()]
	if !ok {
		return errors.New("Provider not found")
	}
	pubKey, err := g.RegisteredProvidersMap[offer.NodeID.ToString()].GetSigningKey()
	if err != nil {
		return err
	}
	// First verify the message
	ok, err = request.VerifySignature(func(sig string, msg interface{}) (bool, error) {
		return fcrcrypto.VerifyMessage(pubKey, sig, msg)
	})
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("Fail to verify the request")
	}

	// Second Need to verify the offer
	ok, err = offer.VerifySignature(func(sig string, msg interface{}) (bool, error) {
		return fcrcrypto.VerifyMessage(pubKey, sig, msg)
	})
	// Error in verifying message
	if err != nil {
		logging.Error("Internal error in verifying group cid offer.")
		// Ignored.
		return err
	}
	// Offer does not pass verification
	if !ok {
		logging.Info("Offer does not pass verification.")
		// Ignored.
		return errors.New("Offer does not pass verification")
	}

	// Store the offer
	if g.Offers.Add(offer) != nil {
		logging.Error("Internal error in adding group cid offer.")
		return err
	}

	logging.Info("Stored offers: %+v", g.Offers)

	logging.Info("Encode provider publish group CID response: %+v", offer)
	response, err := fcrmessages.EncodeProviderPublishGroupCIDResponse(
		*g.GatewayID,
		offer.GetMessageDigest(),
	)
	if err != nil {
		logging.Error("Internal error in encoding publish group cid response.")
		return err
	}

	// Sign the response
	response.SignMessage(func(msg interface{}) (string, error) {
		return fcrcrypto.SignMessage(g.GatewayPrivateKey, g.GatewayPrivateKeyVersion, msg)
	})
	return fcrtcpcomms.SendTCPMessage(conn, response, settings.DefaultTCPInactivityTimeout)
}
