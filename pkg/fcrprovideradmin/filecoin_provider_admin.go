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
	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/register"
	"github.com/ConsenSys/fc-retrieval-provider-admin/internal/control"
	"github.com/ConsenSys/fc-retrieval-provider-admin/internal/settings"
)

// FilecoinRetrievalProviderAdminClient holds information about the interaction of
// the Filecoin Retrieval Provider Admin Client with Filecoin Retrieval Providers.
type FilecoinRetrievalProviderAdminClient struct {
	providerManager *control.ProviderManager
	// TODO have a list of provider objects of all the current providers being interacted with
}

var singleInstance *FilecoinRetrievalProviderAdminClient
var initialised = false

// InitFilecoinRetrievalProviderAdminClient initialise the Filecoin Retreival Client library
func InitFilecoinRetrievalProviderAdminClient(conf Settings) *FilecoinRetrievalProviderAdminClient {
	if initialised {
		logging.ErrorAndPanic("Attempt to init Filecoin Retrieval Provider Admin Client a second time")
	}
	var c = FilecoinRetrievalProviderAdminClient{}
	singleInstance = &c
	initialised = true

	clientSettings := conf.(*settings.ClientProviderAdminSettings)
	c.providerManager = control.NewProviderManager(clientSettings)
	return singleInstance

}

// GetFilecoinRetrievalProviderAdminClient creates a Filecoin Retrieval Provider Admin Client
func GetFilecoinRetrievalProviderAdminClient() *FilecoinRetrievalProviderAdminClient {
	if !initialised {
		logging.ErrorAndPanic("Filecoin Retrieval Provider Admin Client not initialised")
	}

	return singleInstance
}

// InitialiseProvider initialise a given provider
func (c *FilecoinRetrievalProviderAdminClient) InitialiseProvider(providerInfo *register.ProviderRegister, providerPrivKey *fcrcrypto.KeyPair, providerPrivKeyVer *fcrcrypto.KeyVersion) error {
	return c.providerManager.InitialiseProvider(providerInfo, providerPrivKey, providerPrivKeyVer)
}

// PublishGroupCID publish a group cid offer to a given provider
func (c *FilecoinRetrievalProviderAdminClient) PublishGroupCID(providerID *nodeid.NodeID, cids []cid.ContentID, price uint64, expiry int64, qos uint64) error {
	return c.providerManager.PublishGroupCID(providerID, cids, price, expiry, qos)
}

// PublishDHTCID publish a dht cid offer to a given provider
func (c *FilecoinRetrievalProviderAdminClient) PublishDHTCID(providerID *nodeid.NodeID, cids []cid.ContentID, price []uint64, expiry []int64, qos []uint64) error {
	return c.providerManager.PublishDHTCID(providerID, cids, price, expiry, qos)
}

// GetGroupCIDOffer checks the group offer stored in the provider
func (c *FilecoinRetrievalProviderAdminClient) GetGroupCIDOffer(providerID *nodeid.NodeID, gatewayIDs []nodeid.NodeID) (bool, []fcrmessages.CIDGroupInformation, error) {
	return c.providerManager.GetGroupCIDOffer(providerID, gatewayIDs)
}
