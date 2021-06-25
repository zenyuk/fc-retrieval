import axios from 'axios'
import { FilecoinRetrievalClient } from '../src/FilecoinRetrievalClient'
import { NodeID } from '../src/nodeid/nodeid.interface'
import { Settings } from '../src/config/settings.config'

describe('Client', () => {
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

  it('check FCR_GATEWAYS value', async () => {
    console.log('process.env.FCR_GATEWAYS ' + process.env.FCR_GATEWAYS)
    expect(process.env.FCR_GATEWAYS).toBeDefined()
  })

  it('check FCR_PROVIDERS value', async () => {
    console.log('process.env.FCR_PROVIDERS ' + process.env.FCR_PROVIDERS)
    expect(process.env.FCR_PROVIDERS).toBeDefined()
  })

  it.only('Call register', async () => {
    const defaultRegisterURL: string | undefined = process.env.FCR_REGISTER_API_URL
    const url = defaultRegisterURL + '/registers/gateway'
    const response = await axios.get(url)
    const gateways = response.data as any[]
    console.log(url + ' gateways ', gateways)
    expect(gateways.length).toBeGreaterThanOrEqual(2)
  })

  describe('on addGatewaysToUse', () => {
    it.only('succeeds', async () => {
      const client = new FilecoinRetrievalClient(
        new Settings({
          clientId: '101112131415161718191A1B1C1D1E3F',
          registerURL: 'http://localhost:9020',
        }),
      )

      // const gateways: string[] = (process.env.FCR_GATEWAYS as string).split(',')
      const nodeIDS = [new NodeID('101112131415161718191A1B1C1D1E3F202122232425262728292A2B2C2D2E1F')]
      const used = await client.addGatewaysToUse(nodeIDS)
      const active = await client.addActiveGateways(nodeIDS)

      expect(used).toBeDefined()
      // expect(used).toBeGreaterThanOrEqual(1);
      expect(active).toBeDefined()
      // expect(active).toBeGreaterThanOrEqual(1);
      expect(used).toEqual(active)
    })
  })
})
