import { secp256k1Sign } from './secp256k1.service'

export const SigAlgEcdsaSecP256K1Blake2b: KeySigAlg = { algorithm: 1 } as KeySigAlg

export type KeySigAlg = {
  algorithm: number
}

export class KeyPair {
  pKey: Uint8Array
  pubKey: Uint8Array
  alg: KeySigAlg

  constructor({ alg, pKey = new Uint8Array(), pubKey = new Uint8Array() }: any) {
    this.pKey = pKey
    this.pubKey = pubKey
    this.alg = alg
  }

  sign(toBeSigned: Uint8Array): Uint8Array {
    if (this.alg !== SigAlgEcdsaSecP256K1Blake2b) {
      throw new Error('unsupported key algorithm: ' + this.alg.algorithm)
    }

    return secp256k1Sign(this.pKey, toBeSigned)
  }
}
