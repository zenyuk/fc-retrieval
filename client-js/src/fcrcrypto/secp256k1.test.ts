describe('Crypto dependencies test', () => {
  it('secp256k1', async () => {
    const secp256k1 = require('secp256k1')

    expect(secp256k1).toBeDefined()
  })
})
