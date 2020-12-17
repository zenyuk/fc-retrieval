package gateway

import (
	"crypto/ecdsa"
	"errors"
	"sync"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway/clients"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/offers"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)

const (
	protocolVersion   = 1 // Main protocol version
	protocolSupported = 1 // Alternative protocol version
)

// CommunicationChannels holds the node id and internal communication channels for a live tcp connection
type CommunicationChannels struct {
	NodeID                              nodeid.NodeID
	CommsRequestChan, CommsResponseChan chan []byte
}

// Gateway holds the main data structure for the whole gateway.
type Gateway struct {
	ProtocolVersion   int32
	ProtocolSupported []int32

	// GatewayID of this gateway
	GatewayID *nodeid.NodeID

	// Gateway Private Key ID of this gateway
	GatewayPrivateKey *ecdsa.PrivateKey

	// GatewayPrivateKeyVersion is the key version number of the private key.
	GatewayPrivateKeyVersion *fcrcrypto.KeyVersion

	// GatewayPrivateKeySigAlg is the signing algorithm to be used with the private key.
	GatewayPrivateKeySigAlg *fcrcrypto.SigAlg

	// ActiveGateways store connected active gateways
	// A map from gateway id (big int in string repr)
	// to a CommunicationThread
	ActiveGateways     map[string](*CommunicationChannels)
	ActiveGatewaysLock sync.RWMutex

	// ActiveProviders store connected active providers
	// A map from provider id (big in in string repr)
	// to a CommunicationThread
	ActiveProviders     map[string](*CommunicationChannels)
	ActiveProvidersLock sync.RWMutex

	GatewayClient *clients.GatewayClientInteraction

	// Offers, it is threadsafe.
	Offers *offers.Offers
}

// Single instance of the gateway
var instance *Gateway
var doOnce sync.Once

// GetSingleInstance returns the single instance of the gateway
func GetSingleInstance(confs ...*settings.AppSettings) *Gateway {
	doOnce.Do(func() {
		if len(confs) == 0 {
			panic("No settings supplied to Gateway start-up")
		}
		if len(confs) != 1 {
			panic("More than one sets of settings supplied to Gateway start-up")
		}
		conf := confs[0]

		gatewayPrivateKey := fcrcrypto.DecodePrivateKey(conf.GatewayPrivKey)
		gatewayID, err2 := nodeid.NewNodeIDFromString(conf.GatewayID) 
		if err2 != nil {
			panic(err2)
		}

		gatewayPrivateKeyVersion := fcrcrypto.DecodeKeyVersion(conf.GatewayPrivKeyVersion)
		gatewayPrivateKeySigAlg := fcrcrypto.DecodeSigAlg(conf.GatewaySigAlg)

		instance = &Gateway{
			ProtocolVersion:     protocolVersion,
			ProtocolSupported:   []int32{protocolVersion, protocolSupported},
			ActiveGateways:      make(map[string](*CommunicationChannels)),
			ActiveGatewaysLock:  sync.RWMutex{},
			ActiveProviders:     make(map[string](*CommunicationChannels)),
			ActiveProvidersLock: sync.RWMutex{},
			GatewayClient:		 &clients.GatewayClientInteraction{},
			GatewayPrivateKey:	 gatewayPrivateKey,
			GatewayPrivateKeyVersion: gatewayPrivateKeyVersion,
			GatewayPrivateKeySigAlg: gatewayPrivateKeySigAlg,
			GatewayID:			 gatewayID, 
			Offers:              offers.GetSingleInstance(),
		}
	})
	return instance
}

// RegisterGatewayCommunication registers a gateway communication
func RegisterGatewayCommunication(id nodeid.NodeID, gComms *CommunicationChannels) error {
	if instance == nil {
		return errors.New("Error: instance not created")
	}
	instance.ActiveGatewaysLock.Lock()
	defer instance.ActiveGatewaysLock.Unlock()
	_, exist := instance.ActiveGateways[id.ToString()]
	if exist {
		return errors.New("Error: connection existed")
	}
	instance.ActiveGateways[id.ToString()] = gComms
	return nil
}

// DeregisterGatewayCommunication deregisters a gateway communication
// Fail silently
func DeregisterGatewayCommunication(id nodeid.NodeID) {
	if instance != nil {
		instance.ActiveGatewaysLock.Lock()
		defer instance.ActiveGatewaysLock.Unlock()
		_, exist := instance.ActiveGateways[id.ToString()]
		if exist {
			delete(instance.ActiveGateways, id.ToString())
		}
	}
}

// RegisterProviderCommunication registers a provider communication
func RegisterProviderCommunication(id nodeid.NodeID, pComms *CommunicationChannels) error {
	if instance == nil {
		return errors.New("Error: instance not created")
	}
	instance.ActiveProvidersLock.Lock()
	defer instance.ActiveProvidersLock.Unlock()
	_, exist := instance.ActiveProviders[id.ToString()]
	if exist {
		return errors.New("Error: connection existed")
	}
	instance.ActiveProviders[id.ToString()] = pComms
	return nil
}

// DeregisterProviderCommunication deregisters a provider communication
func DeregisterProviderCommunication(id nodeid.NodeID) error {
	if instance == nil {
		return errors.New("Error: instance not created")
	}
	instance.ActiveProvidersLock.Lock()
	defer instance.ActiveProvidersLock.Unlock()
	_, exist := instance.ActiveProviders[id.ToString()]
	if !exist {
		return errors.New("Error: connection not existed")
	}
	delete(instance.ActiveProviders, id.ToString())
	return nil
}
