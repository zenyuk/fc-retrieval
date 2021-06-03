export interface KeyPair {
  pKey: string
  pubKey: string
  alg: KeySigAlg
}

export type KeySigAlg = {
  algorithm: number
}
