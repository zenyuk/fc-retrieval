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
    it.skip('succeeds', async () => {
      const gatewayInfo = new GatewayRegister({
        nodeId: '9876543210',
        address: 'f01234',
        networkInfoAdmin: '127.0.0.1:80',
        networkInfoClient: '127.0.0.1:80',
        networkInfoGateway: '127.0.0.1:80',
        networkInfoProvider: '127.0.0.1:80',
        regionCode: 'FR',
        signingKey:
          '010472e9ab95ed0171cc9f07e9ac0cde6ad23040a97f079ac5702c39867c59149c7c071415ca41c565ef7dda4ebf3cfeb4d52703329e06234720c2e3d25211737ad5',
        signature:
          '0000000106ccdf77b9f655f7f61ca64a219b91891799bbce0373402b2aba763694aed6834b7b369efb39717370d689ac1ac25b45b760cc777653b56fcef6854527e28e2e01',
      })
      const challenge = Buffer.from('B8ydoaJhQrS8Bm36nZuDdhfQFisogUV9BdnoSPgze/Y=', 'base64')
      const clientID = new NodeID('101112131415161718191A1B1C1D1E3F')
      const ttl = 1624519937238

      const done = await requestEstablishment(gatewayInfo, challenge, clientID, ttl)
      expect(done).toEqual(false)
    })
  })
})
