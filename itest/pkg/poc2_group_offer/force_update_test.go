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
	var rm = fcrregistermgr.NewFCRRegisterMgr(gatewayConfig.GetString("REGISTER_API_URL"), true, true, 10*time.Second)
	if err := rm.Start(); err != nil {
		logging.Error("error starting Register Manager: %s", err.Error())
		t.FailNow()
	}
	defer rm.ShutdownAndWait()

	// Configure gateway admin
	gConfBuilder := fcrgatewayadmin.CreateSettings()
	gConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	gConfBuilder.SetRegisterURL(gatewayConfig.GetString("REGISTER_API_URL"))
	gConf := gConfBuilder.Build()
	gwAdmin := fcrgatewayadmin.NewFilecoinRetrievalGatewayAdmin(*gConf)

	// Initialise gateways
	var gwIDs []*nodeid.NodeID
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
			panic(err)
		}
		// Enroll the gateway in the Register srv.
		if err := rm.RegisterGateway(gatewayRegistrar); err != nil {
			logging.Error("gateway registering error: %s", err.Error())
			t.FailNow()
		}
	}

	// Initialise providers
	pConfBuilder := fcrprovideradmin.CreateSettings()
	pConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	pConfBuilder.SetRegisterURL(providerConfig.GetString("REGISTER_API_URL"))
	pConf := pConfBuilder.Build()
	pAdmin := fcrprovideradmin.NewFilecoinRetrievalProviderAdmin(*pConf)
	var pIDs []*nodeid.NodeID
	for i := 0; i < 3; i++ {
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
		// Initialise the provider using provider admin
		err = pAdmin.InitialiseProviderV2(providerRegistrar, providerRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1), walletKey, lotusAP, lotusToken)
		if err != nil {
			panic(err)
		}
		// Enroll the provider in the Register srv.
		if err := rm.RegisterProvider(providerRegistrar); err != nil {
			logging.Error("error registering provider: %s", err.Error())
			t.FailNow()
		}
	}

	walletKey := privateKeys[0]
	// walletAddress := accountAddrs[0]
	privateKeys = privateKeys[1:]
	accountAddrs = accountAddrs[1:]

	// Configure client
	clientConfBuilder := fcrclient.CreateSettings()
	clientConfBuilder.SetEstablishmentTTL(101)
	clientConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	clientConfBuilder.SetRegisterURL(gatewayConfig.GetString("REGISTER_API_URL"))
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

	added := client.AddGatewaysToUse(gwIDs)
	if !assert.Equal(t, 32, added, "32 gateways should be added") {
		t.FailNow()
	}

	added = client.AddActiveGateways(gwIDs)
	if !assert.Equal(t, 32, added, "32 gateways should be active") {
		t.FailNow()
	}

	// Force providers and gateways to update
	for _, p := range pIDs {
		err := pAdmin.ForceUpdate(p)
		if err != nil {
			logging.Error("provider update error: %s", err)
			t.FailNow()
		}
	}
	for _, g := range gwIDs {
		err := gwAdmin.ForceUpdate(g)
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
