package poc2_dht_offer

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ConsenSys/fc-retrieval/client/pkg/fcrclient"
	"github.com/ConsenSys/fc-retrieval/common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrregistermgr"
	"github.com/ConsenSys/fc-retrieval/common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval/common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval/common/pkg/register"
	"github.com/ConsenSys/fc-retrieval/gateway-admin/pkg/fcrgatewayadmin"
	cr "github.com/ConsenSys/fc-retrieval/itest/pkg/util/crypto-facade"
	fil "github.com/ConsenSys/fc-retrieval/itest/pkg/util/filecoin-facade"
	"github.com/ConsenSys/fc-retrieval/provider-admin/pkg/fcrprovideradmin"
)

// Here we have gw 6-21 storing 688... and gw25-31&gw0-8 storing 008

func TestPublishDHTOffer(t *testing.T) {
	t.Log("/*******************************************************/")
	t.Log("/*              Start TestPublishDHTOffer              */")
	t.Log("/*******************************************************/")

	ctx := context.Background()

	lotusToken, superAcct := fil.GetLotusToken()
	logging.Info("Lotus token is: %s", lotusToken)
	logging.Info("Super Acct is %s", superAcct)

	lotusDaemonEndpoint, _ := containers.Lotus.GetLostHostApiEndpoints()
	lotusAP := "http://" + lotusDaemonEndpoint + "/rpc/v0"

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

	registerHostPort := containers.Register.GetRegisterHostApiEndpoint()
	registerApiEndpoint := "http://" + registerHostPort

	// Create and start register manager
	var rm = fcrregistermgr.NewFCRRegisterMgr(registerApiEndpoint, true, true, 10*time.Second)
	if err := rm.Start(); err != nil {
		logging.Error("error starting Register Manager: %s", err.Error())
		t.FailNow()
	}
	defer rm.ShutdownAndWait()

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
		t.FailNow()
	}

	// Configure gateway admin
	gConfBuilder := fcrgatewayadmin.CreateSettings()
	gConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	gConfBuilder.SetRegisterURL(registerApiEndpoint)
	gConf := gConfBuilder.Build()
	gwAdmin := fcrgatewayadmin.NewFilecoinRetrievalGatewayAdmin(*gConf)

	// Initialise gateways

	// map between gateway ID and gateway name
	var gateways = make(map[string]*nodeid.NodeID)
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
			logging.Error("node ID generation error: %s", err.Error())
			t.FailNow()
		}
		fmt.Printf("+++ TestPublishDHTOffer, new gatewayID %+v %+v\n", gatewayID, *gatewayID)

		gatewayName := fmt.Sprintf("gateway-%v", i)
		gateways[gatewayName] = gatewayID
		_, _, gatewayClientApiEndpoint, gatewayAdminApiEndpoint := containers.Gateways[gatewayName].GetGatewayHostApiEndpoints()

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

		// add to client
		added := client.AddGatewaysToUse([]*nodeid.NodeID{gatewayID})
		if !assert.Equal(t, 1, added, "1 gateway should be added") {
			t.FailNow()
		}
		activatedCount := client.AddActiveGateways(gatewayClientApiEndpoint, []*nodeid.NodeID{gatewayID})
		if !assert.Equal(t, 1, activatedCount, "1 gateway should be active") {
			t.FailNow()
		}
	}

	// Initialise providers
	pConfBuilder := fcrprovideradmin.CreateSettings()
	pConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	pConfBuilder.SetRegisterURL(registerApiEndpoint)
	pConf := pConfBuilder.Build()
	pAdmin := fcrprovideradmin.NewFilecoinRetrievalProviderAdmin(*pConf)

	// map between provider ID and provider name
	var providers = make(map[string]*nodeid.NodeID)
	for i := 0; i < 3; i++ {
		walletKey := privateKeys[0]
		walletAddress := accountAddrs[0]
		privateKeys = privateKeys[1:]
		accountAddrs = accountAddrs[1:]

		providerRootPubKey, providerRetrievalPubKey, providerRetrievalPrivateKey, err := cr.GenerateKeys()
		providerID := nodeid.NewRandomNodeID()

		providerName := fmt.Sprintf("provider-%v", i)
		providers[providerName] = providerID
		// get endpoint for Docker host
		_, _, providerAdminApiEndpoint := containers.Providers[providerName].GetProviderHostApiEndpoints()

		providerRegistrar := register.NewProviderRegister(
			providerID.ToString(),
			walletAddress,
			providerRootPubKey,
			providerRetrievalPubKey,
			providerConfig.GetString("PROVIDER_REGION_CODE"),
			providerName+":"+providerConfig.GetString("BIND_GATEWAY_API"),
			providerName+":"+providerConfig.GetString("BIND_REST_API"),
			providerName+":"+providerConfig.GetString("BIND_ADMIN_API"),
		)
		// Initialise the provider using provider admin
		err = pAdmin.InitialiseProviderV2(providerAdminApiEndpoint, providerRegistrar, providerRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1), walletKey, lotusAP, lotusToken)
		if err != nil {
			panic(err)
		}
		// Enroll the provider in the Register srv.
		if err := rm.RegisterProvider(providerRegistrar); err != nil {
			logging.Error("error registering provider: %s", err.Error())
			t.FailNow()
		}
	}

	// Force providers and gateways to update
	for providerName, providerID := range providers {
		_, _, providerAdminApiEndpoint := containers.Providers[providerName].GetProviderHostApiEndpoints()
		err := pAdmin.ForceUpdate(providerAdminApiEndpoint, providerID)
		if err != nil {
			logging.Error("provider update error: %s", err)
			t.FailNow()
		}
	}
	for gatewayName, gatewayID := range gateways {
		_, _, _, gatewayAdminApiEndpoint := containers.Gateways[gatewayName].GetGatewayHostApiEndpoints()
		err := gwAdmin.ForceUpdate(gatewayAdminApiEndpoint, gatewayID)
		if err != nil {
			logging.Error("gateway update error: %s", err)
			t.FailNow()
		}
	}

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
	_, _, providerAdminApiEndpoint := containers.Providers["provider-0"].GetProviderHostApiEndpoints()
	err = pAdmin.PublishGroupCID(providerAdminApiEndpoint, providers["provider-0"], []cid.ContentID{*contentID01, *contentID02, *contentID03}, 42, expiryDate, 42)
	if err != nil {
		logging.Error("error to publish group offer for content IDs: %s, %s, %s; error: %s", contentID01.ToString(), contentID02.ToString(), contentID03.ToString(), err.Error())
		t.FailNow()
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
	_, _, providerAdminApiEndpoint = containers.Providers["provider-1"].GetProviderHostApiEndpoints()
	err = pAdmin.PublishGroupCID(providerAdminApiEndpoint, providers["provider-1"], []cid.ContentID{*contentID11, *contentID12, *contentID13}, 42, expiryDate, 42)
	if err != nil {
		logging.Error("error to publish group offer for content IDs: %s, %s, %s; error: %s", contentID11.ToString(), contentID12.ToString(), contentID13.ToString(), err.Error())
		t.FailNow()
	}

	// Publish DHT Offer from pvd0
	// It will be published to gateway 6 -  gateway 21
	contentID0, err := cid.NewContentIDFromHexString("6880000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		t.Fatal(err)
	}
	expiryDate = time.Now().Local().Add(time.Hour * time.Duration(24)).Unix()
	_, _, providerAdminApiEndpoint = containers.Providers["provider-0"].GetProviderHostApiEndpoints()
	err = pAdmin.PublishDHTCID(providerAdminApiEndpoint, providers["provider-0"], []cid.ContentID{*contentID0}, []uint64{42}, []int64{expiryDate}, []uint64{42})
	if err != nil {
		logging.Error("error to publish DHT offer for content ID: %s; error: %s", contentID0.ToString(), err.Error())
		t.FailNow()
	}

	// Publish DHT Offer from pvd1
	// It will be published to gateway 25 to gateway 31 and gateway 0 to gateway 8
	contentID1, err := cid.NewContentIDFromHexString("0080000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		t.Fatal(err)
	}
	expiryDate = time.Now().Local().Add(time.Hour * time.Duration(24)).Unix()
	_, _, providerAdminApiEndpoint = containers.Providers["provider-0"].GetProviderHostApiEndpoints()
	err = pAdmin.PublishDHTCID(providerAdminApiEndpoint, providers["provider-0"], []cid.ContentID{*contentID1}, []uint64{42}, []int64{expiryDate}, []uint64{42})
	if err != nil {
		logging.Error("error to publish DHT offer for content ID: %s; error: %s", contentID1.ToString(), err.Error())
		t.FailNow()
	}

	// Try Standard Discovery for contentID0
	_, _, gatewayClientApiEndpoint, _ := containers.Gateways["gateway-5"].GetGatewayHostApiEndpoints()
	offers, err := client.FindOffersStandardDiscoveryV2(gatewayClientApiEndpoint, contentID0, gateways["gateway-5"], 1)
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 0, len(offers), "Shouldn't find offer with cid 0 from gateway 5, outside of published ring.") {
		t.Fatal()
	}

	_, _, gatewayClientApiEndpoint, _ = containers.Gateways["gateway-6"].GetGatewayHostApiEndpoints()
	offers, err = client.FindOffersStandardDiscoveryV2(gatewayClientApiEndpoint, contentID0, gateways["gateway-6"], 1)
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 6, boundary of published ring.") {
		t.FailNow()
	}

	_, _, gatewayClientApiEndpoint, _ = containers.Gateways["gateway-12"].GetGatewayHostApiEndpoints()
	offers, err = client.FindOffersStandardDiscoveryV2(gatewayClientApiEndpoint, contentID0, gateways["gateway-12"], 1)
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 12, within published ring.") {
		t.FailNow()
	}

	_, _, gatewayClientApiEndpoint, _ = containers.Gateways["gateway-21"].GetGatewayHostApiEndpoints()
	offers, err = client.FindOffersStandardDiscoveryV2(gatewayClientApiEndpoint, contentID0, gateways["gateway-21"], 1)
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 21, boundary of published ring.") {
		t.FailNow()
	}

	_, _, gatewayClientApiEndpoint, _ = containers.Gateways["gateway-22"].GetGatewayHostApiEndpoints()
	offers, err = client.FindOffersStandardDiscoveryV2(gatewayClientApiEndpoint, contentID0, gateways["gateway-22"], 1)
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 0, len(offers), "Shouldn't find offer with cid 0 from gateway 22, outside of published ring.") {
		t.FailNow()
	}

	// Try Standard Discovery for contentID1
	_, _, gatewayClientApiEndpoint, _ = containers.Gateways["gateway-24"].GetGatewayHostApiEndpoints()
	offers, err = client.FindOffersStandardDiscoveryV2(gatewayClientApiEndpoint, contentID1, gateways["gateway-24"], 1)
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 0, len(offers), "Shouldn't find offer with cid 1 from gateway 24, outside of published ring.") {
		t.FailNow()
	}

	_, _, gatewayClientApiEndpoint, _ = containers.Gateways["gateway-25"].GetGatewayHostApiEndpoints()
	offers, err = client.FindOffersStandardDiscoveryV2(gatewayClientApiEndpoint, contentID1, gateways["gateway-25"], 1)
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 1 from gateway 25, boundary of published ring.") {
		t.FailNow()
	}

	_, _, gatewayClientApiEndpoint, _ = containers.Gateways["gateway-31"].GetGatewayHostApiEndpoints()
	offers, err = client.FindOffersStandardDiscoveryV2(gatewayClientApiEndpoint, contentID1, gateways["gateway-31"], 1)
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 31, within published ring.") {
		t.FailNow()
	}

	_, _, gatewayClientApiEndpoint, _ = containers.Gateways["gateway-0"].GetGatewayHostApiEndpoints()
	offers, err = client.FindOffersStandardDiscoveryV2(gatewayClientApiEndpoint, contentID1, gateways["gateway-0"], 1)
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 0, within published ring.") {
		t.FailNow()
	}

	_, _, gatewayClientApiEndpoint, _ = containers.Gateways["gateway-8"].GetGatewayHostApiEndpoints()
	offers, err = client.FindOffersStandardDiscoveryV2(gatewayClientApiEndpoint, contentID1, gateways["gateway-8"], 1)
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 8, boundary of published ring.") {
		t.FailNow()
	}

	_, _, gatewayClientApiEndpoint, _ = containers.Gateways["gateway-9"].GetGatewayHostApiEndpoints()
	offers, err = client.FindOffersStandardDiscoveryV2(gatewayClientApiEndpoint, contentID1, gateways["gateway-9"], 1)
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 0, len(offers), "Shouldn't find offer with cid 0 from gateway 9, outside of published ring.") {
		t.FailNow()
	}

	// Try DHT Search for content 0
	_, _, gatewayClientApiEndpoint, _ = containers.Gateways["gateway-0"].GetGatewayHostApiEndpoints()
	offers, err = client.FindOffersStandardDiscoveryV2(gatewayClientApiEndpoint, contentID0, gateways["gateway-0"], 1)
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 0, len(offers), "Shouldn't find offer with cid 0 from gateway 0, outside of published ring.") {
		t.FailNow()
	}

	_, _, gatewayClientApiEndpoint, _ = containers.Gateways["gateway-0"].GetGatewayHostApiEndpoints()
	offersMap, err := client.FindOffersDHTDiscoveryV2(gatewayClientApiEndpoint, contentID0, gateways["gateway-0"], 4, 4)
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 4, len(offersMap), "Should have a map of 4 entries.") {
		t.FailNow()
	}

	// It should contact gateway 12, 13, 14, 15
	_, exists := offersMap[gateways["gateway-12"].ToString()]
	if !assert.True(t, exists, "Should query gateway 12.") {
		t.FailNow()
	}
	offers = *(offersMap[gateways["gateway-12"].ToString()])
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 12.") {
		t.FailNow()
	}

	_, exists = offersMap[gateways["gateway-13"].ToString()]
	if !assert.True(t, exists, "Should query gateway 13.") {
		t.FailNow()
	}
	offers = *(offersMap[gateways["gateway-13"].ToString()])
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 13.") {
		t.FailNow()
	}

	_, exists = offersMap[gateways["gateway-14"].ToString()]
	if !assert.True(t, exists, "Should query gateway 14.") {
		t.FailNow()
	}
	offers = *(offersMap[gateways["gateway-14"].ToString()])
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 14.") {
		t.FailNow()
	}

	_, exists = offersMap[gateways["gateway-15"].ToString()]
	if !assert.True(t, exists, "Should query gateway 15.") {
		t.FailNow()
	}
	offers = *(offersMap[gateways["gateway-15"].ToString()])
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 15.") {
		t.FailNow()
	}

	// Try DHT Search for content 1
	_, _, gatewayClientApiEndpoint, _ = containers.Gateways["gateway-15"].GetGatewayHostApiEndpoints()
	offers, err = client.FindOffersStandardDiscoveryV2(gatewayClientApiEndpoint, contentID1, gateways["gateway-15"], 1)
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 0, len(offers), "Shouldn't find offer with cid 0 from gateway 15, outside of published ring.") {
		t.FailNow()
	}

	_, _, gatewayClientApiEndpoint, _ = containers.Gateways["gateway-15"].GetGatewayHostApiEndpoints()
	offersMap, err = client.FindOffersDHTDiscoveryV2(gatewayClientApiEndpoint, contentID1, gateways["gateway-15"], 3, 3)
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 3, len(offersMap), "Should have a map of 3 entries.") {
		t.FailNow()
	}

	// It should contact gateway 0, 1 and 31
	_, exists = offersMap[gateways["gateway-0"].ToString()]
	if !assert.True(t, exists, "Should query gateway 0.") {
		t.FailNow()
	}
	offers = *(offersMap[gateways["gateway-0"].ToString()])
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 0.") {
		t.FailNow()
	}

	_, exists = offersMap[gateways["gateway-1"].ToString()]
	if !assert.True(t, exists, "Should query gateway 1.") {
		t.FailNow()
	}
	offers = *(offersMap[gateways["gateway-1"].ToString()])
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 1.") {
		t.FailNow()
	}

	_, exists = offersMap[gateways["gateway-31"].ToString()]
	if !assert.True(t, exists, "Should query gateway 31.") {
		t.FailNow()
	}
	offers = *(offersMap[gateways["gateway-31"].ToString()])
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 31.") {
		t.FailNow()
	}

	t.Log("/*******************************************************/")
	t.Log("/*               End TestPublishDHTOffer               */")
	t.Log("/*******************************************************/")
}
