package settings

import (
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	log "github.com/ConsenSys/fc-retrieval-common/pkg/logging"
)

// BuilderImpl holds the library configuration
type BuilderImpl struct {
	logLevel                   string
	logTarget                  string
	logServiceName             string
	registerURL                string
	blockchainPrivateKey       *fcrcrypto.KeyPair
	providerAdminPrivateKey    *fcrcrypto.KeyPair
	providerAdminPrivateKeyVer *fcrcrypto.KeyVersion
}

// CreateSettings creates an object with the default settings
func CreateSettings() *BuilderImpl {
	f := BuilderImpl{}
	f.logLevel = defaultLogLevel
	f.logTarget = defaultLogTarget
	f.logServiceName = defaultLogServiceName
	f.registerURL = defaultRegisterURL
	return &f
}

// SetLogging sets the log level and target.
func (f *BuilderImpl) SetLogging(logLevel string, logTarget string, logServiceName string) {
	f.logLevel = logLevel
	f.logTarget = logTarget
	f.logServiceName = logServiceName
}

// SetRegisterURL sets the register URL.
func (f *BuilderImpl) SetRegisterURL(url string) {
	f.registerURL = url
}

// SetBlockchainPrivateKey sets the blockchain private key.
func (f *BuilderImpl) SetBlockchainPrivateKey(bcPkey *fcrcrypto.KeyPair) {
	f.blockchainPrivateKey = bcPkey
}

// SetProviderAdminPrivateKey sets the retrieval private key.
func (f *BuilderImpl) SetProviderAdminPrivateKey(key *fcrcrypto.KeyPair, ver *fcrcrypto.KeyVersion) {
	f.providerAdminPrivateKey = key
	f.providerAdminPrivateKeyVer = ver
}

// Build creates a settings object and initialise the logging system
func (f *BuilderImpl) Build() *ClientProviderAdminSettings {

	log.Init1(f.logLevel, f.logTarget, f.logServiceName)

	c := &ClientProviderAdminSettings{}
	c.registerURL = f.registerURL

	if f.blockchainPrivateKey == nil {
		log.ErrorAndPanic("Settings: Blockchain Private Key not set")
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
