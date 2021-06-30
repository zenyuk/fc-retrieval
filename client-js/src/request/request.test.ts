import axios, { AxiosResponse } from 'axios'

import { sendMessage } from './request'
import { FCRMessage } from '../fcrMessages/fcrMessage.class'

jest.mock('axios')
const mockedAxios = axios as jest.Mocked<typeof axios>

const mockedGatewayMessageRequest = new FCRMessage({
  message_type: 205,
  message_body: Buffer.from(JSON.stringify({ nonce: 42, isAlive: true })).toString('base64'),
})

const mockedGatewayResponse: AxiosResponse = {
  data: [
    {
      message_type: 206,
      protocol_version: 1,
      protocol_supported: [1, 1],
      message_body: 'eyJub25jZSI6NDIsImlzQWxpdmUiOnRydWV9', // base64-encoded {"nonce":42,"isAlive":true}
      message_signature: '',
    },
  ],
  status: 200,
  statusText: 'OK',
  headers: {},
  config: {},
}

const mockedGatewayUrl = 'http://gateway'
const mockedErrorMessage = 'error'

describe('Request', () => {
  describe('when successfully calls sendMessage', () => {
    it('send message request and get message response', async () => {
      mockedAxios.post.mockResolvedValue(mockedGatewayResponse)
      const gateways = await sendMessage(mockedGatewayUrl, mockedGatewayMessageRequest)
      expect(gateways).toEqual(new FCRMessage(mockedGatewayResponse.data))
    })
  })
  describe('when Gateway returns error', () => {
    it('throw error', async () => {
      mockedAxios.post.mockImplementationOnce(() => Promise.reject(new Error(mockedErrorMessage)))
      await expect(sendMessage(mockedGatewayUrl, mockedGatewayMessageRequest)).rejects.toThrow(mockedErrorMessage)
    })
  })
})
