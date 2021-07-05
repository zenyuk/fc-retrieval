package poc2_group_offer

import (
	"context"
	"testing"

	fil "github.com/ConsenSys/fc-retrieval/itest/pkg/util/filecoin-facade"
)

func TestNewAccounts(t *testing.T) {
	t.Log("/*******************************************************/")
	t.Log("/*                 Start TestNewAccounts               */")
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
	t.Logf("privateKeys: %#v", privateKeys)
	t.Logf("accountAddrs: %#v", accountAddrs)

	t.Log("/*******************************************************/")
	t.Log("/*                  End TestNewAccounts                */")
	t.Log("/*******************************************************/")
}
