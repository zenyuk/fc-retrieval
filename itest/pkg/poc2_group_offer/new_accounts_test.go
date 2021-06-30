package poc2_group_offer

import (
	"testing"

	"github.com/ConsenSys/fc-retrieval/itest/pkg/util"
)

func TestNewAccounts(t *testing.T) {
	t.Log("/*******************************************************/")
	t.Log("/*                 Start TestNewAccounts               */")
	t.Log("/*******************************************************/")

	var err error
	privateKeys, accountAddrs, err := util.GenerateAccount(lotusAP, lotusToken, superAcct, 37)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("privateKeys: %#v", privateKeys)
	t.Logf("accountAddrs: %#v", accountAddrs)

	t.Log("/*******************************************************/")
	t.Log("/*                  End TestNewAccounts                */")
	t.Log("/*******************************************************/")
}
