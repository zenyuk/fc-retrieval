import { ContentID } from '../cid/cid.interface'
import { SubCIDOffer } from '../cidoffer/cidoffer.interface'
import { GatewayRegister } from '../register/register.class'

export const requestStandardDiscoverOffer = (
  gatewayInfo: GatewayRegister,
  contentID: ContentID,
  nonce: number,
  ttl: number,
  offerDigests: string[],
  paychAddr: string,
  voucher: string,
): SubCIDOffer[] => {
  const offers = [] as SubCIDOffer[]
  return offers
}
