package poc2_new_gateway

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ConsenSys/fc-retrieval-client/pkg/fcrclient"
	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrregistermgr"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/register"
	"github.com/ConsenSys/fc-retrieval-gateway-admin/pkg/fcrgatewayadmin"
	"github.com/ConsenSys/fc-retrieval-itest/pkg/util"
	"github.com/ConsenSys/fc-retrieval-provider-admin/pkg/fcrprovideradmin"
)

// Here we have gw 6-21 storing 688... and gw25-31&gw0-8 storing 008
// Plus the gw33 storing both 688 and 008

func TestNewGateway(t *testing.T) {
	t.Log("/*******************************************************/")
	t.Log("/*                 Start TestNewGateway                */")
	t.Log("/*******************************************************/")

	var err error
	privateKeys, accountAddrs, err := util.GenerateAccount(lotusAP, lotusToken, superAcct, 37)
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

		gatewayRootPubKey, gatewayRetrievalPubKey, gatewayRetrievalPrivateKey, err := generateKeys()
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
			logging.Error("error registering gateway: %s", err.Error())
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

		providerRootPubKey, providerRetrievalPubKey, providerRetrievalPrivateKey, err := generateKeys()
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
	err = pAdmin.PublishGroupCID(pIDs[1], []cid.ContentID{*contentID11, *contentID12, *contentID13}, 42, expiryDate, 42)
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
	err = pAdmin.PublishDHTCID(pIDs[0], []cid.ContentID{*contentID0}, []uint64{42}, []int64{expiryDate}, []uint64{42})
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
	err = pAdmin.PublishDHTCID(pIDs[0], []cid.ContentID{*contentID1}, []uint64{42}, []int64{expiryDate}, []uint64{42})
	if err != nil {
		logging.Error("error to publish DHT offer for content ID: %s; error: %s", contentID1.ToString(), err.Error())
		t.FailNow()
	}

	///////////////////////////////////////////////////////////////
	// unique part below

	newGatewayWalletKey := privateKeys[0]
	newGatwayWalletAddress := accountAddrs[0]
	privateKeys = privateKeys[1:]
	accountAddrs = accountAddrs[1:]

	gatewayRootPubKey, gatewayRetrievalPubKey, gatewayRetrievalPrivateKey, err := generateKeys()
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

	//////////////////////////////
	// unique assertion part

	offers, err := client.FindOffersStandardDiscoveryV2(contentID0, newGatewayID, 1)
	if err != nil {
		logging.Error("standard offers discovery error: %s", err.Error())
		t.FailNow()
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 0 from gateway 32.") {
		t.FailNow()
	}

	offers, err = client.FindOffersStandardDiscoveryV2(contentID1, newGatewayID, 1)
	if err != nil {
		logging.Error("standard offers discovery error: %s", err.Error())
		t.FailNow()
	}
	if !assert.Equal(t, 1, len(offers), "Should find offer with cid 1 from gateway 32.") {
		t.FailNow()
	}

	t.Log("/*******************************************************/")
	t.Log("/*                  End TestNewGateway                 */")
	t.Log("/*******************************************************/")
}
