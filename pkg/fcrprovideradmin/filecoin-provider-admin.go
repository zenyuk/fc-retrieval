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
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/register"
	"github.com/ConsenSys/fc-retrieval-provider-admin/internal/control"
	"github.com/ConsenSys/fc-retrieval-provider-admin/internal/settings"
)

// FilecoinRetrievalProviderAdminClient holds information about the interaction of
// the Filecoin Retrieval Provider Admin Client with Filecoin Retrieval Providers.
type FilecoinRetrievalProviderAdminClient struct {
	Settings        settings.ClientSettings
	providerManager *control.ProviderManager
	// TODO have a list of provider objects of all the current providers being interacted with
}

var singleInstance *FilecoinRetrievalProviderAdminClient
var initialised = false

// InitFilecoinRetrievalProviderAdminClient initialise the Filecoin Retreival Client library
func InitFilecoinRetrievalProviderAdminClient(settings Settings) *FilecoinRetrievalProviderAdminClient {
	if initialised {
		logging.ErrorAndPanic("Attempt to init Filecoin Retrieval Provider Admin Client a second time")
	}
	var c = FilecoinRetrievalProviderAdminClient{}
	c.startUp(settings)
	singleInstance = &c
	initialised = true
	return singleInstance

}

// GetFilecoinRetrievalProviderAdminClient creates a Filecoin Retrieval Provider Admin Client
func GetFilecoinRetrievalProviderAdminClient() *FilecoinRetrievalProviderAdminClient {
	if !initialised {
		logging.ErrorAndPanic("Filecoin Retrieval Provider Admin Client not initialised")
	}

	return singleInstance
}

func (c *FilecoinRetrievalProviderAdminClient) startUp(conf Settings) {
	logging.Info("Filecoin Retrieval Provider Admin Client started")
	clientSettings := conf.(*settings.ClientSettings)
	c.Settings = *clientSettings
	c.providerManager = control.NewProviderManager(*clientSettings)
}

// Shutdown releases all resources used by the library
func (c *FilecoinRetrievalProviderAdminClient) Shutdown() {
	logging.Info("Filecoin Retrieval Provider Admin Client shutting down")
	c.providerManager.Shutdown()
}

// RegisterProvider to register a provider
func (c *FilecoinRetrievalProviderAdminClient) RegisterProvider() error {
	logging.Info("Filecoin Retrieval Provider Admin Client sending message")
	url := c.Settings.RegisterURL()
	reg := c.Settings.ProviderRegister()
	err := reg.RegisterProvider(url)
	if err != nil {
		logging.Error("Unable to register provider")
		return err
	}
	return nil
}

// SendMessage send message to providers
func (c *FilecoinRetrievalProviderAdminClient) SendMessage(message *fcrmessages.FCRMessage) (
	*fcrmessages.FCRMessage,
	error,
) {
	logging.Info("Filecoin Retrieval Provider Admin Client sending message")
	return c.providerManager.SendMessage(message)
}

// GetRegisteredGateways send message to providers
func (c *FilecoinRetrievalProviderAdminClient) GetRegisteredGateways() ([]register.GatewayRegister, error) {
	logging.Info("Filecoin Retrieval Provider Admin Client getting registered gateways")
	gateways, err := register.GetRegisteredGateways(c.Settings.RegisterURL())
	if err != nil {
		logging.Error("Unable to get registered gateways: %v", err)
		return []register.GatewayRegister{}, err
	}
	return gateways, nil
}
