package settings

// Copyright (C) 2020 ConsenSys Software Inc

// Filecoin Retrieval Client Settings

import (
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/register"
)

// BuilderImpl holds the library configuration
type BuilderImpl struct {
	logLevel               string
	logTarget              string
	establishmentTTL       int64
	tcpInactivityTimeout   int64
	clientID               *nodeid.NodeID
	registerURL            string
	providerRegister       *register.ProviderRegister
	blockchainPrivateKey   *fcrcrypto.KeyPair
	retrievalPrivateKey    *fcrcrypto.KeyPair
	retrievalPrivateKeyVer *fcrcrypto.KeyVersion
}

// CreateSettings creates an object with the default settings.
func CreateSettings() *BuilderImpl {
	f := BuilderImpl{}
	f.logLevel = defaultLogLevel
	f.logTarget = defaultLogTarget
	f.establishmentTTL = defaultEstablishmentTTL
	f.tcpInactivityTimeout = defaultTcpInactivityTimeout
	f.registerURL = defaultRegisterURL
	return &f
}

// SetLogging sets the log level and target.
func (f *BuilderImpl) SetLogging(logLevel string, logTarget string) {
	f.logLevel = logLevel
	f.logTarget = logTarget
}

// SetEstablishmentTTL sets the time to live for the establishment message between client and provider.
func (f *BuilderImpl) SetEstablishmentTTL(ttl int64) {
	f.establishmentTTL = ttl
}

// SetTcpInactivityTimeout sets the tcp inactivity timeout.
func (f *BuilderImpl) SetTcpInactivityTimeout(tcpInactivityTimeout int64) {
	f.tcpInactivityTimeout = tcpInactivityTimeout
}

// SetRegisterURL sets the register URL.
func (f *BuilderImpl) SetRegisterURL(url string) {
	f.registerURL = url
}

// SetProviderRegister sets the provider network info.
func (f *BuilderImpl) SetProviderRegister(info *register.ProviderRegister) {
	f.providerRegister = info
}

// SetBlockchainPrivateKey sets the blockchain private key.
func (f *BuilderImpl) SetBlockchainPrivateKey(bcPkey *fcrcrypto.KeyPair) {
	f.blockchainPrivateKey = bcPkey
}

// SetRetrievalPrivateKey sets the retrieval private key.
func (f *BuilderImpl) SetRetrievalPrivateKey(rPkey *fcrcrypto.KeyPair, ver *fcrcrypto.KeyVersion) {
	f.retrievalPrivateKey = rPkey
	f.retrievalPrivateKeyVer = ver
}

// Build creates a settings object and initialises the logging system.
func (f *BuilderImpl) Build() *ClientSettings {
	var err error

	logging.Init1(f.logLevel, f.logTarget)
	// logging.SetLogLevel(f.logLevel)
	// logging.SetLogTarget(f.logTarget)

	g := ClientSettings{}
	g.establishmentTTL = f.establishmentTTL
	g.tcpInactivityTimeout = f.tcpInactivityTimeout
	g.registerURL = f.registerURL
	g.providerRegister = f.providerRegister

	if f.blockchainPrivateKey == nil {
		logging.ErrorAndPanic("Settings: Blockchain Private Key not set")
	}
	g.blockchainPrivateKey = f.blockchainPrivateKey

	if f.clientID == nil {
		logging.Info("Settings: No Client ID set. Generating random client ID")
		// TODO replace once NewRandomNodeID becomes available.
		g.clientID, err = nodeid.NewRandomNodeID()
		if err != nil {
			logging.ErrorAndPanic("Settings: Error while generating random client ID: %s", err)
		}
	} else {
		g.clientID = f.clientID
	}

	if f.retrievalPrivateKey == nil {
		pKey, err := fcrcrypto.GenerateRetrievalV1KeyPair()
		if err != nil {
			logging.ErrorAndPanic("Settings: Error while generating random retrieval key pair: %s", err)
		}
		g.retrievalPrivateKey = pKey
		g.retrievalPrivateKeyVer = fcrcrypto.DecodeKeyVersion(1)
	} else {
		g.retrievalPrivateKey = f.retrievalPrivateKey
		g.retrievalPrivateKeyVer = f.retrievalPrivateKeyVer
	}

	return &g
}
