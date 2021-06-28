package fcrprovideradmin

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

import "github.com/ConsenSys/fc-retrieval/common/pkg/fcrcrypto"

// ProviderAdminSettings holds the library configuration
type ProviderAdminSettings struct {
	registerURL string

	blockchainPrivateKey *fcrcrypto.KeyPair

	providerAdminPrivateKey    *fcrcrypto.KeyPair
	providerAdminPrivateKeyVer *fcrcrypto.KeyVersion
}

// RegisterURL returns the register url
func (c *ProviderAdminSettings) RegisterURL() string {
	return c.registerURL
}

// BlockchainPrivateKey returns the blockchain private key
func (c *ProviderAdminSettings) BlockchainPrivateKey() *fcrcrypto.KeyPair {
	return c.blockchainPrivateKey
}

// ProviderAdminPrivateKey returns the provider admin private key
func (c *ProviderAdminSettings) ProviderAdminPrivateKey() *fcrcrypto.KeyPair {
	return c.providerAdminPrivateKey
}

// ProviderAdminPrivateKeyVer returns the provider admin private key version
func (c *ProviderAdminSettings) ProviderAdminPrivateKeyVer() *fcrcrypto.KeyVersion {
	return c.providerAdminPrivateKeyVer
}
