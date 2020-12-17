package fcrcrypto

import (
	"log"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"reflect"
	"strconv"
	"math/big"
)

/**
 * Message sigining. This file signs and verifies signatures for messages in JSON format. 
 *
 * TODO: This code will need to be revisited / we may need to have a parallel version for the binary format
 * we end up using between provider and gateway.
 *
 * The process is:
 * 1. Generate a key pair using a certain ECC curve. Note that the curve must match the planned signature algorithm.
 * 2. Associate a key version with the key pair. This should be a monotomically increasing number starting
 *    at zero. An entity can only ever issue 256 keys / / roll-over their key 255 times. 
 * 3. Publish the public key to the Gateway or Provider registration contract, or for Clients, in a 
 *    Client - Gateway Establishment message. Publishing the key requires publishing the key version and
 *    signing algorithm that will be used with the key.
 * 4. Populate all fields of the message. Have Signature field = ""
 * 5. Sign a message, supplying the private key, the key version and algorithm.
 * 6. Receipients of the message extract the key version, and from there determine the public key and 
 *    signature algorithm. This may involve fetching information from the registration contract.
 * 7. Extract the signature from the signature field. Set the Signature to "".
 * 8. Verify the signature.
 * 
 * 
 *
 */


const (
	// SigAlgEcdsaP256Sha512_256 indicates the signature algorithm ECDSA P256 with SHA512/256
	SigAlgEcdsaP256Sha512_256 = 1


	// Offsets within a signature string for SigAlgEcdsaP256Sha512_256
	sigOfsKeyVersionStart = 0
	sigOfsKeyVersionEnd = 2
	sigOfsRStart = sigOfsKeyVersionEnd
	lenOf256BitInteger = 64 // hex characters
	sigOfsREnd = sigOfsRStart + lenOf256BitInteger
	sigOfsSStart = sigOfsREnd
	sigOfsSEnd = sigOfsSStart + lenOf256BitInteger

	hexadecimal = 16
	bitsInKeyVersion = 8
)

// SigAlg is holds a signature algorithm
type SigAlg struct {
	Alg uint8
}

// DecodeSigAlg converts a number to an object.
func DecodeSigAlg(alg uint8) *SigAlg {
	return &SigAlg{Alg: alg}
}

// KeyVersion wraps a key version number.
type KeyVersion struct {
	Ver uint8
}

// DecodeKeyVersion converts a number to an object.
func DecodeKeyVersion(ver uint8) *KeyVersion {
	return &KeyVersion{ver}
}


// Sign signs a message using the specified private key.
// Note that the struct must contain a field "Signature"
func Sign(pKey *ecdsa.PrivateKey, keyVersion KeyVersion, sigAlg SigAlg, msg interface{}) (*string, error) {
	if sigAlg.Alg != SigAlgEcdsaP256Sha512_256 {
		return nil, fmt.Errorf("Unknown signature algorithm: %d", sigAlg.Alg)
	}
	
	tbs := getToBeSigned(msg)
	log.Printf("Sign tbs: %s", tbs)

	hash := sha512.Sum512_256([]byte(tbs))
	r, s, err := ecdsa.Sign(rand.Reader, pKey, hash[:])
	if err != nil {
		return nil, err
	}

	alg := strconv.FormatInt(int64(keyVersion.Ver), hexadecimal)
	// TODO zero fill so that r and s are precisely 32 bytes (64 ascii characters) long
	rawSig := r.Text(hexadecimal) + s.Text(hexadecimal) 
	sig := alg + rawSig
	return &sig, nil
}

// ExtractKeyVersion extracts the key version from a signature string
func ExtractKeyVersion(signature *string) (*KeyVersion, error) {
	runes := []rune(*signature)
	keyVersionStr := string(runes[sigOfsKeyVersionStart:sigOfsKeyVersionEnd])
	keyVer, err := strconv.ParseUint(keyVersionStr, hexadecimal, bitsInKeyVersion)
	if err != nil {
		return nil, err
	}
	return &KeyVersion{Ver: uint8(keyVer)}, nil

}

// Verify verifies a message using the specified public key.
// Note that the struct must contain a field "Signature"
func Verify(pubKey *ecdsa.PublicKey, sigAlg SigAlg, signature *string, msg interface{}) (bool, error) {
	if sigAlg.Alg != SigAlgEcdsaP256Sha512_256 {
		return false, fmt.Errorf("Unknown signature algorithm: %d", sigAlg.Alg)
	}
	
	r, s, err := decodeSignature(signature)
	if err != nil {
		return false, err
	}
	tbs := getToBeSigned(msg)
	log.Printf("Verify tbs: %s", tbs)

	hash := sha512.Sum512_256([]byte(tbs))

	verified := ecdsa.Verify(pubKey, hash[:], r, s)
	return verified, nil
}

func decodeSignature(signature *string) (r, s *big.Int, err error) {
	runes := []rune(*signature)
	rStr := string(runes[sigOfsRStart:sigOfsREnd])
	r = new(big.Int)
    _, ok := r.SetString(rStr, hexadecimal)
    if !ok {
		return nil, nil, fmt.Errorf("Could not decode signature R value")
	}
	sStr := string(runes[sigOfsSStart:sigOfsSEnd])
	s = new(big.Int)
    _, ok = s.SetString(sStr, hexadecimal)
    if !ok {
		return nil, nil, fmt.Errorf("Could not decode signature S value")
	}
	return r, s, nil
}



func getToBeSigned(msg interface{}) string {
	v := reflect.ValueOf(msg)

	var allFields string
    for i := 0; i < v.NumField(); i++ {
		fieldAsString := v.Field(i).String()
		allFields = allFields + fieldAsString
	}
	return allFields
}

