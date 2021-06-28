describe('Crypto dependencies test', () => {
  it('blake2b', async () => {
    const blake2b = require('blake2b')

    let output = new Uint8Array(32)
    const input = Buffer.alloc(2048)

    expect(blake2b(output.length).update(input).digest('hex')).toEqual(
      '200823e5158b3774c11b5c61850ada762f8264144a9bebec3ebac5a2adde67b8',
    )

    blake2b.ready(function () {
      expect(blake2b(output.length).update(input).digest('hex')).toEqual(
        '200823e5158b3774c11b5c61850ada762f8264144a9bebec3ebac5a2adde67b8',
      )
    })
  })

  it('blake2b empty', async () => {
    const blake2b = require('blake2b')

    const output = new Uint8Array(32)
    const empty = Buffer.from('00')
    const blake2bEmpty = blake2b(output.length).update(empty).digest('binary')
    expect(blake2bEmpty.length).toEqual(32)
    expect(blake2bEmpty[0]).toStrictEqual(203)
    expect(blake2bEmpty[1]).toStrictEqual(198)
    expect(blake2bEmpty[31]).toStrictEqual(66)
  })
  it('blake2b', async () => {
    const blake2b = require('blake2b')

    const msg = `{"message_type":101,"protocol_version":1,"protocol_supported":[1,1],"message_body":"eyJnYXRld2F5X2lkIjoiMDgwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMCIsImNoYWxsZW5nZSI6Ik1lZitXMGJKR2w3OW8zcjRmNXN1OUdZSGpDd3RtY2pGUmQ2aXNmRk1kMU09In0=","message_signature":""}`
    const buffer = Buffer.from(msg)
    const digest = blake2b(32).update(buffer).digest('binary')
    expect(digest).toEqual(
      new Uint8Array([
        250, 228, 212, 173, 107, 45, 174, 93, 0, 32, 37, 177, 189, 29, 53, 70, 230, 241, 197, 82, 249, 53, 26, 178, 9,
        220, 224, 114, 19, 56, 65, 230,
      ]),
    )
  })
})
