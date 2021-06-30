package poc2v2

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	tc "github.com/wcgcyx/testcontainers-go"

	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrregistermgr"
	"github.com/ConsenSys/fc-retrieval/common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval/common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval/common/pkg/register"
	"github.com/ConsenSys/fc-retrieval/gateway-admin/pkg/fcrgatewayadmin"
	"github.com/ConsenSys/fc-retrieval/itest/config"
	"github.com/ConsenSys/fc-retrieval/itest/pkg/util"
	"github.com/ConsenSys/fc-retrieval/provider-admin/pkg/fcrprovideradmin"
)

var lotusAP = "http://lotus-full-node:1234/rpc/v0"
var lotusToken string
var superAcct string

var gatewayConfig = config.NewConfig(".env.gateway")
var providerConfig = config.NewConfig(".env.provider")

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
	tag := util.GetCurrentBranch()

	// Get env
	rgEnv := util.GetEnvMap("../../.env.register")
	gwEnv := util.GetEnvMap("../../.env.gateway")
	pvEnv := util.GetEnvMap("../../.env.provider")

	// Create shared net
	ctx := context.Background()
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

	reloadJsTests := os.Getenv("RELOAD_JS_TESTS")
	// Get lotus token
	lotusToken, superAcct = util.GetLotusToken()
	logging.Info("Lotus token is: %s", lotusToken)
	logging.Info("Super Acct is %s", superAcct)

	// Start itest
	done := make(chan bool)
	itestContainer := util.StartItest(ctx, tag, networkName, util.ColorGreen, lotusToken, superAcct, done, true, reloadJsTests)
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
	t.Log("/*                 Start POC2-JS TestNewAccounts               */")
	t.Log("/*******************************************************/")

	var err error
	privateKeys, accountAddrs, err := util.GenerateAccount(lotusAP, lotusToken, superAcct, nGateways+nProviderContainers+3)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("privateKeys: %#v", privateKeys)
	t.Logf("accountAddrs: %#v", accountAddrs)

	t.Log("/*******************************************************/")
	t.Log("/*                  End  POC2-JS TestNewAccounts                */")
	t.Log("/*******************************************************/")
}

func TestInitialiseProviders(t *testing.T) {
	//t.Skip(true)
	t.Log("/*******************************************************/")
	t.Log("/*             Start  POC2-JS TestInitialiseProviders           */")
	t.Log("/*******************************************************/")

	var err error
	privateKeys, accountAddrs, err := util.GenerateAccount(lotusAP, lotusToken, superAcct, 37)
	if err != nil {
		t.Fatal(err)
	}

	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		panic(err)
	}

	confBuilder := fcrprovideradmin.CreateSettings()
	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	confBuilder.SetRegisterURL(providerConfig.GetString("REGISTER_API_URL"))
	conf := confBuilder.Build()
	pAdmin := fcrprovideradmin.NewFilecoinRetrievalProviderAdmin(*conf)

	var pIDs []*nodeid.NodeID

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
		providerRegistrar := register.NewProviderRegister(
			providerID.ToString(),
			walletAddress,
			providerRootSigningKey,
			providerRetrievalSigningKey,
			providerConfig.GetString("PROVIDER_REGION_CODE"),
			providerConfig.GetString("NETWORK_INFO_GATEWAY")[:8]+identifier+providerConfig.GetString("NETWORK_INFO_GATEWAY")[8:],
			providerConfig.GetString("NETWORK_INFO_CLIENT")[:8]+identifier+providerConfig.GetString("NETWORK_INFO_CLIENT")[8:],
			providerConfig.GetString("NETWORK_INFO_ADMIN")[:8]+identifier+providerConfig.GetString("NETWORK_INFO_ADMIN")[8:],
		)

		err = pAdmin.InitialiseProviderV2(providerRegistrar, providerRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1), walletKey, lotusAP, lotusToken)
		if err != nil {
			logging.Error("error initialising provider: %s", err.Error())
			t.FailNow()
		}
	}

	t.Log("/*******************************************************/")
	t.Log("/*               End  POC2-JS TestInitialiseProviders           */")
	t.Log("/*******************************************************/")
}

func TestInitialiseGateways(t *testing.T) {
	t.Log("/*******************************************************/")
	t.Log("/*             Start  POC2-JS TestInitialiseGateways            */")
	t.Log("/*******************************************************/")

	var err error
	privateKeys, accountAddrs, err := util.GenerateAccount(lotusAP, lotusToken, superAcct, 37)
	if err != nil {
		t.Fatal(err)
	}

	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		panic(err)
	}

	confBuilder := fcrgatewayadmin.CreateSettings()
	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	confBuilder.SetRegisterURL(gatewayConfig.GetString("REGISTER_API_URL"))
	conf := confBuilder.Build()
	gwAdmin := fcrgatewayadmin.NewFilecoinRetrievalGatewayAdmin(*conf)

	var gwIDs []*nodeid.NodeID

	for i := 0; i < nGateways; i++ {
		walletKey := privateKeys[0]
		walletAddress := accountAddrs[0]
		privateKeys = privateKeys[1:]
		accountAddrs = accountAddrs[1:]

		gatewayRootPubKey, gatewayRetrievalPubKey, gatewayRetrievalPrivateKey, err := generateKeys()
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
		gatewayRegistrar := register.NewGatewayRegister(
			gatewayID.ToString(),
			walletAddress,
			gatewayRootPubKey,
			gatewayRetrievalPubKey,
			gatewayConfig.GetString("GATEWAY_REGION_CODE"),
			gatewayConfig.GetString("NETWORK_INFO_GATEWAY")[:7]+identifier+gatewayConfig.GetString("NETWORK_INFO_GATEWAY")[7:],
			gatewayConfig.GetString("NETWORK_INFO_PROVIDER")[:7]+identifier+gatewayConfig.GetString("NETWORK_INFO_PROVIDER")[7:],
			gatewayConfig.GetString("NETWORK_INFO_CLIENT")[:7]+identifier+gatewayConfig.GetString("NETWORK_INFO_CLIENT")[7:],
			gatewayConfig.GetString("NETWORK_INFO_ADMIN")[:7]+identifier+gatewayConfig.GetString("NETWORK_INFO_ADMIN")[7:],
		)

		err = gwAdmin.InitialiseGatewayV2(gatewayRegistrar, gatewayRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1), walletKey, lotusAP, lotusToken)
		if err != nil {
			logging.Error("error initialising gateway: %s", err.Error())
			t.FailNow()
		}
	}

	t.Log("/*******************************************************/")
	t.Log("/*               End  POC2-JS TestInitialiseGateway            */")
	t.Log("/*******************************************************/")
}

// Test client JS
func TestClientJS(t *testing.T) {
	t.Log("/*******************************************************/")
	t.Log("/*             Start  POC2-JS TestClientJS              */")
	t.Log("/*******************************************************/")

	assert.Nil(t, util.CallClientJsInstall())

	var err error
	privateKeys, accountAddrs, err := util.GenerateAccount(lotusAP, lotusToken, superAcct, 37)
	if err != nil {
		t.Fatal(err)
	}
	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		panic(err)
	}

	// Create and start register manager
	var rm = fcrregistermgr.NewFCRRegisterMgr(gatewayConfig.GetString("REGISTER_API_URL"), true, true, 10*time.Second)
	if err := rm.Start(); err != nil {
		logging.Error("error starting Register Manager: %s", err.Error())
		t.FailNow()
	}
	defer rm.ShutdownAndWait()

	// Init providers
	pConfBuilder := fcrprovideradmin.CreateSettings()
	pConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	pConfBuilder.SetRegisterURL(providerConfig.GetString("REGISTER_API_URL"))
	pCconf := pConfBuilder.Build()
	pAdmin := fcrprovideradmin.NewFilecoinRetrievalProviderAdmin(*pCconf)
	var pIDs []*nodeid.NodeID
	for i := 0; i < nProviderContainers; i++ {

		walletKey := privateKeys[0]
		walletAddress := accountAddrs[0]
		privateKeys = privateKeys[1:]
		accountAddrs = accountAddrs[1:]

		providerRootPubKey, providerRetrievalPubKey, providerRetrievalPrivateKey, err := generateKeys()
		providerID := nodeid.NewRandomNodeID()
		pIDs = append(pIDs, providerID)

		identifier := fmt.Sprintf("-%v", i)
		providerRegistrar := register.NewProviderRegister(
			providerID.ToString(),
			walletAddress,
			providerRootPubKey,
			providerRetrievalPubKey,
			providerConfig.GetString("PROVIDER_REGION_CODE"),
			providerConfig.GetString("NETWORK_INFO_GATEWAY")[:8]+identifier+providerConfig.GetString("NETWORK_INFO_GATEWAY")[8:],
			providerConfig.GetString("NETWORK_INFO_CLIENT")[:8]+identifier+providerConfig.GetString("NETWORK_INFO_CLIENT")[8:],
			providerConfig.GetString("NETWORK_INFO_ADMIN")[:8]+identifier+providerConfig.GetString("NETWORK_INFO_ADMIN")[8:],
		)

		err = pAdmin.InitialiseProviderV2(providerRegistrar, providerRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1), walletKey, lotusAP, lotusToken)
		if err != nil {
			logging.Error("error initialising provider: %s", err.Error())
			t.FailNow()
		}
		// Enroll the provider in the Register srv.
		if err := rm.RegisterProvider(providerRegistrar); err != nil {
			logging.Error("error registering provider: %s", err.Error())
			t.FailNow()
		}
	}

	// Init gateways
	gConfBuilder := fcrgatewayadmin.CreateSettings()
	gConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	gConfBuilder.SetRegisterURL(gatewayConfig.GetString("REGISTER_API_URL"))
	gConf := gConfBuilder.Build()
	gwAdmin := fcrgatewayadmin.NewFilecoinRetrievalGatewayAdmin(*gConf)
	var gwIDs []*nodeid.NodeID
	for i := 0; i < nGateways; i++ {
		walletKey := privateKeys[0]
		walletAddress := accountAddrs[0]
		privateKeys = privateKeys[1:]
		accountAddrs = accountAddrs[1:]

		gatewayRootPubKey, gatewayRetrievalPubKey, gatewayRetrievalPrivateKey, err := generateKeys()
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
		gatewayRegistrar := register.NewGatewayRegister(
			gatewayID.ToString(),
			walletAddress,
			gatewayRootPubKey,
			gatewayRetrievalPubKey,
			gatewayConfig.GetString("GATEWAY_REGION_CODE"),
			gatewayConfig.GetString("NETWORK_INFO_GATEWAY")[:7]+identifier+gatewayConfig.GetString("NETWORK_INFO_GATEWAY")[7:],
			gatewayConfig.GetString("NETWORK_INFO_PROVIDER")[:7]+identifier+gatewayConfig.GetString("NETWORK_INFO_PROVIDER")[7:],
			gatewayConfig.GetString("NETWORK_INFO_CLIENT")[:7]+identifier+gatewayConfig.GetString("NETWORK_INFO_CLIENT")[7:],
			gatewayConfig.GetString("NETWORK_INFO_ADMIN")[:7]+identifier+gatewayConfig.GetString("NETWORK_INFO_ADMIN")[7:],
		)
		err = gwAdmin.InitialiseGatewayV2(gatewayRegistrar, gatewayRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1), walletKey, lotusAP, lotusToken)
		if err != nil {
			logging.Error("error initialising gateway: %s", err.Error())
			t.FailNow()
		}
		// Enroll the gateway in the Register srv.
		if err := rm.RegisterGateway(gatewayRegistrar); err != nil {
			logging.Error("error registering gateway: %s", err.Error())
			t.FailNow()
		}
	}

	key, err := blockchainPrivateKey.EncodePublicKey()
	if err != nil {
		panic(err)
	}
	walletKey := privateKeys[0]
	privateKeys = privateKeys[1:]
	accountAddrs = accountAddrs[1:]

	err = util.CallClientJsE2E(key, walletKey, gatewayConfig.GetString("REGISTER_API_URL"), lotusAP, lotusToken)
	if err != nil {
		logging.Error("error calling JS client E2E: %s", err.Error())
		t.FailNow()
	}

	t.Log("/*******************************************************/")
	t.Log("/*               End  POC2-JS TestClientJS              */")
	t.Log("/*******************************************************/")
}

// Helper function to generate set of keys
func generateKeys() (rootPubKey string, retrievalPubKey string, retrievalPrivateKey *fcrcrypto.KeyPair, err error) {
	rootKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		return "", "", nil, fmt.Errorf("error generating blockchain key: %s", err.Error())
	}
	if rootKey == nil {
		return "", "", nil, errors.New("error generating blockchain key")
	}

	rootPubKey, err = rootKey.EncodePublicKey()
	if err != nil {
		return "", "", nil, fmt.Errorf("error encoding public key: %s", err.Error())
	}

	retrievalPrivateKey, err = fcrcrypto.GenerateRetrievalV1KeyPair()
	if err != nil {
		return "", "", nil, fmt.Errorf("error generating retrieval key: %s", err.Error())
	}
	if retrievalPrivateKey == nil {
		return "", "", nil, errors.New("error generating retrieval key")
	}

	retrievalPubKey, err = retrievalPrivateKey.EncodePublicKey()
	if err != nil {
		return "", "", nil, fmt.Errorf("error encoding retrieval pub key: %s", err.Error())
	}
	return
}
