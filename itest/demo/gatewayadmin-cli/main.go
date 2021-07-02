package main

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"
	"time"

	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrpaymentmgr"
	"github.com/ConsenSys/fc-retrieval/common/pkg/register"
	"github.com/ConsenSys/fc-retrieval/gateway-admin/pkg/fcrgatewayadmin"
	"github.com/c-bata/go-prompt"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/api/apistruct"
	"github.com/filecoin-project/lotus/chain/types"
	cid2 "github.com/ipfs/go-cid"
)

var lotusAP = "http://127.0.0.1:1234/rpc/v0"
var gwAdmin *fcrgatewayadmin.FilecoinRetrievalGatewayAdmin
var initialised bool

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "init", Description: "Initialise the gateway admin"},
		{Text: "register-gw", Description: "Register and initialise a given gateway"},
		{Text: "inspect-gw-offer-traffic", Description: "Inspect the offer query traffic for a given gateway"},
		{Text: "cache-offer-content", Description: "Ask a given gateway to cache a popular offer content and start serving"},
		{Text: "exit", Description: "Exit the program"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func executor(in string) {
	in = strings.TrimSpace(in)

	blocks := strings.Split(in, " ")
	switch blocks[0] {
	case "init":
		if initialised {
			fmt.Println("Admin has already been initialised.")
			return
		}
		fmt.Println("Initialise gateway admin (dev)...")
		blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
		if err != nil {
			fmt.Println("Error in generating client key.")
			return
		}
		confBuilder := fcrgatewayadmin.CreateSettings()
		confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
		confBuilder.SetRegisterURL("http://127.0.0.1:9020")
		conf := confBuilder.Build()
		gwAdmin = fcrgatewayadmin.NewFilecoinRetrievalGatewayAdmin(*conf)
		initialised = true
		fmt.Println("Done.")
	case "register-gw":
		if !initialised {
			fmt.Println("Admin hasn't been initialised yet.")
			return
		}
		fmt.Println("Initialise gateway (dev)")
		token, acct := getLotusToken()
		keys, addresses, err := generateAccount(lotusAP, token, acct, 20)
		if err != nil {
			fmt.Printf("Fail to initialise gateway: %s\n", err.Error())
			return
		}
		for i := 0; i < 20; i++ {

			key := keys[i]
			address := addresses[i]
			fmt.Printf("(Dev) Created an account for this gateway to use: %s\n", address)

			gatewayRootKey, err := fcrcrypto.GenerateBlockchainKeyPair()
			if err != nil {
				fmt.Printf("Fail to initialise gateway: %s\n", err.Error())
				return
			}
			gatewayRootSigningKey, err := gatewayRootKey.EncodePublicKey()
			if err != nil {
				fmt.Printf("Fail to initialise gateway: %s\n", err.Error())
				return
			}
			gatewayRetrievalPrivateKey, err := fcrcrypto.GenerateRetrievalV1KeyPair()
			if err != nil {
				fmt.Printf("Fail to initialise gateway: %s\n", err.Error())
				return
			}
			gatewayRetrievalSigningKey, err := gatewayRetrievalPrivateKey.EncodePublicKey()
			if err != nil {
				fmt.Printf("Fail to initialise gateway: %s\n", err.Error())
				return
			}

			var idStr string
			if i%2 == 0 {
				idStr = fmt.Sprintf("%X000000000000000000000000000000000000000000000000000000000000000", i/2)
			} else {
				idStr = fmt.Sprintf("%X800000000000000000000000000000000000000000000000000000000000000", i/2)
			}

			gatewayRegister := &register.GatewayRegister{
				NodeID:              idStr,
				Address:             address,
				RootSigningKey:      gatewayRootSigningKey,
				SigningKey:          gatewayRetrievalSigningKey,
				RegionCode:          "au",
				NetworkInfoGateway:  fmt.Sprintf("gateway%v:9012", i),
				NetworkInfoProvider: fmt.Sprintf("gateway%v:9011", i),
				NetworkInfoClient:   fmt.Sprintf("127.0.0.1:%v", 8018+i),
				NetworkInfoAdmin:    fmt.Sprintf("127.0.0.1:%v", 7013+i),
			}
			err = gwAdmin.InitialiseGatewayV2(gatewayRegister, gatewayRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1), key, "http://lotus:1234/rpc/v0", token)
			if err != nil {
				fmt.Printf("Fail to initialise gateway: %s\n", err.Error())
				return
			}
			fmt.Printf("Successfully initialised gateway: %v\n", gatewayRegister.NodeID)
		}
	case "inspect-gw-offer-traffic":
		if !initialised {
			fmt.Println("Admin hasn't been initialised yet.")
			return
		}
		fmt.Println("Hasn't been implemented yet.")
	case "cache-offer-content":
		if !initialised {
			fmt.Println("Admin hasn't been initialised yet.")
			return
		}
		fmt.Println("Hasn't been implemented yet.")
	case "exit":
		fmt.Println("Bye!")
		os.Exit(0)
	default:
		fmt.Printf("Unknown command: %s\n", blocks[0])
	}
}

func main() {
	initialised = false
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			debug.PrintStack()
		}
		handleExit()
	}()
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(">>> "),
	)
	p.Run()
}

func getLotusToken() (string, string) {
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

// The following helper method is used to generate a new filecoin account with 10 filecoins of balance
func generateAccount(lotusAP string, token string, superAcct string, num int) ([]string, []string, error) {
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
	cids := make([]cid2.Cid, 0)

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
		signedMsg, err := fillMsg(mainAddress, &api, msg)
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

// fillMsg will fill the gas and sign a given message
func fillMsg(fromAddress address.Address, api *apistruct.FullNodeStruct, msg *types.Message) (*types.SignedMessage, error) {
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

// wait receipt will wait until receipt is received for a given cid
func waitReceipt(cid *cid2.Cid, api *apistruct.FullNodeStruct) *types.MessageReceipt {
	// Return until recipient is returned (transaction is processed)
	var receipt *types.MessageReceipt
	var err error
	for {
		receipt, err = api.StateGetReceipt(context.Background(), *cid, types.EmptyTSK)
		if err != nil {
			fmt.Printf("Payment manager has error getting recipient of cid: %s\n", cid.String())
		}
		if receipt != nil {
			break
		}
		// TODO, Make the interval configurable
		time.Sleep(1 * time.Second)
	}
	return receipt
}

func generateKeyPair() ([]byte, []byte, error) {
	// Generate Private-Public pairs. Public key will be used as address
	var signer fcrpaymentmgr.SecpSigner
	privateKey, err := signer.GenPrivate()
	if err != nil {
		fmt.Printf("Error generating private key, while creating address %s\n", err.Error())
		return nil, nil, err
	}

	publicKey, err := signer.ToPublic(privateKey)
	if err != nil {
		fmt.Printf("Error generating public key, while creating address %s\n", err.Error())
		return nil, nil, err
	}
	return privateKey, publicKey, err
}

// handleExit fixes the problem of broken terminal when exit in Linux
// ref: https://www.gitmemory.com/issue/c-bata/go-prompt/228/820639887
func handleExit() {
	if _, err := os.Stat("/bin/stty"); os.IsNotExist(err) {
		return
	}
	rawModeOff := exec.Command("/bin/stty", "-raw", "echo")
	rawModeOff.Stdin = os.Stdin
	_ = rawModeOff.Run()
	rawModeOff.Wait()
}
