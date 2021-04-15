package fcrgatewayadmin

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
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
)

// GatewayAdminSettings holds the library configuration
type GatewayAdminSettings struct {
	establishmentTTL int64

	blockchainPrivateKey *fcrcrypto.KeyPair

	gatewayAdminPrivateKey    *fcrcrypto.KeyPair
	gatewayAdminPrivateKeyVer *fcrcrypto.KeyVersion

	registerURL string
}

// EstablishmentTTL returns the establishmentTTL
func (c GatewayAdminSettings) EstablishmentTTL() int64 {
	return c.establishmentTTL
}

// BlockchainPrivateKey returns the BlockchainPrivateKey
func (c GatewayAdminSettings) BlockchainPrivateKey() *fcrcrypto.KeyPair {
	return c.blockchainPrivateKey
}

// GatewayAdminPrivateKey returns the GatewayAdminPrivateKey
func (c GatewayAdminSettings) GatewayAdminPrivateKey() *fcrcrypto.KeyPair {
	return c.gatewayAdminPrivateKey
}

// GatewayAdminPrivateKeyVer returns the GatewayAdminKeyVer
func (c GatewayAdminSettings) GatewayAdminPrivateKeyVer() *fcrcrypto.KeyVersion {
	return c.gatewayAdminPrivateKeyVer
}

// RegisterURL is the URL to the register service
func (c GatewayAdminSettings) RegisterURL() string {
	return c.registerURL
}
