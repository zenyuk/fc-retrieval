package gateway

import (
	"errors"
	"log"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/api"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/api/clientapi"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)

// Single instance of the gateway
var instance *api.Gateway

// Create a new instance
func newInstance(settings util.AppSettings) (*api.Gateway, error) {
	var err error
	var g = api.Gateway{}

	// Set-up the REST API
	g.ClientAPI, err = clientapi.StartClientRestAPI(settings)
	if err != nil {
		log.Println("Error starting server: REST API: " + err.Error())
		return nil, err
	}

	return &g, nil
}

// Create is a factory method to get the single instance of the gateway
func Create(settings util.AppSettings) (*api.Gateway, error) {
	if instance != nil {
		return nil, errors.New("Error: instance already created")
	}
	return newInstance(settings)
}

// GetSingleInstance returns the single instance of the gateway
func GetSingleInstance() (*api.Gateway, error) {
	if instance == nil {
		return nil, errors.New("Error: instance not created")
	}
	return instance, nil
}

// RegisterGatewayCommunication registers a gateway communication
func RegisterGatewayCommunication(id nodeid.NodeID, gComms *api.CommunicationThread) error {
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
func RegisterProviderCommunication(id nodeid.NodeID, pComms *api.CommunicationThread) error {
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
