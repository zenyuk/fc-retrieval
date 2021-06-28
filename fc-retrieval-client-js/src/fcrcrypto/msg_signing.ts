import { decodeKeyVersionFromBytes, encodeKeyVersionAsBytes, KeyVersion } from './key_version.class';
import { decodeSecP256K1PrivateKey } from './secp256k1.service';
import { KeyPair, SigAlgEcdsaSecP256K1Blake2b } from './key_pair.class';
import { decodeStringAsHexArray, str2array } from './encode.helper';

const sigOfsKeyVersionStart: number = 0;
const sigOfsKeyVersionEnd: number = sigOfsKeyVersionStart + KeyVersion.LengthOfKeyVersionInBytes;

// SignMessage signs a message using the specified private key.
// Note that the struct must contain a field "Signature"
export const signMessage = (pKey: KeyPair, keyVersion: KeyVersion, msg: any): Uint8Array => {
  const toBeSigned = getToBeSigned(msg);
  const rawSign = pKey.sign(toBeSigned);
  const keyVerBytes = encodeKeyVersionAsBytes(keyVersion);

  return new Uint8Array([...keyVerBytes, ...rawSign, 0]);
};

export const verifyMessage = async (pubKey: string, signature: string, msg: any): Promise<boolean> => {
  const secp256k1Async = require('js-secp256k1/dist/node-bundle');
  const secp256k1 = await secp256k1Async();
  return secp256k1.ecdsaVerify(str2array(signature), getToBeSigned(msg), str2array(pubKey));
};

// ExtractKeyVersionFromMessage extracts the key version from a signature string
export const extractKeyVersionFromMessage = (signature: string): KeyVersion => {
  const sigBytes = decodeStringAsHexArray(signature);
  return decodeKeyVersionFromBytes(sigBytes.slice(sigOfsKeyVersionStart, sigOfsKeyVersionEnd));
};

// DecodePrivateKey decodes the algorithm and private key from a hex string.
export const decodePrivateKey = (encoded: string): KeyPair => {
  const algKeyBytes = decodeStringAsHexArray(encoded);
  const alg = algKeyBytes[0];

  switch (alg) {
    case SigAlgEcdsaSecP256K1Blake2b.algorithm:
      return decodeSecP256K1PrivateKey(algKeyBytes.slice(1));
    default:
      throw new Error('unknown private key algorithm:' + alg);
  }
};

export const getToBeSigned = (msg: any): Uint8Array => {
  let raw: string = '';
  for (let key of Object.keys(msg)) {
    raw += msg[key];
  }

  return Buffer.from(raw);
};
