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
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
)

// SettingsBuilder holds the library configuration
type SettingsBuilder struct {
	logLevel       string
	logTarget      string
	logServiceName string
	registerURL    string

	blockchainPrivateKey *fcrcrypto.KeyPair

	gatewayAdminPrivateKey    *fcrcrypto.KeyPair
	gatewayAdminPrivateKeyVer *fcrcrypto.KeyVersion
}

// CreateSettings creates an object with the default settings.
func CreateSettings() *SettingsBuilder {
	f := SettingsBuilder{}
	f.logLevel = defaultLogLevel
	f.logTarget = defaultLogTarget
	f.logServiceName = defaultLogServiceName
	return &f
}

// SetLogging sets the log level and target.
func (f *SettingsBuilder) SetLogging(logLevel string, logTarget string, logServiceName string) {
	f.logLevel = logLevel
	f.logTarget = logTarget
	f.logServiceName = logServiceName
}

// SetBlockchainPrivateKey sets the blockchain private key.
func (f *SettingsBuilder) SetBlockchainPrivateKey(bcPkey *fcrcrypto.KeyPair) {
	f.blockchainPrivateKey = bcPkey
}

// SetGatewayAdminPrivateKey sets the private key used for authenticating to the gateway
func (f *SettingsBuilder) SetGatewayAdminPrivateKey(rPkey *fcrcrypto.KeyPair, ver *fcrcrypto.KeyVersion) {
	f.gatewayAdminPrivateKey = rPkey
	f.gatewayAdminPrivateKeyVer = ver
}

// SetRegisterURL sets the URL of the register service
func (f *SettingsBuilder) SetRegisterURL(regURL string) {
	f.registerURL = regURL
}

// Build creates a settings object and initialises the logging system.
func (f *SettingsBuilder) Build() *GatewayAdminSettings {
	logging.Init1(f.logLevel, f.logTarget, f.logServiceName)

	g := GatewayAdminSettings{}
	g.registerURL = f.registerURL

	if f.blockchainPrivateKey == nil {
		logging.ErrorAndPanic("Settings: Blockchain Private Key not set")
	}
	g.blockchainPrivateKey = f.blockchainPrivateKey

	if f.gatewayAdminPrivateKey == nil {
		pKey, err := fcrcrypto.GenerateRetrievalV1KeyPair()
		if err != nil {
			logging.ErrorAndPanic("Settings: Error while generating random retrieval key pair: %s" + err.Error())
		}
		g.gatewayAdminPrivateKey = pKey
		g.gatewayAdminPrivateKeyVer = fcrcrypto.DecodeKeyVersion(1)
	} else {
		g.gatewayAdminPrivateKey = f.gatewayAdminPrivateKey
		g.gatewayAdminPrivateKeyVer = f.gatewayAdminPrivateKeyVer
	}

	return &g
}
