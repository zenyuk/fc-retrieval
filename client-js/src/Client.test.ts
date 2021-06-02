import { Client } from './index'
import { defaults } from './defaults'

describe('Client', () => {
  it('findOffersStandardDiscoveryV2', async () => {
    const client = new Client(defaults)

    const offers = await client.findOffersStandardDiscoveryV2()
    expect(offers).toEqual(['hello'])
  })
})
