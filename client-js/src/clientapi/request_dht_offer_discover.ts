import { ContentID } from '../cid/cid.interface'
import { GatewayRegister } from '../register/register.class'
import { NodeID } from '../nodeid/nodeid.interface'
import { SubCIDOffer } from '../cidoffer/subcidoffer.class'

export interface GatewaySubOffers {
  gatewayID: NodeID
  subOffers: SubCIDOffer[]
}

export const requestDHTOfferDiscover = (
  gw: GatewayRegister,
  contactedGateways: NodeID[],
  contentID: ContentID,
  nonce: number,
  offersDigestsFromAllGateways: string[][],
  paychAddrs: string,
  voucher: string,
): GatewaySubOffers[] => {
  return [] as GatewaySubOffers[]
}
