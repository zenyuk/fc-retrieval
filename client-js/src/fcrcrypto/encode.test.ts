import { bin2String, decodeStringAsHexArray, string2Bin } from './encode.helper'

describe('encode/decode ', () => {
  it('string2Bin/bin2String ', async () => {
    const data = 'F8F9FAFBFCFDFEFF'
    const out = string2Bin(data)
    const outDec = bin2String(out)

    expect(outDec).toEqual(data)
  })

  it('decodeStringAsArray', async () => {
    const data = 'abcdefghijklmn'
    const out = decodeStringAsHexArray(data)
    console.log(`string2Bin ${out}`)

    expect(out).toEqual(new Uint8Array([171, 205, 239, 0, 0, 0, 0]))
  })
})
