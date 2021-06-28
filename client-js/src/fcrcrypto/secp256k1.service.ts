import { KeyPair, SigAlgEcdsaSecP256K1Blake2b } from './key_pair.class';

const secp256k1 = require('secp256k1');

export const Secp256k1PrivateKeyBytes: number = 32;

export const secp256k1Sign = (pKey: Uint8Array, digest: Uint8Array): Uint8Array => {
  const sign: { signature: Uint8Array } = secp256k1.ecdsaSign(digest, pKey);
  return sign.signature;
};

export const decodeSecP256K1PrivateKey = (keyBytes: Uint8Array): KeyPair => {
  const key = new KeyPair(SigAlgEcdsaSecP256K1Blake2b, keyBytes, undefined);

  if (key.pKey.length != Secp256k1PrivateKeyBytes) {
    throw new Error(`incorrect secp256k1 private key length: ${key.pKey}`);
  }
  return key;
};
