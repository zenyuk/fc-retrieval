import { ContentID } from '../cid/cid.interface'
import { NodeID } from '../nodeid/nodeid.interface'
import { FCRMessage } from './fcrMessage.class'

export interface DHTOfferAckRequest {
  pieceCID: ContentID
  gatewayID: NodeID
}

export interface DHTOfferAckResponse {
  cid: ContentID // piece cid
  gatewayID: NodeID // gateway id
  found: boolean // found
  offerRequest: FCRMessage // publish dht offer request
  offerResponse: FCRMessage // publish dht offer resposne
}

export const encodeClientDHTOfferAckRequest = (pieceCID: ContentID, gatewayID: NodeID): FCRMessage => {
  return {} as FCRMessage
}

export const decodeClientDHTOfferAckRequest = (fcrMsg: FCRMessage): DHTOfferAckRequest => {
  return {} as DHTOfferAckRequest
}
export const encodeClientDHTOfferAckResponse = (pieceCID: ContentID, gatewayID: NodeID): FCRMessage => {
  return {} as FCRMessage
}

export const decodeClientDHTOfferAckResponse = (fcrMsg: FCRMessage): DHTOfferAckResponse => {
  return {} as DHTOfferAckResponse
}
