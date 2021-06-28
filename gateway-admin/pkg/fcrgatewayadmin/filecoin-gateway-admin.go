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

  "github.com/ConsenSys/fc-retrieval/common/pkg/cidoffer"
  "github.com/ConsenSys/fc-retrieval/common/pkg/fcrcrypto"
  "github.com/ConsenSys/fc-retrieval/common/pkg/nodeid"
  "github.com/ConsenSys/fc-retrieval/common/pkg/register"
  "github.com/ConsenSys/fc-retrieval/gateway-admin/pkg/api/adminapi"
)

// FilecoinRetrievalGatewayAdmin is an example implementation using the api,
// which holds information about the interaction of the Filecoin
// Retrieval Gateway Admin with Filecoin Retrieval Gateways.
type FilecoinRetrievalGatewayAdmin struct {
	Settings GatewayAdminSettings

	// List of gateways this admin is in use
	ActiveGateways     map[string]register.GatewayRegistrar
	ActiveGatewaysLock sync.RWMutex

  AdminApiCaller      adminapi.AdminApi
}

// NewFilecoinRetrievalGatewayAdmin initialise the Filecoin Retreival Gateway Admin library
func NewFilecoinRetrievalGatewayAdmin(settings GatewayAdminSettings) *FilecoinRetrievalGatewayAdmin {
	return &FilecoinRetrievalGatewayAdmin{
		Settings:           settings,
		ActiveGateways:     make(map[string]register.GatewayRegistrar),
		ActiveGatewaysLock: sync.RWMutex{},
    AdminApiCaller:     adminapi.NewAdminApi(),
	}
}

// InitialiseGateway initialise a given gateway
func (c *FilecoinRetrievalGatewayAdmin) InitialiseGateway(gatewayRegistrar register.GatewayRegistrar, gatewayPrivKey *fcrcrypto.KeyPair, gatewayPrivKeyVer *fcrcrypto.KeyVersion) error {
	err := c.AdminApiCaller.RequestInitialiseKey(gatewayRegistrar, gatewayPrivKey, gatewayPrivKeyVer, c.Settings.gatewayAdminPrivateKey, c.Settings.gatewayAdminPrivateKeyVer)
	if err != nil {
		return err
	}

	// Add this gateway to the active gateways list
	c.ActiveGatewaysLock.Lock()
	c.ActiveGateways[gatewayRegistrar.GetNodeID()] = gatewayRegistrar
  c.ActiveGatewaysLock.Unlock()
	return nil
}

// InitialiseGatewayV2 initialise a given v2 gateway
func (c *FilecoinRetrievalGatewayAdmin) InitialiseGatewayV2(
  gatewayRegistrar register.GatewayRegistrar,
	gatewayPrivKey *fcrcrypto.KeyPair,
	gatewayPrivKeyVer *fcrcrypto.KeyVersion,
	lotusWalletPrivateKey string,
	lotusAP string,
	lotusAuthToken string,
) error {
	err := c.AdminApiCaller.RequestInitialiseKeyV2(
    gatewayRegistrar,
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

	// Add this gateway to the active gateways list
	c.ActiveGatewaysLock.Lock()
	c.ActiveGateways[gatewayRegistrar.GetNodeID()] = gatewayRegistrar
	c.ActiveGatewaysLock.Unlock()
	return nil
}

// ResetClientReputation requests a Gateway to initialise a client's reputation to the default value.
func (c *FilecoinRetrievalGatewayAdmin) ResetClientReputation(_ *nodeid.NodeID) error {
	return errors.New("not implemented yet")
}

// SetClientReputation requests a Gateway to set a client's reputation to a specified value.
func (c *FilecoinRetrievalGatewayAdmin) SetClientReputation(_ *nodeid.NodeID, _ int64) error {
	return errors.New("not implemented yet")
}

// GetCIDOffersList requests a Gateway's current list of CID Offers.
func (c *FilecoinRetrievalGatewayAdmin) GetCIDOffersList() ([]cidoffer.CIDOffer, error) {
	return nil, errors.New("not implemented yet")
}

// ForceUpdate forces the provider to update its internal register
func (c *FilecoinRetrievalGatewayAdmin) ForceUpdate(gatewayID *nodeid.NodeID) error {
	c.ActiveGatewaysLock.RLock()
	defer c.ActiveGatewaysLock.RUnlock()
  gatewayRegistrar, exists := c.ActiveGateways[gatewayID.ToString()]
	if !exists {
		return errors.New("unable to find the gateway in admin storage")
	}
	return c.AdminApiCaller.RequestForceRefresh(gatewayRegistrar, c.Settings.gatewayAdminPrivateKey, c.Settings.gatewayAdminPrivateKeyVer)
}

// ListDHTOffer asks the gateway to list dht offer from providers
func (c *FilecoinRetrievalGatewayAdmin) ListDHTOffer(gatewayID *nodeid.NodeID) error {
	c.ActiveGatewaysLock.RLock()
	defer c.ActiveGatewaysLock.RUnlock()
  gatewayRegistrar, exists := c.ActiveGateways[gatewayID.ToString()]
	if !exists {
		return errors.New("unable to find the gateway in admin storage")
	}
	return c.AdminApiCaller.RequestListDHTOffer(gatewayRegistrar, c.Settings.gatewayAdminPrivateKey, c.Settings.gatewayAdminPrivateKeyVer)
}

func (c *FilecoinRetrievalGatewayAdmin) UpdateGatewaySupportedFeatures(gatewayRegistrar register.GatewayRegistrar, providers []nodeid.NodeID) error {
	c.ActiveGatewaysLock.Lock()
	defer c.ActiveGatewaysLock.Unlock()
	err := c.AdminApiCaller.SetGroupCIDOfferSupportedForProviders(gatewayRegistrar, providers, c.Settings.gatewayAdminPrivateKey, c.Settings.gatewayAdminPrivateKeyVer)
	if err != nil {
		return err
	}

	c.ActiveGateways[gatewayRegistrar.GetNodeID()] = gatewayRegistrar
	return nil
}
