import axios, { AxiosResponse } from 'axios'
import { getGateways } from './register.service';

jest.mock('axios')
const mockedAxios = axios as jest.Mocked<typeof axios>

const mockedRegisterResponse: AxiosResponse = {
  data: [
    {
      nodeId: '9876543210',
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

const mockedRegisterUrl = 'http://register'
const mockedErrorMessage = 'error'

describe('Request', () => {
  describe('when successfully calls getGateways', () => {
    it('getGateways', async () => {
      mockedAxios.get.mockResolvedValue(mockedRegisterResponse)
      const gateways = await getGateways(mockedRegisterUrl)
      expect(gateways).toEqual(mockedRegisterResponse.data)
    })
  })
  describe('when Register returns error', () => {
    it('throw error', async () => {
      mockedAxios.get.mockImplementationOnce(() => Promise.reject(new Error(mockedErrorMessage)))
      await expect(getGateways(mockedRegisterUrl)).rejects.toThrow(mockedErrorMessage)
    })
  })
})
