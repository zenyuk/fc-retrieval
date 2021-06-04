import { ContentID } from '../cid/cid.interface'
import { SubCIDOffer } from '../cidoffer/subcidoffer.class'
import { GatewayRegister } from '../register/register.class'
import { sendMessage } from '../request/request'
import { FCRMessage } from '../fcrMessages/fcrMessage.class'

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

  const request = {} as FCRMessage
  const response = sendMessage(gatewayInfo.networkInfoClient, request)

  return offers
}
