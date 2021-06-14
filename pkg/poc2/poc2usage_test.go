/*
Package poc2 - set of end-to-end tests, designed to demonstrate functionality required for Proof of Concept stage 2.
*/
package poc2

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	tc "github.com/wcgcyx/testcontainers-go"

	"github.com/ConsenSys/fc-retrieval-client/pkg/fcrclient"
	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/register"
	"github.com/ConsenSys/fc-retrieval-gateway-admin/pkg/fcrgatewayadmin"
	"github.com/ConsenSys/fc-retrieval-itest/config"
	"github.com/ConsenSys/fc-retrieval-itest/pkg/util"
	"github.com/ConsenSys/fc-retrieval-provider-admin/pkg/fcrprovideradmin"
)

var gatewayConfig = config.NewConfig(".env.gateway")
var providerConfig = config.NewConfig(".env.provider")
var gwAdmin *fcrgatewayadmin.FilecoinRetrievalGatewayAdmin
var pAdmin *fcrprovideradmin.FilecoinRetrievalProviderAdmin
var client *fcrclient.FilecoinRetrievalClient
var gwIDs []*nodeid.NodeID
var pIDs []*nodeid.NodeID

func TestMain(m *testing.M) {
	// Need to make sure this env is not set in host machine
	itestEnv := os.Getenv("ITEST_CALLING_FROM_CONTAINER")

	if itestEnv != "" {
		// Env is set, we are calling from docker container
		// This logging should be only called after all tests finished.
		m.Run()
		return
	}
	// Env is not set, we are calling from host
	// We need a redis, a register, 17 gateways and 3 providers
	tag := util.GetCurrentBranch()

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
	registerContainer := util.StartRegister(ctx, tag, networkName, util.ColorYellow, rgEnv, true)

	// Start 3 providers
	var providerContainers []*tc.Container
	for i := 0; i < 3; i++ {
		c := util.StartProvider(ctx, fmt.Sprintf("provider-%v", i), tag, networkName, util.ColorBlue, pvEnv, true)
		providerContainers = append(providerContainers, &c)
	}

	// Start 33 gateways
	var gatewayContainers []*tc.Container
	for i := 0; i < 33; i++ {
		c := util.StartGateway(ctx, fmt.Sprintf("gateway-%v", i), tag, networkName, util.ColorCyan, gwEnv, true)
		gatewayContainers = append(gatewayContainers, &c)
	}

	// Start itest
	done := make(chan bool)
	itestContainer := util.StartItest(ctx, tag, networkName, util.ColorGreen, "", "", done, true, "")
	if _, err := itestContainer.Exec(ctx, []string{"export"}); err != nil {
		logging.Error("can't execute 'export' command in test container")
	}

	// Block until done.
	if <-done {
		logging.Info("Tests passed, shutdown...")
	} else {
		logging.Error("Tests failed, shutdown...")
	}

	if err := itestContainer.Terminate(ctx); err != nil {
		logging.Error("error while terminating test container: %s", err.Error())
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
	if err :=  registerContainer.Terminate(ctx); err != nil {
		logging.Error("error while terminating test container: %s", err.Error())
	}
	if err :=  redisContainer.Terminate(ctx); err != nil {
		logging.Error("error while terminating test container: %s", err.Error())
	}
	if err :=  (*network).Remove(ctx); err != nil {
		logging.Error("error while terminating test container network: %s", err.Error())
	}
}

func TestInitialiseProviders(t *testing.T) {
	t.Log("/*******************************************************/")
	t.Log("/*             Start TestInitialiseProviders           */")
	t.Log("/*******************************************************/")

	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		panic(err)
	}

	confBuilder := fcrprovideradmin.CreateSettings()
	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	confBuilder.SetRegisterURL(providerConfig.GetString("REGISTER_API_URL"))
	conf := confBuilder.Build()
	pAdmin = fcrprovideradmin.NewFilecoinRetrievalProviderAdmin(*conf)

	for i := 0; i < 3; i++ {
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
		providerRegister := &register.ProviderRegister{
			NodeID:             providerID.ToString(),
			Address:            providerConfig.GetString("PROVIDER_ADDRESS"),
			RootSigningKey:     providerRootSigningKey,
			SigningKey:         providerRetrievalSigningKey,
			RegionCode:         providerConfig.GetString("PROVIDER_REGION_CODE"),
			NetworkInfoGateway: providerConfig.GetString("NETWORK_INFO_GATEWAY")[:8] + identifier + providerConfig.GetString("NETWORK_INFO_GATEWAY")[8:],
			NetworkInfoClient:  providerConfig.GetString("NETWORK_INFO_CLIENT")[:8] + identifier + providerConfig.GetString("NETWORK_INFO_CLIENT")[8:],
			NetworkInfoAdmin:   providerConfig.GetString("NETWORK_INFO_ADMIN")[:8] + identifier + providerConfig.GetString("NETWORK_INFO_ADMIN")[8:],
		}

		err = pAdmin.InitialiseProvider(providerRegister, providerRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1))
		if err != nil {
			panic(err)
		}
	}

	t.Log("/*******************************************************/")
	t.Log("/*               End TestInitialiseProviders           */")
	t.Log("/*******************************************************/")
}

func TestInitialiseGateways(t *testing.T) {
	t.Log("/*******************************************************/")
	t.Log("/*             Start TestInitialiseGateways            */")
	t.Log("/*******************************************************/")

	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		panic(err)
	}

	confBuilder := fcrgatewayadmin.CreateSettings()
	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	confBuilder.SetRegisterURL(gatewayConfig.GetString("REGISTER_API_URL"))
	conf := confBuilder.Build()
	gwAdmin = fcrgatewayadmin.NewFilecoinRetrievalGatewayAdmin(*conf)

	// Only initialise 32 gateways, with one extra to initialise later to test list single cid offer
	for i := 0; i < 32; i++ {
		gatewayRootKey, err := fcrcrypto.GenerateBlockchainKeyPair()
		if err != nil {
			panic(err)
		}
		gatewayRootSigningKey, err := gatewayRootKey.EncodePublicKey()
		if err != nil {
			panic(err)
		}
		gatewayRetrievalPrivateKey, err := fcrcrypto.GenerateRetrievalV1KeyPair()
		if err != nil {
			panic(err)
		}
		gatewayRetrievalSigningKey, err := gatewayRetrievalPrivateKey.EncodePublicKey()
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
		gatewayRegister := &register.GatewayRegister{
			NodeID:              gatewayID.ToString(),
			Address:             gatewayConfig.GetString("GATEWAY_ADDRESS"),
			RootSigningKey:      gatewayRootSigningKey,
			SigningKey:          gatewayRetrievalSigningKey,
			RegionCode:          gatewayConfig.GetString("GATEWAY_REGION_CODE"),
			NetworkInfoGateway:  gatewayConfig.GetString("NETWORK_INFO_GATEWAY")[:7] + identifier + gatewayConfig.GetString("NETWORK_INFO_GATEWAY")[7:],
			NetworkInfoProvider: gatewayConfig.GetString("NETWORK_INFO_PROVIDER")[:7] + identifier + gatewayConfig.GetString("NETWORK_INFO_PROVIDER")[7:],
			NetworkInfoClient:   gatewayConfig.GetString("NETWORK_INFO_CLIENT")[:7] + identifier + gatewayConfig.GetString("NETWORK_INFO_CLIENT")[7:],
			NetworkInfoAdmin:    gatewayConfig.GetString("NETWORK_INFO_ADMIN")[:7] + identifier + gatewayConfig.GetString("NETWORK_INFO_ADMIN")[7:],
		}

		err = gwAdmin.InitialiseGateway(gatewayRegister, gatewayRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1))
		if err != nil {
			panic(err)
		}
	}

	t.Log("/*******************************************************/")
	t.Log("/*               End TestInitialiseGateway            */")
	t.Log("/*******************************************************/")
}

func TestInitialiseClient(t *testing.T) {
	t.Log("/*******************************************************/")
	t.Log("/*             Start TestInitialiseClient              */")
	t.Log("/*******************************************************/")

	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		panic(err)
	}

	confBuilder := fcrclient.CreateSettings()
	confBuilder.SetEstablishmentTTL(101)
	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	confBuilder.SetRegisterURL(gatewayConfig.GetString("REGISTER_API_URL"))
	conf := confBuilder.Build()
	client, err = fcrclient.NewFilecoinRetrievalClient(*conf)
	if !assert.Nil(t, err, "Error should be nil") {
		t.Fatal(err)
	}

	added := client.AddGatewaysToUse(gwIDs)
	if !assert.Equal(t, 32, added, "32 gateways should be added") {
		t.Fatal()
	}

	added = client.AddActiveGateways(gwIDs)
	assert.Equal(t, 32, added, "32 gateways should be active")

	t.Log("/*******************************************************/")
	t.Log("/*               End TestInitialiseClient              */")
	t.Log("/*******************************************************/")
}

func TestForceUpdate(t *testing.T) {
	t.Log("/*******************************************************/")
	t.Log("/*                Start TestForceUpdate                */")
	t.Log("/*******************************************************/")

	for i := 0; i < 3; i++ {
		err := pAdmin.ForceUpdate(pIDs[i])
		if err != nil {
			t.Fatal(err)
		}
	}

	for i := 0; i < 32; i++ {
		err := gwAdmin.ForceUpdate(gwIDs[i])
		if err != nil {
			t.Fatal(err)
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

// TODO: Add tests to configure the 32 gateways that do not ask for group offer
// func TestSubscriptionTurnOff()
// func TestPublishGroupOfferFail() -> try publish offer, use client to search, return no result
// func TestSubscriptionTurnOn() -> turn on gateway 0 from pvd0, gateway 1 from pvd1.
// need to modify the following, I've commented out part of the code that needs to be updated

func TestPublishGroupOffer(t *testing.T) {
	t.Log("/*******************************************************/")
	t.Log("/*             Start TestPublishGroupOffer             */")
	t.Log("/*******************************************************/")

	// Publish offer 0 from provider 0
	contentID01, err := cid.NewContentIDFromHexString("101112131415161718191A1B1C1D1E1F202122232425262728292A2B2C2D2E21")
	if err != nil {
		t.Fatal(err)
	}
	contentID02, err := cid.NewContentIDFromHexString("101112131415161718191A1B1C1D1E1F202122232425262728292A2B2C2D2E22")
	if err != nil {
		t.Fatal(err)
	}
	contentID03, err := cid.NewContentIDFromHexString("101112131415161718191A1B1C1D1E1F202122232425262728292A2B2C2D2E23")
	if err != nil {
		t.Fatal(err)
	}
	expiryDate := time.Now().Local().Add(time.Hour * time.Duration(24)).Unix()
	err = pAdmin.PublishGroupCID(pIDs[0], []cid.ContentID{*contentID01, *contentID02, *contentID03}, 42, expiryDate, 42)
	if err != nil {
		t.Fatal(err)
	}

	// Publish offer 1 from provider 1
	contentID11, err := cid.NewContentIDFromHexString("101112131415161718191A1B1C1D1E1F202122232425262728292A2B2C2D2E24")
	if err != nil {
		t.Fatal(err)
	}
	contentID12, err := cid.NewContentIDFromHexString("101112131415161718191A1B1C1D1E1F202122232425262728292A2B2C2D2E25")
	if err != nil {
		t.Fatal(err)
	}
	contentID13, err := cid.NewContentIDFromHexString("101112131415161718191A1B1C1D1E1F202122232425262728292A2B2C2D2E26")
	if err != nil {
		t.Fatal(err)
	}
	expiryDate = time.Now().Local().Add(time.Hour * time.Duration(24)).Unix()
	err = pAdmin.PublishGroupCID(pIDs[1], []cid.ContentID{*contentID11, *contentID12, *contentID13}, 42, expiryDate, 42)
	if err != nil {
		t.Fatal(err)
	}

	// Query gateway 0 for offer 0
	offers, err := client.FindOffersStandardDiscovery(contentID01, gwIDs[0])
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 01.") {
		t.Fatal()
	}

	// Query gateway 0 for offer 1
	offers, err = client.FindOffersStandardDiscovery(contentID11, gwIDs[0])
	if err != nil {
		t.Fatal(err)
	}
	// if !assert.Equal(t, 0, len(offers), "Should not find offer with cid 11.") {
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 11.") {
		t.Fatal()
	}

	// Query gateway 1 for offer 0
	offers, err = client.FindOffersStandardDiscovery(contentID01, gwIDs[1])
	if err != nil {
		t.Fatal(err)
	}
	// if !assert.Equal(t, 0, len(offers), "Should not find offer with cid 01.") {
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 01.") {
		t.Fatal()
	}

	// Query gateway 1 for offer 1
	offers, err = client.FindOffersStandardDiscovery(contentID11, gwIDs[1])
	if err != nil {
		panic(err)
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 11.") {
		t.Fatal()
	}

	// Query gateway 2 for offer 0
	offers, err = client.FindOffersStandardDiscovery(contentID01, gwIDs[2])
	if err != nil {
		panic(err)
	}
	// if !assert.Equal(t, 0, len(offers), "Should not find offer with cid 01.") {
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 01.") {
		t.Fatal()
	}

	// Query gateway 2 for offer 1
	offers, err = client.FindOffersStandardDiscovery(contentID11, gwIDs[2])
	if err != nil {
		t.Fatal(err)
	}
	// if !assert.Equal(t, 0, len(offers), "Should not find offer with cid 11.") {
	assert.Equal(t, 1, len(offers), "Should find offer with cid 11.")

	t.Log("/*******************************************************/")
	t.Log("/*              End TestPublishGroupOffer              */")
	t.Log("/*******************************************************/")
}

func TestPublishDHTOffer(t *testing.T) {
	t.Log("/*******************************************************/")
	t.Log("/*              Start TestPublishDHTOffer              */")
	t.Log("/*******************************************************/")

	// Publish DHT Offer from pvd0
	// It will be published to gateway 6 -  gateway 21
	contentID0, err := cid.NewContentIDFromHexString("6880000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		t.Fatal(err)
	}
	expiryDate := time.Now().Local().Add(time.Hour * time.Duration(24)).Unix()
	err = pAdmin.PublishDHTCID(pIDs[0], []cid.ContentID{*contentID0}, []uint64{42}, []int64{expiryDate}, []uint64{42})
	if err != nil {
		t.Fatal(err)
	}

	// Publish DHT Offer from pvd1
	// It will be published to gateway 25 to gateway 31 and gateway 0 to gateway 8
	contentID1, err := cid.NewContentIDFromHexString("0080000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		t.Fatal(err)
	}
	expiryDate = time.Now().Local().Add(time.Hour * time.Duration(24)).Unix()
	err = pAdmin.PublishDHTCID(pIDs[0], []cid.ContentID{*contentID1}, []uint64{42}, []int64{expiryDate}, []uint64{42})
	if err != nil {
		t.Fatal(err)
	}

	// Try Standard Discovery for contentID0
	offers, err := client.FindOffersStandardDiscovery(contentID0, gwIDs[5])
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 0, len(offers), "Shouldn't find offer with cid 0 from gateway 5, outside of published ring.") {
		t.Fatal()
	}

	offers, err = client.FindOffersStandardDiscovery(contentID0, gwIDs[6])
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 6, boundary of published ring.") {
		t.Fatal()
	}

	offers, err = client.FindOffersStandardDiscovery(contentID0, gwIDs[12])
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 12, within published ring.") {
		t.Fatal()
	}

	offers, err = client.FindOffersStandardDiscovery(contentID0, gwIDs[21])
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 21, boundary of published ring.") {
		t.Fatal()
	}

	offers, err = client.FindOffersStandardDiscovery(contentID0, gwIDs[22])
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 0, len(offers), "Shouldn't find offer with cid 0 from gateway 22, outside of published ring.") {
		t.Fatal()
	}

	// Try Standard Discovery for contentID1
	offers, err = client.FindOffersStandardDiscovery(contentID1, gwIDs[24])
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 0, len(offers), "Shouldn't find offer with cid 1 from gateway 24, outside of published ring.") {
		t.Fatal()
	}

	offers, err = client.FindOffersStandardDiscovery(contentID1, gwIDs[25])
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 1 from gateway 25, boundary of published ring.") {
		t.Fatal()
	}

	offers, err = client.FindOffersStandardDiscovery(contentID1, gwIDs[31])
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 31, within published ring.") {
		t.Fatal()
	}

	offers, err = client.FindOffersStandardDiscovery(contentID1, gwIDs[0])
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 0, within published ring.") {
		t.Fatal()
	}

	offers, err = client.FindOffersStandardDiscovery(contentID1, gwIDs[8])
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 8, boundary of published ring.") {
		t.Fatal()
	}

	offers, err = client.FindOffersStandardDiscovery(contentID1, gwIDs[9])
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 0, len(offers), "Shouldn't find offer with cid 0 from gateway 9, outside of published ring.") {
		t.Fatal()
	}

	// Try DHT Search for content 0
	offers, err = client.FindOffersStandardDiscovery(contentID0, gwIDs[0])
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 0, len(offers), "Shouldn't find offer with cid 0 from gateway 0, outside of published ring.") {
		t.Fatal()
	}

	offersMap, err := client.FindOffersDHTDiscovery(contentID0, gwIDs[0], 4)
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 4, len(offersMap), "Should have a map of 4 entries.") {
		t.Fatal()
	}

	// It should contact gateway 12, 13, 14, 15
	_, exists := offersMap[gwIDs[12].ToString()]
	if !assert.True(t, exists, "Should query gateway 12.") {
		t.Fatal()
	}
	offers = *(offersMap[gwIDs[12].ToString()])
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 12.") {
		t.Fatal()
	}

	_, exists = offersMap[gwIDs[13].ToString()]
	if !assert.True(t, exists, "Should query gateway 13.") {
		t.Fatal()
	}
	offers = *(offersMap[gwIDs[13].ToString()])
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 13.") {
		t.Fatal()
	}

	_, exists = offersMap[gwIDs[14].ToString()]
	if !assert.True(t, exists, "Should query gateway 14.") {
		t.Fatal()
	}
	offers = *(offersMap[gwIDs[14].ToString()])
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 14.") {
		t.Fatal()
	}

	_, exists = offersMap[gwIDs[15].ToString()]
	if !assert.True(t, exists, "Should query gateway 15.") {
		t.Fatal()
	}
	offers = *(offersMap[gwIDs[15].ToString()])
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 15.") {
		t.Fatal()
	}

	// Try DHT Search for content 1
	offers, err = client.FindOffersStandardDiscovery(contentID1, gwIDs[15])
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 0, len(offers), "Shouldn't find offer with cid 0 from gateway 15, outside of published ring.") {
		t.Fatal()
	}

	offersMap, err = client.FindOffersDHTDiscovery(contentID1, gwIDs[15], 3)
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 3, len(offersMap), "Should have a map of 3 entries.") {
		t.Fatal()
	}

	// It should contact gateway 0, 1 and 31
	_, exists = offersMap[gwIDs[0].ToString()]
	if !assert.True(t, exists, "Should query gateway 0.") {
		t.Fatal()
	}
	offers = *(offersMap[gwIDs[0].ToString()])
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 0.") {
		t.Fatal()
	}

	_, exists = offersMap[gwIDs[1].ToString()]
	if !assert.True(t, exists, "Should query gateway 1.") {
		t.Fatal()
	}
	offers = *(offersMap[gwIDs[1].ToString()])
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 1.") {
		t.Fatal()
	}

	_, exists = offersMap[gwIDs[31].ToString()]
	if !assert.True(t, exists, "Should query gateway 31.") {
		t.Fatal()
	}
	offers = *(offersMap[gwIDs[31].ToString()])
	assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 31.")

	t.Log("/*******************************************************/")
	t.Log("/*               End TestPublishDHTOffer               */")
	t.Log("/*******************************************************/")
}

// At this point, we have gw 6-21 storing 688... and gw25-31&gw0-8 storing 008

// TODO: Add test
// 1. Turn on subscription gateway 0 from pvd2.
// 2. Publish group cid offer contains cid 8080000000000000000000000000000000000000000000000000000000000000 and 8080000000000000000000000000000000000000000000000000000000000001
// Try to do dht search from gateway2 for 8080000000000000000000000000000000000000000000000000000000000000 with dht num 4
// 3. It should query gateway 15, 16, 17, 18, but no one gets the offer
// 4. Try to do std search for gateway 0, it has offer. All good.

func TestNewGateway(t *testing.T) {
	t.Log("/*******************************************************/")
	t.Log("/*                 Start TestNewGateway                */")
	t.Log("/*******************************************************/")

	gatewayRootKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		panic(err)
	}
	gatewayRootSigningKey, err := gatewayRootKey.EncodePublicKey()
	if err != nil {
		panic(err)
	}
	gatewayRetrievalPrivateKey, err := fcrcrypto.GenerateRetrievalV1KeyPair()
	if err != nil {
		panic(err)
	}
	gatewayRetrievalSigningKey, err := gatewayRetrievalPrivateKey.EncodePublicKey()
	if err != nil {
		panic(err)
	}
	idStr := "3880000000000000000000000000000000000000000000000000000000000000"

	gatewayID, err := nodeid.NewNodeIDFromHexString(idStr)
	if err != nil {
		panic(err)
	}
	gwIDs = append(gwIDs, gatewayID)

	identifier := fmt.Sprintf("-32")
	gatewayRegister := &register.GatewayRegister{
		NodeID:              gatewayID.ToString(),
		Address:             gatewayConfig.GetString("GATEWAY_ADDRESS"),
		RootSigningKey:      gatewayRootSigningKey,
		SigningKey:          gatewayRetrievalSigningKey,
		RegionCode:          gatewayConfig.GetString("GATEWAY_REGION_CODE"),
		NetworkInfoGateway:  gatewayConfig.GetString("NETWORK_INFO_GATEWAY")[:7] + identifier + gatewayConfig.GetString("NETWORK_INFO_GATEWAY")[7:],
		NetworkInfoProvider: gatewayConfig.GetString("NETWORK_INFO_PROVIDER")[:7] + identifier + gatewayConfig.GetString("NETWORK_INFO_PROVIDER")[7:],
		NetworkInfoClient:   gatewayConfig.GetString("NETWORK_INFO_CLIENT")[:7] + identifier + gatewayConfig.GetString("NETWORK_INFO_CLIENT")[7:],
		NetworkInfoAdmin:    gatewayConfig.GetString("NETWORK_INFO_ADMIN")[:7] + identifier + gatewayConfig.GetString("NETWORK_INFO_ADMIN")[7:],
	}

	err = gwAdmin.InitialiseGateway(gatewayRegister, gatewayRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1))
	if err != nil {
		panic(err)
	}

	// Force update
	for i := 0; i < 3; i++ {
		err := pAdmin.ForceUpdate(pIDs[i])
		if err != nil {
			t.Fatal(err)
		}
	}

	for i := 0; i < 33; i++ {
		err := gwAdmin.ForceUpdate(gwIDs[i])
		if err != nil {
			t.Fatal(err)
		}
	}

	err = gwAdmin.ListDHTOffer(gwIDs[32])
	if err != nil {
		t.Fatal(err)
	}

	added := client.AddGatewaysToUse([]*nodeid.NodeID{gwIDs[32]})
	if !assert.Equal(t, 1, added, "1 gateway should be added") {
		t.Fatal()
	}

	added = client.AddActiveGateways([]*nodeid.NodeID{gwIDs[32]})
	if !assert.Equal(t, 1, added, "1 gateway should be active") {
		t.Fatal()
	}

	// This new gateway should have used list cid offer to get both cid0 and cid1
	contentID0, err := cid.NewContentIDFromHexString("6880000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		t.Fatal(err)
	}
	contentID1, err := cid.NewContentIDFromHexString("0080000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		t.Fatal(err)
	}

	offers, err := client.FindOffersStandardDiscovery(contentID0, gwIDs[32])
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 32.") {
		t.Fatal(err)
	}

	offers, err = client.FindOffersStandardDiscovery(contentID1, gwIDs[32])
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(offers), "Should find offer with cid 1 from gateway 32.")

	t.Log("/*******************************************************/")
	t.Log("/*                  End TestNewGateway                 */")
	t.Log("/*******************************************************/")
}

// At this point, we have gw 6-21 storing 688... and gw25-31&gw0-8 storing 008
// Plus the gw33 storing both 688 and 008

func TestPublishDHTOfferWithNewGateway(t *testing.T) {
	t.Log("/*******************************************************/")
	t.Log("/*       Start TestPublishDHTOfferWithNewGateway       */")
	t.Log("/*******************************************************/")

	// First, force all gws and pvds to refresh
	for i := 0; i < 3; i++ {
		err := pAdmin.ForceUpdate(pIDs[i])
		if err != nil {
			t.Fatal(err)
		}
	}

	for i := 0; i < 33; i++ {
		err := gwAdmin.ForceUpdate(gwIDs[i])
		if err != nil {
			t.Fatal(err)
		}
	}

	// Publish DHT Offer from pvd3
	// It used to be published to gateway 7 - gateway 22, but now, it should be published to gateway 7 - 21 and gateway 32
	contentID0, err := cid.NewContentIDFromHexString("7080000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		t.Fatal(err)
	}
	expiryDate := time.Now().Local().Add(time.Hour * time.Duration(24)).Unix()
	err = pAdmin.PublishDHTCID(pIDs[2], []cid.ContentID{*contentID0}, []uint64{42}, []int64{expiryDate}, []uint64{42})
	if err != nil {
		t.Fatal(err)
	}

	// Try Standard Discovery for contentID0
	offers, err := client.FindOffersStandardDiscovery(contentID0, gwIDs[6])
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 0, len(offers), "Shouldn't find offer with cid 0 from gateway 6, outside of published ring.") {
		t.Fatal()
	}

	offers, err = client.FindOffersStandardDiscovery(contentID0, gwIDs[32])
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 32, boundary of published ring.") {
		t.Fatal()
	}

	offers, err = client.FindOffersStandardDiscovery(contentID0, gwIDs[15])
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 15, within published ring.") {
		t.Fatal()
	}

	offers, err = client.FindOffersStandardDiscovery(contentID0, gwIDs[21])
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 21, boundary of published ring.") {
		t.Fatal()
	}

	offers, err = client.FindOffersStandardDiscovery(contentID0, gwIDs[22])
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 0, len(offers), "Shouldn't find offer with cid 0 from gateway 22, outside of published ring.")

	t.Log("/*******************************************************/")
	t.Log("/*        End TestPublishDHTOfferWithNewGateway        */")
	t.Log("/*******************************************************/")
}

// To this point, we have
// At this point, we have gw 6-21 storing 688... and gw25-31&gw0-8 storing 008
// Plus the gw33 storing both 688 and 008
// Also, gw7-21 and gw32 storing 708

func TestDHTOfferAck(t *testing.T) {
	t.Log("/*******************************************************/")
	t.Log("/*                Start TestDHTOfferAck                */")
	t.Log("/*******************************************************/")

	cidValid, err := cid.NewContentIDFromHexString("7080000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		t.Fatal(err)
	}

	cidInValid, err := cid.NewContentIDFromHexString("7080000000000000000000000000000000000000000000000000000000000001")
	if err != nil {
		t.Fatal(err)
	}

	// Test a valid cid, valid gateway pair
	exists, err := client.FindDHTOfferAck(cidValid, gwIDs[10], pIDs[2])
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, exists, "Offer should exist for gateway 10 and valid cid")

	// Test an invalid cid, valid gateway pair
	exists, err = client.FindDHTOfferAck(cidInValid, gwIDs[10], pIDs[2])
	if err != nil {
		t.Fatal(err)
	}
	assert.False(t, exists, "Offer shouldn't exist for gateway 10 and invalid cid")

	// Test a valid cid, invalid gateway pair
	exists, err = client.FindDHTOfferAck(cidValid, gwIDs[30], pIDs[2])
	if err != nil {
		t.Fatal(err)
	}
	assert.False(t, exists, "Offer shouldn't exist for gateway 30 and valid cid")

	// Test an invalid cid, invalid gateway pair
	exists, err = client.FindDHTOfferAck(cidInValid, gwIDs[30], pIDs[2])
	if err != nil {
		t.Fatal(err)
	}
	assert.False(t, exists, "Offer shouldn't exist for gateway 30 and invalid cid")

	t.Log("/*******************************************************/")
	t.Log("/*                 End TestDHTOfferAck                 */")
	t.Log("/*******************************************************/")
}
