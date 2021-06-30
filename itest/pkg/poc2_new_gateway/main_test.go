/*
Package poc2_new_gateway - set of end-to-end tests, designed to demonstrate functionality required for Proof of Concept stage 2.
*/
package poc2_new_gateway

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	tc "github.com/wcgcyx/testcontainers-go"

	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval/common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval/itest/config"
	"github.com/ConsenSys/fc-retrieval/itest/pkg/util"
)

const lotusAP = "http://lotus-full-node:1234/rpc/v0"

var lotusToken string
var superAcct string
var gatewayConfig = config.NewConfig(".env.gateway")
var providerConfig = config.NewConfig(".env.provider")

func TestMain(m *testing.M) {
	// Need to make sure this env is not set in host machine
	itestEnv := os.Getenv("ITEST_CALLING_FROM_CONTAINER")

	if itestEnv != "" {
		lotusToken = os.Getenv("LOTUS_TOKEN")
		superAcct = os.Getenv("SUPER_ACCT")
		// Env is set, we are calling from docker container
		// This logging should be only called after all tests finished.
		m.Run()
		return
	}
	// Env is not set, we are calling from host
	// We need a redis, a register, 17 gateways and 3 providers

	// Get env
	rgEnv := util.GetEnvMap("../../.env.register")
	gwEnv := util.GetEnvMap("../../.env.gateway")
	pvEnv := util.GetEnvMap("../../.env.provider")

	// Create shared net
	ctx := context.Background()
	network, networkName := util.CreateNetwork(ctx)

	// Start redis
	redisContainer := util.StartRedis(ctx, networkName, true)

	// Start register
	registerContainer := util.StartRegister(ctx, networkName, util.ColorYellow, rgEnv, true)

	// Start 3 providers
	var providerContainers []*tc.Container
	for i := 0; i < 3; i++ {
		c := util.StartProvider(ctx, fmt.Sprintf("provider-%v", i), networkName, util.ColorBlue, pvEnv, true)
		providerContainers = append(providerContainers, &c)
	}

	// Start 33 gateways
	var gatewayContainers []*tc.Container
	for i := 0; i < 33; i++ {
		c := util.StartGateway(ctx, fmt.Sprintf("gateway-%v", i), networkName, util.ColorCyan, gwEnv, true)
		gatewayContainers = append(gatewayContainers, &c)
	}

	// Start lotus
	lotusContainer := util.StartLotusFullNode(ctx, networkName, false)

	// Get lotus token
	lotusToken, superAcct = util.GetLotusToken()
	logging.Info("Lotus token is: %s", lotusToken)
	logging.Info("Super Acct is %s", superAcct)

	// Start itest
	done := make(chan bool)
	itestContainer := util.StartItest(ctx, networkName, util.ColorGreen, lotusToken, superAcct, done, true, "")
	// Block until done.
	if <-done {
		logging.Info("Tests passed, shutdown...")
	} else {
		logging.Error("Tests failed, shutdown...")
	}

	if err := itestContainer.Terminate(ctx); err != nil {
		logging.Error("error while terminating test container: %s", err.Error())
	}
	if err := lotusContainer.Terminate(ctx); err != nil {
		logging.Error("error while terminating Lotus test container: %s", err.Error())
	}
	for _, c := range gatewayContainers {
		if err := (*c).Terminate(ctx); err != nil {
			logging.Error("error while terminating gateway test container: %s", err.Error())
		}
	}
	for _, c := range providerContainers {
		if err := (*c).Terminate(ctx); err != nil {
			logging.Error("error while terminating provider test container: %s", err.Error())
		}
	}
	if err := registerContainer.Terminate(ctx); err != nil {
		logging.Error("error while terminating test container: %s", err.Error())
	}
	if err := redisContainer.Terminate(ctx); err != nil {
		logging.Error("error while terminating test container: %s", err.Error())
	}
	if err := (*network).Remove(ctx); err != nil {
		logging.Error("error while terminating test container network: %s", err.Error())
	}
}

// Helper function to generate set of keys
func generateKeys() (rootPubKey string, retrievalPubKey string, retrievalPrivateKey *fcrcrypto.KeyPair, err error) {
	rootKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		return "", "", nil, fmt.Errorf("error generating blockchain key: %s", err.Error())
	}
	if rootKey == nil {
		return "", "", nil, errors.New("error generating blockchain key")
	}

	rootPubKey, err = rootKey.EncodePublicKey()
	if err != nil {
		return "", "", nil, fmt.Errorf("error encoding public key: %s", err.Error())
	}

	retrievalPrivateKey, err = fcrcrypto.GenerateRetrievalV1KeyPair()
	if err != nil {
		return "", "", nil, fmt.Errorf("error generating retrieval key: %s", err.Error())
	}
	if retrievalPrivateKey == nil {
		return "", "", nil, errors.New("error generating retrieval key")
	}

	retrievalPubKey, err = retrievalPrivateKey.EncodePublicKey()
	if err != nil {
		return "", "", nil, fmt.Errorf("error encoding retrieval pub key: %s", err.Error())
	}
	return
}
