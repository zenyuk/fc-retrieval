package poc2_group_offer

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ConsenSys/fc-retrieval/client/pkg/fcrclient"
	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrregistermgr"
	"github.com/ConsenSys/fc-retrieval/common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval/common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval/common/pkg/register"
	"github.com/ConsenSys/fc-retrieval/gateway-admin/pkg/fcrgatewayadmin"
	cr "github.com/ConsenSys/fc-retrieval/itest/pkg/util/crypto-facade"
	fil "github.com/ConsenSys/fc-retrieval/itest/pkg/util/filecoin-facade"
	"github.com/ConsenSys/fc-retrieval/provider-admin/pkg/fcrprovideradmin"
)

func TestForceUpdate(t *testing.T) {
	t.Log("/*******************************************************/")
	t.Log("/*                Start TestForceUpdate                */")
	t.Log("/*******************************************************/")

	ctx := context.Background()
	lotusToken, superAcct := fil.GetLotusToken()
	lotusDaemonApiEndpoint, _ := containers.Lotus.GetLostHostApiEndpoints()
	var lotusAP = "http://" + lotusDaemonApiEndpoint + "/rpc/v0"
	var err error
	privateKeys, accountAddrs, err := fil.GenerateAccount(ctx, lotusAP, lotusToken, superAcct, 37)
	if err != nil {
		t.Fatal(err)
	}

	// Main key used across the chain
	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		logging.Error("error generating blockchain key: %s", err.Error())
		t.FailNow()
	}

	// Create and start register manager
	registerApiEndpoint := "http://" + containers.Register.GetRegisterHostApiEndpoint()
	var rm = fcrregistermgr.NewFCRRegisterMgr(registerApiEndpoint, false, true, 10*time.Second)
	if err := rm.Start(); err != nil {
		logging.Error("error starting Register Manager: %s", err.Error())
		t.FailNow()
	}
	defer rm.ShutdownAndWait()

	walletKey := privateKeys[0]
	// walletAddress := accountAddrs[0]
	privateKeys = privateKeys[1:]
	accountAddrs = accountAddrs[1:]

	// Configure client
	clientConfBuilder := fcrclient.CreateSettings()
	clientConfBuilder.SetEstablishmentTTL(101)
	clientConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	clientConfBuilder.SetRegisterURL(registerApiEndpoint)
	clientConfBuilder.SetWalletPrivateKey(walletKey)
	clientConfBuilder.SetLotusAP(lotusAP)
	clientConfBuilder.SetLotusAuthToken(lotusToken)
	clientConf := clientConfBuilder.Build()
	client, err := fcrclient.NewFilecoinRetrievalClient(*clientConf, rm)
	if !assert.Nil(t, err, "Error should be nil") {
		t.Fatal(err)
	}
	res := client.PaymentMgr()
	if !assert.NotNil(t, res, "Fail to initialise a payment manager") {
		t.FailNow()
	}

	// Configure gateway admin
	gConfBuilder := fcrgatewayadmin.CreateSettings()
	gConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	gConfBuilder.SetRegisterURL("http://" + containers.Register.GetRegisterHostApiEndpoint())
	gConf := gConfBuilder.Build()
	gwAdmin := fcrgatewayadmin.NewFilecoinRetrievalGatewayAdmin(*gConf)

	// Initialise gateways

	// map between gateway ID and gateway name
	var gateways = make(map[string]*nodeid.NodeID)
	// Only initialise 32 gateways, with one extra to initialise later to test list single cid offer
	for i := 0; i < 32; i++ {
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

		gatewayName := fmt.Sprintf("gateway-%v", i)
		gateways[gatewayName] = gatewayID
		_, _, gatewayClientApiEndpoint, gatewayAdminApiEndpoint := containers.Gateways[gatewayName].GetGatewayHostApiEndpoints()
		gatewayRegistrar := register.NewGatewayRegister(
			gatewayID.ToString(),
			walletAddress,
			gatewayRootPubKey,
			gatewayRetrievalPubKey,
			gatewayConfig.GetString("GATEWAY_REGION_CODE"),
			gatewayName+":"+gatewayConfig.GetString("BIND_GATEWAY_API"),
			gatewayName+":"+gatewayConfig.GetString("BIND_PROVIDER_API"),
			gatewayName+":"+gatewayConfig.GetString("BIND_REST_API"),
			gatewayName+":"+gatewayConfig.GetString("BIND_ADMIN_API"),
		)

		err = gwAdmin.InitialiseGatewayV2(gatewayAdminApiEndpoint, gatewayRegistrar, gatewayRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1), walletKey, lotusAP, lotusToken)
		if err != nil {
			panic(err)
		}
		// Enroll the gateway in the Register srv.
		if err := rm.RegisterGateway(gatewayRegistrar); err != nil {
			logging.Error("gateway registering error: %s", err.Error())
			t.FailNow()
		}

		// add to client
		added := client.AddGatewaysToUse([]*nodeid.NodeID{gatewayID})
		if !assert.Equal(t, 1, added, "1 gateway should be added") {
			t.FailNow()
		}
		activatedCount := client.AddActiveGateways(gatewayClientApiEndpoint, []*nodeid.NodeID{gatewayID})
		if !assert.Equal(t, 1, activatedCount, "1 gateway should be active") {
			t.FailNow()
		}
	}

	// Initialise providers
	pConfBuilder := fcrprovideradmin.CreateSettings()
	pConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	pConfBuilder.SetRegisterURL("http://" + containers.Register.GetRegisterHostApiEndpoint())
	pConf := pConfBuilder.Build()
	pAdmin := fcrprovideradmin.NewFilecoinRetrievalProviderAdmin(*pConf)

	// map between provider ID and provider name
	var providers = make(map[string]*nodeid.NodeID)
	for i := 0; i < 3; i++ {
		walletKey := privateKeys[0]
		walletAddress := accountAddrs[0]
		privateKeys = privateKeys[1:]
		accountAddrs = accountAddrs[1:]

		providerRootPubKey, providerRetrievalPubKey, providerRetrievalPrivateKey, err := cr.GenerateKeys()
		providerID := nodeid.NewRandomNodeID()

		providerName := fmt.Sprintf("provider-%v", i)
		providers[providerName] = providerID
		// get endpoint for Docker host
		_, _, providerAdminApiEndpoint := containers.Providers[providerName].GetProviderHostApiEndpoints()

		providerRegistrar := register.NewProviderRegister(
			providerID.ToString(),
			walletAddress,
			providerRootPubKey,
			providerRetrievalPubKey,
			providerConfig.GetString("PROVIDER_REGION_CODE"),
			providerName+":"+providerConfig.GetString("BIND_GATEWAY_API"),
			providerName+":"+providerConfig.GetString("BIND_REST_API"),
			providerName+":"+providerConfig.GetString("BIND_ADMIN_API"),
		)
		// Initialise the provider using provider admin
		err = pAdmin.InitialiseProviderV2(providerAdminApiEndpoint, providerRegistrar, providerRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1), walletKey, lotusAP, lotusToken)
		if err != nil {
			panic(err)
		}
		// Enroll the provider in the Register srv.
		if err := rm.RegisterProvider(providerRegistrar); err != nil {
			logging.Error("error registering provider: %s", err.Error())
			t.FailNow()
		}
	}

	// Force providers and gateways to update
	for providerName, providerID := range providers {
		_, _, providerAdminApiEndpoint := containers.Providers[providerName].GetProviderHostApiEndpoints()
		err := pAdmin.ForceUpdate(providerAdminApiEndpoint, providerID)
		if err != nil {
			logging.Error("provider update error: %s", err)
			t.FailNow()
		}
	}
	for gatewayName, gatewayID := range gateways {
		_, _, _, gatewayAdminApiEndpoint := containers.Gateways[gatewayName].GetGatewayHostApiEndpoints()
		err := gwAdmin.ForceUpdate(gatewayAdminApiEndpoint, gatewayID)
		if err != nil {
			logging.Error("gateway update error: %s", err)
			t.FailNow()
		}
	}

	// Now List DHT Offers
	// for i := 0; i < 32; i++ {
	// 	err := gwAdmin.ListDHTOffer(gwIDs[i])
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// }

	t.Log("/*******************************************************/")
	t.Log("/*                  End TestForceUpdate                */")
	t.Log("/*******************************************************/")
}
