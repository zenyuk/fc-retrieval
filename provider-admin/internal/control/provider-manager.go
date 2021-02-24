package control

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
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	// "time"

	// "github.com/ConsenSys/fc-retrieval-gateway/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrtcpcomms"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/register"

	"github.com/ConsenSys/fc-retrieval-provider-admin/internal/providerapi"
	"github.com/ConsenSys/fc-retrieval-provider-admin/internal/settings"
)

// ProviderManager managers the pool of providers and the connections to them.
type ProviderManager struct {
	settings  settings.ClientSettings
	providers []ActiveProvider
	// RegisteredProvidersMap stores mapping from provider id (big int in string repr) to its registration info
	RegisteredProvidersMap     map[string]register.RegisteredNode
	RegisteredProvidersMapLock sync.RWMutex

	// ProviderCommPool manages connection for outgoing request to providers
	ProviderCommPool *fcrtcpcomms.CommunicationPool
}

// ActiveProvider contains information for a single provider
type ActiveProvider struct {
	info  register.ProviderRegister
	comms *providerapi.Comms
}

// NewProviderManager returns an initialised instance of the provider manager.
func NewProviderManager(settings settings.ClientSettings) *ProviderManager {
	p := ProviderManager{
		RegisteredProvidersMap:     make(map[string]register.RegisteredNode),
		RegisteredProvidersMapLock: sync.RWMutex{},
	}
	p.settings = settings
	p.ProviderCommPool = fcrtcpcomms.NewCommunicationPool(p.RegisteredProvidersMap, &p.RegisteredProvidersMapLock)
	p.providerManagerRunner()
	return &p
}

// TODO this should be in a go routine and loop for ever.
func (p *ProviderManager) providerManagerRunner() {
	logging.Info("Provider Manager: Management thread started")

	providers, err := register.GetRegisteredProviders(p.settings.RegisterURL())
	if err != nil {
		logging.Error("Not provider registered")
		return
	}
	for _, provider := range providers {
		providerID, err := nodeid.NewNodeIDFromString(provider.NodeID)
		if err != nil {
			logging.Error("Error with nodeID %v: %v", provider.NodeID, err)
			continue
		}
		p.RegisteredProvidersMap[providerID.ToString()] = &provider
	}

	// TODO this loop is where the managing of providers that the client is using happens.
	logging.Info("Provider Manager: GetProviders returned %d providers", len(providers))
	// for _, provider := range providers {
	// 	logging.Info("Setting-up comms with: %+v", provider)
	// 	comms, err := providerapi.NewProviderAPIComms(&provider, &p.settings)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	activeProvider := ActiveProvider{provider, comms}
	// 	p.providers = append(p.providers, activeProvider)
	// }

	// logging.Info("Provider Manager using %d providers", len(p.providers))
}

// BlockProvider adds a host to disallowed list of providers
func (p *ProviderManager) BlockProvider(hostName string) {
	// TODO
}

// UnblockProvider add a host to allowed list of providers
func (p *ProviderManager) UnblockProvider(hostName string) {
	// TODO

}

// FindOffersStandardDiscovery finds offers using the standard discovery mechanism.
// func (p *ProviderManager) FindOffersStandardDiscovery(contentID *cid.ContentID) ([]cidoffer.CidGroupOffer, error) {
// 	if len(g.providers) == 0 {
// 		return nil, fmt.Errorf("No providers available")
// 	}

// 	var aggregateOffers []cidoffer.CidGroupOffer
// 	for _, gw := range g.providers {
// 		// TODO need to do nonce management
// 		// TODO need to do requests to all providers in parallel, rather than serially
// 		offers, err := gw.comms.ProviderStdCIDDiscovery(contentID, 1)
// 		if err != nil {
// 			logging.Warn("ProviderStdDiscovery error. Provider: %s, Error: %s", gw.info.NodeID, err)
// 		}
// 		// TODO: probably should remove duplicate offers at this point
// 		aggregateOffers = append(aggregateOffers, offers...)
// 	}
// 	return aggregateOffers, nil
// }

// GetConnectedProviders returns the list of domain names of providers that the client
// is currently connected to.
func (p *ProviderManager) GetConnectedProviders() []string {
	urls := make([]string, len(p.providers))
	for i, provider := range p.providers {
		urls[i] = provider.comms.ApiURL
	}
	return urls
}

// Shutdown stops go routines and closes sockets. This should be called as part
// of the graceful library shutdown
func (p *ProviderManager) Shutdown() {
	// TODO
}

// SendMessage send message to providers
func (p *ProviderManager) SendMessage(message *fcrmessages.FCRMessage) (
	*fcrmessages.FCRMessage,
	error,
) {
	mJSON, _ := json.Marshal(message)
	logging.Info("JSON sent: %s", string(mJSON))
	contentReader := bytes.NewReader(mJSON)
	providerRegister := p.settings.ProviderRegister()
	fmt.Printf("Provider registered: %+v\n", providerRegister)
	req, _ := http.NewRequest("POST", "http://"+providerRegister.NetworkInfoAdmin+"/v1", contentReader)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		logging.Fatal("Error: %+v", err)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	var data fcrmessages.FCRMessage
	json.NewDecoder(res.Body).Decode(&data)
	return &data, nil
}
