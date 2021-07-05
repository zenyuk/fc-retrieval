package clientapi

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
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/ant0ine/go-json-rest/rest"

	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval/common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval/common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval/gateway/internal/core"
)

const defaultPaymentLane = 0

// HandleClientDHTCIDDiscoverRequestV2 is used to handle client request for cid offer
func HandleClientDHTCIDDiscoverRequestV2(w rest.ResponseWriter, request *fcrmessages.FCRMessage) {
	// Get core structure
	c := core.GetSingleInstance()

	cid, nonce, ttl, numDHT, _, paymentChannelAddress, voucher, err := fcrmessages.DecodeClientDHTDiscoverRequestV2(request)
	if err != nil {
		s := "Fail to decode message."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusBadRequest)
		return
	}

	// First check if the message can be discarded
	if time.Now().Unix() > ttl {
		// Message expired.
		return
	}
	// Get a list of gatewayIDs to contact
	gateways, err := c.RegisterMgr.GetGatewaysNearCID(cid, int(numDHT), c.GatewayID)
	if err != nil {
		s := "Fail to obtain peers."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusBadRequest)
		return
	}

	amount, err := c.PaymentMgr.Receive(paymentChannelAddress, voucher)
	if err != nil {
		s := "Internal error in payment manager Receive."
		logging.Error(s)
		rest.Error(w, s, http.StatusBadRequest)
		return
	}

	expectedAmount := new(big.Int).Mul(c.Settings.SearchPrice, big.NewInt(numDHT))
	if amount.Cmp(expectedAmount) < 0 {
		s := "Insufficient Funds, received " + amount.String() + ", expected: " + expectedAmount.String()
		logging.Error(s)
		rest.Error(w, s, http.StatusInternalServerError)
	}

	// Construct response
	// TODO: Right now, it ignores the incremental result filed.
	// Will return all in one message.
	// Now requesting gateways.
	contacted := make([]nodeid.NodeID, 0)
	contactedResp := make([]fcrmessages.FCRMessage, 0)
	unContactable := make([]nodeid.NodeID, 0)
	for _, gw := range gateways {
		id, err := nodeid.NewNodeIDFromHexString(gw.GetNodeID())
		if err != nil {
			s := "Fail to generate node id."
			logging.Error(s + err.Error())
			rest.Error(w, s, http.StatusBadRequest)
			return
		}
		// Pay this gateway
		var needTopup bool
		paymentChannelAddress, voucher, needTopup, err = c.PaymentMgr.Pay(gw.GetAddress(), 0, c.Settings.SearchPrice)
		if err != nil {
			s := "Fail to pay recipient."
			logging.Error(s + err.Error())
			rest.Error(w, s, http.StatusBadRequest)
			return
		}
		if needTopup {
			paymentChannelAddress, voucher, err = topupAndPay(c, gw.GetAddress(), defaultPaymentLane, c.Settings.SearchPrice)
			if err != nil {
				s := "Internal error: Failure to topup the gateway: " + gw.GetAddress()
				logging.Error(s + err.Error())
				rest.Error(w, s, http.StatusInternalServerError)
				return
			}
		}
		res, err := c.P2PServer.RequestGatewayFromGateway(id, fcrmessages.GatewayDHTDiscoverRequestV2Type, cid, id, paymentChannelAddress, voucher)
		if err != nil {
			unContactable = append(unContactable, *id)
		} else {
			contacted = append(contacted, *id)
			contactedResp = append(contactedResp, *res)
		}
	}

	response, err := fcrmessages.EncodeClientDHTDiscoverResponseV2(contacted, contactedResp, unContactable, nonce, false, "")
	if err != nil {
		s := "Internal error: Fail to encode message."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusInternalServerError)
		return
	}

	// Sign message
	err = response.Sign(c.GatewayPrivateKey, c.GatewayPrivateKeyVersion)
	if err != nil {
		s := "Internal error: Fail to sign message."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusInternalServerError)
		return
	}
	if err := w.WriteJson(response); err != nil {
		logging.Error("can't write JSON during HandleClientDHTCIDDiscoverRequestV2 %s", err.Error())
	}
}

func topupAndPay(c *core.Core, recipient string, paymentLane uint64, amount *big.Int) (paymentChannel string, voucher string, err error) {
	if err := c.PaymentMgr.Topup(recipient, c.Settings.TopupAmount); err != nil {
		return "", "", fmt.Errorf("Unable to topup, error: %s ", err.Error())
	}
	var needTopup bool
	paymentChannel, voucher, needTopup, err = c.PaymentMgr.Pay(recipient, paymentLane, amount)
	if err != nil {
		return paymentChannel, voucher, fmt.Errorf("Unable to make payment, did topup first; error: %s ", err.Error())
	}
	if needTopup {
		err = c.PaymentMgr.Topup(recipient, c.Settings.TopupAmount)
		if err != nil {
			return paymentChannel, voucher, fmt.Errorf("Unable to topup second time, error: %s ", err.Error())
		}
		paymentChannel, voucher, needTopup, err = c.PaymentMgr.Pay(recipient, paymentLane, amount)
		if err != nil {
			return paymentChannel, voucher, fmt.Errorf("Unable to make payment, did topup 2 times; error: %s ", err.Error())
		}
		if needTopup {
			return paymentChannel, voucher, fmt.Errorf("unable to make payment, it's keep saying 'need topup', tried to topup 2 times")
		}
	}
	return paymentChannel, voucher, nil
}
