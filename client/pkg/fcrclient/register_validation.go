package fcrclient

/*
 * Copyright 2020 ConsenSys Software Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

import (
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/register"
)

// Validate the information coming from the Register.
// Return true if the information is valid.
func validateGatewayInfo(gateway *register.GatewayRegister) bool {
	// All of the fields must have a value in them.
	if gateway.NodeID == "" {
		logging.Warn("Gateway registration issue: NodeID not set")
		return false
	}
	if gateway.Address == "" {
		logging.Warn("Gateway registration issue: Gateway IP address or domain name not set")
		return false
	}
	if gateway.NetworkInfoGateway == "" {
		logging.Warn("Gateway registration issue: Port for Gateway to Gateway communications not set")
		return false
	}
	if gateway.NetworkInfoProvider == "" {
		logging.Warn("Gateway registration issue: Port for Provider to Gateway communications not set")
		return false
	}
	if gateway.NetworkInfoClient == "" {
		logging.Warn("Gateway registration issue: Port for Client to Gateway communications not set")
		return false
	}
	if gateway.NetworkInfoAdmin == "" {
		logging.Warn("Gateway registration issue: Port for Admin to Gateway communications not set")
		return false
	}
	if gateway.RegionCode == "" {
		logging.Warn("Gateway registration issue: Region Code not set")
		return false
	}
	_, err := gateway.GetRootSigningKey()
	if err != nil {
		logging.Warn("Gateway registration issue: Root Signing Public Key error: %+v", err)
		return false
	}
	_, err = gateway.GetSigningKey()
	if err != nil {
		logging.Warn("Gateway registration issue: Retrieval Signing Key error: %+v", err)
		return false
	}
	return true
}

// Validate the information coming from the Register.
// Return true if the information is valid.
func validateProviderInfo(provider *register.ProviderRegister) bool {
	// All of the fields must have a value in them.
	if provider.NodeID == "" {
		logging.Warn("Provider registration issue: NodeID not set")
		return false
	}
	if provider.Address == "" {
		logging.Warn("Provider registration issue: Provider IP address or domain name not set")
		return false
	}
	if provider.NetworkInfoGateway == "" {
		logging.Warn("Provider registration issue: Port for Gateway to Provider communications not set")
		return false
	}
	if provider.NetworkInfoClient == "" {
		logging.Warn("Provider registration issue: Port for Client to Provider communications not set")
		return false
	}
	if provider.NetworkInfoAdmin == "" {
		logging.Warn("Provider registration issue: Port for Admin to Provider communications not set")
		return false
	}
	if provider.RegionCode == "" {
		logging.Warn("Provider registration issue: Region Code not set")
		return false
	}
	_, err := provider.GetRootSigningKey()
	if err != nil {
		logging.Warn("Provider registration issue: Root Signing Public Key error: %+v", err)
		return false
	}
	_, err = provider.GetSigningKey()
	if err != nil {
		logging.Warn("Provider registration issue: Retrieval Signing Key error: %+v", err)
		return false
	}
	return true
}
