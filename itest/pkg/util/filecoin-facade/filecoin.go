package filecoin_facade

import (
	"context"
	"encoding/hex"
	"errors"
	"net/http"
	"os/exec"
	"time"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/api/apistruct"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/ipfs/go-cid"

	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrpaymentmgr"
	"github.com/ConsenSys/fc-retrieval/common/pkg/logging"
)

// GenerateAccount - helper method, used to generate a new filecoin account with 10 filecoins of balance
func GenerateAccount(ctx context.Context, lotusAP string, token string, superAcct string, num int) ([]string, []string, error) {
	// Get API
	var api apistruct.FullNodeStruct
	headers := http.Header{"Authorization": []string{"Bearer " + token}}
	closer, err := jsonrpc.NewMergeClient(ctx, lotusAP, "Filecoin", []interface{}{&api.Internal, &api.CommonStruct.Internal}, headers)
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
		amt, err := types.ParseFIL("100")
		if err != nil {
			return nil, nil, err
		}

		msg := &types.Message{
			To:     address1,
			From:   mainAddress,
			Value:  types.BigInt(amt),
			Method: 0,
		}
		signedMsg, err := fillMsg(ctx, mainAddress, &api, msg)
		if err != nil {
			return nil, nil, err
		}

		// Send request to lotus
		contentID, err := api.MpoolPush(ctx, signedMsg)
		if err != nil {
			return nil, nil, err
		}
		cids = append(cids, contentID)

		// Add to result
		privateKeys = append(privateKeys, privKeyStr)
		addresses = append(addresses, address1.String())
	}

	// Finally check receipts
	for _, contentID := range cids {
		receipt := waitReceipt(&contentID, &api)
		if receipt.ExitCode != 0 {
			return nil, nil, errors.New("transaction fail to execute")
		}
	}

	return privateKeys, addresses, nil
}

// GetLotusToken gets the lotus token and the super account from the lotus container
func GetLotusToken() (string, string) {
	cmd := exec.Command("docker", "ps", "--filter", "ancestor=consensys/lotus-full-node:latest", "--format", "{{.ID}}")
	stdout, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	id := string(stdout[:len(stdout)-1])
	cmd = exec.Command("docker", "exec", id, "bash", "-c", "cd ~/.lotus; cat token")
	stdout, err = cmd.Output()
	if err != nil {
		panic(err)
	}
	token := string(stdout)

	cmd = exec.Command("docker", "exec", id, "bash", "-c", "./lotus wallet default")
	stdout, err = cmd.Output()
	if err != nil {
		panic(err)
	}
	acct := string(stdout[:len(stdout)-1])
	return token, acct
}

func generateKeyPair() ([]byte, []byte, error) {
	// Generate Private-Public pairs. Public key will be used as address
	var signer fcrpaymentmgr.SecpSigner
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
func fillMsg(ctx context.Context, fromAddress address.Address, api *apistruct.FullNodeStruct, msg *types.Message) (*types.SignedMessage, error) {
	// Get nonce
	nonce, err := api.MpoolGetNonce(ctx, msg.From)
	if err != nil {
		return nil, err
	}
	msg.Nonce = nonce

	// Calculate gas
	limit, err := api.GasEstimateGasLimit(ctx, msg, types.EmptyTSK)
	if err != nil {
		return nil, err
	}
	msg.GasLimit = int64(float64(limit) * 1.25)

	premium, err := api.GasEstimateGasPremium(ctx, 10, msg.From, msg.GasLimit, types.EmptyTSK)
	if err != nil {
		return nil, err
	}
	msg.GasPremium = premium

	feeCap, err := api.GasEstimateFeeCap(ctx, msg, 20, types.EmptyTSK)
	if err != nil {
		return nil, err
	}
	msg.GasFeeCap = feeCap

	// Sign message
	return api.WalletSignMessage(ctx, fromAddress, msg)
}

// wait receipt will wait until receipt is received for a given cid
func waitReceipt(cid *cid.Cid, api *apistruct.FullNodeStruct) *types.MessageReceipt {
	// Return until recipient is returned (transaction is processed)
	var receipt *types.MessageReceipt
	var err error
	for {
		receipt, err = api.StateGetReceipt(context.Background(), *cid, types.EmptyTSK)
		if err != nil {
			logging.Warn("Payment manager has error getting recipient of cid: %s", cid.String())
		}
		if receipt != nil {
			break
		}
		// TODO, Make the interval configurable
		time.Sleep(1 * time.Second)
	}
	return receipt
}
