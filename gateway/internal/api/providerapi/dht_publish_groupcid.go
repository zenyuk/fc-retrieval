package providerapi

import (
	"errors"
	"net"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrtcpcomms"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
)

func handleProviderDHTPublishGroupCIDRequest(conn net.Conn, request *fcrmessages.FCRMessage) error {
	// Get the core structure
	g := gateway.GetSingleInstance()

	nonce, providerID, offers, err := fcrmessages.DecodeProviderDHTPublishGroupCIDRequest(request)
	if err != nil {
		logging.Info("Provider publish cid offer dht request fail to decode.")
		return nil
	}

	// Get the provider's signing key
	g.RegisteredProvidersMapLock.RLock()
	defer g.RegisteredProvidersMapLock.RUnlock()
	_, ok := g.RegisteredProvidersMap[providerID.ToString()]
	if !ok {
		return errors.New("Provider public key not found")
	}
	pubKey, err := g.RegisteredProvidersMap[providerID.ToString()].GetSigningKey()
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

	for _, offer := range offers {
		// Need to verify the offer one by one
		ok, err = offer.VerifySignature(func(sig string, msg interface{}) (bool, error) {
			return fcrcrypto.VerifyMessage(pubKey, sig, msg)
		})

		if err != nil {
			logging.Error("Internal error in verifying the cid offer.")
			continue
		}

		if !ok {
			logging.Info("Offer does not pass verification.")
			continue
		}

		if g.Offers.Add(&offer) != nil {
			// Ignored.
			logging.Error("Internal error in adding single cid offer.")
		}
	}

	// Sign the request
	sig, err := fcrcrypto.SignMessage(g.GatewayPrivateKey, g.GatewayPrivateKeyVersion, request)
	if err != nil {
		// Ignored.
		logging.Error("Internal error in signing message.")
		return nil
	}

	response, err := fcrmessages.EncodeProviderDHTPublishGroupCIDAck(nonce, sig)
	if err != nil {
		logging.Error("Internal error in encoding message.")
		return nil
	}
	// Sign the response
	response.SignMessage(func(msg interface{}) (string, error) {
		return fcrcrypto.SignMessage(g.GatewayPrivateKey, g.GatewayPrivateKeyVersion, msg)
	})
	return fcrtcpcomms.SendTCPMessage(conn, response, settings.DefaultTCPInactivityTimeout)
}
