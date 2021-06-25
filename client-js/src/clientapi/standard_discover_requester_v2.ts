import { ContentID } from '../cid/cid.interface'
import { FCRMessage } from '../fcrMessages/fcrMessage.class'
import { FCRMessageType } from '../fcrMessages/type.enum'
import { GatewayRegister } from '../register/register.class'
import { sendMessage } from '../request/request'

// RequestStandardDiscoverV2 requests a standard discover to a given gateway for a given contentID, nonce and ttl.
export const requestStandardDiscoverV2 = async (
  gatewayInfo: GatewayRegister,
  contentID: ContentID,
  nonce: number,
  ttl: number,
  paychAddr: string,
  voucher: string,
): Promise<string[]> => {
  const request = encodeClientStandardDiscoverRequestV2(contentID, nonce, ttl, paychAddr, voucher)

  const response = await sendMessage(gatewayInfo.networkInfoClient, request)

  const pubKey = gatewayInfo.getRootSigningKeyPair()

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
  response: FCRMessage,
): DecodedClientStandardDiscoverResponseV2 => {
  return {} as DecodedClientStandardDiscoverResponseV2
}

export const encodeClientStandardDiscoverRequestV2 = (
  contentID: ContentID,
  nonce: number,
  ttl: number,
  paychAddr: string,
  voucher: string,
): FCRMessage => {
  let body = ''

  return new FCRMessage({
    message_type: FCRMessageType.ClientStandardDiscoverRequestV2Type,
    message_body: Buffer.from(body).toString('base64'),
  })
}
