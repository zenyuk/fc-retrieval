import { bin2String, decodeString, decodeStringAsHexArray, encodeToString, string2Bin } from './encode.helper';

describe('encode/decode ', () => {
  it('encodeToString', async () => {
    const data = 'F8F9FAFBFCFDFEFF';
    const out = encodeToString(data);
    const outDec = decodeString(out);

    expect(outDec).toEqual(data);
  });

  it('string2Bin/bin2String ', async () => {
    const data = 'F8F9FAFBFCFDFEFF';
    const out = string2Bin(data);
    const outDec = bin2String(out);

    expect(outDec).toEqual(data);
  });

  it('decodeStringAsArray', async () => {
    const data = 'abcdefghijklmn';
    const out = decodeStringAsHexArray(data);
    console.log(`string2Bin ${out}`);

    expect(out).toEqual(new Uint8Array([171, 205, 239, 0, 0, 0, 0]));
  });
});
