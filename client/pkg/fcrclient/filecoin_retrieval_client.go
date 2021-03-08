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
	"github.com/ConsenSys/fc-retrieval-client/internal/control"
	"github.com/ConsenSys/fc-retrieval-client/internal/settings"
	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
)

// FilecoinRetrievalClient holds information about the interaction of
// the Filecoin Retrieval Client with Filecoin Retrieval Gateways.
type FilecoinRetrievalClient struct {
	clientManager *control.ClientManager
	// TODO have a list of gateway objects of all the current gateways being interacted with
}

// NewFilecoinRetrievalClient initialise the Filecoin Retreival Client library
func NewFilecoinRetrievalClient(conf Settings) (*FilecoinRetrievalClient, error) {
	logging.Info("Filecoin Retrieval Client started")
	var c = FilecoinRetrievalClient{}
	clientSettings := conf.(*settings.ClientSettings)
	c.clientManager = control.NewClientManager(*clientSettings)
	return &c, nil
}

// RefreshLatestProviderInfo refreshes the provider info
func (c *FilecoinRetrievalClient) RefreshLatestProviderInfo() {
	logging.Info("Filecoin Retrieval Client refresh provider info")
	c.clientManager.RefreshLatestProviderInfo()
}

// RefreshLatestGatewayInfo refreshes the gateway info
func (c *FilecoinRetrievalClient) RefreshLatestGatewayInfo() {
	logging.Info("Filecoin Retrieval Client refresh gateway info")
	c.clientManager.RefreshLatestGatewayInfo()
}

// FindGateways find gateways located near too the specified location. Use AddGateways
// to use these gateways.
func (c *FilecoinRetrievalClient) FindGateways(location string, maxNumToLocate int) ([]*nodeid.NodeID, error) {
	logging.Info("Find gateways")
	return c.clientManager.FindGateways(location, maxNumToLocate)
}

// AddGateways adds one or more gateways to use.
func (c *FilecoinRetrievalClient) AddGateways(gwNodeIDs []*nodeid.NodeID) int {
	logging.Info("Add gateways")
	return c.clientManager.AddGateways(gwNodeIDs)
}

// RemoveGateways removes one or more gateways from the list of Gateways to use.
func (c *FilecoinRetrievalClient) RemoveGateways(gwNodeIDs []*nodeid.NodeID) int {
	logging.Info("Remove gateways")
	return c.clientManager.RemoveGateways(gwNodeIDs)
}

// RemoveAllGateways removes all gateways from the list of Gateways to use.
func (c *FilecoinRetrievalClient) RemoveAllGateways() int {
	logging.Info("Remove all gateways")
	return c.clientManager.RemoveAllGateways()
}

// GetGateways returns the list of gateways that are being used.
func (c *FilecoinRetrievalClient) GetGateways() []*nodeid.NodeID {
	logging.Info("Get gateways")
	return c.clientManager.GetGateways()
}

// FindOffersStandardDiscovery finds offer using standard discovery
func (c *FilecoinRetrievalClient) FindOffersStandardDiscovery(contentID *(cid.ContentID)) ([]cidoffer.CidGroupOffer, error) {
	logging.Info("Find offers std discovery")
	return c.clientManager.FindOffersStandardDiscovery(contentID)
}

// // FindBestOffers locates offsers for supplying the content associated with the pieceCID
// func (c *FilecoinRetrievalClient) FindBestOffers(pieceCID [32]byte, maxPrice uint64, maxExpectedLatency int64) ([]cidoffer.CidGroupOffer, error) {
// 	cid := cid.NewContentIDFromBytes(pieceCID[:])
// 	logging.Trace("FindBestOffers(pieceCID: %s, maxPrice: %d, maxExpectedLatency: %d",
// 		cid.ToString(), maxPrice, maxExpectedLatency)

// 	rawOffers, err := c.gatewayManager.FindOffersStandardDiscovery(cid)
// 	if err != nil {
// 		return nil, err
// 	}
// 	logging.Trace("FindBestOffers(pieceCID: %s) offers found before filtering: %d", cid.ToString(), len(rawOffers))
// 	var offers []cidoffer.CidGroupOffer
// 	for _, offer := range rawOffers {
// 		if offer.Price < maxPrice {
// 			offers = append(offers, offer)
// 		}
// 		// TODO: need to have latency filter.
// 	}

// 	logging.Info("FindBestOffers(pieceCID: %s) found %d offers", cid.ToString(), len(offers))
// 	return offers, nil
// }

// // ConnectedGateways returns a slice of the URLs for the gateways this client is connected to.
// func (c *FilecoinRetrievalClient) ConnectedGateways() []string {
// 	return c.gatewayManager.GetConnectedGateways()
// }

// // Shutdown releases all resources used by the library
// func (c *FilecoinRetrievalClient) Shutdown() {
// 	logging.Info("Filecoin Retrieval Client shutting down")
// 	c.gatewayManager.Shutdown()
// }
