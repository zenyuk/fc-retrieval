package settings

// Copyright (C) 2020 ConsenSys Software Inc

// Filecoin Retrieval Client Settings

import (
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
)

// BuilderImpl holds the library configuration
type BuilderImpl struct {
	logLevel         string
	logTarget        string
	logServiceName   string
	establishmentTTL int64
	clientID         *nodeid.NodeID
	registerURL      string

	blockchainPrivateKey *fcrcrypto.KeyPair

	retrievalPrivateKey    *fcrcrypto.KeyPair
	retrievalPrivateKeyVer *fcrcrypto.KeyVersion
}

// CreateSettings creates an object with the default settings.
func CreateSettings() *BuilderImpl {
	f := BuilderImpl{}
	f.logLevel = defaultLogLevel
	f.logTarget = defaultLogTarget
	f.logServiceName = defaultLogServiceName
	f.establishmentTTL = defaultEstablishmentTTL
	f.registerURL = defaultRegisterURL
	return &f
}

// SetLogging sets the log level and target.
func (f *BuilderImpl) SetLogging(logLevel string, logTarget string, logServiceName string) {
	f.logLevel = logLevel
	f.logTarget = logTarget
	f.logServiceName = logServiceName
}

// SetEstablishmentTTL sets the time to live for the establishment message between client and gateway.
func (f *BuilderImpl) SetEstablishmentTTL(ttl int64) {
	f.establishmentTTL = ttl
}

// SetRegisterURL sets the register URL.
func (f *BuilderImpl) SetRegisterURL(url string) {
	f.registerURL = url
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

	logging.Init1(f.logLevel, f.logTarget, f.logServiceName)
	// logging.SetLogLevel(f.logLevel)
	// logging.SetLogTarget(f.logTarget)

	g := ClientSettings{}
	g.establishmentTTL = f.establishmentTTL
	g.registerURL = f.registerURL

	if f.blockchainPrivateKey == nil {
		logging.ErrorAndPanic("Settings: Blockchain Private Key not set")
	}
	g.blockchainPrivateKey = f.blockchainPrivateKey

	if f.clientID == nil {
		logging.Info("Settings: No Client ID set. Generating random client ID")
		// TODO replace once NewRandomNodeID becomes available.
		g.clientID = nodeid.NewRandomNodeID()
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
