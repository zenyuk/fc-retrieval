package fcrpaymentmgr

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

import "errors"

// To implement payment manager.
// payment manager will need to have to
// 1. Initialise a map[recipient string](payment channel string) <- this can be a db
// 2. Topup() will call lotus api to top up a payment channel, searching by the recipient address,
//		Searching is done through lotus API, if there is a payment channel existing but not in
//		local memory, the mapping will need to be stored in the map.
//		If there isn't any payment channel existed, call lotus API to create a payment channel,
//		and the mapping again will need to be stored into the map.
// 3. Pay() will first get the payment channel, searching by the recipient address, first search
//		the map, if no such payment channel exists, then search through lotus API (update the map
//		if found). On finding a channel, use again lotus API to create voucher. If not finding any
//		channel, returns error.
// 4. Receive() will just be a wrapper of the lotus paych receive. It returns the amount received.
// Note to handle access from multiple threads.

// Payment manager manages all payment related functions
type FCRPaymentMgr struct {
}

// NewFCRPaymentMgr creates a new payment manager.
func NewFCRPaymentMgr(lotusAccount string, lotusAPI string) (*FCRPaymentMgr, error) {
	return nil, errors.New("Not yet implemented")
}

// Topup will topup a payment channel to recipient with given amount.
func (mgr *FCRPaymentMgr) Topup(recipient string, amount int) error {
	return errors.New("Not yet implemented")
}

// Pay will pay the recipient a given amount and return the voucher.
func (mgr *FCRPaymentMgr) Pay(recipient string, amount int) (string, error) {
	return "", errors.New("Not yet implemented")
}

// Receive will receive a given voucher at a given payment channel and return the amount received.
func (mgr *FCRPaymentMgr) Receive(paychAddr string, voucher string) (int, error) {
	return 0, errors.New("Not yet implemented")
}
