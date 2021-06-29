import { KeyPair, SigAlgEcdsaSecP256K1Blake2b } from './key_pair.class'
import { retrievalV1Hash } from './blakeb.service'

const secp256k1 = require('secp256k1')
const blake2b = require('blake2b')

export const Secp256k1PrivateKeyBytes: number = 32

export const secp256k1Sign = (pKey: Uint8Array, toBeSigned: Uint8Array): Uint8Array => {
  const digest = retrievalV1Hash(toBeSigned)
  const sign: { signature: Uint8Array } = secp256k1.ecdsaSign(digest, pKey)
  return sign.signature
}

export const decodeSecP256K1PrivateKey = (keyBytes: Uint8Array): KeyPair => {
  const key = new KeyPair({ alg: SigAlgEcdsaSecP256K1Blake2b, pKey: keyBytes })

  if (key.pKey.length != Secp256k1PrivateKeyBytes) {
    throw new Error(`incorrect secp256k1 private key length: ${key.pKey}`)
  }
  return key
}

export const secp256k1Verify = (signature: Uint8Array, msg: string, pubKey: Uint8Array): boolean => {
  try {
    const digest = blake2b(32).update(Buffer.from(msg)).digest('binary')

    return secp256k1.ecdsaVerify(signature, digest, pubKey)
  } catch (e) {
    console.error(`Signature verify ${e}`)
    return false
  }
}
