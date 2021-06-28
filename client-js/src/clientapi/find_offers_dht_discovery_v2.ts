import { ContentID } from '../cid/cid.interface'
import { GatewayRegister } from '../register/register.class'
import { NodeID } from '../nodeid/nodeid.interface'
import { FCRMessage } from '../fcrMessages/fcrMessage.class'

export interface requestedDHTDiscoverV2 {
  contactedGateways: NodeID[]
  contactedResp: FCRMessage[]
  uncontactable: NodeID[]
}

export const requestDHTDiscoverV2 = (
  gw: GatewayRegister,
  contentID: ContentID,
  nonce: number,
  ttl: number,
  numDHT: number,
  incrementalResult: boolean,
  paymentChannel: string,
  voucher: string,
): requestedDHTDiscoverV2 => {
  return {
    contactedGateways: [] as NodeID[],
    contactedResp: [] as FCRMessage[],
    uncontactable: [] as NodeID[],
  } as requestedDHTDiscoverV2
}

export interface decodedGatewayDHTDiscoverResponseV2 {
  pieceCID: ContentID
  nonce: number
  found: boolean
  subCidOffersDigest: string[]
  fundedPaymentChannel: boolean[]
}

export const decodeGatewayDHTDiscoverResponseV2 = (resp: FCRMessage): decodedGatewayDHTDiscoverResponseV2 => {
  return {} as decodedGatewayDHTDiscoverResponseV2
}
