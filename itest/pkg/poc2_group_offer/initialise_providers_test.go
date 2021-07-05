package poc2_group_offer

import (
	"context"
	"fmt"
	"testing"

	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval/common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval/common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval/common/pkg/register"
	fil "github.com/ConsenSys/fc-retrieval/itest/pkg/util/filecoin-facade"
	"github.com/ConsenSys/fc-retrieval/provider-admin/pkg/fcrprovideradmin"
)

func TestInitialiseProviders(t *testing.T) {
	t.Log("/*******************************************************/")
	t.Log("/*             Start TestInitialiseProviders           */")
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
	confBuilder := fcrprovideradmin.CreateSettings()
	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	confBuilder.SetRegisterURL(registerApiEndpoint)
	conf := confBuilder.Build()
	pAdmin := fcrprovideradmin.NewFilecoinRetrievalProviderAdmin(*conf)

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

		providerName := fmt.Sprintf("provider-%v", i)
		// get endpoint for Docker host
		_, _, providerAdminApiEndpoint := containers.Providers[providerName].GetProviderHostApiEndpoints()
		providerRegistrar := register.NewProviderRegister(
			providerID.ToString(),
			walletAddress,
			providerRootSigningKey,
			providerRetrievalSigningKey,
			providerConfig.GetString("PROVIDER_REGION_CODE"),
			providerName+":"+providerConfig.GetString("BIND_GATEWAY_API"),
			providerName+":"+providerConfig.GetString("BIND_REST_API"),
			providerName+":"+providerConfig.GetString("BIND_ADMIN_API"),
		)

		err = pAdmin.InitialiseProviderV2(providerAdminApiEndpoint, providerRegistrar, providerRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1), walletKey, lotusAP, lotusToken)
		if err != nil {
			logging.Error("error initialising provider: %s", err.Error())
			t.FailNow()
		}
	}

	t.Log("/*******************************************************/")
	t.Log("/*               End TestInitialiseProviders           */")
	t.Log("/*******************************************************/")
}
