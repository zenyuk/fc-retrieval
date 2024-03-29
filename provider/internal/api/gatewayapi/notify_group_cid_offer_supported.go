package gatewayapi

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
	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrp2pserver"
	"github.com/ConsenSys/fc-retrieval/common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval/common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval/common/pkg/slice"
	"github.com/ConsenSys/fc-retrieval/provider/internal/core"
)

func HandleGatewayNotifyProviderGroupCIDOfferSupportRequest(_ *fcrp2pserver.FCRServerReader, writer *fcrp2pserver.FCRServerWriter, request *fcrmessages.FCRMessage) error {
	// Get core structure
	c := core.GetSingleInstance()

	gatewayID, groupCIDOfferSupported, err := fcrmessages.DecodeGatewayNotifyProviderGroupCIDOfferSupportRequest(request)
	if err != nil {
		// Reply with invalid message
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}

	// Get the gateway's signing key
	gatewayInfo := c.RegisterMgr.GetGateway(gatewayID)
	if gatewayInfo == nil {
		logging.Warn("Gateway information not found for %s.", gatewayID.ToString())
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}
	pubKey, err := gatewayInfo.GetSigningKey()
	if err != nil {
		logging.Warn("Fail to obtain the public key for %s", gatewayID.ToString())
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}

	// Verify the message
	if request.Verify(pubKey) != nil {
		logging.Warn("Fail to verify the request from %s", gatewayID.ToString())
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}

	updateGatewaysSupportingGroupCIDOffer(groupCIDOfferSupported, c.GroupOfferGatewayIDs, gatewayID)

	// Construct response
	response, err := fcrmessages.EncodeGatewayNotifyProviderGroupCIDOfferSupportResponse(true)
	if err != nil {
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}
	// Sign the response
	if response.Sign(c.ProviderPrivateKey, c.ProviderPrivateKeyVersion) != nil {
		logging.Error("Internal error in signing message.")
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}
	// Respond
	err = writer.Write(response, c.Settings.TCPInactivityTimeout)
	if err != nil {
		return err
	}

	return nil
}

// updateGatewaysSupportingGroupCIDOffer changes locally stored collection of Gateway IDs supporting Group CID Offer
func updateGatewaysSupportingGroupCIDOffer(groupCIDOfferSupported bool, savedGateways []nodeid.NodeID, gatewayID *nodeid.NodeID) {
	exists, foundIndex, err := slice.Exists(savedGateways, gatewayID)
	if err != nil {
		logging.Error("Error while searching for existing gatewayID %v", gatewayID)
	}

	if groupCIDOfferSupported {
		if !exists {
			// add new
			savedGateways = append(savedGateways, *gatewayID)
		}
		return
	}

	// groupCIDOfferSupported == false
	// Group CID Offer should NOT be supported for this Provider
	if exists {
		// remove existing
		savedGateways[foundIndex] = savedGateways[len(savedGateways)-1]
		savedGateways = savedGateways[:len(savedGateways)-1]
	}
}
