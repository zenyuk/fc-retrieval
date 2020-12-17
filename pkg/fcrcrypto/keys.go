package fcrcrypto

import (
	"crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "crypto/x509"
    "encoding/pem"
)

// GenKeyPair generates a key pair
func GenKeyPair() (*ecdsa.PrivateKey, error) {
    curve := elliptic.P256()
    return ecdsa.GenerateKey(curve, rand.Reader)
}


// EncodePrivateKey converts a private key to a string
func EncodePrivateKey(privateKey *ecdsa.PrivateKey) string {
    x509Encoded, err := x509.MarshalECPrivateKey(privateKey)
    if err != nil {
        panic(err)
    }
    pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})
    return string(pemEncoded)
}


// EncodePublicKey converts a public key to a string
func EncodePublicKey(publicKey *ecdsa.PublicKey) string {
    x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)
    pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})
    return string(pemEncodedPub)
}


// DecodePrivateKey converts a string to a private key
func DecodePrivateKey(pemEncoded string) *ecdsa.PrivateKey {
    block, _ := pem.Decode([]byte(pemEncoded))
    x509Encoded := block.Bytes
    privateKey, _ := x509.ParseECPrivateKey(x509Encoded)
    return privateKey
}

// DecodePublicKey converts a string to a public key
func DecodePublicKey(pemEncodedPub string) *ecdsa.PublicKey {
    blockPub, _ := pem.Decode([]byte(pemEncodedPub))
    x509EncodedPub := blockPub.Bytes
    genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
    publicKey := genericPublicKey.(*ecdsa.PublicKey)
    return publicKey
}