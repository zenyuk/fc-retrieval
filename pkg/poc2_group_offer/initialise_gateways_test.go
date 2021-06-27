package poc2_group_offer

import (
	"fmt"
	"testing"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/register"
	"github.com/ConsenSys/fc-retrieval-gateway-admin/pkg/fcrgatewayadmin"
	"github.com/ConsenSys/fc-retrieval-itest/pkg/util"
)

func TestInitialiseGateways(t *testing.T) {
	t.Log("/*******************************************************/")
	t.Log("/*             Start TestInitialiseGateways            */")
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

	// Only initialise 32 gateways, with one extra to initialise later to test list single cid offer
	for i := 0; i < 32; i++ {
		walletKey := privateKeys[0]
		walletAddress := accountAddrs[0]
		privateKeys = privateKeys[1:]
		accountAddrs = accountAddrs[1:]

		gatewayRootPubKey, gatewayRetrievalPubKey, gatewayRetrievalPrivateKey, err := generateKeys()
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
	t.Log("/*               End TestInitialiseGateway            */")
	t.Log("/*******************************************************/")
}
