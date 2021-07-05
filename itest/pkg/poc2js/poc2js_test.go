package poc2v2

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrregistermgr"
	"github.com/ConsenSys/fc-retrieval/common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval/common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval/common/pkg/register"
	"github.com/ConsenSys/fc-retrieval/gateway-admin/pkg/fcrgatewayadmin"
	"github.com/ConsenSys/fc-retrieval/itest/config"
	cr "github.com/ConsenSys/fc-retrieval/itest/pkg/util/crypto-facade"
	fil "github.com/ConsenSys/fc-retrieval/itest/pkg/util/filecoin-facade"
	js "github.com/ConsenSys/fc-retrieval/itest/pkg/util/jsclient-facade"
	tc "github.com/ConsenSys/fc-retrieval/itest/pkg/util/test-containers"

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
	lotusToken, superAcct = fil.GetLotusToken()
	logging.Info("Lotus token is: %s", lotusToken)
	logging.Info("Super Acct is %s", superAcct)

	const testName = "poc2-new-gateway"
	ctx := context.Background()
	containersAndPorts, network, err := tc.StartContainers(ctx, nGateways, nProviderContainers, testName, true)
	if err != nil {
		logging.Error("%s failed, container starting error: %s", testName, err.Error())
		tc.StopContainers(ctx, containersAndPorts, network)
		os.Exit(1)
	}
	defer tc.StopContainers(ctx, containersAndPorts, network)
	m.Run()
}

func TestNewAccounts(t *testing.T) {
	t.Log("/*******************************************************/")
	t.Log("/*                 Start POC2-JS TestNewAccounts               */")
	t.Log("/*******************************************************/")

	ctx := context.Background()
	var err error
	privateKeys, accountAddrs, err := fil.GenerateAccount(ctx, lotusAP, lotusToken, superAcct, nGateways+nProviderContainers+3)
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

	ctx := context.Background()
	var err error
	privateKeys, accountAddrs, err := fil.GenerateAccount(ctx, lotusAP, lotusToken, superAcct, 37)
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

	ctx := context.Background()
	var err error
	privateKeys, accountAddrs, err := fil.GenerateAccount(ctx, lotusAP, lotusToken, superAcct, 37)
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

		gatewayRootPubKey, gatewayRetrievalPubKey, gatewayRetrievalPrivateKey, err := cr.GenerateKeys()
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
			logging.Error("gateway initialising error: %s", err.Error())
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

	assert.Nil(t, js.CallClientJsInstall())

	ctx := context.Background()
	var err error
	privateKeys, accountAddrs, err := fil.GenerateAccount(ctx, lotusAP, lotusToken, superAcct, 37)
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

		providerRootPubKey, providerRetrievalPubKey, providerRetrievalPrivateKey, err := cr.GenerateKeys()
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

		gatewayRootPubKey, gatewayRetrievalPubKey, gatewayRetrievalPrivateKey, err := cr.GenerateKeys()
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
			logging.Error("gateway initialising error: %s", err.Error())
			t.FailNow()
		}
		// Enroll the gateway in the Register srv.
		if err := rm.RegisterGateway(gatewayRegistrar); err != nil {
			logging.Error("gateway registering error: %s", err.Error())
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

	err = js.CallClientJsE2E(key, walletKey, gatewayConfig.GetString("REGISTER_API_URL"), lotusAP, lotusToken)
	if err != nil {
		logging.Error("error calling JS client E2E: %s", err.Error())
		t.FailNow()
	}

	t.Log("/*******************************************************/")
	t.Log("/*               End  POC2-JS TestClientJS              */")
	t.Log("/*******************************************************/")
}
