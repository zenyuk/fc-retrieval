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

	err := clientapi.StartClientRestAPI(settings)
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

	// Get all registerd Gateways
	go updateRegisteredGateways(settings.RegisterAPIURL, g)

	// Get all registered Providers
	go updateRegisteredProviders(settings.RegisterAPIURL, g)

	// Configure what should be called if Control-C is hit.
	util.SetUpCtrlCExit(gracefulExit)

	logging.Info("Filecoin Gateway Start-up Complete")

	// Wait forever.
	select {}
}

func updateRegisteredGateways(url string, g *gateway.Gateway) {
	for {
		gateways, err := register.GetRegisteredGateways(url)
		if err != nil {
			logging.Error("Error in getting registered gateways: %s", err.Error())
		} else {
			// Check if nothing is changed.
			update := false
			g.RegisteredGatewaysMapLock.RLock()
			if len(gateways) != len(g.RegisteredGatewaysMap) {
				update = true
			} else {
				for _, gateway := range gateways {
					// Skip itself
					if gateway.NodeID == g.GatewayID.ToString() {
						continue
					}
					storedInfo, exist := g.RegisteredGatewaysMap[strings.ToLower(gateway.NodeID)]
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
						if gateway.Address != storedInfo.GetAddress() ||
							gateway.NetworkInfoAdmin != storedInfo.GetNetworkInfoAdmin() ||
							gateway.NetworkInfoClient != storedInfo.GetNetworkInfoClient() ||
							gateway.NetworkInfoProvider != storedInfo.GetNetworkInfoProvider() ||
							gateway.NetworkInfoGateway != storedInfo.GetNetworkInfoGateway() ||
							gateway.RegionCode != storedInfo.GetRegionCode() ||
							gateway.RootSigningKey != rootSigningKey ||
							gateway.SigningKey != signingKey {
							update = true
							break
						}
					}
				}
			}
			g.RegisteredGatewaysMapLock.RUnlock()
			if update {
				g.RegisteredGatewaysMapLock.Lock()
				g.RegisteredGatewaysMap = make(map[string]register.RegisteredNode)
				for _, gateway := range gateways {
					// Skip itself
					if gateway.NodeID == g.GatewayID.ToString() {
						continue
					}
					logging.Info("Add to registered gateways map: nodeID=%+v", gateway.NodeID)
					g.RegisteredGatewaysMap[strings.ToLower(gateway.NodeID)] = &gateway
				}
				g.RegisteredGatewaysMapLock.Unlock()
			}
		}
		// Sleep for 5 seconds, refresh every 5 seconds
		time.Sleep(5 * time.Second)
	}
}

func updateRegisteredProviders(url string, g *gateway.Gateway) {
	for {
		providers, err := register.GetRegisteredProviders(url)
		if err != nil {
			logging.Error("Error in getting registered providers: %s", err.Error())
		} else {
			// Check if nothing is changed.
			update := false
			g.RegisteredProvidersMapLock.RLock()
			if len(providers) != len(g.RegisteredProvidersMap) {
				update = true
			} else {
				for _, provider := range providers {
					storedInfo, exist := g.RegisteredProvidersMap[strings.ToLower(provider.NodeID)]
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
			g.RegisteredProvidersMapLock.RUnlock()
			if update {
				g.RegisteredProvidersMapLock.Lock()
				g.RegisteredProvidersMap = make(map[string]register.RegisteredNode)
				for _, provider := range providers {
					logging.Info("Add to registered providers map: nodeID=%+v", provider.NodeID)
					g.RegisteredProvidersMap[strings.ToLower(provider.NodeID)] = &provider
				}
				g.RegisteredProvidersMapLock.Unlock()
			}
		}
		// Sleep for 5 seconds, refresh every 5 seconds
		time.Sleep(5 * time.Second)
	}
}

func gracefulExit() {
	logging.Info("Filecoin Gateway Shutdown: Start")

	logging.Error("graceful shutdown code not written yet!")
	// TODO

	logging.Info("Filecoin Gateway Shutdown: Completed")
}
