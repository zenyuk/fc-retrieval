import { FCRMessage } from "../fcrMessages/fcrMessage.class";
import { FCRMessageType } from "../fcrMessages/type.enum";
import { NodeID } from "../nodeid/nodeid.interface";
import { GatewayRegister } from "../register/register.class";
import { requestEstablishment } from "./establishment_requester";

jest.mock('../fcrMessages/client_establishment', () => {
  return {
    encodeClientEstablishmentRequest: jest.fn().mockImplementation(() => new FCRMessage(
      FCRMessageType.ClientEstablishmentRequestType,
			'{"client_id":"101112131415161718191A1B1C1D1E3F","challenge":"B8ydoaJhQrS8Bm36nZuDdhfQFisogUV9BdnoSPgze/Y=","ttl":"1624519937238"}',
    )),
		decodeClientEstablishmentResponse: jest.fn().mockImplementation(() => ({
      gateway_id: new NodeID('9876543210'),
			challenge: 'B8ydoaJhQrS8Bm36nZuDdhfQFisogUV9BdnoSPgze/Y=',
    })),
  };
});
jest.mock('../request/request', () => {
  return {
    sendMessage: jest.fn().mockImplementation(() => Promise.resolve(new FCRMessage(
			FCRMessageType.ClientEstablishmentResponseType,
			'{"gateway_id":"9876543210","challenge":"B8ydoaJhQrS8Bm36nZuDdhfQFisogUV9BdnoSPgze/Y="}'
		))),
  };
});

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
				rootSigningKey: '0xABCDE123456789',
				signingKey: '0x987654321EDCBA',
			})
			const challenge = Buffer.from('B8ydoaJhQrS8Bm36nZuDdhfQFisogUV9BdnoSPgze/Y=', 'base64')
			const clientID = new NodeID('101112131415161718191A1B1C1D1E3F')
			const ttl = 1624519937238

			const done = await requestEstablishment(gatewayInfo, challenge, clientID, ttl)

      expect(done).toEqual(true)
    });
  })
});
