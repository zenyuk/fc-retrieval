package crypto_facade

import (
	"errors"
	"fmt"

	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrcrypto"
)

// GenerateKeys - helper function to generate set of keys, relying on file coin retrieval crypto package
func GenerateKeys() (rootPubKey string, retrievalPubKey string, retrievalPrivateKey *fcrcrypto.KeyPair, err error) {
	rootKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		return "", "", nil, fmt.Errorf("error generating blockchain key: %s", err.Error())
	}
	if rootKey == nil {
		return "", "", nil, errors.New("error generating blockchain key")
	}

	rootPubKey, err = rootKey.EncodePublicKey()
	if err != nil {
		return "", "", nil, fmt.Errorf("error encoding public key: %s", err.Error())
	}

	retrievalPrivateKey, err = fcrcrypto.GenerateRetrievalV1KeyPair()
	if err != nil {
		return "", "", nil, fmt.Errorf("error generating retrieval key: %s", err.Error())
	}
	if retrievalPrivateKey == nil {
		return "", "", nil, errors.New("error generating retrieval key")
	}

	retrievalPubKey, err = retrievalPrivateKey.EncodePublicKey()
	if err != nil {
		return "", "", nil, fmt.Errorf("error encoding retrieval pub key: %s", err.Error())
	}
	return
}
