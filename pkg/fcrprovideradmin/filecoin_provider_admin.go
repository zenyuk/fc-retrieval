package fcrprovideradmin

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
	"sync"

	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/register"
	"github.com/ConsenSys/fc-retrieval-provider-admin/pkg/api/adminapi"
)

// providerInfo stores the information of a remote provider instance
type providerInfo struct {
	adminAP string

	id      *nodeid.NodeID
	privKey *fcrcrypto.KeyPair
	pubKey  *fcrcrypto.KeyPair
}

// FilecoinRetrievalProviderAdmin is an example implementation using the api,
// which holds information about the interaction of the Filecoin
// Retrieval Provider Admin with Filecoin Retrieval Providers.
type FilecoinRetrievalProviderAdmin struct {
	Settings ProviderAdminSettings

	// List of providers this admin is in use
	ActiveProviders     map[string]*providerInfo
	ActiveProvidersLock sync.RWMutex
}

// NewFilecoinRetrievalProviderAdmin initialise the Filecoin Retreival Provider Admin library
func NewFilecoinRetrievalProviderAdmin(settings ProviderAdminSettings) *FilecoinRetrievalProviderAdmin {
	return &FilecoinRetrievalProviderAdmin{
		Settings:            settings,
		ActiveProviders:     make(map[string]*providerInfo),
		ActiveProvidersLock: sync.RWMutex{},
	}
}

// InitialiseProvider initialise a given provider
func (c *FilecoinRetrievalProviderAdmin) InitialiseProvider(adminAP string, nodeID *nodeid.NodeID, providerPrivKey *fcrcrypto.KeyPair, providerPrivKeyVer *fcrcrypto.KeyVersion) error {
	err := adminapi.RequestInitialiseKey(adminAP, nodeID, providerPrivKey, providerPrivKeyVer, c.Settings.providerAdminPrivateKey, c.Settings.providerAdminPrivateKeyVer)
	if err != nil {
		return err
	}

	// Get pub key
	encoded, err := providerPrivKey.EncodePublicKey()
	if err != nil {
		logging.Error("Error in encoding public key")
		return err
	}
	pubKey, err := fcrcrypto.DecodePublicKey(encoded)
	if err != nil {
		logging.Error("Error in generating signing key.")
		return err
	}

	// Add this provider to the active providers list
	c.ActiveProvidersLock.Lock()
	c.ActiveProviders[nodeID.ToString()] = &providerInfo{
		adminAP: adminAP,
		id:      nodeID,
		privKey: providerPrivKey,
		pubKey:  pubKey,
	}
	c.ActiveProvidersLock.Unlock()
	return nil
}

// RegisterProvider registers the given providerInfo to the register
func (c *FilecoinRetrievalProviderAdmin) RegisterProvider(providerInfo *register.ProviderRegister) error {
	// Register this provider
	err := providerInfo.RegisterProvider(c.Settings.RegisterURL())
	if err != nil {
		logging.Error("Error in register the provider.")
	}
	return err
}

// PublishGroupCID publish a group cid offer to a given provider
func (c *FilecoinRetrievalProviderAdmin) PublishGroupCID(providerID *nodeid.NodeID, cids []cid.ContentID, price uint64, expiry int64, qos uint64) error {
	c.ActiveProvidersLock.RLock()
	defer c.ActiveProvidersLock.RUnlock()
	providerInfo, exists := c.ActiveProviders[providerID.ToString()]
	if !exists {
		return errors.New("Unable to find the provider in admin storage")
	}
	return adminapi.RequestPublishGroupOffer(providerInfo.adminAP, providerInfo.pubKey, cids, price, expiry, qos, c.Settings.providerAdminPrivateKey, c.Settings.providerAdminPrivateKeyVer)
}

// PublishDHTCID publish a dht cid offer to a given provider
func (c *FilecoinRetrievalProviderAdmin) PublishDHTCID(providerID *nodeid.NodeID, cids []cid.ContentID, price []uint64, expiry []int64, qos []uint64) error {
	c.ActiveProvidersLock.RLock()
	defer c.ActiveProvidersLock.RUnlock()
	providerInfo, exists := c.ActiveProviders[providerID.ToString()]
	if !exists {
		return errors.New("Unable to find the provider in admin storage")
	}
	return adminapi.RequestPublishDHTOffer(providerInfo.adminAP, providerInfo.pubKey, cids, price, expiry, qos, c.Settings.providerAdminPrivateKey, c.Settings.providerAdminPrivateKeyVer)
}

// GetGroupCIDOffer checks the group offer stored in the provider
func (c *FilecoinRetrievalProviderAdmin) GetGroupCIDOffer(providerID *nodeid.NodeID, gatewayIDs []nodeid.NodeID) (bool, []cidoffer.CIDOffer, error) {
	c.ActiveProvidersLock.RLock()
	defer c.ActiveProvidersLock.RUnlock()
	providerInfo, exists := c.ActiveProviders[providerID.ToString()]
	if !exists {
		return false, nil, errors.New("Unable to find the provider in admin storage")
	}
	return adminapi.RequestGetPublishedOffer(providerInfo.adminAP, providerInfo.pubKey, gatewayIDs, c.Settings.providerAdminPrivateKey, c.Settings.providerAdminPrivateKeyVer)
}

// ForceUpdate forces the provider to update its internal register
func (c *FilecoinRetrievalProviderAdmin) ForceUpdate(providerID *nodeid.NodeID) error {
	c.ActiveProvidersLock.RLock()
	defer c.ActiveProvidersLock.RUnlock()
	providerInfo, exists := c.ActiveProviders[providerID.ToString()]
	if !exists {
		return errors.New("Unable to find the provider in admin storage")
	}
	return adminapi.RequestForceRefresh(providerInfo.adminAP, providerInfo.pubKey, c.Settings.providerAdminPrivateKey, c.Settings.providerAdminPrivateKeyVer)
}
