import { Client } from './index'
import { defaults } from './defaults'
import axios, { AxiosResponse } from 'axios'
jest.mock('axios')
const mockedAxios = axios as jest.Mocked<typeof axios>

const mockedResponse: AxiosResponse = {
  data: [
    {
      nodeID: '9876543210',
      address: 'f01234',
      networkInfoAdmin: '127.0.0.1:80',
      networkInfoClient: '127.0.0.1:80',
      networkInfoGateway: '127.0.0.1:80',
      networkInfoProvider: '127.0.0.1:80',
      regionCode: 'FR',
      rootSigningKey: '0xABCDE123456789',
      signingKey: '0x987654321EDCBA',
    },
  ],
  status: 200,
  statusText: 'OK',
  headers: {},
  config: {},
}

describe('Client', () => {
  // it('findGateways', async () => {
  //   const client = new Client(defaults)
  //   mockedAxios.get.mockResolvedValue(mockedResponse)

  //   const gateways = await client.findGateways()
  //   expect(gateways).toEqual(mockedResponse.data)
  // })

  it('addGatewaysToUse', async () => {
    const client = new Client(defaults)
    // mockedAxios.get.mockResolvedValue(mockedResponse)

    await client.addGatewaysToUse(['9876543210'])
    const gateways = await client.getGatewaysToUse()
    expect(gateways.nodeID).toEqual({ '9876543210': mockedResponse.data[0].nodeID })
  })
})
