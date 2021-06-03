import { ContentID } from '../cid/cid.interface'
import { GatewayRegister } from '../register/register.class'

// RequestStandardDiscoverV2 requests a standard discover to a given gateway for a given contentID, nonce and ttl.
export const requestStandardDiscoverV2 = (
  gatewayInfo: GatewayRegister,
  contentID: ContentID,
  nonce: number,
  ttl: number,
  paychAddr: string,
  voucher: string,
): string[] => {
  //
  const offerDigests = [] as string[]
  return offerDigests
}
