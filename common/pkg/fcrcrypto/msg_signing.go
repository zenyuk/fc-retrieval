package fcrcrypto

import (
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
)

/**
 * Message sigining. This file contains code to sign and verifie signatures for messages in JSON format.
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
 */

const (
	// Offsets within a signature string for SigAlgEcdsaP256Sha512_256
	sigOfsKeyVersionStart = 0
	sigOfsKeyVersionEnd   = sigOfsKeyVersionStart + lengthOfKeyVersionInBytes
	sigOfsRawSig          = sigOfsKeyVersionEnd
)

// SignMessage signs a message using the specified private key.
// Note that the struct must contain a field "Signature"
func SignMessage(pKey *KeyPair, keyVersion *KeyVersion, msg interface{}) (string, error) {
	rawSig, err := pKey.Sign(getToBeSigned(msg))
	if err != nil {
		return "", err
	}
	keyVerBytes := keyVersion.EncodeKeyVersionAsBytes()
	sigBytes := append(keyVerBytes, rawSig...)
	return hex.EncodeToString(sigBytes), nil
}

// ExtractKeyVersionFromMessage extracts the key version from a signature string
func ExtractKeyVersionFromMessage(signature string) (*KeyVersion, error) {
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return nil, err
	}
	return DecodeKeyVersionFromBytes(sigBytes[sigOfsKeyVersionStart:sigOfsKeyVersionEnd])
}

// VerifyMessage verifies a message using the specified public key.
// Note that the struct must contain a field "Signature"
func VerifyMessage(pubKey *KeyPair, signature string, msg interface{}) (bool, error) {
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return false, err
	}

	if len(sigBytes) < sigOfsRawSig {
		err := errors.New("sigBytes is empty, unable to verify, please check signature")
		return false, err
	}
	return pubKey.Verify(sigBytes[sigOfsRawSig:], getToBeSigned(msg))
}

func getToBeSigned(msg interface{}) []byte {
	allFields := DumpStructPayload(msg)

	return []byte(allFields)
}

func DumpStructPayloadV(val reflect.Value) string {

	if val.Kind() == reflect.Interface && !val.IsNil() {
		elm := val.Elem()
		if elm.Kind() == reflect.Ptr && !elm.IsNil() && elm.Elem().Kind() == reflect.Ptr {
			val = elm
		}
	}
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	var out string
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)

		if valueField.Kind() == reflect.Interface && !valueField.IsNil() {
			elm := valueField.Elem()
			if elm.Kind() == reflect.Ptr && !elm.IsNil() && elm.Elem().Kind() == reflect.Ptr {
				valueField = elm
			}
		}

		if valueField.Kind() == reflect.Ptr {
			valueField = valueField.Elem()
		}

		var value interface{}
		if valueField.CanInterface() {
			value = valueField.Interface()
		}

		if value != nil {
			out = fmt.Sprintf("%v%v", out, value)
		}

		if valueField.Kind() == reflect.Struct {
			DumpStructPayloadV(valueField)
		}
	}

	for i := 0; i < val.NumMethod(); i++ {
		out += val.Method(i).Call([]reflect.Value{})[0].String()
	}
	return out
}

func DumpStructPayload(v interface{}) string {
	return DumpStructPayloadV(reflect.ValueOf(v))
}
