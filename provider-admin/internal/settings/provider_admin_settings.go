package settings

import "github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"

// ClientProviderAdminSettings holds the library configuration
type ClientProviderAdminSettings struct {
	registerURL string

	blockchainPrivateKey       *fcrcrypto.KeyPair
	providerAdminPrivateKey    *fcrcrypto.KeyPair
	providerAdminPrivateKeyVer *fcrcrypto.KeyVersion
}

// RegisterURL returns the register url
func (c *ClientProviderAdminSettings) RegisterURL() string {
	return c.registerURL
}

// BlockchainPrivateKey returns the blockchain private key
func (c *ClientProviderAdminSettings) BlockchainPrivateKey() *fcrcrypto.KeyPair {
	return c.blockchainPrivateKey
}

// ProviderAdminPrivateKey returns the provider admin private key
func (c *ClientProviderAdminSettings) ProviderAdminPrivateKey() *fcrcrypto.KeyPair {
	return c.providerAdminPrivateKey
}

// ProviderAdminPrivateKeyVer returns the provider admin private key version
func (c *ClientProviderAdminSettings) ProviderAdminPrivateKeyVer() *fcrcrypto.KeyVersion {
	return c.providerAdminPrivateKeyVer
}
