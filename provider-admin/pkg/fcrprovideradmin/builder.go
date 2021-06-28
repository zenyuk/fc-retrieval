/*
Package fcrprovideradmin - is a Provider Admin library for FileCoin Secondary Retrieval.

It's designed to be imported by an application which requires to communicate with and administer Retrieval Providers
in FileCoin blockchain network.
*/
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

import (
	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval/common/pkg/logging"
)

// SettingsBuilder holds the library configuration
type SettingsBuilder struct {
	logLevel       string
	logTarget      string
	logServiceName string
	registerURL    string

	blockchainPrivateKey *fcrcrypto.KeyPair

	providerAdminPrivateKey    *fcrcrypto.KeyPair
	providerAdminPrivateKeyVer *fcrcrypto.KeyVersion

	lotusWalletPrivateKey string
	lotusAP               string
	lotusAuthToken        string
}

// CreateSettings creates an object with the default settings
func CreateSettings() *SettingsBuilder {
	f := SettingsBuilder{}
	f.logLevel = defaultLogLevel
	f.logTarget = defaultLogTarget
	f.logServiceName = defaultLogServiceName
	f.registerURL = defaultRegisterURL
	return &f
}

// SetLogging sets the log level and target.
func (f *SettingsBuilder) SetLogging(logLevel string, logTarget string, logServiceName string) {
	f.logLevel = logLevel
	f.logTarget = logTarget
	f.logServiceName = logServiceName
}

// SetRegisterURL sets the register URL.
func (f *SettingsBuilder) SetRegisterURL(url string) {
	f.registerURL = url
}

// SetBlockchainPrivateKey sets the blockchain private key.
func (f *SettingsBuilder) SetBlockchainPrivateKey(bcPkey *fcrcrypto.KeyPair) {
	f.blockchainPrivateKey = bcPkey
}

// SetProviderAdminPrivateKey sets the retrieval private key.
func (f *SettingsBuilder) SetProviderAdminPrivateKey(key *fcrcrypto.KeyPair, ver *fcrcrypto.KeyVersion) {
	f.providerAdminPrivateKey = key
	f.providerAdminPrivateKeyVer = ver
}

// Build creates a settings object and initialise the logging system
func (f *SettingsBuilder) Build() *ProviderAdminSettings {

	logging.Init1(f.logLevel, f.logTarget, f.logServiceName)

	c := &ProviderAdminSettings{}
	c.registerURL = f.registerURL

	if f.blockchainPrivateKey == nil {
		logging.ErrorAndPanic("Settings: Blockchain Private Key not set")
	}
	c.blockchainPrivateKey = f.blockchainPrivateKey

	if f.providerAdminPrivateKey == nil {
		pKey, err := fcrcrypto.GenerateRetrievalV1KeyPair()
		if err != nil {
			logging.ErrorAndPanic("Settings: Error while generating random retrieval key pair: %s", err)
		}
		c.providerAdminPrivateKey = pKey
		c.providerAdminPrivateKeyVer = fcrcrypto.DecodeKeyVersion(1)
	} else {
		c.providerAdminPrivateKey = f.providerAdminPrivateKey
		c.providerAdminPrivateKeyVer = f.providerAdminPrivateKeyVer
	}

	return c
}
