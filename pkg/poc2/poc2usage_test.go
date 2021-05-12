package poc2

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/ConsenSys/fc-retrieval-client/pkg/fcrclient"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/register"
	"github.com/ConsenSys/fc-retrieval-gateway-admin/pkg/fcrgatewayadmin"
	"github.com/ConsenSys/fc-retrieval-itest/config"
	"github.com/ConsenSys/fc-retrieval-itest/pkg/util"
	"github.com/ConsenSys/fc-retrieval-provider-admin/pkg/fcrprovideradmin"
)

var gatewayConfig = config.NewConfig(".env.gateway")
var providerConfig = config.NewConfig(".env.provider")
var gwAdmin *fcrgatewayadmin.FilecoinRetrievalGatewayAdmin
var pAdmin *fcrprovideradmin.FilecoinRetrievalProviderAdmin
var client *fcrclient.FilecoinRetrievalClient
var gwIDs []*nodeid.NodeID
var pIDs []*nodeid.NodeID

func TestMain(m *testing.M) {
	// Need to make sure this env is not set in host machine
	itestEnv := os.Getenv("ITEST_CALLING_FROM_CONTAINER")

	if itestEnv != "" {
		// Env is set, we are calling from docker container
		// This logging should be only called after all tests finished.
		m.Run()
		return
	}
	// Env is not set, we are calling from host
	// We need a redis, a register, 17 gateways and 3 providers
	tag := util.GetCurrentBranch()
	network := "itest-shared"
	util.CleanContainers(network)

	// Get env
	rgEnv := util.GetEnvMap("../../.env.register")
	gwEnv := util.GetEnvMap("../../.env.gateway")
	pvEnv := util.GetEnvMap("../../.env.provider")

	// Create shared net
	ctx := context.Background()
	net := *util.CreateNetwork(ctx, network)
	defer net.Remove(ctx)

	// Start redis
	util.StartRedis(ctx, network, true)

	// Start register
	util.StartRegister(ctx, tag, network, util.ColorYellow, rgEnv, true)

	// Start 3 providers
	for i := 0; i < 3; i++ {
		util.StartProvider(ctx, fmt.Sprintf("provider-%v", i), tag, network, util.ColorBlue, pvEnv, true)
	}

	// Start 17 gateways
	for i := 0; i < 17; i++ {
		util.StartGateway(ctx, fmt.Sprintf("gateway-%v", i), tag, network, util.ColorCyan, gwEnv, true)
	}

	// Start itest
	done := make(chan bool)
	util.StartItest(ctx, tag, network, util.ColorGreen, done, true)

	// Block until done.
	if <-done {
		logging.Info("Tests passed, shutdown...")
	} else {
		logging.Fatal("Tests failed, shutdown...")
	}
	// Clean containers to shutdown
	util.CleanContainers(network)
}

func TestInitialiseProviders(t *testing.T) {
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

	for i := 0; i < 3; i++ {
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
		if err != nil {
			panic(err)
		}
		pIDs = append(pIDs, providerID)

		identifier := fmt.Sprintf("-%v", i)
		providerRegister := &register.ProviderRegister{
			NodeID:             providerID.ToString(),
			Address:            providerConfig.GetString("PROVIDER_ADDRESS"),
			RootSigningKey:     providerRootSigningKey,
			SigningKey:         providerRetrievalSigningKey,
			RegionCode:         providerConfig.GetString("PROVIDER_REGION_CODE"),
			NetworkInfoGateway: providerConfig.GetString("NETWORK_INFO_GATEWAY")[:8] + identifier + providerConfig.GetString("NETWORK_INFO_GATEWAY")[8:],
			NetworkInfoClient:  providerConfig.GetString("NETWORK_INFO_CLIENT")[:8] + identifier + providerConfig.GetString("NETWORK_INFO_CLIENT")[8:],
			NetworkInfoAdmin:   providerConfig.GetString("NETWORK_INFO_ADMIN")[:8] + identifier + providerConfig.GetString("NETWORK_INFO_ADMIN")[8:],
		}

		err = pAdmin.InitialiseProvider(providerRegister, providerRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1))
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

	// Only initialise 16 gateways, with one extra to initialise later to test list single cid offer
	for i := 0; i < 16; i++ {
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
		gatewayID := nodeid.NewRandomNodeID()
		if err != nil {
			panic(err)
		}
		gwIDs = append(gwIDs, gatewayID)

		identifier := fmt.Sprintf("-%v", i)
		gatewayRegister := &register.GatewayRegister{
			NodeID:              gatewayID.ToString(),
			Address:             gatewayConfig.GetString("GATEWAY_ADDRESS"),
			RootSigningKey:      gatewayRootSigningKey,
			SigningKey:          gatewayRetrievalSigningKey,
			RegionCode:          gatewayConfig.GetString("GATEWAY_REGION_CODE"),
			NetworkInfoGateway:  gatewayConfig.GetString("NETWORK_INFO_GATEWAY")[:7] + identifier + gatewayConfig.GetString("NETWORK_INFO_GATEWAY")[7:],
			NetworkInfoProvider: gatewayConfig.GetString("NETWORK_INFO_PROVIDER")[:7] + identifier + gatewayConfig.GetString("NETWORK_INFO_PROVIDER")[7:],
			NetworkInfoClient:   gatewayConfig.GetString("NETWORK_INFO_CLIENT")[:7] + identifier + gatewayConfig.GetString("NETWORK_INFO_CLIENT")[7:],
			NetworkInfoAdmin:    gatewayConfig.GetString("NETWORK_INFO_ADMIN")[:7] + identifier + gatewayConfig.GetString("NETWORK_INFO_ADMIN")[7:],
		}

		err = gwAdmin.InitialiseGateway(gatewayRegister, gatewayRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1))
		if err != nil {
			panic(err)
		}
	}

	t.Log("/*******************************************************/")
	t.Log("/*               End TestInitialiseGateway            */")
	t.Log("/*******************************************************/")
}
