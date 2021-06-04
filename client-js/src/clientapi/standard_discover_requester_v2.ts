import { ContentID } from '../cid/cid.interface'
import { GatewayRegister } from '../register/register.class'
import { createFCRMessage, FcrMessage } from '../fcrMessages/fcrMessage'

// RequestStandardDiscoverV2 requests a standard discover to a given gateway for a given contentID, nonce and ttl.
export const requestStandardDiscoverV2 = (
  gatewayInfo: GatewayRegister,
  contentID: ContentID,
  nonce: number,
  ttl: number,
  paychAddr: string,
  voucher: string,
): string[] => {
  const request = encodeClientStandardDiscoverRequestV2(contentID, nonce, ttl, paychAddr, voucher)

  const response = sendMessage(gatewayInfo.networkInfoClient, request)

  const pubKey = gatewayInfo.getSigningKey()

  response.verify(pubKey)

  const responseV2 = decodeClientStandardDiscoverResponseV2(response)

  return responseV2.subCIDOfferDigests
}

export interface DecodedClientStandardDiscoverResponseV2 {
  contentID: ContentID
  nonce: number
  found: boolean
  subCIDOfferDigests: string[]
  fundedPaymentChannel: boolean[]
}

export const decodeClientStandardDiscoverResponseV2 = (
  response: FcrMessage,
): DecodedClientStandardDiscoverResponseV2 => {
  return {} as DecodedClientStandardDiscoverResponseV2
}

export const encodeClientStandardDiscoverRequestV2 = (
  contentID: ContentID,
  nonce: number,
  ttl: number,
  paychAddr: string,
  voucher: string,
): FcrMessage => {
  let ClientStandardDiscoverRequestV2Type = 0
  let body = ''

  return createFCRMessage(ClientStandardDiscoverRequestV2Type, body)
}

export const sendMessage = (url: string, request: FcrMessage): FcrMessage => {
  return {} as FcrMessage
}
