package poc2v2

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/register"
	"github.com/ConsenSys/fc-retrieval-gateway-admin/pkg/fcrgatewayadmin"
	"github.com/ConsenSys/fc-retrieval-itest/config"
	"github.com/ConsenSys/fc-retrieval-itest/pkg/util"
	"github.com/ConsenSys/fc-retrieval-provider-admin/pkg/fcrprovideradmin"
	"github.com/stretchr/testify/assert"
	tc "github.com/wcgcyx/testcontainers-go"
)

var lotusAP = "http://lotus-full-node:1234/rpc/v0"
var lotusToken string
var superAcct string
var gatewayConfig = config.NewConfig(".env.gateway")
var providerConfig = config.NewConfig(".env.provider")
var gwAdmin *fcrgatewayadmin.FilecoinRetrievalGatewayAdmin
var pAdmin *fcrprovideradmin.FilecoinRetrievalProviderAdmin

var ctx context.Context
var tag string

var gwIDs []*nodeid.NodeID
var pIDs []*nodeid.NodeID
var privateKeys []string
var accountAddrs []string

const (
	nGateways           = 2 // before 33
	nProviderContainers = 2 // before 3
)

func TestMain(m *testing.M) {

	// Need to make sure this env is not set in host machine
	itestEnv := os.Getenv("ITEST_CALLING_FROM_CONTAINER")

	if itestEnv != "" {
		lotusToken = os.Getenv("LOTUS_TOKEN")
		superAcct = os.Getenv("SUPER_ACCT")
		// Env is set, we are calling from docker container
		// This logging should be only called after all tests finished.
		m.Run()
		return
	}
	// Env is not set, we are calling from host
	// We need a redis, a register, 17 gateways and 3 providers
	tag = util.GetCurrentBranch()

	// Get env
	rgEnv := util.GetEnvMap("../../.env.register")
	gwEnv := util.GetEnvMap("../../.env.gateway")
	pvEnv := util.GetEnvMap("../../.env.provider")

	// Create shared net
	ctx = context.Background()
	network, networkName := util.CreateNetwork(ctx)
	defer (*network).Remove(ctx)

	// Start redis
	redisContainer := util.StartRedis(ctx, networkName, true)
	defer redisContainer.Terminate(ctx)

	// Start register
	registerContainer := util.StartRegister(ctx, tag, networkName, util.ColorYellow, rgEnv, true)
	defer registerContainer.Terminate(ctx)
	// Start providers
	var providerContainers []*tc.Container
	for i := 0; i < nProviderContainers; i++ {
		c := util.StartProvider(ctx, fmt.Sprintf("provider-%v", i), tag, networkName, util.ColorBlue, pvEnv, true)
		providerContainers = append(providerContainers, &c)
	}
	defer func() {
		for _, c := range providerContainers {
			(*c).Terminate(ctx)
		}
	}()
	// Start gateways
	var gatewayContainers []*tc.Container

	for i := 0; i < nGateways; i++ {
		c := util.StartGateway(ctx, fmt.Sprintf("gateway-%v", i), tag, networkName, util.ColorCyan, gwEnv, true)
		gatewayContainers = append(gatewayContainers, &c)
	}
	defer func() {
		for _, c := range gatewayContainers {
			(*c).Terminate(ctx)
		}
	}()

	// Start lotus
	lotusContainer := util.StartLotusFullNode(ctx, networkName, false)
	defer lotusContainer.Terminate(ctx)

	// Get lotus token
	lotusToken, superAcct = util.GetLotusToken()
	logging.Info("Lotus token is: %s", lotusToken)
	logging.Info("Super Acct is %s", superAcct)

	// Start itest
	done := make(chan bool)
	itestContainer := util.StartItest(ctx, tag, networkName, util.ColorGreen, lotusToken, superAcct, done, true)
	defer itestContainer.Terminate(ctx)
	// Block until done.
	if <-done {
		logging.Info("Tests passed, shutdown...")
	} else {
		logging.Fatal("Tests failed, shutdown...")
	}
}

func TestNewAccounts(t *testing.T) {
	t.Log("/*******************************************************/")
	t.Log("/*                 Start TestNewAccounts               */")
	t.Log("/*******************************************************/")

	var err error
	privateKeys, accountAddrs, err = util.GenerateAccount(lotusAP, lotusToken, superAcct, nGateways+nProviderContainers+3)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("/*******************************************************/")
	t.Log("/*                  End TestNewAccounts                */")
	t.Log("/*******************************************************/")
}

func TestInitialiseProviders(t *testing.T) {
	t.Skip(true)
	t.Log("/*******************************************************/")
	t.Log("/*             Start TestInitialiseProviders           */")
	t.Log("/*******************************************************/")

	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		panic(err)
	}

	confBuilder := fcrprovideradmin.CreateSettings()
	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	confBuilder.SetRegisterURL(providerConfig.GetString("REGISTER_API_URL"))
	conf := confBuilder.Build()
	pAdmin = fcrprovideradmin.NewFilecoinRetrievalProviderAdmin(*conf)

	for i := 0; i < nProviderContainers; i++ {

		walletKey := privateKeys[0]
		walletAddress := accountAddrs[0]
		privateKeys = privateKeys[1:]
		accountAddrs = accountAddrs[1:]

		providerRootKey, err := fcrcrypto.GenerateBlockchainKeyPair()
		if err != nil {
			panic(err)
		}
		providerRootSigningKey, err := providerRootKey.EncodePublicKey()
		if err != nil {
			panic(err)
		}
		providerRetrievalPrivateKey, err := fcrcrypto.GenerateRetrievalV1KeyPair()
		if err != nil {
			panic(err)
		}
		providerRetrievalSigningKey, err := providerRetrievalPrivateKey.EncodePublicKey()
		if err != nil {
			panic(err)
		}
		providerID := nodeid.NewRandomNodeID()
		pIDs = append(pIDs, providerID)

		identifier := fmt.Sprintf("-%v", i)
		providerRegister := &register.ProviderRegister{
			NodeID:             providerID.ToString(),
			Address:            walletAddress,
			RootSigningKey:     providerRootSigningKey,
			SigningKey:         providerRetrievalSigningKey,
			RegionCode:         providerConfig.GetString("PROVIDER_REGION_CODE"),
			NetworkInfoGateway: providerConfig.GetString("NETWORK_INFO_GATEWAY")[:8] + identifier + providerConfig.GetString("NETWORK_INFO_GATEWAY")[8:],
			NetworkInfoClient:  providerConfig.GetString("NETWORK_INFO_CLIENT")[:8] + identifier + providerConfig.GetString("NETWORK_INFO_CLIENT")[8:],
			NetworkInfoAdmin:   providerConfig.GetString("NETWORK_INFO_ADMIN")[:8] + identifier + providerConfig.GetString("NETWORK_INFO_ADMIN")[8:],
		}

		err = pAdmin.InitialiseProviderV2(providerRegister, providerRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1), walletKey, lotusAP, lotusToken)
		if err != nil {
			panic(err)
		}
	}

	t.Log("/*******************************************************/")
	t.Log("/*               End TestInitialiseProviders           */")
	t.Log("/*******************************************************/")
}

func TestInitialiseGateways(t *testing.T) {
	t.Log("/*******************************************************/")
	t.Log("/*             Start TestInitialiseGateways            */")
	t.Log("/*******************************************************/")

	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		panic(err)
	}

	confBuilder := fcrgatewayadmin.CreateSettings()
	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	confBuilder.SetRegisterURL(gatewayConfig.GetString("REGISTER_API_URL"))
	conf := confBuilder.Build()
	gwAdmin = fcrgatewayadmin.NewFilecoinRetrievalGatewayAdmin(*conf)

	// Only initialise 32 gateways, with one extra to initialise later to test list single cid offer
	for i := 0; i < nGateways; i++ {
		walletKey := privateKeys[0]
		walletAddress := accountAddrs[0]
		privateKeys = privateKeys[1:]
		accountAddrs = accountAddrs[1:]

		gatewayRootKey, err := fcrcrypto.GenerateBlockchainKeyPair()
		if err != nil {
			panic(err)
		}
		gatewayRootSigningKey, err := gatewayRootKey.EncodePublicKey()
		if err != nil {
			panic(err)
		}
		gatewayRetrievalPrivateKey, err := fcrcrypto.GenerateRetrievalV1KeyPair()
		if err != nil {
			panic(err)
		}
		gatewayRetrievalSigningKey, err := gatewayRetrievalPrivateKey.EncodePublicKey()
		if err != nil {
			panic(err)
		}
		var idStr string
		if i%2 == 0 {
			idStr = fmt.Sprintf("%X000000000000000000000000000000000000000000000000000000000000000", i/2)
		} else {
			idStr = fmt.Sprintf("%X800000000000000000000000000000000000000000000000000000000000000", i/2)
		}
		t.Log(idStr)

		gatewayID, err := nodeid.NewNodeIDFromHexString(idStr)
		if err != nil {
			panic(err)
		}
		gwIDs = append(gwIDs, gatewayID)

		identifier := fmt.Sprintf("-%v", i)
		gatewayRegister := &register.GatewayRegister{
			NodeID:              gatewayID.ToString(),
			Address:             walletAddress,
			RootSigningKey:      gatewayRootSigningKey,
			SigningKey:          gatewayRetrievalSigningKey,
			RegionCode:          gatewayConfig.GetString("GATEWAY_REGION_CODE"),
			NetworkInfoGateway:  gatewayConfig.GetString("NETWORK_INFO_GATEWAY")[:7] + identifier + gatewayConfig.GetString("NETWORK_INFO_GATEWAY")[7:],
			NetworkInfoProvider: gatewayConfig.GetString("NETWORK_INFO_PROVIDER")[:7] + identifier + gatewayConfig.GetString("NETWORK_INFO_PROVIDER")[7:],
			NetworkInfoClient:   gatewayConfig.GetString("NETWORK_INFO_CLIENT")[:7] + identifier + gatewayConfig.GetString("NETWORK_INFO_CLIENT")[7:],
			NetworkInfoAdmin:    gatewayConfig.GetString("NETWORK_INFO_ADMIN")[:7] + identifier + gatewayConfig.GetString("NETWORK_INFO_ADMIN")[7:],
		}

		err = gwAdmin.InitialiseGatewayV2(gatewayRegister, gatewayRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1), walletKey, lotusAP, lotusToken)
		if err != nil {
			panic(err)
		}
	}

	t.Log("/*******************************************************/")
	t.Log("/*               End TestInitialiseGateway            */")
	t.Log("/*******************************************************/")
}

func TestInitialiseClient(t *testing.T) {
	t.Log("/*******************************************************/")
	t.Log("/*             Start TestInitialiseClient              */")
	t.Log("/*******************************************************/")

	cmd := exec.Command("npm", "install")
	cmd.Dir = "/usr/src/github.com/ConsenSys/fc-retrieval-client-js/"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		t.Log("os.call ExitError ", err.(*exec.ExitError).String())
		t.Log("os.call error ", string(err.(*exec.ExitError).Stderr))
	}
	assert.Nil(t, err)

	cmd = exec.Command("npm", "run", "test-e2e")
	//cmd = exec.Command("npm", "run", "test-e2e-watch")
	cmd.Dir = "/usr/src/github.com/ConsenSys/fc-retrieval-client-js/"

	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		panic(err)
	}
	key, err := blockchainPrivateKey.EncodePublicKey()
	if err != nil {
		panic(err)
	}
	walletKey := privateKeys[0]
	privateKeys = privateKeys[1:]
	accountAddrs = accountAddrs[1:]

	cmd.Env = append(os.Environ(),
		"ESTABLISHMENT_TTL=101",
		fmt.Sprintf("BLOCKCHAIN_PUBLIC_KEY=%s", key),
		fmt.Sprintf("REGISTER_API_URL=%s", gatewayConfig.GetString("REGISTER_API_URL")),
		fmt.Sprintf("WALLET_PRIVATE_KEY=%s", walletKey),
		fmt.Sprintf("LOTUS_AP=%s", lotusAP),
		fmt.Sprintf("LOTUS_AUTH_TOKEN=%s", lotusToken),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		t.Log("os.call ExitError ", err.(*exec.ExitError).String())
		t.Log("os.call error ", string(err.(*exec.ExitError).Stderr))
	}
	assert.Nil(t, err)

	t.Log("/*******************************************************/")
	t.Log("/*               End TestInitialiseClient              */")
	t.Log("/*******************************************************/")
}
