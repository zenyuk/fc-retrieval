package gateway


import (
	"log"
	"errors"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/util"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/api/clientapi"
)


// Single instance of the gateway
var instance *Gateway


// Gateway holds the main data structure for the whole gateway.
type Gateway struct {
	clientAPI *clientapi.ClientAPI
}


// Create a new instance
func newInstance(settings util.AppSettings) (*Gateway, error) {
	var err error 
	var g = Gateway{}

	// Set-up the REST API
	g.clientAPI, err = clientapi.StartClientRestAPI(settings)
	if err != nil {
		log.Println("Error starting server: REST API: " + err.Error())
		return nil, err
	}


	return &g, nil
}


// Create is a factory method to get the single instance of the gateway
func Create(settings util.AppSettings) (*Gateway, error) {
	if (instance != nil) {
		return nil, errors.New("Error: instance already created")
	}
	return newInstance(settings)
}

// GetSingleInstance returns the single instance of the gateway
func GetSingleInstance(settings util.AppSettings) (*Gateway, error) {
	if (instance == nil) {
		return nil, errors.New("Error: instance not created")
	}
	return instance, nil
}


