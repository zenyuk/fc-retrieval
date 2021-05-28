package fcrgatewayadmin

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

	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/register"
	"github.com/ConsenSys/fc-retrieval-gateway-admin/pkg/api/adminapi"
)

// FilecoinRetrievalGatewayAdmin is an example implementation using the api,
// which holds information about the interaction of the Filecoin
// Retrieval Gateway Admin with Filecoin Retrieval Gateways.
type FilecoinRetrievalGatewayAdmin struct {
	Settings GatewayAdminSettings

	// List of gateways this admin is in use
	ActiveGateways     map[string]register.GatewayRegister
	ActiveGatewaysLock sync.RWMutex
}

// NewFilecoinRetrievalGatewayAdmin initialise the Filecoin Retreival Gateway Admin library
func NewFilecoinRetrievalGatewayAdmin(settings GatewayAdminSettings) *FilecoinRetrievalGatewayAdmin {
	return &FilecoinRetrievalGatewayAdmin{
		Settings:           settings,
		ActiveGateways:     make(map[string]register.GatewayRegister),
		ActiveGatewaysLock: sync.RWMutex{},
	}
}

// InitialiseGateway initialise a given gateway
func (c *FilecoinRetrievalGatewayAdmin) InitialiseGateway(gatewayInfo *register.GatewayRegister, gatewayPrivKey *fcrcrypto.KeyPair, gatewayPrivKeyVer *fcrcrypto.KeyVersion) error {
	err := adminapi.RequestInitialiseKey(gatewayInfo, gatewayPrivKey, gatewayPrivKeyVer, c.Settings.gatewayAdminPrivateKey, c.Settings.gatewayAdminPrivateKeyVer)
	if err != nil {
		return err
	}

	// Register this gateway
	err = gatewayInfo.RegisterGateway(c.Settings.RegisterURL())
	if err != nil {
		logging.Error("Error in register the gateway.")
		return err
	}

	// Add this provider to the active gateways list
	c.ActiveGatewaysLock.Lock()
	c.ActiveGateways[gatewayInfo.NodeID] = *gatewayInfo
	c.ActiveGatewaysLock.Unlock()
	return nil
}

// InitialiseGatewayV2 initialise a given v2 gateway
func (c *FilecoinRetrievalGatewayAdmin) InitialiseGatewayV2(
	gatewayInfo *register.GatewayRegister,
	gatewayPrivKey *fcrcrypto.KeyPair,
	gatewayPrivKeyVer *fcrcrypto.KeyVersion,
	lotusWalletPrivateKey string,
	lotusAP string,
	lotusAuthToken string,
) error {
	err := adminapi.RequestInitialiseKeyV2(
		gatewayInfo,
		gatewayPrivKey,
		gatewayPrivKeyVer,
		c.Settings.gatewayAdminPrivateKey,
		c.Settings.gatewayAdminPrivateKeyVer,
		lotusWalletPrivateKey,
		lotusAP,
		lotusAuthToken,
	)
	if err != nil {
		return err
	}

	// Register this gateway
	err = gatewayInfo.RegisterGateway(c.Settings.RegisterURL())
	if err != nil {
		logging.Error("Error in register the gateway.")
		return err
	}

	// Add this gateway to the active gateways list
	c.ActiveGatewaysLock.Lock()
	c.ActiveGateways[gatewayInfo.NodeID] = *gatewayInfo
	c.ActiveGatewaysLock.Unlock()
	return nil
}

// ResetClientReputation requests a Gateway to initialise a client's reputation to the default value.
func (c *FilecoinRetrievalGatewayAdmin) ResetClientReputation(clientID *nodeid.NodeID) error {
	return errors.New("Not implemented yet")
}

// SetClientReputation requests a Gateway to set a client's reputation to a specified value.
func (c *FilecoinRetrievalGatewayAdmin) SetClientReputation(clientID *nodeid.NodeID, rep int64) error {
	return errors.New("Not implemented yet")
}

// GetCIDOffersList requests a Gateway's current list of CID Offers.
func (c *FilecoinRetrievalGatewayAdmin) GetCIDOffersList() ([]cidoffer.CIDOffer, error) {
	return nil, errors.New("Not implemented yet")
}

// ForceUpdate forces the provider to update its internal register
func (c *FilecoinRetrievalGatewayAdmin) ForceUpdate(gatewayID *nodeid.NodeID) error {
	c.ActiveGatewaysLock.RLock()
	defer c.ActiveGatewaysLock.RUnlock()
	gatewayInfo, exists := c.ActiveGateways[gatewayID.ToString()]
	if !exists {
		return errors.New("Unable to find the gateway in admin storage")
	}
	return adminapi.RequestForceRefresh(&gatewayInfo, c.Settings.gatewayAdminPrivateKey, c.Settings.gatewayAdminPrivateKeyVer)
}

// ListDHTOffer asks the gateway to list dht offer from providers
func (c *FilecoinRetrievalGatewayAdmin) ListDHTOffer(gatewayID *nodeid.NodeID) error {
	c.ActiveGatewaysLock.RLock()
	defer c.ActiveGatewaysLock.RUnlock()
	gatewayInfo, exists := c.ActiveGateways[gatewayID.ToString()]
	if !exists {
		return errors.New("Unable to find the gateway in admin storage")
	}
	return adminapi.RequestListDHTOffer(&gatewayInfo, c.Settings.gatewayAdminPrivateKey, c.Settings.gatewayAdminPrivateKeyVer)
}

func (c *FilecoinRetrievalGatewayAdmin) UpdateGatewaySupportedFeatures(gatewayInfo *register.GatewayRegister, providers []nodeid.NodeID) error {
	c.ActiveGatewaysLock.Lock()
	defer c.ActiveGatewaysLock.Unlock()
	err := adminapi.SetGroupCIDOfferSupportedForProviders(gatewayInfo, providers, c.Settings.gatewayAdminPrivateKey, c.Settings.gatewayAdminPrivateKeyVer)
	if err != nil {
		return err
	}

	c.ActiveGateways[gatewayInfo.NodeID] = *gatewayInfo
	return nil
}
