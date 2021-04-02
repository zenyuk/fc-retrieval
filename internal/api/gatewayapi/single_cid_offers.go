package gatewayapi

import (
	"errors"
	"math/big"
	"net"
	"strconv"
	"strings"

	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrtcpcomms"
	log "github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-provider/internal/core"
	"github.com/ConsenSys/fc-retrieval-provider/internal/util/settings"
)

func handleSingleCIDOffersPublishRequest(conn net.Conn, request *fcrmessages.FCRMessage, settings settings.AppSettings) error {
	// Get core structure
	c := core.GetSingleInstance()
	gatewayID, cidMin, cidMax, registrationBlock, registrationTransactionReceipt, registrationMerkleRoot, registrationMerkleProof, err := fcrmessages.DecodeGatewayListDHTOfferRequest(request)
	if err != nil {
		return err
	}
	// TODO, Need to check registration info
	log.Info("Registration info: %v, %v, %v, %v", registrationBlock, registrationTransactionReceipt, registrationMerkleRoot, registrationMerkleProof)

	// Get the gateways's signing key
	c.RegisteredGatewaysMapLock.RLock()
	defer c.RegisteredGatewaysMapLock.RUnlock()
	_, ok := c.RegisteredGatewaysMap[strings.ToLower(gatewayID.ToString())]
	if !ok {
		return errors.New("Gateway register not found")
	}
	pubKey, err := c.RegisteredGatewaysMap[strings.ToLower(gatewayID.ToString())].GetSigningKey()
	if err != nil {
		return err
	}
	// Verify the incoming request
	if request.Verify(pubKey) != nil {
		return errors.New("Fail to verify the request")
	}

	// Search offers
	min, err := strconv.ParseInt(cidMin.ToString(), 16, 32) // TODO, CHECK IF THIS IS CORRECT
	if err != nil {
		return err
	}
	max, err := strconv.ParseInt(cidMax.ToString(), 16, 32) // TODO, CHECK IF THIS IS CORRECT
	if err != nil {
		return err
	}
	if max < min {
		return errors.New("Invalid parameters")
	}

	maxOffers := 500
	offers := make([]cidoffer.CIDOffer, 0)
	for i := min; i <= max; i++ {
		id, err := cid.NewContentID(big.NewInt(i))
		if err != nil {
			return err
		}
		offers, exists := c.SingleOffers.GetOffers(id)
		if exists {
			for _, offer := range offers {
				offers = append(offers, offer)
				if len(offers) >= maxOffers {
					break
				}
			}
		}
		if len(offers) >= maxOffers {
			break
		}
	}
	maxOffersPerMsg := 50
	msgs := make([]fcrmessages.FCRMessage, 0)
	for {
		if len(offers) > maxOffersPerMsg {
			msg, err := fcrmessages.EncodeProviderPublishDHTOfferRequest(c.ProviderID, 1, offers[:50]) //TODO: Add nonce
			if err != nil {
				return err
			}
			msgs = append(msgs, *msg)
			offers = offers[50:]
		} else {
			msg, err := fcrmessages.EncodeProviderPublishDHTOfferRequest(c.ProviderID, 1, offers) //TODO: Add nonce
			if err != nil {
				return err
			}
			msgs = append(msgs, *msg)
			break
		}
	}

	// Construct response
	response, err := fcrmessages.EncodeGatewayListDHTOfferResponse(msgs)
	if err != nil {
		return err
	}
	// Sign the response
	if response.Sign(c.ProviderPrivateKey, c.ProviderPrivateKeyVersion) != nil {
		return errors.New("Error in signing the response")
	}
	// Respond
	err = fcrtcpcomms.SendTCPMessage(conn, response, settings.TCPLongInactivityTimeout)
	if err != nil {
		return err
	}

	// Get acks
	acks, err := fcrtcpcomms.ReadTCPMessage(conn, settings.TCPLongInactivityTimeout)
	if err != nil {
		return err
	}
	// Verify the acks
	if acks.Verify(pubKey) != nil {
		return errors.New("Fail to verify the acks")
	}

	acknowledgements, err := fcrmessages.DecodeGatewayListDHTOfferAck(acks)
	if len(acknowledgements) != len(offers) {
		return errors.New("Invalid response")
	}
	for i, acknowledgement := range acknowledgements {
		// TODO: Check nonce.
		_, signature, err := fcrmessages.DecodeProviderPublishDHTOfferResponse(&acknowledgement)
		if err != nil {
			return err
		}
		ok, err := fcrcrypto.VerifyMessage(pubKey, signature, msgs[i])
		if err != nil {
			return err
		}
		if !ok {
			return errors.New("Verification failed")
		}
		// It's okay, add to acknowledgements map
		c.AcknowledgementMapLock.Lock()
		c.AcknowledgementMap[offers[i].GetCIDs()[0].ToString()][gatewayID.ToString()] = core.DHTAcknowledgement{
			Msg:    msgs[i],
			MsgAck: acknowledgement,
		}
		c.AcknowledgementMapLock.Unlock()
	}
	return nil
}
