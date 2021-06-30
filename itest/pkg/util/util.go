/*
Package util - common functions used in end-to-end and integration testing. Allowing to start different types of
Retrieval network nodes for testing.
*/
package util

import (
	"bufio"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/api/apistruct"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/google/uuid"
	"github.com/ipfs/go-cid"
	"github.com/testcontainers/testcontainers-go/wait"
	tc "github.com/wcgcyx/testcontainers-go"

	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrpaymentmgr"
	"github.com/ConsenSys/fc-retrieval/common/pkg/logging"
)

const ColorRed = "\033[31m"
const ColorGreen = "\033[32m"
const ColorYellow = "\033[33m"
const ColorBlue = "\033[34m"
const ColorPurple = "\033[35m"
const ColorCyan = "\033[36m"
const ColorWhite = "\033[37m" // Used by redis
const ColorGray = "\033[90m"  // Used by lotus
const ColorBrightRed = "\033[91m"
const ColorBrightGreen = "\033[91m"
const ColorBrightYellow = "\033[91m"
const ColorBrightBlue = "\033[91m"
const ColorBrightPurple = "\033[91m"
const ColorBrightCyan = "\033[91m"
const ColorBrightWhite = "\033[91m"

const lotusDaemonWaitFor = "retrieval client"
const lotusFullNodeWaitFor = "starting winning PoSt warmup"
const networkMode = "default"

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

// CreateNetwork creates a network
func CreateNetwork(ctx context.Context) (*tc.Network, string) {
	randomUuid, _ := uuid.NewRandom()
	networkName := randomUuid.String()
	net, err := tc.GenericNetwork(ctx, tc.GenericNetworkRequest{
		NetworkRequest: tc.NetworkRequest{
			Name:           networkName,
			CheckDuplicate: true,
		},
	})
	if err != nil {
		panic(err)
	}
	return &net, networkName
}

// GetEnvMap gets the env map from a given env file
func GetEnvMap(envFile string) map[string]string {
	env := make(map[string]string)
	file, err := os.Open(envFile)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		res := strings.Split(scanner.Text(), "=")
		if len(res) >= 2 {
			env[res[0]] = res[1]
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	if err := file.Close(); err != nil {
		logging.Error("error closing file: %s", err.Error())
	}
	return env
}

// StartLotusDaemon starts Lotus local development network, daemon only (miner is missing)
func StartLotusDaemon(ctx context.Context, network string, verbose bool) tc.Container {
	// Start lotus
	req := tc.ContainerRequest{
		Image:          "consensys/lotus-daemon:latest",
		Networks:       []string{network},
		NetworkMode:    container.NetworkMode(networkMode),
		NetworkAliases: map[string][]string{network: {"lotus-daemon"}},
		WaitingFor:     wait.ForLog(lotusDaemonWaitFor),
		ExposedPorts:   []string{"1234"},
		AutoRemove:     true,
	}
	lotusC, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}
	if verbose {
		g := &logConsumer{name: fmt.Sprintf("lotus-daemon"), color: ColorGray}
		err = lotusC.StartLogProducer(ctx)
		if err != nil {
			panic(err)
		}
		lotusC.FollowOutput(g)
	}
	return lotusC
}

// StartLotusFullNode starts Lotus local development network, two services: miner and daemon in one container
func StartLotusFullNode(ctx context.Context, network string, verbose bool) tc.Container {
	// Start lotus
	req := tc.ContainerRequest{
		Image:          "consensys/lotus-full-node:latest",
		Networks:       []string{network},
		NetworkMode:    container.NetworkMode(networkMode),
		NetworkAliases: map[string][]string{network: {"lotus-full-node"}},
		WaitingFor:     wait.ForLog(lotusFullNodeWaitFor).WithStartupTimeout(10 * time.Minute),
		ExposedPorts:   []string{"1234", "2345"},
		AutoRemove:     true,
		// --cpus=<value>
	}
	lotusC, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}
	if verbose {
		g := &logConsumer{name: fmt.Sprintf("lotus-full-node"), color: ColorGray}
		err = lotusC.StartLogProducer(ctx)
		if err != nil {
			panic(err)
		}
		lotusC.FollowOutput(g)
	}
	return lotusC
}

// StartRedis - starts redis; used by register
func StartRedis(ctx context.Context, network string, verbose bool) tc.Container {
	// Start redis
	req := tc.ContainerRequest{
		Image:          "redis:alpine",
		Cmd:            []string{"redis-server", "--requirepass", "xxxx"},
		Networks:       []string{network},
		NetworkMode:    container.NetworkMode(networkMode),
		NetworkAliases: map[string][]string{network: {"redis"}},
		WaitingFor:     wait.ForLog("Ready to accept connections"),
		AutoRemove:     true,
	}
	redisC, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}
	if verbose {
		g := &logConsumer{name: fmt.Sprintf("redis"), color: ColorWhite}
		err = redisC.StartLogProducer(ctx)
		if err != nil {
			panic(err)
		}
		redisC.FollowOutput(g)
	}
	return redisC
}

// StartRegister - starts the register
func StartRegister(ctx context.Context, network string, color string, env map[string]string, verbose bool) tc.Container {
	// Start a register container
	req := tc.ContainerRequest{
		Image:          "consensys/fc-retrieval/register:latest",
		Networks:       []string{network},
		Env:            env,
		NetworkMode:    container.NetworkMode(networkMode),
		NetworkAliases: map[string][]string{network: {"register"}},
		WaitingFor:     wait.ForLog("Serving register at"),
		AutoRemove:     true,
	}
	registerC, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}
	if verbose {
		g := &logConsumer{name: fmt.Sprintf("register"), color: color}
		err = registerC.StartLogProducer(ctx)
		if err != nil {
			panic(err)
		}
		registerC.FollowOutput(g)
	}
	return registerC
}

// StartGateway - start a gateway of specific id, network, log color and env
func StartGateway(ctx context.Context, id string, network string, color string, env map[string]string, verbose bool) tc.Container {
	// Start a gateway container
	req := tc.ContainerRequest{
		Image:          "consensys/fc-retrieval/gateway:latest",
		Networks:       []string{network},
		Env:            env,
		NetworkMode:    container.NetworkMode(networkMode),
		NetworkAliases: map[string][]string{network: {id}},
		WaitingFor:     wait.ForLog("Filecoin Gateway Start-up Complete"),
		AutoRemove:     true,
	}
	gatewayC, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}
	if verbose {
		g := &logConsumer{name: id, color: color}
		err = gatewayC.StartLogProducer(ctx)
		if err != nil {
			panic(err)
		}
		gatewayC.FollowOutput(g)
	}
	return gatewayC
}

// StartProvider - start a provider of specific id, network, log color and env
func StartProvider(ctx context.Context, id string, network string, color string, env map[string]string, verbose bool) tc.Container {
	// Start a provider container
	req := tc.ContainerRequest{
		Image:          "consensys/fc-retrieval/provider:latest",
		Networks:       []string{network},
		Env:            env,
		NetworkMode:    container.NetworkMode(networkMode),
		NetworkAliases: map[string][]string{network: {id}},
		WaitingFor:     wait.ForLog("Filecoin Provider Start-up Complete"),
		AutoRemove:     true,
	}
	providerC, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}
	if verbose {
		g := &logConsumer{name: id, color: color}
		err = providerC.StartLogProducer(ctx)
		if err != nil {
			panic(err)
		}
		providerC.FollowOutput(g)
	}
	return providerC
}

// StartItest - start the itest, must only be called in host
func StartItest(ctx context.Context, network string, color string, lotusToken string, superAcct string, done chan bool, verbose bool, reloadJsTests string) tc.Container {
	// Start a itest container
	// Mount testdir
	absPath, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}
	// Mount common, client, gw-admin, pvd-admin
	commonPath, err := filepath.Abs("../../../common/pkg")
	if err != nil {
		panic(err)
	}
	clientPath, err := filepath.Abs("../../../client/pkg")
	if err != nil {
		panic(err)
	}
	gwAdminPath, err := filepath.Abs("../../../gateway-admin/pkg")
	if err != nil {
		panic(err)
	}
	pvdAdminPath, err := filepath.Abs("../../../provider-admin/pkg")
	if err != nil {
		panic(err)
	}
	clientJsPath, err := filepath.Abs("../../../client-js")
	if err != nil {
		panic(err)
	}

	req := tc.ContainerRequest{
		Image:          "consensys/fc-retrieval/itest",
		Name:           "itest",
		Networks:       []string{network},
		Env:            map[string]string{"ITEST_CALLING_FROM_CONTAINER": "yes", "LOTUS_TOKEN": lotusToken, "SUPER_ACCT": superAcct, "RELOAD_JS_TESTS": reloadJsTests},
		NetworkMode:    container.NetworkMode(networkMode),
		NetworkAliases: map[string][]string{network: {"itest"}},
		BindMounts: map[string]string{
			clientJsPath: "/go/src/github.com/ConsenSys/fc-retrieval/client-js/",
			absPath:      "/go/src/github.com/ConsenSys/fc-retrieval/itest/pkg/temp/",
			commonPath:   "/go/src/github.com/ConsenSys/fc-retrieval/common/pkg/",
			clientPath:   "/go/src/github.com/ConsenSys/fc-retrieval/client/pkg/",
			gwAdminPath:  "/go/src/github.com/ConsenSys/fc-retrieval/gateway-admin/pkg/",
			pvdAdminPath: "/go/src/github.com/ConsenSys/fc-retrieval/provider-admin/pkg/"},
		Cmd:        []string{"go", "test", "-v", "--count=1", "/go/src/github.com/ConsenSys/fc-retrieval/itest/pkg/temp/"},
		AutoRemove: true,
	}

	if reloadJsTests == "yes" {
		req.Cmd = append(req.Cmd, "-timeout=0")
	}

	itestC, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}
	if verbose {
		g := &logConsumer{name: fmt.Sprintf("itest"), color: color, done: done}
		err = itestC.StartLogProducer(ctx)
		if err != nil {
			panic(err)
		}
		itestC.FollowOutput(g)
	}
	return itestC
}

type logConsumer struct {
	name  string
	color string
	// The following are used by itest only
	done chan bool
}

func (g *logConsumer) Accept(l tc.Log) {
	log := string(l.Content)
	fmt.Print(g.color, "[", strings.ToUpper(g.name), "]", "\033[0m:", log)
	if g.done != nil {
		if strings.Contains(log, "--- FAIL:") {
			// Tests have falied.
			g.done <- false
		} else if strings.Contains(log, "ok") && strings.Contains(log, "github.com/ConsenSys/fc-retrieval/itest/pkg/") {
			// Tests have all passed.
			g.done <- true
		}
	}
}

// GenerateAccount - helper method, used to generate a new filecoin account with 10 filecoins of balance
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
		contentID, err := api.MpoolPush(context.Background(), signedMsg)
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

func CallClientJsInstall() error {
	cmd := exec.Command("npm", "install")
	cmd.Dir = "../../../client-js/"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func CallClientJsE2E(key string, walletKey string, registerAPI string, lotusAP string, lotusToken string, gtws string, prvs string) error {
	cmd := exec.Command("npm", "run", "test-e2e")
	if os.Getenv("RELOAD_JS_TESTS") == "yes" {
		cmd = exec.Command("npm", "run", "test-e2e-watch")
	}
	cmd.Dir = "/usr/src/github.com/ConsenSys/client-js"

	cmd.Env = append(os.Environ(),
		"ESTABLISHMENT_TTL=101",
		fmt.Sprintf("FCR_BLOCKCHAIN_PUBLIC_KEY=%s", key),
		fmt.Sprintf("FCR_REGISTER_API_URL=%s", registerAPI),
		fmt.Sprintf("FCR_WALLET_PRIVATE_KEY=%s", walletKey),
		fmt.Sprintf("FCR_LOTUS_AP=%s", lotusAP),
		fmt.Sprintf("FCR_LOTUS_AUTH_TOKEN=%s", lotusToken),
		fmt.Sprintf("FCR_GATEWAYS=%s", gtws),
		fmt.Sprintf("FCR_PROVIDERS=%s", prvs),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
