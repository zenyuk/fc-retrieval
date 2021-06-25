import { FCRMessage } from '../fcrMessages/fcrMessage.class'
import { FCRMessageType } from '../fcrMessages/type.enum'
import { NodeID } from '../nodeid/nodeid.interface'
import { GatewayRegister } from '../register/register.class'
import { requestEstablishment } from './establishment_requester'

jest.mock('../fcrMessages/client_establishment', () => {
  return {
    encodeClientEstablishmentRequest: jest.fn().mockImplementation(
      () =>
        new FCRMessage({
          message_type: FCRMessageType.ClientEstablishmentRequestType,
          message_body: Buffer.from(
            '{"client_id":"101112131415161718191A1B1C1D1E3F","challenge":"B8ydoaJhQrS8Bm36nZuDdhfQFisogUV9BdnoSPgze/Y=","ttl":"1624519937238"}',
          ).toString('base64'),
        }),
    ),
    decodeClientEstablishmentResponse: jest.fn().mockImplementation(() => ({
      gateway_id: new NodeID('9876543210'),
      challenge: 'B8ydoaJhQrS8Bm36nZuDdhfQFisogUV9BdnoSPgze/Y=',
    })),
  }
})
jest.mock('../request/request', () => {
  return {
    sendMessage: jest.fn().mockImplementation(() =>
      Promise.resolve(
        new FCRMessage({
          message_type: FCRMessageType.ClientEstablishmentResponseType,
          message_body: Buffer.from(
            '{"gateway_id":"9876543210","challenge":"B8ydoaJhQrS8Bm36nZuDdhfQFisogUV9BdnoSPgze/Y="}',
          ).toString('base64'),
        }),
      ),
    ),
  }
})

describe('establishment_requester.test', () => {
  describe('on requestEstablishment', () => {
    it('succeeds', async () => {
      const gatewayInfo = new GatewayRegister({
        nodeId: '9876543210',
        address: 'f01234',
        networkInfoAdmin: '127.0.0.1:80',
        networkInfoClient: '127.0.0.1:80',
        networkInfoGateway: '127.0.0.1:80',
        networkInfoProvider: '127.0.0.1:80',
        regionCode: 'FR',
        rootSigningKey:
          '01047799f37b014564e23578447d718e5c70a786b0e4e58ca25cb2a086b822434594d910b9b8c0fcbfe9f4c2db321e874819e0614be5b57fbb5080accd69adb2eaad',
      })
      const challenge = Buffer.from('B8ydoaJhQrS8Bm36nZuDdhfQFisogUV9BdnoSPgze/Y=', 'base64')
      const clientID = new NodeID('101112131415161718191A1B1C1D1E3F')
      const ttl = 1624519937238

      const done = await requestEstablishment(gatewayInfo, challenge, clientID, ttl)

      expect(done).toEqual(true)
    })
  })
})
