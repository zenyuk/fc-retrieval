package main

// Copyright (C) 2020 ConsenSys Software Inc

import (
	"strings"
	"time"
	_ "github.com/joho/godotenv/autoload"

	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-register/pkg/register"
	"github.com/ConsenSys/fc-retrieval-gateway/config"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/api/adminapi"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/api/clientapi"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/api/gatewayapi"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/api/providerapi"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util"
)

func main() {
	conf := config.NewConfig()
	settings := config.Map(conf)
	logging.Init(conf)
	logging.Info("Filecoin Gateway Start-up: Started")

	logging.Info("Settings: %+v", settings)

	// Initialise a dummy gateway instance.
	g := gateway.GetSingleInstance(&settings)

	// Get all registerd Gateways
	gateways, err := register.GetRegisteredGateways(settings.RegisterAPIURL)
	if err != nil {
		logging.Error("Unable to get registered gateways: %v", err)
	}
	g.RegisteredGatewaysMapLock.Lock()
	logging.Info("All registered gateways: %+v", gateways)
	for _, gateway := range gateways {
		g.RegisteredGatewaysMap[strings.ToLower(gateway.NodeID)] = &gateway
	}
	g.RegisteredGatewaysMapLock.Unlock()

	// Get all registered Providers
	// providers, err := register.GetRegisteredProviders(settings.RegisterAPIURL)
	// if err != nil {
	// 	logging.Error("Unable to get registered providers: %v", err)
	// }
	// g.RegisteredProvidersMapLock.Lock()
	// logging.Info("All registered providers: %+v", providers)
	// for _, provider := range providers {
	// 	g.RegisteredProviders[provider.NodeID] = &provider
	// }
	// g.RegisteredProvidersMapLock.Unlock()

	err = clientapi.StartClientRestAPI(settings)
	if err != nil {
		logging.Error("Error starting server: Client REST API: %s", err.Error())
		return
	}

	err = gatewayapi.StartGatewayAPI(settings)
	if err != nil {
		logging.Error("Error starting gateway tcp server: %s", err.Error())
		return
	}

	err = providerapi.StartProviderAPI(settings)
	if err != nil {
		logging.Error("Error starting provider tcp server: %s", err.Error())
		return
	}

	err = adminapi.StartAdminAPI(settings, g)
	if err != nil {
		logging.Error("Error starting admin tcp server: %s", err.Error())
		return
	}

	go updateRegisteredProviders(settings.RegisterAPIURL, g)

	// Configure what should be called if Control-C is hit.
	util.SetUpCtrlCExit(gracefulExit)

	logging.Info("Filecoin Gateway Start-up Complete")

	// Wait forever.
	select {}
}

func updateRegisteredProviders(url string, c *gateway.Gateway) {
	logging.Info("Update registered Providers")
	for {
		providers, err := register.GetRegisteredProviders(url)
		if err != nil {
			logging.Error("Error in getting registered providers: %v", err.Error())
		} else {
			// Check if nothing is changed.
			update := false
			c.RegisteredProvidersMapLock.RLock()
			if len(providers) != len(c.RegisteredProvidersMap) {
				update = true
			} else {
				for _, provider := range providers {
					storedInfo, exist := c.RegisteredProvidersMap[strings.ToLower(provider.NodeID)]
					if !exist {
						update = true
						break
					} else {
						key, err := storedInfo.GetRootSigningKey()
						rootSigningKey, err2 := key.EncodePublicKey()
						key, err3 := storedInfo.GetSigningKey()
						signingKey, err4 := key.EncodePublicKey()
						if err != nil || err2 != nil || err3 != nil || err4 != nil {
							logging.Error("Error in generating key string")
							break
						}
						if provider.Address != storedInfo.GetAddress() ||
							provider.NetworkInfoAdmin != storedInfo.GetNetworkInfoAdmin() ||
							provider.NetworkInfoClient != storedInfo.GetNetworkInfoClient() ||
							provider.NetworkInfoGateway != storedInfo.GetNetworkInfoGateway() ||
							provider.RegionCode != storedInfo.GetRegionCode() ||
							provider.RootSigningKey != rootSigningKey ||
							provider.SigningKey != signingKey {
							update = true
							break
						}
					}
				}
			}
			c.RegisteredProvidersMapLock.RUnlock()
			if update {
				c.RegisteredProvidersMapLock.Lock()
				c.RegisteredProvidersMap = make(map[string]register.RegisteredNode)
				for _, provider := range providers {
					logging.Info("Add to registered providers map: nodeID=%+v", provider.NodeID)
					c.RegisteredProvidersMap[strings.ToLower(provider.NodeID)] = &provider
				}
				c.RegisteredProvidersMapLock.Unlock()
			}
		}
		// Sleep for 5 seconds, refresh every 5 seconds
		//time.Sleep(600 * time.Second)
		time.Sleep(5 * time.Second)
	}
}

func gracefulExit() {
	logging.Info("Filecoin Gateway Shutdown: Start")

	logging.Error("graceful shutdown code not written yet!")
	// TODO

	logging.Info("Filecoin Gateway Shutdown: Completed")
}
