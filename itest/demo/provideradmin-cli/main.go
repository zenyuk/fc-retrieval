package main

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/ConsenSys/fc-retrieval/common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrpaymentmgr"
	"github.com/ConsenSys/fc-retrieval/common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval/common/pkg/register"
	"github.com/ConsenSys/fc-retrieval/provider-admin/pkg/fcrprovideradmin"
	"github.com/c-bata/go-prompt"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/api/apistruct"
	"github.com/filecoin-project/lotus/chain/types"
	cid2 "github.com/ipfs/go-cid"
)

var lotusAP = "http://127.0.0.1:1234/rpc/v0"
var pAdmin *fcrprovideradmin.FilecoinRetrievalProviderAdmin
var initialised bool

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "init", Description: "Initialise the provider admin"},
		{Text: "register-pvd", Description: "Register and initialise a given provider"},
		{Text: "list-pvds", Description: "list currently managed providers"},
		{Text: "publish-dht-offer", Description: "Ask a given provider to publish a dht offer uring 1 given cid"},
		{Text: "publish-offer-given-cids", Description: "Ask a given provider to publish an offer using 3 given cids"},
		{Text: "publish-offer-from-file", Description: "Ask a given provider to publish an offer for a given file"},
		{Text: "publish-random-offer", Description: "(Dev) Ask a given provider to publish random offer"},
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
		fmt.Println("Initialise provider admin (dev)...")
		blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
		if err != nil {
			fmt.Println("Error in generating client key.")
			return
		}
		confBuilder := fcrprovideradmin.CreateSettings()
		confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
		confBuilder.SetRegisterURL("http://127.0.0.1:9020")
		conf := confBuilder.Build()
		pAdmin = fcrprovideradmin.NewFilecoinRetrievalProviderAdmin(*conf)
		initialised = true
		fmt.Println("Done.")
	case "register-pvd":
		if !initialised {
			fmt.Println("Admin hasn't been initialised yet.")
			return
		}
		fmt.Println("Initialise provider (dev)")
		token, acct := getLotusToken()
		keys, addresses, err := generateAccount(lotusAP, token, acct, 1)
		if err != nil {
			fmt.Printf("Fail to initialise provider: %s\n", err.Error())
			return
		}
		key := keys[0]
		address := addresses[0]
		fmt.Printf("(Dev) Created an account for this provider to use: %s\n", address)

		providerRootKey, err := fcrcrypto.GenerateBlockchainKeyPair()
		if err != nil {
			fmt.Printf("Fail to initialise provider: %s\n", err.Error())
			return
		}
		providerRootSigningKey, err := providerRootKey.EncodePublicKey()
		if err != nil {
			fmt.Printf("Fail to initialise provider: %s\n", err.Error())
			return
		}
		providerRetrievalPrivateKey, err := fcrcrypto.GenerateRetrievalV1KeyPair()
		if err != nil {
			fmt.Printf("Fail to initialise provider: %s\n", err.Error())
			return
		}
		providerRetrievalSigningKey, err := providerRetrievalPrivateKey.EncodePublicKey()
		if err != nil {
			fmt.Printf("Fail to initialise provider: %s\n", err.Error())
			return
		}
		providerRegister := &register.ProviderRegister{
			NodeID:             nodeid.NewRandomNodeID().ToString(),
			Address:            address,
			RootSigningKey:     providerRootSigningKey,
			SigningKey:         providerRetrievalSigningKey,
			RegionCode:         "au",
			NetworkInfoGateway: "provider:9032",
			NetworkInfoClient:  "127.0.0.1:9030",
			NetworkInfoAdmin:   "127.0.0.1:9033",
		}
		err = pAdmin.InitialiseProviderV2(providerRegister, providerRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1), key, "http://lotus:1234/rpc/v0", token)
		if err != nil {
			fmt.Printf("Fail to initialise provider: %s\n", err.Error())
			return
		}
		fmt.Printf("Successfully initialised provider: %v\n", providerRegister.NodeID)
	case "list-pvds":
		if !initialised {
			fmt.Println("Admin hasn't been initialised yet.")
			return
		}
		fmt.Println("Providers currently in use:")
		pAdmin.ActiveProvidersLock.RLock()
		for _, pvd := range pAdmin.ActiveProviders {
			fmt.Println(pvd.GetNodeID())
		}
		pAdmin.ActiveProvidersLock.RUnlock()
	case "publish-offer-from-file":
		if !initialised {
			fmt.Println("Admin hasn't been initialised yet.")
			return
		}
		fmt.Println("Hasn't been implemented yet.")
	case "publish-random-offer":
		if !initialised {
			fmt.Println("Admin hasn't been initialised yet.")
			return
		}
		if len(blocks) != 2 {
			fmt.Println("Usage: publish-random-offer $nodeID")
			return
		}
		id, err := nodeid.NewNodeIDFromHexString(blocks[1])
		if err != nil {
			fmt.Printf("Invalid nodeID: %s\n", blocks[1])
			return
		}
		cid1 := cid.NewRandomContentID()
		cid2 := cid.NewRandomContentID()
		cid3 := cid.NewRandomContentID()
		expiryDate := time.Now().Local().Add(time.Hour * 24).Unix()
		price := rand.Intn(100-1) + 1
		err = pAdmin.PublishGroupCID(id, []cid.ContentID{*cid1, *cid2, *cid3}, uint64(price), expiryDate, 42)
		if err != nil {
			fmt.Println("Error in publishing offer.")
			return
		}
		fmt.Printf("Published offer for cid: [\n%v\n%v\n%v\n] at a price of %v\n", cid1.ToString(), cid2.ToString(), cid3.ToString(), price)
	case "publish-offer-given-cids":
		if !initialised {
			fmt.Println("Admin hasn't been initialised yet.")
			return
		}
		if len(blocks) != 5 {
			fmt.Println("Usage: publish-offer-given-cids $nodeID $cid1 $cid2 $cid3")
			return
		}
		id, err := nodeid.NewNodeIDFromHexString(blocks[1])
		if err != nil {
			fmt.Printf("Invalid nodeID: %s\n", blocks[1])
			return
		}
		cid1, err := cid.NewContentIDFromHexString(blocks[2])
		if err != nil {
			fmt.Printf("Invalid cid: %s\n", blocks[2])
			return
		}
		cid2, err := cid.NewContentIDFromHexString(blocks[3])
		if err != nil {
			fmt.Printf("Invalid cid: %s\n", blocks[3])
			return
		}
		cid3, err := cid.NewContentIDFromHexString(blocks[4])
		if err != nil {
			fmt.Printf("Invalid cid: %s\n", blocks[4])
			return
		}
		expiryDate := time.Now().Local().Add(time.Hour * 24).Unix()
		price := 44
		err = pAdmin.PublishGroupCID(id, []cid.ContentID{*cid1, *cid2, *cid3}, uint64(price), expiryDate, 42)
		if err != nil {
			fmt.Println("Error in publishing offer.")
			return
		}
		fmt.Printf("Published offer for cid: [\n%v\n%v\n%v\n] at a price of %v\n", cid1.ToString(), cid2.ToString(), cid3.ToString(), price)
	case "publish-dht-offer":
		if !initialised {
			fmt.Println("Admin hasn't been initialised yet.")
			return
		}
		if len(blocks) != 3 {
			fmt.Println("Usage: publish-offer-given-cid $nodeID $cid")
			return
		}
		id, err := nodeid.NewNodeIDFromHexString(blocks[1])
		if err != nil {
			fmt.Printf("Invalid nodeID: %s\n", blocks[1])
			return
		}
		ccid, err := cid.NewContentIDFromHexString(blocks[2])
		if err != nil {
			fmt.Printf("Invalid cid: %s\n", blocks[2])
			return
		}
		expiryDate := time.Now().Local().Add(time.Hour * 24).Unix()
		err = pAdmin.PublishDHTCID(id, []cid.ContentID{*ccid}, []uint64{45}, []int64{expiryDate}, []uint64{43})
		if err != nil {
			fmt.Println("Error in publishing DHT offer.")
			return
		}
		fmt.Printf("Published DHT offer for cid: [\n%v\n]\n", ccid.ToString())
	case "exit":
		fmt.Println("Bye!")
		os.Exit(0)
	default:
		fmt.Printf("Unknown command: %s\n", blocks[0])
	}
}

func main() {
	initialised = false
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
