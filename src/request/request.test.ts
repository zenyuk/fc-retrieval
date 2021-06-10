import axios, { AxiosResponse } from 'axios'

import { getGateways, sendMessage } from './request'
import { FCRMessage } from '../fcrMessages/fcrMessage.class'

jest.mock('axios')
const mockedAxios = axios as jest.Mocked<typeof axios>

const mockedRegisterResponse: AxiosResponse = {
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

const mockedGatewayMessageRequest = new FCRMessage(205, { nonce: 42, isAlive: true })

const mockedGatewayResponse: AxiosResponse = {
  data: [
    {
      messageType: 206,
      protocolVersion: 1,
      protocolSupported: [1, 1],
      messageBody: { nonce: 42, isAlive: true },
      signature: '',
    },
  ],
  status: 200,
  statusText: 'OK',
  headers: {},
  config: {},
}

const mockedRegisterUrl = 'http://register'
const mockedGatewayUrl = 'http://gateway'
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

  describe('when successfully calls sendMessage', () => {
    it('send message request and get message response', async () => {
      mockedAxios.post.mockResolvedValue(mockedGatewayResponse)
      const gateways = await sendMessage(mockedGatewayUrl, mockedGatewayMessageRequest)
      expect(gateways).toEqual(mockedGatewayResponse.data)
    })
  })
  describe('when Gateway returns error', () => {
    it('throw error', async () => {
      mockedAxios.post.mockImplementationOnce(() => Promise.reject(new Error(mockedErrorMessage)))
      await expect(sendMessage(mockedGatewayUrl, mockedGatewayMessageRequest)).rejects.toThrow(mockedErrorMessage)
    })
  })
})
