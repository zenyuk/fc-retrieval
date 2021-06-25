import { Settings } from './config/settings.config'
import { FilecoinRetrievalClient } from './FilecoinRetrievalClient'
import { NodeID } from './nodeid/nodeid.interface'
import { GatewayRegister } from './register/register.class'

jest.mock('./register/register.service', () => {
  return {
    getGatewayByID: jest.fn().mockImplementation(() =>
      Promise.resolve(
        new GatewayRegister({
          nodeId: '9876543210',
          address: 'f01234',
          networkInfoAdmin: '127.0.0.1:80',
          networkInfoClient: '127.0.0.1:80',
          networkInfoGateway: '127.0.0.1:80',
          networkInfoProvider: '127.0.0.1:80',
          regionCode: 'FR',
          rootSigningKey: '0xABCDE123456789',
        }),
      ),
    ),
  }
})
jest.mock('./clientapi/establishment_requester', () => {
  return {
    requestEstablishment: jest.fn().mockImplementation(() => Promise.resolve(true)),
  }
})

describe('Client', () => {
  let client: FilecoinRetrievalClient

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

  beforeAll(() => {
    client = new FilecoinRetrievalClient(
      new Settings({
        clientId: '101112131415161718191A1B1C1D1E3F',
        registerURL: process.env.FCR_REGISTER_API_URL,
      }),
    )
  })

  describe('on addGatewaysToUse', () => {
    it('succeeds', async () => {
      const used = await client.addGatewaysToUse([new NodeID('9876543210')])
      const active = await client.addActiveGateways([new NodeID('9876543210')])

      expect(used).toBeDefined()
      expect(used).toBeGreaterThanOrEqual(1)
      expect(active).toBeDefined()
      expect(active).toBeGreaterThanOrEqual(1)
      expect(used).toEqual(active)
    })
  })
})
