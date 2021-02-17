package settings

// Copyright (C) 2020 ConsenSys Software Inc

// Filecoin Retrieval Gateway Admin Client Settings

import (
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
)

// ClientGatewayAdminSettings holds the library configuration
type ClientGatewayAdminSettings struct {
	establishmentTTL int64

	blockchainPrivateKey *fcrcrypto.KeyPair

	gatewayAdminPrivateKey    *fcrcrypto.KeyPair
	gatewayAdminPrivateKeyVer *fcrcrypto.KeyVersion
}

// EstablishmentTTL returns the establishmentTTL
func (c ClientGatewayAdminSettings) EstablishmentTTL() int64 {
	return c.establishmentTTL
}

// BlockchainPrivateKey returns the BlockchainPrivateKey
func (c ClientGatewayAdminSettings) BlockchainPrivateKey() *fcrcrypto.KeyPair {
	return c.blockchainPrivateKey
}

// GatewayAdminPrivateKey returns the GatewayAdminPrivateKey
func (c ClientGatewayAdminSettings) GatewayAdminPrivateKey() *fcrcrypto.KeyPair {
	return c.gatewayAdminPrivateKey
}

// GatewayAdminPrivateKeyVer returns the GatewayAdminKeyVer
func (c ClientGatewayAdminSettings) GatewayAdminPrivateKeyVer() *fcrcrypto.KeyVersion {
	return c.gatewayAdminPrivateKeyVer
}
