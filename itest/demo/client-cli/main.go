package main

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/ConsenSys/fc-retrieval/client/pkg/fcrclient"
	"github.com/ConsenSys/fc-retrieval/common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval/common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrpaymentmgr"
	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrregistermgr"
	"github.com/ConsenSys/fc-retrieval/common/pkg/nodeid"
	"github.com/c-bata/go-prompt"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/api/apistruct"
	"github.com/filecoin-project/lotus/chain/types"
	cid2 "github.com/ipfs/go-cid"
)

var lotusAP = "http://127.0.0.1:1234/rpc/v0"
var client *fcrclient.FilecoinRetrievalClient
var offerMap map[string]*cidoffer.SubCIDOffer
var initialised bool
var registerURL = "http://127.0.0.1:9020"
var rm = fcrregistermgr.NewFCRRegisterMgr(registerURL, true, true, 2*time.Second)

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "init", Description: "Initialise the client"},
		{Text: "ls-registered", Description: "List registered gateways"},
		{Text: "ls-active", Description: "List gateways currently in use"},
		{Text: "add-active", Description: "Add active gateway"},
		{Text: "find-offer", Description: "Find offers for given cid"},
		{Text: "find-offer-dht", Description: "Find offers for given cid using DHT discovery"},
		{Text: "list-offers", Description: "List obtained offers"},
		{Text: "retrieve", Description: "Retrieve data using an offer"},
		{Text: "retrieve-fast", Description: "Fast-retrieve data (automated offer discovery, selection and data retrieval)"},
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
			fmt.Println("Client has already been initialised.")
			return
		}
		fmt.Println("Initialise client (dev)...")
		token, acct := getLotusToken()
		keys, addresses, err := generateAccount(lotusAP, token, acct, 1)
		if err != nil {
			fmt.Printf("Fail to initialise client: %s\n", err.Error())
			return
		}
		key := keys[0]
		address := addresses[0]
		fmt.Printf("(Dev) Created an account for this client to use: %s\n", address)
		blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
		if err != nil {
			fmt.Println("Error in generating client key.")
			return
		}
		confBuilder := fcrclient.CreateSettings()
		confBuilder.SetEstablishmentTTL(101)
		confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
		confBuilder.SetWalletPrivateKey(key)
		confBuilder.SetLotusAP(lotusAP)
		confBuilder.SetLotusAuthToken(token)
		conf := confBuilder.Build()
		err = rm.Start()
		if err != nil {
			fmt.Printf("Fail to start rm for client: %s\n", err.Error())
			return
		}
		client, err = fcrclient.NewFilecoinRetrievalClient(*conf, rm)
		if err != nil {
			fmt.Printf("Fail to initialise client: %s\n", err.Error())
			return
		}
		res := client.PaymentMgr()
		if res == nil {
			fmt.Println("Fail to initialise client's payment manager.")
			return
		}
		offerMap = make(map[string]*cidoffer.SubCIDOffer)
		initialised = true
	case "ls-registered":
		if !initialised {
			fmt.Println("Client hasn't been initialised yet.")
			return
		}
		gws := rm.GetAllGateways()
		fmt.Println("Registered gateways:")
		for _, gw := range gws {
			fmt.Printf("%v\n", gw.GetNodeID())
		}

	case "ls-active":
		if !initialised {
			fmt.Println("Client hasn't been initialised yet.")
			return
		}
		gws := client.GetActiveGateways()
		fmt.Println("Current gateways in use:")
		for _, gw := range gws {
			fmt.Printf("%v\n", gw.ToString())
		}
	case "add-active":
		if !initialised {
			fmt.Println("Client hasn't been initialised yet.")
			return
		}
		if len(blocks) != 2 {
			fmt.Println("Usage: add-active $nodeID")
			return
		}
		id, err := nodeid.NewNodeIDFromHexString(blocks[1])
		if err != nil {
			fmt.Printf("Invalid nodeID: %s\n", blocks[1])
			return
		}
		client.AddGatewaysToUse([]*nodeid.NodeID{id})
		added := client.AddActiveGateways([]*nodeid.NodeID{id})
		if added != 1 {
			fmt.Println("Fail to use gateway.")
			return
		}
		info := rm.GetGateway(id)
		err = client.PaymentMgr().Topup(info.GetAddress(), client.Settings.TopUpAmount())
		if err != nil {
			fmt.Println("Error in creating payment channel to given gateway")
			return
		}
		fmt.Println("Add gateway successful.")
	case "find-offer":
		if !initialised {
			fmt.Println("Client hasn't been initialised yet.")
			return
		}
		if len(blocks) != 2 {
			fmt.Println("Usage: find-offer $contentID")
			return
		}
		id, err := cid.NewContentIDFromHexString(blocks[1])
		if err != nil {
			fmt.Printf("Invalid contentID: %s\n", blocks[1])
			return
		}
		temp := make(map[string]*cidoffer.SubCIDOffer)
		gws := client.GetActiveGateways()
		for _, gw := range gws {
			res, err := client.FindOffersStandardDiscoveryV2(id, gw, 5)
			if err != nil {
				fmt.Printf("Error querying gateway %v: %v\n", gw.ToString(), err.Error())
				continue
			}
			for _, offer := range res {
				temp[offer.GetSignature()] = &offer
			}
		}
		fmt.Println("Find offers:")
		for key, val := range temp {
			offerMap[key] = val
			fmt.Printf("Offer %v: provider-%v, cid-%v, price-%v, expiry-%v, qos-%v\n", key, val.GetProviderID().ToString(), val.GetSubCID().ToString(), val.GetPrice(), val.GetExpiry(), val.GetQoS())
		}
	case "find-offer-dht":
		if !initialised {
			fmt.Println("Client hasn't been initialised yet.")
			return
		}
		if len(blocks) != 3 {
			fmt.Println("Usage: find-offer-dht $nodeID $contentID")
			return
		}
		node, err := nodeid.NewNodeIDFromHexString(blocks[1])
		if err != nil {
			fmt.Printf("Invalid nodeID: %s\n", blocks[1])
			return
		}
		id, err := cid.NewContentIDFromHexString(blocks[2])
		if err != nil {
			fmt.Printf("Invalid contentID: %s\n", blocks[2])
			return
		}
		res, err := client.FindOffersDHTDiscoveryV2(id, node, 4, 4)
		if err != nil {
			fmt.Printf("Error in dht discovery: %s\n", err.Error())
			return
		}
		fmt.Println("Find offers:")
		for key, val := range res {
			fmt.Printf("Gateway: %v\n", key)
			for _, o := range *val {
				fmt.Printf("\tOffer %v: provider-%v, cid-%v, price-%v, expiry-%v, qos-%v\n", key, o.GetProviderID().ToString(), o.GetSubCID().ToString(), o.GetPrice(), o.GetExpiry(), o.GetQoS())
			}
		}
		fmt.Println("Note <dev>: Offers found are not saved.")
	case "list-offers":
		if !initialised {
			fmt.Println("Client hasn't been initialised yet.")
			return
		}
		fmt.Println("Current offers:")
		for key, val := range offerMap {
			fmt.Printf("Offer %v: provider-%v, cid-%v, price-%v, expiry-%v, qos-%v\n", key, val.GetProviderID().ToString(), val.GetSubCID().ToString(), val.GetPrice(), val.GetExpiry(), val.GetQoS())
		}
	case "retrieve":
		if !initialised {
			fmt.Println("Client hasn't been initialised yet.")
			return
		}
		fmt.Println("Retrieve hasn't been implemented yet.")
	case "retrieve-fast":
		if !initialised {
			fmt.Println("Client hasn't been initialised yet.")
			return
		}
		fmt.Println("Fast retrieve hasn't been implemented yet.")
	case "exit":
		rm.Shutdown()
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
