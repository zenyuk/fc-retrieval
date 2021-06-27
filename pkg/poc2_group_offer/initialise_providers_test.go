package poc2_group_offer

import (
	"fmt"
	"testing"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/register"
	"github.com/ConsenSys/fc-retrieval-itest/pkg/util"
	"github.com/ConsenSys/fc-retrieval-provider-admin/pkg/fcrprovideradmin"
)

func TestInitialiseProviders(t *testing.T) {
	t.Log("/*******************************************************/")
	t.Log("/*             Start TestInitialiseProviders           */")
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

	for i := 0; i < 3; i++ {

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
	t.Log("/*               End TestInitialiseProviders           */")
	t.Log("/*******************************************************/")
}
