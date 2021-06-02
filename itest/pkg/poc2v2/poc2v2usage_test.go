package poc2v2

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-itest/pkg/util"
)

//83003f298382e3211b441438b610da078166b162972f0e2f024c6663671c3fb6
//04f76ccfe30ba7c8c586bd619ed3bdb181c56d559a04088e7920ac1efad2a1affcf2236caa378700af5a33611975cc84e0c4ded41df64f65156a524aec33a54ed1

func TestMain(m *testing.M) {
	// Need to make sure this env is not set in host machine
	itestEnv := os.Getenv("ITEST_CALLING_FROM_CONTAINER")

	if itestEnv != "" {
		// Env is set, we are calling from docker container
		m.Run()
		return
	}
	// Env is not set, we are calling from host
	// We need a lotus
	tag := util.GetCurrentBranch()

	// Create shared net
	bgCtx := context.Background()
	ctx, _ := context.WithTimeout(bgCtx, time.Minute*2)
	network, networkName := util.CreateNetwork(ctx)
	defer (*network).Remove(ctx)

	// Start lotus
	lotusContainer := util.StartLotusFullNode(ctx, networkName, false)
	defer lotusContainer.Terminate(ctx)

	// Start itest
	done := make(chan bool)
	itestContainer := util.StartItest(ctx, tag, networkName, util.ColorGreen, done, true)
	defer itestContainer.Terminate(ctx)
	// Block until done.
	if <-done {
		logging.Info("Tests passed, shutdown...")
	} else {
		logging.Fatal("Tests failed, shutdown...")
	}
}

func TestNewAccounts(t *testing.T) {
	ap := "http://lotus-full-node:1234/rpc/v0"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.nuS2Ueb0llOeX3q4K53AZXZNJkh8thYk3MXFGg4VSXA"

	t.Log("Start" + time.Now().String())
	_, accounts, err := util.GenerateAccount(ap, token, 40)
	if err != nil {
		t.Fatal(err)
	}

	for _, acc := range accounts {
		t.Log(acc)
	}

	t.Log("End" + time.Now().String())

	time.Sleep(120 * time.Second)

	// account1 := accounts[0]
	// account2 := accounts[1]
	// t.Log(account1)
	// t.Log(account2)

	// mgr1, err := fcrpaymentmgr.NewFCRPaymentMgr(privateKeys[0], ap, token)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// mgr2, err := fcrpaymentmgr.NewFCRPaymentMgr(privateKeys[1], ap, token)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// _, _, topup, err := mgr1.Pay(account2, 0, "0.001")
	// assert.Nil(t, err)
	// assert.True(t, topup)

	// err = mgr1.Topup(account2, "0.1")
	// assert.Nil(t, err)

	// chanAddr, voucher, topup, err := mgr1.Pay(account2, 0, "0.001")
	// assert.Nil(t, err)
	// assert.False(t, topup)

	// received, err := mgr2.Receive(chanAddr, voucher)
	// assert.Nil(t, err)

	// t.Log(received)
}
