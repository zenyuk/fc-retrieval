package gatewayapi

import (
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrp2pserver"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/core"
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

// HandleGatewayDHTOfferRequest handles the gateway dht discover request
func HandleGatewayDHTOfferRequest(reader *fcrp2pserver.FCRServerReader, writer *fcrp2pserver.FCRServerWriter, request *fcrmessages.FCRMessage) error {
	// Get the core structure
	c := core.GetSingleInstance()

	pieceCID, nonce, _, suboffers, fundedPaymentChannel, err := fcrmessages.DecodeGatewayDHTDiscoverOfferResponse(request)
	if err != nil {
		// Reply with invalid message
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}

	// Respond to the request
	_, exists := c.OffersMgr.GetOffers(pieceCID)
	if exists {
		for i, offer := range suboffers {

			// what means fundedPaymentChannel true or false
			if fundedPaymentChannel[i] { // this or not this
				continue
			}

			// should validate ??
			offer.HasExpired()
			offer.GetSignature()

			// Now Verify the merkle proof
			if offer.VerifyMerkleProof() != nil {
				logging.Error("Merkle proof verification failed.")
				continue
			}

			// how can I get chn and voucher?
			var chanAddr string
			var voucher string

			received, err := c.PaymentMgr.Receive(chanAddr, voucher)
			if err != nil {
				// return writer.Write ????
				continue
			}
			// is this condition right?
			fundedPaymentChannel[i] = received.Uint64() >= offer.GetPrice()
		}
	}

	// Construct response
	response, err := fcrmessages.EncodeGatewayDHTDiscoverOfferResponse(pieceCID, nonce, exists, suboffers, fundedPaymentChannel)
	if err != nil {
		// TODO: Do we need a response of internal error?
		// There are three possible errors, 1. Protocol errors (request is not correct) 2. Communication errors (lost connection) and 3. Internal errors.
		// Need to do error management.
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}

	// Sign the response
	if response.Sign(c.GatewayPrivateKey, c.GatewayPrivateKeyVersion) != nil {
		logging.Error("Internal error in signing message.")
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}

	return writer.Write(response, c.Settings.TCPInactivityTimeout)
}
