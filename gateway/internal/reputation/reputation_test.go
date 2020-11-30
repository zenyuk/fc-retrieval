package reputation
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
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
    "math/big"
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestClientRepInitial(t *testing.T) {
	id := big.NewInt(2)
	n := nodeid.NewNodeID(id)
	r := GetSingleInstance()
	r.EstablishClientReputation(n)

	rep := r.GetClientReputation(n)
	assert.Equal(t, clientInitialReputation, rep, "Initial reputation not set correctly")
}

func TestClientRepDeposit(t *testing.T) {
	testClientReputationChange(t, GetSingleInstance().OnChainDeposit, clientOnChainDeposit)
}

func TestClientRepEstablishmentChallenge(t *testing.T) {
	testClientReputationChange(t, GetSingleInstance().ClientEstablishmentChallenge, clientEstablishmentChallenge)
}

func TestClientStdDiscOneCidOffer(t *testing.T) {
	testClientReputationChange(t, GetSingleInstance().ClientStdDiscOneCidOffer, clientStdDiscOneCidOffer)
}

func TestClientStdDiscNoCidOffers(t *testing.T) {
	testClientReputationChange(t, GetSingleInstance().ClientStdDiscNoCidOffers, clientStdDiscNoCidOffers)
}

func TestClientStdDiscLateCidOffers(t *testing.T) {
	testClientReputationChange(t, GetSingleInstance().ClientStdDiscLateCidOffers, clientStdDiscLateCidOffers)
}

func TestClientStdDiscNonPayment(t *testing.T) {
	testClientReputationChange(t, GetSingleInstance().ClientStdDiscNonPayment, clientStdDiscNonPayment)
}

func TestClientDhtDiscOneCidOffer(t *testing.T) {
	testClientReputationChange(t, GetSingleInstance().ClientDhtDiscOneCidOffer, clientDhtDiscOneCidOffer)
}

func TestClientDhtDiscNoCidOffers(t *testing.T) {
	testClientReputationChange(t, GetSingleInstance().ClientDhtDiscNoCidOffers, clientDhtDiscNoCidOffers)
}

func TestClientDhtDiscLateCidOffers(t *testing.T) {
	testClientReputationChange(t, GetSingleInstance().ClientDhtDiscLateCidOffers, clientDhtDiscLateCidOffers)
}

func TestClientDhtDiscNonPayment(t *testing.T) {
	testClientReputationChange(t, GetSingleInstance().ClientDhtDiscNonPayment, clientDhtDiscNonPayment)
}

func TestClientMicroPayment(t *testing.T) {
	testClientReputationChange(t, GetSingleInstance().ClientMicroPayment, clientMicroPayment)
}

func TestClientInvalidMessage(t *testing.T) {
	testClientReputationChange(t, GetSingleInstance().ClientInvalidMessage, clientInvalidMessage)
}

func TestClientRepMax(t *testing.T) {
	id := big.NewInt(2)
	n := nodeid.NewNodeID(id)
	r := GetSingleInstance()
	r.EstablishClientReputation(n)

	// Deposit enough to exceed maximum. This code will need to change if the reputation
	// earned by doing a deposit is reduced.
	for i := 0; i < 11; i++ {
		r.OnChainDeposit(n)
	}
	rep := r.GetClientReputation(n)
	assert.Equal(t, clientMaxReputation, rep, "reputation does not equal max")
}


func testClientReputationChange(t *testing.T, f func(clientNodeID *nodeid.NodeID), expectedChange int) {
	id := big.NewInt(2)
	n := nodeid.NewNodeID(id)
	r := GetSingleInstance()
	r.EstablishClientReputation(n)
	f(n)
	rep := r.GetClientReputation(n)
	assert.Equal(t, clientInitialReputation + expectedChange, rep, "reputation not set correctly")
}
