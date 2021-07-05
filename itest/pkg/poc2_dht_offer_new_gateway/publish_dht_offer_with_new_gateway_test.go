package poc2_dht_offer_new_gateway

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
// Plus the gw33 storing both 688 and 008
// Also, gw7-21 and gw32 storing 708

func TestPublishDHTOfferWithNewGateway(t *testing.T) {
	t.Log("/*******************************************************/")
	t.Log("/*       Start TestPublishDHTOfferWithNewGateway       */")
	t.Log("/*******************************************************/")

	ctx := context.Background()
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

	// Create and start register manager
	var rm = fcrregistermgr.NewFCRRegisterMgr(gatewayConfig.GetString("REGISTER_API_URL"), true, true, 10*time.Second)
	if err := rm.Start(); err != nil {
		logging.Error("error starting Register Manager: %s", err.Error())
		t.FailNow()
	}
	defer rm.ShutdownAndWait()

	// Configure gateway admin
	gConfBuilder := fcrgatewayadmin.CreateSettings()
	gConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	gConfBuilder.SetRegisterURL(gatewayConfig.GetString("REGISTER_API_URL"))
	gConf := gConfBuilder.Build()
	gwAdmin := fcrgatewayadmin.NewFilecoinRetrievalGatewayAdmin(*gConf)

	// Initialise gateways
	var gwIDs []*nodeid.NodeID
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
			panic(err)
		}
		// Enroll the gateway in the Register srv.
		if err := rm.RegisterGateway(gatewayRegistrar); err != nil {
			logging.Error("gateway registering error: %s", err.Error())
			t.FailNow()
		}
	}

	// Initialise providers
	pConfBuilder := fcrprovideradmin.CreateSettings()
	pConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	pConfBuilder.SetRegisterURL(providerConfig.GetString("REGISTER_API_URL"))
	pConf := pConfBuilder.Build()
	pAdmin := fcrprovideradmin.NewFilecoinRetrievalProviderAdmin(*pConf)
	var pIDs []*nodeid.NodeID
	for i := 0; i < 3; i++ {
		walletKey := privateKeys[0]
		walletAddress := accountAddrs[0]
		privateKeys = privateKeys[1:]
		accountAddrs = accountAddrs[1:]

		providerRootPubKey, providerRetrievalPubKey, providerRetrievalPrivateKey, err := cr.GenerateKeys()
		providerID := nodeid.NewRandomNodeID()
		pIDs = append(pIDs, providerID)

		identifier := fmt.Sprintf("-%v", i)
		providerRegistrar := register.NewProviderRegister(
			providerID.ToString(),
			walletAddress,
			providerRootPubKey,
			providerRetrievalPubKey,
			providerConfig.GetString("PROVIDER_REGION_CODE"),
			providerConfig.GetString("NETWORK_INFO_GATEWAY")[:8]+identifier+providerConfig.GetString("NETWORK_INFO_GATEWAY")[8:],
			providerConfig.GetString("NETWORK_INFO_CLIENT")[:8]+identifier+providerConfig.GetString("NETWORK_INFO_CLIENT")[8:],
			providerConfig.GetString("NETWORK_INFO_ADMIN")[:8]+identifier+providerConfig.GetString("NETWORK_INFO_ADMIN")[8:],
		)
		// Initialise the provider in the Admin
		err = pAdmin.InitialiseProviderV2(providerRegistrar, providerRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1), walletKey, lotusAP, lotusToken)
		if err != nil {
			panic(err)
		}
		// Enroll the provider in the Register srv.
		if err := rm.RegisterProvider(providerRegistrar); err != nil {
			logging.Error("error registering provider: %s", err.Error())
			t.FailNow()
		}
	}

	clientWalletKey := privateKeys[0]
	// drop used key and account pair
	privateKeys = privateKeys[1:]
	accountAddrs = accountAddrs[1:]

	// Configure client
	clientConfBuilder := fcrclient.CreateSettings()
	clientConfBuilder.SetEstablishmentTTL(101)
	clientConfBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	clientConfBuilder.SetRegisterURL(gatewayConfig.GetString("REGISTER_API_URL"))
	clientConfBuilder.SetWalletPrivateKey(clientWalletKey)
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

	added := client.AddGatewaysToUse(gwIDs)
	if !assert.Equal(t, 32, added, "32 gateways should be added") {
		t.FailNow()
	}

	added = client.AddActiveGateways(gwIDs)
	if !assert.Equal(t, 32, added, "32 gateways should be active") {
		t.FailNow()
	}

	// Force providers and gateways to update
	for _, p := range pIDs {
		err := pAdmin.ForceUpdate(p)
		if err != nil {
			logging.Error("provider update error: %s", err)
			t.FailNow()
		}
	}
	for _, g := range gwIDs {
		err := gwAdmin.ForceUpdate(g)
		if err != nil {
			logging.Error("gateway update error: %s", err)
			t.FailNow()
		}
	}

	newGatewayWalletKey := privateKeys[0]
	newGatwayWalletAddress := accountAddrs[0]
	privateKeys = privateKeys[1:]
	accountAddrs = accountAddrs[1:]

	gatewayRootPubKey, gatewayRetrievalPubKey, gatewayRetrievalPrivateKey, err := cr.GenerateKeys()
	newGatewayID, err := nodeid.NewNodeIDFromHexString("3880000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		panic(err)
	}
	gwIDs = append(gwIDs, newGatewayID)
	identifier := fmt.Sprintf("-32")
	gatewayRegistrar := register.NewGatewayRegister(
		newGatewayID.ToString(),
		newGatwayWalletAddress,
		gatewayRootPubKey,
		gatewayRetrievalPubKey,
		gatewayConfig.GetString("GATEWAY_REGION_CODE"),
		gatewayConfig.GetString("NETWORK_INFO_GATEWAY")[:7]+identifier+gatewayConfig.GetString("NETWORK_INFO_GATEWAY")[7:],
		gatewayConfig.GetString("NETWORK_INFO_PROVIDER")[:7]+identifier+gatewayConfig.GetString("NETWORK_INFO_PROVIDER")[7:],
		gatewayConfig.GetString("NETWORK_INFO_CLIENT")[:7]+identifier+gatewayConfig.GetString("NETWORK_INFO_CLIENT")[7:],
		gatewayConfig.GetString("NETWORK_INFO_ADMIN")[:7]+identifier+gatewayConfig.GetString("NETWORK_INFO_ADMIN")[7:],
	)
	err = gwAdmin.InitialiseGatewayV2(gatewayRegistrar, gatewayRetrievalPrivateKey, fcrcrypto.DecodeKeyVersion(1), newGatewayWalletKey, lotusAP, lotusToken)
	if err != nil {
		logging.Error("error to initialise gateway ID: %s; error: %s", gatewayRegistrar.GetNodeID(), err.Error())
		t.FailNow()
	}
	if err = rm.RegisterGateway(gatewayRegistrar); err != nil {
		logging.Error("error to register gateway ID: %s; error: %s", gatewayRegistrar.GetNodeID(), err.Error())
		t.FailNow()
	}
	// Force update providers and admins
	for _, p := range pIDs {
		err := pAdmin.ForceUpdate(p)
		if err != nil {
			logging.Error("error updating provider ID: %s; error: %s", p.ToString(), err.Error())
			t.FailNow()
		}
	}
	for _, g := range gwIDs {
		err := gwAdmin.ForceUpdate(g)
		if err != nil {
			logging.Error("error updating gateway ID: %s; error: %s", g.ToString(), err.Error())
			t.FailNow()
		}
	}

	err = gwAdmin.ListDHTOffer(newGatewayID)
	if err != nil {
		logging.Error("error listing DHT offer for the new gateway: %s", err.Error())
		t.FailNow()
	}

	added = client.AddGatewaysToUse([]*nodeid.NodeID{newGatewayID})
	if !assert.Equal(t, 1, added, "1 gateway should be added") {
		t.FailNow()
	}

	added = client.AddActiveGateways([]*nodeid.NodeID{newGatewayID})
	if !assert.Equal(t, 1, added, "1 gateway should be active") {
		t.FailNow()
	}

	////////////////////////////////////////////////////////
	// unique assertion part

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
	offers, err := client.FindOffersStandardDiscoveryV2(contentID0, gwIDs[6], 1)
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 0, len(offers), "Shouldn't find offer with cid 0 from gateway 6, outside of published ring.") {
		t.Fatal()
	}

	offers, err = client.FindOffersStandardDiscoveryV2(contentID0, gwIDs[32], 1)
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 32, boundary of published ring.") {
		t.Fatal()
	}

	offers, err = client.FindOffersStandardDiscoveryV2(contentID0, gwIDs[15], 1)
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 15, within published ring.") {
		t.Fatal()
	}

	offers, err = client.FindOffersStandardDiscoveryV2(contentID0, gwIDs[21], 1)
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 21, boundary of published ring.") {
		t.Fatal()
	}

	offers, err = client.FindOffersStandardDiscoveryV2(contentID0, gwIDs[22], 1)
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 0, len(offers), "Shouldn't find offer with cid 0 from gateway 22, outside of published ring.") {
		t.FailNow()
	}

	t.Log("/*******************************************************/")
	t.Log("/*        End TestPublishDHTOfferWithNewGateway        */")
	t.Log("/*******************************************************/")
}
