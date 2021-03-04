package fcrprovideradmin

// Copyright (C) 2020 ConsenSys Software Inc

// Filecoin Retrieval Client Settings

import (
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-provider-admin/internal/settings"
)

// SettingsBuilder holds the library configuration
type SettingsBuilder interface {
	// SetLogging sets the log level and target.
	SetLogging(logLevel string, logTarget string, logServiceName string)

	// SetRegisterURL sets the register URL.
	SetRegisterURL(url string)

	// SetBlockchainPrivateKey sets the blockchain private key.
	SetBlockchainPrivateKey(bcPkey *fcrcrypto.KeyPair)

	// SetProviderAdminPrivateKey sets the retrieval private key.
	SetProviderAdminPrivateKey(rPkey *fcrcrypto.KeyPair, ver *fcrcrypto.KeyVersion)

	// Build creates a settings object and initialises the logging system.
	Build() *Settings
}

// Settings holds the library configuration
type Settings interface {
	RegisterURL() 								string
	BlockchainPrivateKey() 				*fcrcrypto.KeyPair
	ProviderAdminPrivateKey() 		*fcrcrypto.KeyPair
	ProviderAdminPrivateKeyVer() 	*fcrcrypto.KeyVersion
}

// CreateSettings loads up default settings
func CreateSettings() SettingsBuilder {
	f := newBuilderImpl()
	builder := SettingsBuilder(f)
	return builder
}

type settingsBuilderImpl struct {
	impl *settings.BuilderImpl
}

func newBuilderImpl() settingsBuilderImpl {
	return settingsBuilderImpl{settings.CreateSettings()}
}

// SetLogging sets the log level and target.
func (f settingsBuilderImpl) SetLogging(logLevel string, logTarget string, logServiceName string) {
	f.impl.SetLogging(logLevel, logTarget, logServiceName)
}

// SetBlockchainPrivateKey sets the blockchain private key.
func (f settingsBuilderImpl) SetBlockchainPrivateKey(bcPkey *fcrcrypto.KeyPair) {
	f.impl.SetBlockchainPrivateKey(bcPkey)
}

// SetproviderAdminPrivateKey sets the retrieval private key.
func (f settingsBuilderImpl) SetProviderAdminPrivateKey(rPkey *fcrcrypto.KeyPair, ver *fcrcrypto.KeyVersion) {
	f.impl.SetProviderAdminPrivateKey(rPkey, ver)
}

// SetRegisterURL sets the register URL.
func (f settingsBuilderImpl) SetRegisterURL(url string) {
	f.impl.SetRegisterURL(url)
}

// Build generates the settings.
func (f settingsBuilderImpl) Build() *Settings {
	clientSettings := f.impl.Build()
	set := Settings(clientSettings)
	return &set
}
