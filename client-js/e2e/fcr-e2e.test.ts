import axios from 'axios'

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

  it('Call register', async () => {
    const defaultRegisterURL: string = process.env.FCR_REGISTER_API_URL
    const url = defaultRegisterURL + '/registers/gateway'
    const response = await axios.get(url)
    const gateways = response.data as any[]
    expect(gateways.length).toBeGreaterThanOrEqual(2)
  })
})
