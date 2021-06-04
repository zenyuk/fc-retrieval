import { ContentID } from '../cid/cid.interface'
import { FCRMessage } from '../fcrMessages/fcrMessage.class'
import { NodeID } from '../nodeid/nodeid.interface'
import { ProviderRegister } from '../register/register.class'
import {
  encodeClientDHTOfferAckRequest,
  decodeClientDHTOfferAckResponse,
  DHTOfferAckResponse,
} from '../fcrMessages/client_dht_offer_ack'
import { sendMessage } from '../request/request'

export const requestDHTOfferAck = (
  providerInfo: ProviderRegister,
  cid: ContentID,
  gatewayID: NodeID,
): DHTOfferAckResponse => {
  const request = encodeClientDHTOfferAckRequest(cid, gatewayID)

  const response = sendMessage(providerInfo.networkInfoClient, request)

  const pubKey = providerInfo.getSigningKey()

  response.verify(pubKey)

  const dhtOfferAckResponse = decodeClientDHTOfferAckResponse(response)

  return dhtOfferAckResponse
}
