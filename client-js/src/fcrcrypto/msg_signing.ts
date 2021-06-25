import { decodeKeyVersionFromBytes, encodeKeyVersionAsBytes, KeyVersion } from './key_version.class'
import { decodeSecP256K1PrivateKey, secp256k1Verify } from './secp256k1.service'
import { KeyPair, SigAlgEcdsaSecP256K1Blake2b } from './key_pair.class'
import { decodeHexArrayAsString, decodeStringAsHexArray } from './encode.helper'

const secp256k1PublicKeyBytes = 65
const sigOfsKeyVersionStart: number = 0
const sigOfsKeyVersionEnd: number = sigOfsKeyVersionStart + KeyVersion.LengthOfKeyVersionInBytes
const sigOfsRawSig: number = sigOfsKeyVersionEnd
// SignMessage signs a message using the specified private key.
// Note that the struct must contain a field "Signature"
export const signMessage = (pKey: KeyPair, keyVersion: KeyVersion, msg: any): string => {
  const toBeSigned = getToBeSigned(msg)
  const rawSign = pKey.sign(toBeSigned)
  const keyVerBytes = encodeKeyVersionAsBytes(keyVersion)

  return decodeHexArrayAsString(new Uint8Array([...keyVerBytes, ...rawSign, 0]))
}

// DecodePrivateKey decodes the algorithm and private key from a hex string.
export const decodePrivateKey = (encoded: string): KeyPair => {
  const algKeyBytes = decodeStringAsHexArray(encoded)
  const alg = algKeyBytes[0]

  switch (alg) {
    case SigAlgEcdsaSecP256K1Blake2b.algorithm:
      return decodeSecP256K1PrivateKey(algKeyBytes.slice(1))
    default:
      throw new Error('unknown private key algorithm:' + alg)
  }
}

// DecodePublicKey decodes the algorithm and public key from a hex string.
export const decodePublicKey = (encoded: string): KeyPair => {
  const algKeyBytes = decodeStringAsHexArray(encoded)

  const alg = algKeyBytes[0]
  switch (alg) {
    case SigAlgEcdsaSecP256K1Blake2b.algorithm:
      return decodeSecP256K1PublicKey(algKeyBytes.slice(1))
    default:
      throw Error('unknown public key algorithm: ' + alg)
  }
}

export const verifyAnyMessage = (pubKey: KeyPair, signature: string, msg: any): boolean => {
  return verifyMessage(pubKey, signature, JSON.stringify(msg))
}
export const verifyMessage = (pubKey: KeyPair, signature: string, msg: string): boolean => {
  const verify = secp256k1Verify(extractKeyFromMessage(signature), msg, pubKey.pubKey)
  return verify
}

const decodeSecP256K1PublicKey = (keyBytes: Uint8Array): KeyPair => {
  const key = new KeyPair({ alg: SigAlgEcdsaSecP256K1Blake2b.algorithm, pubKey: keyBytes })

  if (key.pubKey.length != secp256k1PublicKeyBytes) {
    throw new Error(`incorrect secp256k1 public key length: ${key.pubKey.length}`)
  }

  return key
}

// extractKeyFromMessage extracts the key from a signature string
export const extractKeyFromMessage = (signature: string): Uint8Array => {
  const sigBytes = decodeStringAsHexArray(signature)
  if (sigBytes.length < sigOfsRawSig) {
    throw new Error('sigBytes is empty, unable to verify, please check signature')
  }
  return sigBytes.slice(KeyVersion.LengthOfKeyVersionInBytes, sigBytes.length - 1)
}

// ExtractKeyVersionFromMessage extracts the key version from a signature string
export const extractKeyVersionFromMessage = (signature: string): KeyVersion => {
  const sigBytes = decodeStringAsHexArray(signature)
  return decodeKeyVersionFromBytes(sigBytes.slice(sigOfsKeyVersionStart, sigOfsKeyVersionEnd))
}

export const getToBeSigned = (msg: any): Uint8Array => {
  let raw: string = ''
  for (let key of Object.keys(msg)) {
    raw += msg[key]
  }

  return Buffer.from(raw)
}
