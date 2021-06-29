import axios, { AxiosResponse } from 'axios'

import { FilecoinRPC } from './filecoin_rpc.class'

jest.mock('axios')
const mockedAxios = axios as jest.Mocked<typeof axios>

const mockedFilecoinUrl = 'http://filecoin'
const mockedToken = 'eyJub25jZSI6NDIsImlzQWxpdmUiOnRydWV9'
const mockedAddress = 'eyJub25jZSI6NDIsImlzQWxpdmUiOnRydWV9'
const mockedCid = 'eyJub25jZSI6NDIsImlzQWxpdmUiOnRydWV9'
const mockedNonceResponse: AxiosResponse = {
  data: [
    {
      result: 42,
    },
  ],
  status: 200,
  statusText: 'OK',
  headers: {},
  config: {},
}

const mockedMpoolPushResponse: AxiosResponse = {
  data: {
    result: {
      cid: '987654321',
    },
  },
  status: 200,
  statusText: 'OK',
  headers: {},
  config: {},
}

const mockedMpoolPushErrorResponse: AxiosResponse = {
  data: {
    error: {
      message: 'custom error',
    },
  },
  status: 200,
  statusText: 'OK',
  headers: {},
  config: {},
}

const mockedWaitMessageResponse: AxiosResponse = {
  data: [
    {
      result: 42,
    },
  ],
  status: 200,
  statusText: 'OK',
  headers: {},
  config: {},
}

const mockedGasEstimationResponse: AxiosResponse = {
  data: [
    {
      result: 42,
    },
  ],
  status: 200,
  statusText: 'OK',
  headers: {},
  config: {},
}

const mockedReadStateResponse: AxiosResponse = {
  data: [
    {
      result: 42,
    },
  ],
  status: 200,
  statusText: 'OK',
  headers: {},
  config: {},
}

mockedAxios.create.mockImplementation(config => axios)

describe('FilecoinRPC', () => {
  describe('getNonce', () => {
    it('get Nonce', async () => {
      const filecoinRPC = new FilecoinRPC(mockedFilecoinUrl, mockedToken)
      mockedAxios.post.mockResolvedValue(mockedNonceResponse)
      const nonceData = await filecoinRPC.getNonce(mockedAddress)
      expect(nonceData).toEqual(mockedNonceResponse.data)
    })
  })

  describe('sendSignedMessage', () => {
    describe('sendSignedMessage', () => {
      it('send signed message', async () => {
        const filecoinRPC = new FilecoinRPC(mockedFilecoinUrl, mockedToken)
        mockedAxios.post.mockResolvedValue(mockedMpoolPushResponse)
        const nonceData = await filecoinRPC.sendSignedMessage(mockedAddress)
        expect(nonceData).toEqual(mockedMpoolPushResponse.data.result)
      })
    })
    describe('sendSignedMessage error', () => {
      it('rejects', async () => {
        const filecoinRPC = new FilecoinRPC(mockedFilecoinUrl, mockedToken)
        mockedAxios.post.mockResolvedValue(mockedMpoolPushErrorResponse)
        await expect(filecoinRPC.sendSignedMessage(mockedAddress)).rejects.toThrow(Error)
      })
    })
  })

  describe('waitMessage', () => {
    it('wait for message response', async () => {
      const filecoinRPC = new FilecoinRPC(mockedFilecoinUrl, mockedToken)
      mockedAxios.post.mockResolvedValue(mockedWaitMessageResponse)
      const nonceData = await filecoinRPC.waitMessage(mockedCid)
      expect(nonceData).toEqual(mockedWaitMessageResponse.data)
    })
  })

  describe('getGasEstimation', () => {
    it('get gas estimation', async () => {
      const filecoinRPC = new FilecoinRPC(mockedFilecoinUrl, mockedToken)
      mockedAxios.post.mockResolvedValue(mockedGasEstimationResponse)
      const nonceData = await filecoinRPC.getGasEstimation(mockedAddress)
      expect(nonceData).toEqual(mockedGasEstimationResponse.data)
    })
  })

  describe('readState', () => {
    it('read state', async () => {
      const filecoinRPC = new FilecoinRPC(mockedFilecoinUrl, mockedToken)
      mockedAxios.post.mockResolvedValue(mockedReadStateResponse)
      const nonceData = await filecoinRPC.readState(mockedAddress)
      expect(nonceData).toEqual(mockedReadStateResponse.data)
    })
  })
})
