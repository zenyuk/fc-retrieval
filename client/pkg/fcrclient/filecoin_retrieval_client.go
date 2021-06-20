package fcrclient

/*
 * Copyright 2020 ConsenSys Software Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrpaymentmgr"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrregistermgr"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/register"

	"github.com/ConsenSys/fc-retrieval-client/pkg/api/clientapi"
)

// FilecoinRetrievalClient is an example implementation using the api,
// which holds information about the interaction of the Filecoin
// Retrieval Client with Filecoin Retrieval Gateways.
type FilecoinRetrievalClient struct {
	Settings ClientSettings

	// List of gateways this client can potentially use
	GatewaysToUse     map[string]register.GatewayRegistrar
	GatewaysToUseLock sync.RWMutex

	// List of gateway in use. A client may request a node be added to this list
	ActiveGateways     map[string]register.GatewayRegistrar
	ActiveGatewaysLock sync.RWMutex

	// PaymentMgr payment manager
	paymentMgr     *fcrpaymentmgr.FCRPaymentMgr
	paymentMgrLock sync.RWMutex

	clientApi      clientapi.ClientApi
	registerMgr    *fcrregistermgr.FCRRegisterMgr
}

// NewFilecoinRetrievalClient initialise the Filecoin Retrieval Client library
func NewFilecoinRetrievalClient(settings ClientSettings, registerMgr *fcrregistermgr.FCRRegisterMgr) (*FilecoinRetrievalClient, error) {
	f := &FilecoinRetrievalClient{
		Settings:           settings,
		GatewaysToUse:      make(map[string]register.GatewayRegistrar),
		GatewaysToUseLock:  sync.RWMutex{},
		ActiveGateways:     make(map[string]register.GatewayRegistrar),
		ActiveGatewaysLock: sync.RWMutex{},
		clientApi:          clientapi.NewClientApi(),
		registerMgr:        registerMgr,
	}
	return f, nil
}

func (c *FilecoinRetrievalClient) PaymentMgr() *fcrpaymentmgr.FCRPaymentMgr {
	// lazy init
	if c.paymentMgr == nil {
		c.paymentMgrLock.Lock()
		defer c.paymentMgrLock.Unlock()
		if c.paymentMgr == nil {
			mgr, err := fcrpaymentmgr.NewFCRPaymentMgr(c.Settings.walletPrivateKey, c.Settings.lotusAP, c.Settings.lotusAuthToken)
			if err != nil {
				logging.Error("Error initializing payment manager.")
				return nil
			}
			c.paymentMgr = mgr
		}
	}

	return c.paymentMgr
}

// FindGateways find gateways located near to the specified location. Use AddGateways
// to use these gateways.
func (c *FilecoinRetrievalClient) FindGateways(location string, maxNumToLocate int) ([]*nodeid.NodeID, error) {
	// Determine gateways to use. For the moment, this is just "use all of them"
	// TODO: This will have to become, use gateways that this client has FIL registered with.
	gateways := c.registerMgr.GetAllGateways()
	if gateways == nil {
		return nil, errors.New("error in getting all registered gateways")
	}

	res := make([]*nodeid.NodeID, 0)
	count := 0
	for _, info := range gateways {
		if info.GetRegionCode() == location {
			nodeID, err := nodeid.NewNodeIDFromHexString(info.GetNodeID())
			if err != nil {
				logging.Error("Error in generating node id, skipping: %v", info.GetNodeID())
				continue
			}
			res = append(res, nodeID)
			count++
			if count >= maxNumToLocate {
				break
			}
		}
	}
	return res, nil
}

// AddGatewaysToUse adds one or more gateways to use.
func (c *FilecoinRetrievalClient) AddGatewaysToUse(gwNodeIDs []*nodeid.NodeID) int {
	numAdded := 0
	for _, gwToAddID := range gwNodeIDs {
		if gwToAddID == nil {
			logging.Error("can't add a gateway to use, a nil was given")
			continue
		}
		c.GatewaysToUseLock.RLock()
		_, exist := c.GatewaysToUse[gwToAddID.ToString()]
		c.GatewaysToUseLock.RUnlock()
		if exist {
			continue
		}
		gateway := c.registerMgr.GetGateway(gwToAddID)
		if gateway == nil {
			logging.Error("error getting registered gateway %v", gwToAddID.ToString())
			continue
		}
		if !validateGatewayInfo(gateway) {
			logging.Error("Register info not valid.")
			continue
		}
		// Success
		c.GatewaysToUseLock.Lock()
		c.GatewaysToUse[gwToAddID.ToString()] = gateway
		c.GatewaysToUseLock.Unlock()
		numAdded++
	}
	return numAdded
}

// RemoveGatewaysToUse removes one or more gateways from the list of Gateways to use.
// This also removes the gateway from gateways in active map.
func (c *FilecoinRetrievalClient) RemoveGatewaysToUse(gwNodeIDs []*nodeid.NodeID) int {
	c.GatewaysToUseLock.Lock()
	defer c.GatewaysToUseLock.Unlock()

	numRemoved := 0
	for _, gwToRemoveID := range gwNodeIDs {
		_, exist := c.GatewaysToUse[gwToRemoveID.ToString()]
		if exist {
			delete(c.GatewaysToUse, gwToRemoveID.ToString())
			numRemoved++
			c.ActiveGatewaysLock.Lock()
			delete(c.ActiveGateways, gwToRemoveID.ToString())
			c.ActiveGatewaysLock.Unlock()
		}
	}

	return numRemoved
}

// RemoveAllGatewaysToUse removes all gateways from the list of Gateways.
// This also cleared all gateways in active
func (c *FilecoinRetrievalClient) RemoveAllGatewaysToUse() int {
	c.GatewaysToUseLock.Lock()
	defer c.GatewaysToUseLock.Unlock()
	c.ActiveGatewaysLock.Lock()
	defer c.ActiveGatewaysLock.Unlock()

	numRemoved := len(c.GatewaysToUse)
	c.GatewaysToUse = make(map[string]register.GatewayRegistrar)
	c.ActiveGateways = make(map[string]register.GatewayRegistrar)

	return numRemoved
}

// GetGatewaysToUse returns the list of gateways to use.
func (c *FilecoinRetrievalClient) GetGatewaysToUse() []*nodeid.NodeID {
	c.GatewaysToUseLock.RLock()
	defer c.GatewaysToUseLock.RUnlock()

	res := make([]*nodeid.NodeID, 0)
	for key := range c.GatewaysToUse {
		nodeID, err := nodeid.NewNodeIDFromHexString(key)
		if err != nil {
			logging.Error("Error in generating node id.")
			continue
		}
		res = append(res, nodeID)
	}

	return res
}

// AddActiveGateways adds one or more gateways to active gateway map.
// Returns the number of gateways added.
func (c *FilecoinRetrievalClient) AddActiveGateways(gwNodeIDs []*nodeid.NodeID) int {
	numAdded := 0
	for _, gwToAddID := range gwNodeIDs {
		c.ActiveGatewaysLock.RLock()
		_, exist := c.ActiveGateways[gwToAddID.ToString()]
		c.ActiveGatewaysLock.RUnlock()
		if exist {
			continue
		}
		c.GatewaysToUseLock.RLock()
		gatewayRegistrar, exist := c.GatewaysToUse[strings.ToLower(gwToAddID.ToString())]
		c.GatewaysToUseLock.RUnlock()
		if !exist {
			logging.Error("Given node id: %v does not exist in gateways to use map, consider add the gateway first.", gwToAddID.ToString())
			continue
		}
		// Attempt an establishment
		challenge := make([]byte, 32)
		rand.Read(challenge)
		ttl := time.Now().Unix() + c.Settings.EstablishmentTTL()
		err := c.clientApi.RequestEstablishment(gatewayRegistrar, challenge, c.Settings.ClientID(), ttl)
		if err != nil {
			logging.Error("Error in initial establishment: %v", err.Error())
			continue
		}
		// It is success
		c.ActiveGatewaysLock.Lock()
		c.ActiveGateways[gwToAddID.ToString()] = gatewayRegistrar
		c.ActiveGatewaysLock.Unlock()
		numAdded++
	}
	return numAdded
}

// RemoveActiveGateways removes one or more gateways from the list of Gateways in active.
func (c *FilecoinRetrievalClient) RemoveActiveGateways(gwNodeIDs []*nodeid.NodeID) int {
	c.ActiveGatewaysLock.Lock()
	defer c.ActiveGatewaysLock.Unlock()

	numRemoved := 0
	for _, gwToRemoveID := range gwNodeIDs {
		_, exist := c.ActiveGateways[gwToRemoveID.ToString()]
		if exist {
			delete(c.ActiveGateways, gwToRemoveID.ToString())
			numRemoved++
		}
	}

	return numRemoved
}

// RemoveAllActiveGateways removes all gateways from the list of Gateways in active.
func (c *FilecoinRetrievalClient) RemoveAllActiveGateways() int {
	c.ActiveGatewaysLock.Lock()
	defer c.ActiveGatewaysLock.Unlock()

	numRemoved := len(c.ActiveGateways)
	c.ActiveGateways = make(map[string]register.GatewayRegistrar)

	return numRemoved
}

// GetActiveGateways returns the list of gateways that are active.
func (c *FilecoinRetrievalClient) GetActiveGateways() []*nodeid.NodeID {
	c.ActiveGatewaysLock.RLock()
	defer c.ActiveGatewaysLock.RUnlock()

	res := make([]*nodeid.NodeID, 0)
	for key := range c.ActiveGateways {
		nodeID, err := nodeid.NewNodeIDFromHexString(key)
		if err != nil {
			logging.Error("Error in generating node id.")
			continue
		}
		res = append(res, nodeID)
	}

	return res
}

// FindOffersStandardDiscovery finds offer using standard discovery from given gateways
func (c *FilecoinRetrievalClient) FindOffersStandardDiscovery(contentID *cid.ContentID, gatewayID *nodeid.NodeID) ([]cidoffer.SubCIDOffer, error) {
	c.ActiveGatewaysLock.RLock()
	defer c.ActiveGatewaysLock.RUnlock()

	gw, exists := c.ActiveGateways[gatewayID.ToString()]
	if !exists {
		return make([]cidoffer.SubCIDOffer, 0), errors.New("given gatewayID is not in active nodes map")
	}
	// TODO need to do nonce management
	offers, err := c.clientApi.RequestStandardDiscover(gw, contentID, rand.Int63(), time.Now().Unix()+c.Settings.EstablishmentTTL(), "", "")
	if err != nil {
		logging.Warn("GatewayStdDiscovery error. Gateway: %s, Error: %s", gw.GetNodeID(), err)
		return make([]cidoffer.SubCIDOffer, 0), errors.New("error in requesting standard discovery")
	}
	// Verify the offer one by one
	for _, offer := range offers {
		// Get provider's pubkey
		provider := c.registerMgr.GetProvider(offer.GetProviderID())
		if provider == nil {
			logging.Error("Error getting provider info.")
			continue
		}
		if !validateProviderInfo(provider) {
			logging.Error("Provider register info not valid.")
			continue
		}
		pubKey, err := provider.GetSigningKey()
		if err != nil {
			logging.Error("Fail to obtain public key.")
			continue
		}
		// Verify the offer sig
		err = offer.Verify(pubKey)
		if err != nil {
			logging.Error("Offer signature fail to verify 4, with error: " + fmt.Sprintf("%v", err))
			continue
		}
		// Now Verify the merkle proof
		if offer.VerifyMerkleProof() != nil {
			logging.Error("Merkle proof verification failed.")
			continue
		}
		// Offer pass verification
		logging.Info("Offer pass every verification, added to result")
	}
	return offers, nil
}

// FindOffersDHTDiscovery finds offer using dht discovery from given gateways
func (c *FilecoinRetrievalClient) FindOffersDHTDiscovery(contentID *cid.ContentID, gatewayID *nodeid.NodeID, numDHT int64) (map[string]*[]cidoffer.SubCIDOffer, error) {
	c.ActiveGatewaysLock.RLock()
	defer c.ActiveGatewaysLock.RUnlock()

	offersMap := make(map[string]*[]cidoffer.SubCIDOffer)

	gw, exists := c.ActiveGateways[gatewayID.ToString()]
	if !exists {
		return offersMap, errors.New("given gatewayID is not in active nodes map")
	}
	// TODO need to do nonce management
	contacted, contactedResp, uncontactable, err := c.clientApi.RequestDHTDiscover(gw, contentID, rand.Int63(), time.Now().Unix()+c.Settings.EstablishmentTTL(), numDHT, false, "", "")
	if err != nil {
		logging.Warn("GatewayDHTDiscovery error. Gateway: %s, Error: %s", gw.GetNodeID(), err)
		return offersMap, errors.New("error in requesting dht discovery")
	}
	for i := 0; i < len(uncontactable); i++ {
		logging.Warn("Gateway: %v is uncontactable.", uncontactable[i].ToString())
	}

	for i := 0; i < len(contacted); i++ {
		id := contacted[i]
		resp := contactedResp[i]
		// Verify the sub response
		// Get gateway's pubkey
		gateway := c.registerMgr.GetGateway(&id)
		if gateway == nil {
			logging.Error("Error in getting gateway info.")
			continue
		}
		if !validateGatewayInfo(gateway) {
			logging.Error("Gateway register info not valid.")
			continue
		}
		pubKey, err := gateway.GetSigningKey()
		if err != nil {
			logging.Error("Fail to obtain public key.")
			continue
		}
		if resp.Verify(pubKey) != nil {
			logging.Error("Fail to verify sub response.")
			continue
		}
		_, _, _, offers, _, err := fcrmessages.DecodeGatewayDHTDiscoverResponse(&resp)
		if err != nil {
			logging.Error("Fail to decode response")
			continue
		}
		entry := make([]cidoffer.SubCIDOffer, 0)
		for _, offer := range offers {
			// Get provider's pubkey
			provider := c.registerMgr.GetProvider(offer.GetProviderID())
			if provider == nil {
				logging.Error("Error getting provider info.")
				continue
			}
			if !validateProviderInfo(provider) {
				logging.Error("Provider register info not valid.")
				continue
			}
			pubKey, err := provider.GetSigningKey()
			if err != nil {
				logging.Error("Fail to obtain public key.")
				continue
			}
			// Verify the offer sig
			err = offer.Verify(pubKey)
			if err != nil {
				logging.Error("Offer signature fail to verify 5, with error: " + fmt.Sprintf("%v", err))
				continue
			}
			// Now Verify the merkle proof
			if offer.VerifyMerkleProof() != nil {
				logging.Error("Merkle proof verification failed.")
				continue
			}
			// Offer pass verification
			entry = append(entry, offer)
			logging.Info("Offer pass every verification, added to result")
		}
		offersMap[id.ToString()] = &entry
	}

	return offersMap, nil
}

// FindOffersDHTDiscoveryV2 finds offer using dht discovery from given gateway with maximum number of offers
// offersNumberLimit - maximum number of offers the client asking to have
func (c *FilecoinRetrievalClient) FindOffersDHTDiscoveryV2(contentID *cid.ContentID, gatewayID *nodeid.NodeID, numDHT int64, offersNumberLimit int) (map[string]*[]cidoffer.SubCIDOffer, error) {
	offersMap := make(map[string]*[]cidoffer.SubCIDOffer)

	c.ActiveGatewaysLock.RLock()
	defer c.ActiveGatewaysLock.RUnlock()

	// entryGateway - a Gateway which will be an entry point for us to get to other Gateways
	entryGateway, exists := c.ActiveGateways[gatewayID.ToString()]
	if !exists {
		return nil, errors.New("given gatewayID is not in active nodes map")
	}

	defaultPaymentLane := uint64(0)
	initialRequestPaymentAmount := new(big.Int).Mul(big.NewInt(numDHT), big.NewInt(defaultSearchPrice))
	paymentChannel, voucher, needTopup, paymentErr := c.PaymentMgr().Pay(entryGateway.GetAddress(), defaultPaymentLane, initialRequestPaymentAmount)
	if paymentErr != nil {
		return nil, fmt.Errorf("Unable to make payment for initial DHT offers discovery, error: %s ", paymentErr.Error())
	}
	if needTopup {
		if err := c.PaymentMgr().Topup(entryGateway.GetAddress(), c.Settings.topUpAmount); err != nil {
			return nil, fmt.Errorf("Unable to topup while paying for initial offer DHT discovery, error: %s ", err.Error())
		}
		paymentChannel, voucher, needTopup, paymentErr = c.PaymentMgr().Pay(entryGateway.GetAddress(), defaultPaymentLane, initialRequestPaymentAmount)
		if paymentErr != nil {
			return nil, fmt.Errorf("Unable to make payment for initial DHT offers discovery, with topup first, error: %s ", paymentErr.Error())
		}
	}
	logging.Info("Successful initial payment for DHT offers discovery, payment channel: %s, voucher: %s", paymentChannel, voucher)

	// TODO need to do nonce management
	nonce := rand.Int63()
	ttl := time.Now().Unix() + c.Settings.EstablishmentTTL()
	contactedGateways, contactedResp, uncontactable, err := c.clientApi.RequestDHTDiscoverV2(entryGateway, contentID, nonce, ttl, numDHT, false, paymentChannel, voucher)
	if err != nil {
		logging.Warn("GatewayDHTDiscovery error. Gateway: %s, Error: %s", entryGateway.GetNodeID(), err)
		return nil, errors.New("error in requesting dht discovery")
	}
	for i := 0; i < len(uncontactable); i++ {
		logging.Warn("Gateway: %v is uncontactable.", uncontactable[i].ToString())
	}

	var addedSubOffersCount int
	var offersDigestsFromAllGateways [][][cidoffer.CIDOfferDigestSize]byte
	for i := 0; i < len(contactedGateways); i++ {
		contactedGatewayID := contactedGateways[i]
		resp := contactedResp[i]
		// Verify the sub response
		// Get gateway's pubkey
		gateway  := c.registerMgr.GetGateway(&contactedGatewayID)
		if gateway == nil {
			logging.Error("Error in getting gateway info.")
			continue
		}
		if !validateGatewayInfo(gateway) {
			logging.Error("Gateway register info not valid.")
			continue
		}
		pubKey, err := gateway.GetSigningKey()
		if err != nil {
			logging.Error("Fail to obtain public key.")
			continue
		}
		if resp.Verify(pubKey) != nil {
			logging.Error("Fail to verify sub response.")
			continue
		}
		_, _, found, offerDigests, _, err := fcrmessages.DecodeGatewayDHTDiscoverResponseV2(&resp)
		if err != nil {
			logging.Error("Fail to decode response")
			continue
		}
		if !found {
			return offersMap, nil
		}
		// comply with given offers number limit
		if addedSubOffersCount+len(offerDigests) > offersNumberLimit {
			offerDigests = offerDigests[:(offersNumberLimit - addedSubOffersCount)]
		}

		offersDigestsFromAllGateways = append(offersDigestsFromAllGateways, offerDigests)
		addedSubOffersCount += len(offerDigests)

		if addedSubOffersCount >= offersNumberLimit {
			break
		}
	}

	// Need to pay again
	unit := 0
	for _, entry := range offersDigestsFromAllGateways {
		unit += len(entry)
	}
	offerRequestPaymentAmount := new(big.Int).Mul(big.NewInt(int64(unit)), c.Settings.offerPrice)

	paymentChannel, voucher, needTopup, paymentErr = c.PaymentMgr().Pay(entryGateway.GetAddress(), defaultPaymentLane, offerRequestPaymentAmount)
	if paymentErr != nil {
		return nil, fmt.Errorf("Unable to make payment for initial DHT offers discovery, error: %s ", paymentErr.Error())
	}
	if needTopup {
		if err := c.PaymentMgr().Topup(entryGateway.GetAddress(), c.Settings.topUpAmount); err != nil {
			return nil, fmt.Errorf("Unable to topup while paying for initial offer DHT discovery, error: %s ", err.Error())
		}
		paymentChannel, voucher, needTopup, paymentErr = c.PaymentMgr().Pay(entryGateway.GetAddress(), defaultPaymentLane, offerRequestPaymentAmount)
		if paymentErr != nil {
			return nil, fmt.Errorf("Unable to make payment for initial DHT offers discovery, with topup first, error: %s ", paymentErr.Error())
		}
	}
	logging.Info("Successful initial payment for DHT offers discovery, payment channel: %s, voucher: %s", paymentChannel, voucher)

	allGatewaysOffers, discoverError := c.clientApi.RequestDHTOfferDiscover(entryGateway, contactedGateways, contentID, nonce, offersDigestsFromAllGateways, paymentChannel, voucher)
	if discoverError != nil {
		return nil, fmt.Errorf("error getting sub-offers from their digests: %s", discoverError)
	}
	for _, entry := range allGatewaysOffers {
		offersMap[entry.GatewayID.ToString()] = &entry.SubOffers
	}

	return offersMap, nil
}

// FindDHTOfferAck finds offer ack for a cid, gateway pair
func (c *FilecoinRetrievalClient) FindDHTOfferAck(contentID *cid.ContentID, gatewayID *nodeid.NodeID, providerID *nodeid.NodeID) (bool, error) {
	provider := c.registerMgr.GetProvider(providerID)
	if provider == nil {
		logging.Error("Error getting registered provider %v", providerID)
		return false, errors.New("provider not found inside register")
	}
	if !validateProviderInfo(provider) {
		logging.Error("Register info not valid.")
		return false, errors.New("invalid register info")
	}

	found, request, ack, err := c.clientApi.RequestDHTOfferAck(provider, contentID, gatewayID)
	if err != nil {
		return false, err
	}
	if !found {
		return false, nil
	}
	// We need to verify	1, request is indeed a request, ack is indeed an ack.
	//						2, request is signed by provider, ack is signed by gateway
	//						3, ack and request has same nonce
	// TODO: We will have to also include the hash of the request inside the ack, so the hash can be checked

	// Get provider pub key and gw pub key
	// Get gateway's pubkey
	gateway := c.registerMgr.GetGateway(gatewayID)
	if gateway == nil {
		logging.Error("Error in getting gateway info.")
		return false, errors.New("error in getting gateway info")
	}
	if !validateGatewayInfo(gateway) {
		logging.Error("Gateway register info not valid.")
		return false, errors.New("gateway register info not valid")
	}
	gwPubKey, err := gateway.GetSigningKey()
	if err != nil {
		logging.Error("Fail to obtain public key.")
		return false, errors.New("fail to obtain public key")
	}
	// Get provider's pubkey
	pvdPubKey, err := provider.GetSigningKey()
	if err != nil {
		logging.Error("Fail to obtain public key.")
		return false, errors.New("fail to obtain public key")
	}
	// Verify the request.
	if request.Verify(pvdPubKey) != nil {
		return false, errors.New("error in verifying request")
	}
	// Verify the offer indeed contains the given cid
	_, _, offers, err := fcrmessages.DecodeProviderPublishDHTOfferRequest(request)
	if err != nil {
		return false, err
	}
	found = false
	for _, offer := range offers {
		if offer.GetCIDs()[0].ToString() == contentID.ToString() {
			found = true
			break
		}
	}
	if !found {
		return false, errors.New("initial request does not contain the given cid")
	}
	// Verify the ack
	if ack.Verify(gwPubKey) != nil {
		return false, errors.New("error in verifying the ack")
	}
	_, signature, err := fcrmessages.DecodeProviderPublishDHTOfferResponse(ack)
	if err != nil {
		return false, err
	}
	// Verify ack against request
	ok, err := fcrcrypto.VerifyMessage(gwPubKey, signature, request)
	if err != nil {
		return false, err
	}
	return ok, nil
}

// FindOffersStandardDiscoveryV2 finds offer using standard discovery from given gateways
func (c *FilecoinRetrievalClient) FindOffersStandardDiscoveryV2(contentID *cid.ContentID, gatewayID *nodeid.NodeID, maxOffers int) ([]cidoffer.SubCIDOffer, error) {
	c.ActiveGatewaysLock.RLock()
	defer c.ActiveGatewaysLock.RUnlock()

	gw, exists := c.ActiveGateways[gatewayID.ToString()]
	if !exists {
		return make([]cidoffer.SubCIDOffer, 0), errors.New("given gatewayID is not in active nodes map")
	}

	cidOffers := make([]cidoffer.SubCIDOffer, 0)

	// It will call the payment manager to pay for the initial request
	paychAddr, voucher, topup, err := c.paymentMgr.Pay(gw.GetAddress(), 0, c.Settings.searchPrice)
	if err != nil {
		return cidOffers, err
	}

	// There isn't any balance in the payment channel, need to topup (create)
	if topup == true {
		// If topup failed, then there is not enough balance, return error.
		err = c.PaymentMgr().Topup(gw.GetAddress(), c.Settings.topUpAmount)
		if err != nil {
			logging.Warn("Topup. Gateway: %s, Error: %s", gw.GetNodeID(), err)
			return make([]cidoffer.SubCIDOffer, 0), errors.New("error in payment manager topup - is not enough balance")
		}
		paychAddr, voucher, topup, err = c.paymentMgr.Pay(gw.GetAddress(), 0, c.Settings.searchPrice)
		if err != nil {
			return cidOffers, err
		}
	}

	// It pays for the first request to get a list of offer digests.
	// TODO need to do nonce management
	offerDigests, err := c.clientApi.RequestStandardDiscoverV2(gw, contentID, rand.Int63(), time.Now().Unix()+c.Settings.EstablishmentTTL(), paychAddr, voucher)
	if err != nil {
		logging.Warn("GatewayStdDiscovery error. Gateway: %s, Error: %s", gw.GetNodeID(), err)
		return make([]cidoffer.SubCIDOffer, 0), errors.New("error in requesting standard discovery")
	}
	if len(offerDigests) == 0 {
		// No offer found
		return cidOffers, nil
	}

	// Depend on the input (maximum num of offers) It will send further requests to request offers from the same gateway.
	// It will call the payment manager to pay for the initial request
	if len(offerDigests) > maxOffers {
		offerDigests = offerDigests[:maxOffers]
	}

	lenOffers := new(big.Int).SetInt64(int64(len(offerDigests)))
	expectedAmount := lenOffers.Mul(lenOffers, c.Settings.offerPrice)

	paychAddr, voucher, topup, err = c.paymentMgr.Pay(gw.GetAddress(), 0, expectedAmount)
	if err != nil {
		return cidOffers, err
	}

	// There isn't any balance in the payment channel, need to topup (create)
	if topup == true {
		// If topup failed, then there is not enough balance, return error.
		err = c.PaymentMgr().Topup(gw.GetAddress(), c.Settings.TopUpAmount())
		if err != nil {
			logging.Warn("Topup. Gateway: %s, Error: %s", gw.GetNodeID(), err)
			return make([]cidoffer.SubCIDOffer, 0), errors.New("error in payment manager topup - is not enough balance")
		}
		paychAddr, voucher, topup, err = c.paymentMgr.Pay(gw.GetAddress(), 0, expectedAmount)
		if err != nil {
			return cidOffers, err
		}
	}

	offers, err := c.clientApi.RequestStandardDiscoverOffer(gw, contentID, rand.Int63(), time.Now().Unix()+c.Settings.EstablishmentTTL(), offerDigests, paychAddr, voucher)

	var validOffers []cidoffer.SubCIDOffer
	// Verify the offer one by one
	for _, offer := range offers {
		// Get provider's pubkey
		provider := c.registerMgr.GetProvider(offer.GetProviderID())
		if provider == nil {
			logging.Error("Error getting provider info.")
			continue
		}
		if !validateProviderInfo(provider) {
			logging.Error("Provider register info not valid.")
			continue
		}
		if err != nil {
			logging.Error("Offer signature fail to verify 1, with error: " + fmt.Sprintf("%v", err))
			continue
		}
		pubKey, err := provider.GetSigningKey()
		if err != nil {
			logging.Error("Fail to obtain public key.")
			continue
		}
		// Verify the offer sig
		err = offer.Verify(pubKey)
		if err != nil {
			logging.Error("Offer signature fail to verify 2, with error: " + fmt.Sprintf("%v", err))
			continue
		}
		// Now Verify the merkle proof
		if offer.VerifyMerkleProof() != nil {
			logging.Error("Merkle proof verification failed.")
			continue
		}
		// Offer pass verification
		logging.Info("Offer pass every verification, added to result")

		validOffers = append(validOffers, offer)
		if len(validOffers) == maxOffers {
			break
		}
	}

	return validOffers, nil
}
