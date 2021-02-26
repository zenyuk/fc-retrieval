package fcrprovideradmin

// Copyright (C) 2020 ConsenSys Software Inc

// Filecoin Retrieval Client Settings

import (
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/register"
	"github.com/ConsenSys/fc-retrieval-provider-admin/internal/settings"
)

// SettingsBuilder holds the library configuration
type SettingsBuilder interface {
	// SetLogging sets the log level and target.
	SetLogging(logLevel string, logTarget string)

	// SetEstablishmentTTL sets the time to live for the establishment message between client and provider.
	SetEstablishmentTTL(ttl int64)

	// SetTcpInactivityTimeout sets the tcp inactivity timeout.
	SetTcpInactivityTimeout(tcpInactivityTimeout int64)

	// SetBlockchainPrivateKey sets the blockchain private key.
	SetBlockchainPrivateKey(bcPkey *fcrcrypto.KeyPair)

	// SetRetrievalPrivateKey sets the retrieval private key.
	SetRetrievalPrivateKey(rPkey *fcrcrypto.KeyPair, ver *fcrcrypto.KeyVersion)

	// SetRegisterURL sets the register URL.
	SetRegisterURL(url string)

	// SetProviderRegister sets the provider network info.
	SetProviderRegister(info *register.ProviderRegister)

	// Build creates a settings object and initialises the logging system.
	Build() *Settings
}

// Settings holds the library configuration
type Settings interface {
	EstablishmentTTL() int64
	TcpInactivityTimeout() int64
	ClientID() *nodeid.NodeID
	BlockchainPrivateKey() *fcrcrypto.KeyPair
	RetrievalPrivateKey() *fcrcrypto.KeyPair
	RetrievalPrivateKeyVer() *fcrcrypto.KeyVersion
	RegisterURL() string
	ProviderRegister() *register.ProviderRegister
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
func (f settingsBuilderImpl) SetLogging(logLevel string, logTarget string) {
	f.impl.SetLogging(logLevel, logTarget)
}

// SetEstablishmentTTL sets the time to live for the establishment message between client and provider.
func (f settingsBuilderImpl) SetEstablishmentTTL(ttl int64) {
	f.impl.SetEstablishmentTTL(ttl)
}

// SetTcpInactivityTimeout sets the tcp inactivity timeout.
func (f settingsBuilderImpl) SetTcpInactivityTimeout(tcpInactivityTimeout int64) {
	f.impl.SetTcpInactivityTimeout(tcpInactivityTimeout)
}

// SetBlockchainPrivateKey sets the blockchain private key.
func (f settingsBuilderImpl) SetBlockchainPrivateKey(bcPkey *fcrcrypto.KeyPair) {
	f.impl.SetBlockchainPrivateKey(bcPkey)
}

// SetRetrievalPrivateKey sets the retrieval private key.
func (f settingsBuilderImpl) SetRetrievalPrivateKey(rPkey *fcrcrypto.KeyPair, ver *fcrcrypto.KeyVersion) {
	f.impl.SetRetrievalPrivateKey(rPkey, ver)
}

// SetRegisterURL sets the register URL.
func (f settingsBuilderImpl) SetRegisterURL(url string) {
	f.impl.SetRegisterURL(url)
}

// SetProviderURL sets the provider URL.
func (f settingsBuilderImpl) SetProviderRegister(info *register.ProviderRegister) {
	f.impl.SetProviderRegister(info)
}

// Build generates the settings.
func (f settingsBuilderImpl) Build() *Settings {
	clientSettings := f.impl.Build()
	set := Settings(clientSettings)
	return &set
}
