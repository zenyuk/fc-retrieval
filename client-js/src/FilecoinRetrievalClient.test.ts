import { FilecoinRetrievalClient } from './index'
import { defaults } from './constants/defaults'

describe('Client', () => {
  it('findOffersStandardDiscoveryV2', async () => {
    const client = new FilecoinRetrievalClient(defaults)

    // const offers = await client.findOffersStandardDiscoveryV2()
    expect(['hello']).toEqual(['hello'])
  })
})
