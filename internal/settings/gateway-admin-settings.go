package settings

// Copyright (C) 2020 ConsenSys Software Inc

// Filecoin Retrieval Gateway Admin Client Settings

import (
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)

// ClientGatewayAdminSettings holds the library configuration
type ClientGatewayAdminSettings struct {
	establishmentTTL int64
	clientID         *nodeid.NodeID

	blockchainPrivateKey *fcrcrypto.KeyPair

	retrievalPrivateKey    *fcrcrypto.KeyPair
	retrievalPrivateKeyVer *fcrcrypto.KeyVersion
}

// EstablishmentTTL returns the establishmentTTL
func (c ClientGatewayAdminSettings) EstablishmentTTL() int64 {
	return c.establishmentTTL
}

// ClientID returns the ClientID
func (c ClientGatewayAdminSettings) ClientID() *nodeid.NodeID {
	return c.clientID
}

// BlockchainPrivateKey returns the BlockchainPrivateKey
func (c ClientGatewayAdminSettings) BlockchainPrivateKey() *fcrcrypto.KeyPair {
	return c.blockchainPrivateKey
}

// RetrievalPrivateKey returns the RetrievalPrivateKey
func (c ClientGatewayAdminSettings) RetrievalPrivateKey() *fcrcrypto.KeyPair {
	return c.retrievalPrivateKey
}

// RetrievalPrivateKeyVer returns the RetrievalPrivateKeyVer
func (c ClientGatewayAdminSettings) RetrievalPrivateKeyVer() *fcrcrypto.KeyVersion {
	return c.retrievalPrivateKeyVer
}
