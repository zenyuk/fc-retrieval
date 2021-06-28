import { array2str } from './encode.helper';

export class KeyVersion {
  static BitsInKeyVersion: number = 32;
  static LengthOfKeyVersionInBytes: number = 4;
  static InitialKeyVersion: KeyVersion = { keyVersion: 1 } as KeyVersion;
  static KeyVersionIncrement: number = 1;

  keyVersion: number;

  constructor(keyVersion: number) {
    this.keyVersion = keyVersion;
  }
}
export const decodeKeyVersionFromBytes = (version: Uint8Array): KeyVersion => {
  if (version.length < KeyVersion.LengthOfKeyVersionInBytes) {
    throw new Error(`Version bytes incorrect length: ${version.length}`);
  }
  const k = {} as KeyVersion;
  k.keyVersion = Number(array2str(version));

  return k;
};

export const encodeKeyVersionAsBytes = (k: KeyVersion): Uint8Array => {
  const out = new Uint8Array(KeyVersion.LengthOfKeyVersionInBytes);
  out[0] = k.keyVersion >> 24;
  out[1] = k.keyVersion >> 16;
  out[2] = k.keyVersion >> 8;
  out[3] = k.keyVersion >> 0;

  return out;
};
