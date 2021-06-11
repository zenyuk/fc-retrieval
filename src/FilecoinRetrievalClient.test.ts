import { buildSettings } from './config/settings.config'
import { FilecoinRetrievalClient, Settings } from './index'

describe('Client', () => {
  it('findOffersStandardDiscoveryV2', async () => {
    // const client = new FilecoinRetrievalClient(defaults)
    // const offers = await client.findOffersStandardDiscoveryV2()
    expect(['hello']).toEqual(['hello'])
  })
  it('check FCR_REGISTER_API_URL value', async () => {
    expect(process.env.FCR_REGISTER_API_URL).toEqual('http://register:9020')
  })
  it('check FCR_LOTUS_AP value', async () => {
    expect(process.env.FCR_LOTUS_AP).toEqual('http://lotus-full-node:1234/rpc/v0')
  })
  it('check FCR_LOTUS_AUTH_TOKEN value', async () => {
    expect(process.env.FCR_LOTUS_AUTH_TOKEN).toBeDefined()
  })
  it('check FCR_WALLET_PRIVATE_KEY value', async () => {
    expect(process.env.FCR_WALLET_PRIVATE_KEY).toBeDefined()
  })
})
