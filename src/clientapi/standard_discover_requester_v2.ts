import { ContentID } from '../cid/cid.interface'
import { GatewayRegister } from '../register/register.class'
import { sendMessage } from '../request/request'

// RequestStandardDiscoverV2 requests a standard discover to a given gateway for a given contentID, nonce and ttl.
export const requestStandardDiscoverV2 = (
  gatewayInfo: GatewayRegister,
  contentID: ContentID,
  nonce: number,
  ttl: number,
  paychAddr: string,
  voucher: string,
): string[] => {
  const resp = sendMessage(gatewayInfo.networkInfoClient, request)

  //gatewayInfo.signingKey

  //
  const offerDigests = [] as string[]
  return offerDigests
}
