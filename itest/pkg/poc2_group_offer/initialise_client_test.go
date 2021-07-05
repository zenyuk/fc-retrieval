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
)

func TestInitialiseClient(t *testing.T) {
	t.Log("/*******************************************************/")
	t.Log("/*             Start TestInitialiseClient              */")
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

	// Configure gateway admin
	gConfBuilder := fcrgatewayadmin.CreateSettings()
	gConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	gConfBuilder.SetRegisterURL(registerApiEndpoint)
	gConf := gConfBuilder.Build()
	gwAdmin := fcrgatewayadmin.NewFilecoinRetrievalGatewayAdmin(*gConf)

	// map between gateway ID and gateway name
	var gateways = make(map[string]*nodeid.NodeID)

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

		gatewayName := fmt.Sprintf("gateway-%v", i)
		gateways[gatewayName] = gatewayID
		gwIDs = append(gwIDs, gatewayID)
		_, _, _, gatewayAdminApiEndpoint := containers.Gateways[gatewayName].GetGatewayHostApiEndpoints()
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
	}

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
		t.Fatal()
	}

	added := client.AddGatewaysToUse(gwIDs)
	if !assert.Equal(t, 32, added, "32 gateways should be added") {
		t.Fatal()
	}

	for gatewayName, gatewayID := range gateways {
		_, _, gatewayClientApiEndpoint, _ := containers.Gateways[gatewayName].GetGatewayHostApiEndpoints()
		activatedCount := client.AddActiveGateways(gatewayClientApiEndpoint, []*nodeid.NodeID{gatewayID})
		if !assert.Equal(t, 1, activatedCount, "gateway should be active") {
			t.FailNow()
		}
	}

	t.Log("/*******************************************************/")
	t.Log("/*               End TestInitialiseClient              */")
	t.Log("/*******************************************************/")
}
