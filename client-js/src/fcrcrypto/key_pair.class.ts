import { secp256k1Sign } from './secp256k1.service';
const blake2b = require('blake2b');

export const SigAlgEcdsaSecP256K1Blake2b: KeySigAlg = { algorithm: 1 } as KeySigAlg;

export type KeySigAlg = {
  algorithm: number;
};

export class KeyPair {
  pKey: Uint8Array;
  pubKey: Uint8Array | undefined;
  alg: KeySigAlg;

  constructor(alg: KeySigAlg, pKey: Uint8Array, pubKey: Uint8Array | undefined) {
    this.pKey = pKey;
    this.pubKey = pubKey;
    this.alg = alg;
  }

  sign(toBeSigned: Uint8Array): Uint8Array {
    if (this.alg !== SigAlgEcdsaSecP256K1Blake2b) {
      throw new Error('unsupported key algorithm: ' + this.alg.algorithm);
    }
    const digest = retrievalV1Hash(toBeSigned);

    return secp256k1Sign(this.pKey, digest);
  }
}

// RetrievalV1Hash message digests some data using the algorithm used by version one of the
// Filecoin retrieval protocol.
export const retrievalV1Hash = (data: Uint8Array): Uint8Array => {
  return blake2b(32).update(data).digest('binary');
};
