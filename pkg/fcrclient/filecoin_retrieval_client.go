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
func NewFilecoinRetrievalClient(conf Settings) *FilecoinRetrievalClient {
	logging.Info("Filecoin Retrieval Client started")
	var c = FilecoinRetrievalClient{}
	clientSettings := conf.(*settings.ClientSettings)
	c.clientManager = control.NewClientManager(*clientSettings)
	return &c
}

// FindGateways find gateways located near too the specified location. Use AddGateways
// to use these gateways.
func (c *FilecoinRetrievalClient) FindGateways(location string, maxNumToLocate int) ([]*nodeid.NodeID, error) {
	logging.Info("Find gateways")
	return c.clientManager.FindGateways(location, maxNumToLocate)
}

// AddGatewaysToUse adds one or more gateways to use
func (c *FilecoinRetrievalClient) AddGatewaysToUse(gwNodeIDs []*nodeid.NodeID) int {
	logging.Info("Add gateways to use")
	return c.clientManager.AddGatewaysToUse(gwNodeIDs)
}

// RemoveGatewaysToUse removes one or more gateways from the list of gateways to use
func (c *FilecoinRetrievalClient) RemoveGatewaysToUse(gwNodeIDs []*nodeid.NodeID) int {
	logging.Info("Remove gateways to use")
	return c.clientManager.RemoveGatewaysToUse(gwNodeIDs)
}

// RemoveAllGatewaysToUse removes all gateways from the list of gateways to use
func (c *FilecoinRetrievalClient) RemoveAllGatewaysToUse() int {
	logging.Info("Remove all gateways to use")
	return c.clientManager.RemoveAllGatewaysToUse()
}

// GetGatewaysToUse returns the list of gateways that can be used
func (c *FilecoinRetrievalClient) GetGatewaysToUse() []*nodeid.NodeID {
	logging.Info("Get gateways to use")
	return c.clientManager.GetGatewaysToUse()
}

// AddActiveGateways adds one or more gateways to active.
func (c *FilecoinRetrievalClient) AddActiveGateways(gwNodeIDs []*nodeid.NodeID) int {
	logging.Info("Add active gateways")
	return c.clientManager.AddActiveGateways(gwNodeIDs)
}

// RemoveActiveGateways removes one or more gateways from the list of active gateways.
func (c *FilecoinRetrievalClient) RemoveActiveGateways(gwNodeIDs []*nodeid.NodeID) int {
	logging.Info("Remove active gateways")
	return c.clientManager.RemoveActiveGateways(gwNodeIDs)
}

// RemoveAllActiveGateways removes all gateways from the list of active gateways.
func (c *FilecoinRetrievalClient) RemoveAllActiveGateways() int {
	logging.Info("Remove all active gateways")
	return c.clientManager.RemoveAllActiveGateways()
}

// GetActiveGateways returns the list of gateways that are being used.
func (c *FilecoinRetrievalClient) GetActiveGateways() []*nodeid.NodeID {
	logging.Info("Get active gateways")
	return c.clientManager.GetActiveGateways()
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
