/*
Package poc2 - set of end-to-end tests, designed to demonstrate functionality required for Proof of Concept stage 2.
*/
package poc2_group_offer

import (
	"context"
	"os"
	"testing"

	"github.com/ConsenSys/fc-retrieval/common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval/itest/config"
	fil "github.com/ConsenSys/fc-retrieval/itest/pkg/util/filecoin-facade"
	tc "github.com/ConsenSys/fc-retrieval/itest/pkg/util/test-containers"
)

const lotusAP = "http://lotus-full-node:1234/rpc/v0"

var lotusToken string
var superAcct string
var gatewayConfig = config.NewConfig(".env.gateway")
var providerConfig = config.NewConfig(".env.provider")

func TestMain(m *testing.M) {
	lotusToken, superAcct = fil.GetLotusToken()
	logging.Info("Lotus token is: %s", lotusToken)
	logging.Info("Super Acct is %s", superAcct)

	const testName = "poc2-group-offer"
	ctx := context.Background()
	containersAndPorts, network, err := tc.StartContainers(ctx, 33, 3, testName, true)
	if err != nil {
		logging.Error("%s failed, container starting error: %s", testName, err.Error())
		tc.StopContainers(ctx, containersAndPorts, network)
		os.Exit(1)
	}
	defer tc.StopContainers(ctx, containersAndPorts, network)
	m.Run()
}
