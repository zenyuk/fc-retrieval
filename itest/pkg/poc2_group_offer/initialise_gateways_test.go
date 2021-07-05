package poc2_group_offer

import (
	"context"
	"fmt"
	"testing"

	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval/common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval/common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval/common/pkg/register"
	"github.com/ConsenSys/fc-retrieval/gateway-admin/pkg/fcrgatewayadmin"
	cr "github.com/ConsenSys/fc-retrieval/itest/pkg/util/crypto-facade"
	fil "github.com/ConsenSys/fc-retrieval/itest/pkg/util/filecoin-facade"
)

func TestInitialiseGateways(t *testing.T) {
	t.Log("/*******************************************************/")
	t.Log("/*             Start TestInitialiseGateways            */")
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

	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		panic(err)
	}

	registerApiEndpoint := "http://" + containers.Register.GetRegisterHostApiEndpoint()

	confBuilder := fcrgatewayadmin.CreateSettings()
	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	confBuilder.SetRegisterURL(registerApiEndpoint)
	conf := confBuilder.Build()
	gwAdmin := fcrgatewayadmin.NewFilecoinRetrievalGatewayAdmin(*conf)

	var gwIDs []*nodeid.NodeID

	// Only initialise 32 gateways, with one extra to initialise later to test list single cid offer
	for i := 0; i < 32; i++ {
		walletKey := privateKeys[0]
		walletAddress := accountAddrs[0]
		privateKeys = privateKeys[1:]
		accountAddrs = accountAddrs[1:]

		gatewayRootPubKey, gatewayRetrievalPubKey, gatewayRetrievalPrivateKey, err := cr.GenerateKeys()
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

		gatewayName := fmt.Sprintf("gateway-%v", i)
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

		_, _, _, gatewayAdminApiEndpoint := containers.Gateways[gatewayName].GetGatewayHostApiEndpoints()
		err = gwAdmin.InitialiseGatewayV2(gatewayAdminApiEndpoint, gatewayRegistrar, gatewayRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1), walletKey, lotusAP, lotusToken)
		if err != nil {
			logging.Error("gateway initialising error: %s", err.Error())
			t.FailNow()
		}
	}

	t.Log("/*******************************************************/")
	t.Log("/*               End TestInitialiseGateway            */")
	t.Log("/*******************************************************/")
}
