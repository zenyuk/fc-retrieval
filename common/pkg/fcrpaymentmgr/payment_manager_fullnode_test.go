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
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/api/apistruct"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/assert"
)

func TestGenerateLotusAccount(t *testing.T) {
	// lotus token
	lotusAP := fileToString("LOTUS_AP.out")
	lotusToken := fileToString("LOTUS_TOKEN.out")
	superAcct := fileToString("LOTUS_SUPER_ACCT.out")
	if lotusAP == "" || lotusToken == "" || superAcct == "" {
		t.Skip("lotusAP or lotusToken or superAcct is empty string")
	}
	logging.Info("Lotus token is: %s", lotusToken)
	logging.Info("Super Acct is %s", superAcct)

	var err error
	privateKeys, accountAddrs, err := GenerateAccount(lotusAP, lotusToken, superAcct, 2)
	if err != nil {
		t.Fatal(err)
	}
	keysFile, _ := os.Create("LOTUS_PRIVATE_KEYS.out")
	addrsFile, _ := os.Create("LOTUS_ACCOUT_ADDRS.out")
	for i, privateKey := range privateKeys {
		fmt.Fprintln(keysFile, privateKey)
		fmt.Fprintln(addrsFile, accountAddrs[i])
	}
	defer keysFile.Close()
	defer addrsFile.Close()

	// func NewFCRPaymentMgr(privateKey, lotusAPIAddr, authToken string) (*FCRPaymentMgr, error) {
	payMgr, err := NewFCRPaymentMgr(privateKeys[0], lotusAP, lotusToken)
	logging.Info("%#v", err)

	payMgr.Dump()
}

func TestLoutusPaymentManager(t *testing.T) {
	lotusAP := fileToString("LOTUS_AP.out")
	lotusToken := fileToString("LOTUS_TOKEN.out")
	superAcct := fileToString("LOTUS_SUPER_ACCT.out")
	if lotusAP == "" || lotusToken == "" || superAcct == "" {
		t.Skip("lotusToken or superAcct is empty string")
	}

	privateKeys := strings.Fields(fileToString("LOTUS_PRIVATE_KEYS.out"))
	accountAddrs := strings.Fields(fileToString("LOTUS_ACCOUT_ADDRS.out"))
	if len(privateKeys) == 0 || len(accountAddrs) == 0 {
		t.Skip("lotus privateKeys or accountAddrs is empty")
	}

	// func NewFCRPaymentMgr(privateKey, lotusAPIAddr, authToken string) (*FCRPaymentMgr, error) {
	payMgr, err := NewFCRPaymentMgr(privateKeys[0], lotusAP, lotusToken)
	logging.Info("NewFCRPaymentMgr err:%#v", err)

	bigInt01 := new(big.Int)
	bigInt05 := new(big.Int)
	bigInt06 := new(big.Int)
	fmt.Sscan("100_000_000_000_000_000", bigInt01) // 0.1
	fmt.Sscan("500_000_000_000_000_000", bigInt05) // 0.5
	fmt.Sscan("600_000_000_000_000_000", bigInt06) // 0.6

	err = payMgr.topupAndPay(accountAddrs[0], bigInt05, bigInt06) // topup .5 pay .6
	// "Receive=<nil>,topup=500000000000000000,pay=600000000000000000,needTopup=true,err=EOF"
	assert.NotNil(t, err)

	err = payMgr.topupAndPay(accountAddrs[0], big.NewInt(0), bigInt05) // topup 0, pay .5
	// "Receive=500000000000000000,topup=0,pay=500000000000000000,needTopup=false,err=<nil>"
	assert.Equal(t, err, nil)

	err = payMgr.topupAndPay(accountAddrs[0], bigInt05, bigInt05) // topup .5, pay .5
	// "Receive=<nil>,topup=500000000000000000,pay=500000000000000000,needTopup=true,err=EOF"
	assert.NotNil(t, err)

	err = payMgr.topupAndPay(accountAddrs[0], bigInt01, bigInt01) // topup .1, pay .1
	// "Receive=100000000000000000,topup=100000000000000000,pay=100000000000000000,needTopup=false,err=<nil>"
	assert.Equal(t, err, nil)
}

func fileToString(fileName string) string {
	dat, _ := ioutil.ReadFile(fileName)
	return strings.TrimSpace(string(dat))
}

func (payMgr *FCRPaymentMgr) topupAndPay(recipient string, topupAmt *big.Int, payAmt *big.Int) error {
	// func (mgr *FCRPaymentMgr) Topup(recipient string, amount *big.Int) error {
	err := payMgr.Topup(recipient, topupAmt)
	logging.Info("Topup err: %#v", err)

	// func (mgr *FCRPaymentMgr) Pay(recipient string, lane uint64, amount *big.Int) (string, string, bool, error) {
	// Return channel address, voucher, true if needs to top up, and error.
	channel, voucher, needTopUp, err := payMgr.Pay(recipient, 0, payAmt)
	logging.Info("Pay channel: %s", channel)
	logging.Info("voucher: %s", voucher)
	logging.Info("err: %#v", err)

	// func (mgr *FCRPaymentMgr) Receive(channel string, voucher string) (*big.Int, error) {
	receiveAmt, err := payMgr.Receive(channel, voucher)
	logging.Info("Receive=%s,topup=%s,pay=%s,needTopup=%v,err=%v",
		receiveAmt.String(), topupAmt.String(), payAmt.String(), needTopUp,
		err)
	return err
}

// The following helper method is used to generate a new filecoin account with 1 filecoins of balance
func GenerateAccount(lotusAP string, token string, superAcct string, num int) ([]string, []string, error) {
	// Get API
	var api apistruct.FullNodeStruct
	headers := http.Header{"Authorization": []string{"Bearer " + token}}
	closer, err := jsonrpc.NewMergeClient(context.Background(), lotusAP, "Filecoin", []interface{}{&api.Internal, &api.CommonStruct.Internal}, headers)
	if err != nil {
		return nil, nil, err
	}
	defer closer()

	mainAddress, err := address.NewFromString(superAcct)
	if err != nil {
		return nil, nil, err
	}

	privateKeys := make([]string, 0)
	addresses := make([]string, 0)
	cids := make([]cid.Cid, 0)

	// Send messages
	for i := 0; i < num; i++ {
		privKey, pubKey, err := generateKeyPair()
		if err != nil {
			return nil, nil, err
		}
		privKeyStr := hex.EncodeToString(privKey)

		address1, err := address.NewSecp256k1Address(pubKey)
		if err != nil {
			return nil, nil, err
		}

		// Get amount
		amt, err := types.ParseFIL("1")
		if err != nil {
			return nil, nil, err
		}

		msg := &types.Message{
			To:     address1,
			From:   mainAddress,
			Value:  types.BigInt(amt),
			Method: 0,
		}
		signedMsg, err := fillMsgForAddress(mainAddress, &api, msg)
		if err != nil {
			return nil, nil, err
		}

		// Send request to lotus
		cid, err := api.MpoolPush(context.Background(), signedMsg)
		if err != nil {
			return nil, nil, err
		}
		cids = append(cids, cid)

		// Add to result
		privateKeys = append(privateKeys, privKeyStr)
		addresses = append(addresses, address1.String())
	}

	// Finally check receipts
	for _, cid := range cids {
		receipt := waitReceipt(&cid, &api)
		if receipt.ExitCode != 0 {
			return nil, nil, errors.New("Transaction fail to execute")
		}
	}

	return privateKeys, addresses, nil
}

func generateKeyPair() ([]byte, []byte, error) {
	// Generate Private-Public pairs. Public key will be used as address
	var signer SecpSigner
	privateKey, err := signer.GenPrivate()
	if err != nil {
		logging.Error("Error generating private key, while creating address %s", err.Error())
		return nil, nil, err
	}

	publicKey, err := signer.ToPublic(privateKey)
	if err != nil {
		logging.Error("Error generating public key, while creating address %s", err.Error())
		return nil, nil, err
	}
	return privateKey, publicKey, err
}

// fillMsg will fill the gas and sign a given message
func fillMsgForAddress(fromAddress address.Address, api *apistruct.FullNodeStruct, msg *types.Message) (*types.SignedMessage, error) {
	// Get nonce
	nonce, err := api.MpoolGetNonce(context.Background(), msg.From)
	if err != nil {
		return nil, err
	}
	msg.Nonce = nonce

	// Calculate gas
	limit, err := api.GasEstimateGasLimit(context.Background(), msg, types.EmptyTSK)
	if err != nil {
		return nil, err
	}
	msg.GasLimit = int64(float64(limit) * 1.25)

	premium, err := api.GasEstimateGasPremium(context.Background(), 10, msg.From, msg.GasLimit, types.EmptyTSK)
	if err != nil {
		return nil, err
	}
	msg.GasPremium = premium

	feeCap, err := api.GasEstimateFeeCap(context.Background(), msg, 20, types.EmptyTSK)
	if err != nil {
		return nil, err
	}
	msg.GasFeeCap = feeCap

	// Sign message
	return api.WalletSignMessage(context.Background(), fromAddress, msg)
}

// Dump is used for debugging only, it returns the string repr of payment manager.
// This is not threadsafe and should only be called for debugging.
func (mgr *FCRPaymentMgr) Dump() string {
	var sb strings.Builder
	sb.WriteString("Payment Manager status:\n")
	sb.WriteString(fmt.Sprintf("Private key: %v\n", hex.EncodeToString(mgr.privKey)))
	sb.WriteString(fmt.Sprintf("Adress: %v\n", mgr.address.String()))
	sb.WriteString(fmt.Sprintf("Auth token: %v\n", mgr.authToken))
	sb.WriteString("Outbound channels:\n")
	for recipient, cs := range mgr.outboundChs {
		sb.WriteString(fmt.Sprintf("\t\tChannel address: %v (Recipient address %v)\n", cs.addr.String(), recipient))
		sb.WriteString(fmt.Sprintf("\t\tChannel balance: %v\n", cs.balance.String()))
		sb.WriteString(fmt.Sprintf("\t\tChannel redeemed: %v\n", cs.redeemed.String()))
		sb.WriteString("\t\tLane states:\n")
		for lane, ls := range cs.laneStates {
			sb.WriteString(fmt.Sprintf("\t\t\t\tLane: %v, Nonce: %v, Redeemed: %v\n", lane, ls.nonce, ls.redeemed))
			sb.WriteString(fmt.Sprintf("\t\t\t\tVouchers: "))
			for _, voucher := range ls.vouchers {
				sb.WriteString(voucher)
				sb.WriteString(" ")
			}
			sb.WriteString("\n")
		}
	}
	sb.WriteString("Inbound channels:\n")
	for _, cs := range mgr.inboundChs {
		sb.WriteString(fmt.Sprintf("\t\tChannel address: %v\n", cs.addr.String()))
		sb.WriteString(fmt.Sprintf("\t\tChannel balance: %v\n", cs.balance.String()))
		sb.WriteString(fmt.Sprintf("\t\tChannel redeemed: %v\n", cs.redeemed.String()))
		sb.WriteString("\t\tLane states:\n")
		for lane, ls := range cs.laneStates {
			sb.WriteString(fmt.Sprintf("\t\t\t\tLane: %v, Nonce: %v, Redeemed: %v\n", lane, ls.nonce, ls.redeemed))
			sb.WriteString(fmt.Sprintf("\t\t\t\tVouchers: "))
			for _, voucher := range ls.vouchers {
				sb.WriteString(voucher)
				sb.WriteString(" ")
			}
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
