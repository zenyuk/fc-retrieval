package fcrprovideradmin

import "github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"

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
