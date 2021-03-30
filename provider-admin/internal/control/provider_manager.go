package control

import (
	"errors"
	// "fmt"
	"sync"

	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages/fcrmsgpvdadmin"
	log "github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	req "github.com/ConsenSys/fc-retrieval-common/pkg/request"
	"github.com/ConsenSys/fc-retrieval-provider-admin/internal/settings"
	"github.com/ConsenSys/fc-retrieval-register/pkg/register"
)

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

// ProviderManager manages the pool of providers and the connections to them
type ProviderManager struct {
	settings *settings.ClientProviderAdminSettings

	ActiveProviders     map[string]register.ProviderRegister
	ActiveProvidersLock sync.RWMutex
}

// NewProviderManager returns an initialised instance of the provider manager
func NewProviderManager(settings *settings.ClientProviderAdminSettings) *ProviderManager {
	return &ProviderManager{
		settings:            settings,
		ActiveProviders:     make(map[string]register.ProviderRegister),
		ActiveProvidersLock: sync.RWMutex{},
	}
}

// InitialiseProvider initialise a given provider
func (p *ProviderManager) InitialiseProvider(providerInfo *register.ProviderRegister, providerPrivKey *fcrcrypto.KeyPair, providerPrivKeyVer *fcrcrypto.KeyVersion) error {
	// TODO: Check given providerInfo is correct
	// First, Get pubkey
	pubKey, err := providerInfo.GetSigningKey()
	if err != nil {
		log.Error("Error in obtaining signing key from register info.")
		return err
	}

	nodeID, err := nodeid.NewNodeIDFromHexString(providerInfo.NodeID)
	if err != nil {
		log.Error("Error in generating nodeID.")
		return err
	}

	// Second, send key exchange to activate the given provider
	request, err := fcrmsgpvdadmin.EncodeProviderAdminInitialiseKeyRequest(nodeID, providerPrivKey, providerPrivKeyVer)
	if err != nil {
		log.Error("Error in encoding message.")
		return err
	}

	// Sign the request
	if request.Sign(p.settings.ProviderAdminPrivateKey(), p.settings.ProviderAdminPrivateKeyVer()) != nil {
		return errors.New("Error in signing the request")
	}

	response, err := req.SendMessage("http://"+providerInfo.NetworkInfoAdmin+"/v1", request)
	// response, err := SendMessage(providerInfo.NetworkInfoAdmin, request)
	if err != nil {
		log.Error("Error in sending the message.")
		return err
	}

	// Verify the response
	if response.Verify(pubKey) != nil {
		return errors.New("Fail to verify the response")
	}

	ok, err := fcrmsgpvdadmin.DecodeProviderAdminInitialiseKeyResponse(response)
	if err != nil {
		log.Error("Error in decoding the message.")
		return err
	}
	if !ok {
		log.Error("Initialise provider failed.")
		return errors.New("Fail to initialise provider")
	}

	// Finally register the provider
	err = providerInfo.RegisterProvider(p.settings.RegisterURL())
	if err != nil {
		log.Error("Error in register the provider.")
		return err
	}

	// Finally add the provider to the active providers list
	p.ActiveProvidersLock.Lock()
	p.ActiveProviders[providerInfo.NodeID] = *providerInfo
	p.ActiveProvidersLock.Unlock()
	return nil
}

// PublishGroupCID publish a group cid offer to a given provider
func (p *ProviderManager) PublishGroupCID(providerID *nodeid.NodeID, cids []cid.ContentID, price uint64, expiry int64, qos uint64) error {
	request, err := fcrmsgpvdadmin.EncodeProviderAdminPublishGroupOfferRequest(cids, price, expiry, qos)
	// Sign the request
	if request.Sign(p.settings.ProviderAdminPrivateKey(), p.settings.ProviderAdminPrivateKeyVer()) != nil {
		return errors.New("Error in signing the request")
	}

	response, err := p.SendMessage(providerID, request)
	if err != nil {
		log.Error("Error in sending the message.")
		return err
	}
	// Verify the response
	// Get pubKey
	p.ActiveProvidersLock.RLock()
	pubKey, err := fcrcrypto.DecodePublicKey(p.ActiveProviders[providerID.ToString()].SigningKey)
	if err != nil {
		return err
	}
	p.ActiveProvidersLock.RUnlock()
	if response.Verify(pubKey) != nil {
		return errors.New("Fail to verify the response")
	}

	received, err := fcrmsgpvdadmin.DecodeProviderAdminPublishGroupOfferResponse(response)
	if err != nil {
		log.Error("Error in decoding the message.")
		return err
	}
	if !received {
		log.Error("Publish offer failed.")
		return errors.New("Fail to publish offer")
	}
	return nil
}

// PublishDHTCID publish a dht cid offer to a given provider
func (p *ProviderManager) PublishDHTCID(providerID *nodeid.NodeID, cids []cid.ContentID, price []uint64, expiry []int64, qos []uint64) error {
	request, err := fcrmsgpvdadmin.EncodeProviderAdminPublishDHTOfferRequest(cids, price, expiry, qos)
	// Sign the request
	if request.Sign(p.settings.ProviderAdminPrivateKey(), p.settings.ProviderAdminPrivateKeyVer()) != nil {
		return errors.New("Error in signing the request")
	}

	response, err := p.SendMessage(providerID, request)
	if err != nil {
		log.Error("Error in sending the message.")
		return err
	}
	// Verify the response
	// Get pubKey
	p.ActiveProvidersLock.RLock()
	pubKey, err := fcrcrypto.DecodePublicKey(p.ActiveProviders[providerID.ToString()].SigningKey)
	if err != nil {
		return err
	}
	p.ActiveProvidersLock.RUnlock()
	if response.Verify(pubKey) != nil {
		return errors.New("Fail to verify the response")
	}

	received, err := fcrmsgpvdadmin.DecodeProviderAdminPublishDHTOfferResponse(response)
	if err != nil {
		log.Error("Error in decoding the message.")
		return err
	}
	if !received {
		log.Error("Publish offer failed.")
		return errors.New("Fail to publish offer")
	}
	return nil
}

// GetGroupCIDOffer checks the group offer stored in the provider
func (p *ProviderManager) GetGroupCIDOffer(providerID *nodeid.NodeID, gatewayIDs []nodeid.NodeID) (
	bool, // found
	[]cidoffer.CIDOffer, // offers
	error, // error
) {
	request, err := fcrmsgpvdadmin.EncodeProviderAdminGetPublishedOfferRequest(gatewayIDs)
	// Sign the request
	if request.Sign(p.settings.ProviderAdminPrivateKey(), p.settings.ProviderAdminPrivateKeyVer()) != nil {
		log.Error("Error in signing the request.")
		return false, nil, err
	}

	response, err := p.SendMessage(providerID, request)
	if err != nil {
		log.Error("Error in sending the message.")
		return false, nil, err
	}

	// Verify the response
	// Get pubKey
	p.ActiveProvidersLock.RLock()
	pubKey, err := fcrcrypto.DecodePublicKey(p.ActiveProviders[providerID.ToString()].SigningKey)
	if err != nil {
		return false, nil, err
	}
	p.ActiveProvidersLock.RUnlock()
	if response.Verify(pubKey) != nil {
		return false, nil, errors.New("Fail to verify the response")
	}

	found, offers, err := fcrmsgpvdadmin.DecodeProviderAdminGetPublishedOfferResponse(response)
	if err != nil {
		log.Error("Error in decoding the message")
		return false, nil, err
	}
	return found, offers, nil
}

// SendMessage sends a message to a managed provider and obtain a response
func (p *ProviderManager) SendMessage(providerID *nodeid.NodeID, message *fcrmessages.FCRMessage) (*fcrmessages.FCRMessage, error) {
	log.Info("Provider Manager sending message to providerID: %v", providerID.ToString())
	p.ActiveProvidersLock.RLock()
	defer p.ActiveProvidersLock.RUnlock()
	providerRegister, ok := p.ActiveProviders[providerID.ToString()]
	if !ok {
		log.Error("Provider not found in the provider manager, please initialise it first.")
		return nil, errors.New("Provider not found")
	}
	data, err := req.SendMessage("http://"+providerRegister.NetworkInfoAdmin+"/v1", message)
	if err != nil {
		return nil, err
	}
	return data, nil
}
