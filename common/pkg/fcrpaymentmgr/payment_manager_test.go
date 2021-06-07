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

import (
	"fmt"
	"math/big"
)

func ExampleNewFCRPaymentMgr_01() {
	mgr, err := NewFCRPaymentMgr("privateKey", "lotusAPIAddr", "authToken")
	fmt.Printf("%T\n", mgr)
	fmt.Println(err)
	// Output:
	// *fcrpaymentmgr.FCRPaymentMgr
	// encoding/hex: invalid byte: U+0070 'p'
}

func ExampleNewFCRPaymentMgr_02() {
	mgr, err := NewFCRPaymentMgr("", "", "")
	fmt.Printf("%T\n", mgr)
	fmt.Println(err)
	// Output:
	// *fcrpaymentmgr.FCRPaymentMgr
	// Unable to get public key, private key is empty
}

func ExampleNewFCRPaymentMgr_03() {
	mgr, err := NewFCRPaymentMgr("AA", "", "")
	fmt.Printf("%T\n", mgr)
	fmt.Println(err)
	// Output:
	// *fcrpaymentmgr.FCRPaymentMgr
	// <nil>
}

func ExampleTopup_01() {
	mgr, err := NewFCRPaymentMgr("AA", "", "")
	err = mgr.Topup("recipient", big.NewInt(0))
	fmt.Println(err)
	// Output: unknown address network
}

