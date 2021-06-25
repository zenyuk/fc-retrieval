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
})
