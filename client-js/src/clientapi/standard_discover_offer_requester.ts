import { ContentID } from '../cid/cid.interface'
import { SubCIDOffer } from '../cidoffer/subcidoffer.class'
import {
  decodeClientStandardDiscoverOfferResponse,
  encodeClientStandardDiscoverOfferRequest,
} from '../fcrMessages/client_standard_discover_offer.message'
import { GatewayRegister } from '../register/register.class'
import { sendMessage } from '../request/request'

export const requestStandardDiscoverOffer = async (
  gatewayInfo: GatewayRegister,
  contentID: ContentID,
  nonce: number,
  ttl: number,
  offerDigests: string[],
  paychAddr: string,
  voucher: string,
): Promise<SubCIDOffer[]> => {
  const request = encodeClientStandardDiscoverOfferRequest(contentID, nonce, ttl, offerDigests, paychAddr, voucher)
  // TODO: handle errors
  // if err != nil {
  // 	logging.Error("Error encoding Client Standard Discover Request: %+v", err)
  // 	return nil, err
  // }

  // Send request and get response
  const response = await sendMessage(gatewayInfo.networkInfoClient, request)
  // TODO: handle errors
  // if err != nil {
  // 	return nil, err
  // }

  // Get the gateway's public key
  const pubKey = gatewayInfo.getRootSigningKeyPair()
  // TODO: handle errors
  // if err != nil {
  // 	return nil, err
  // }

  // Verify the response
  if (response.verify(pubKey) === false) {
    throw Error('Verification failed')
  }

  // Decode the response, TODO deal with fundedpayment channels and found
  const { pieceCID: cID, nonce: nonceRecv, subCIDOffers: offers } = decodeClientStandardDiscoverOfferResponse(response)
  // TODO: handle errors
  // if err != nil {
  // 	return nil, err
  // }
  if (cID.toString() !== contentID.toString()) {
    throw Error('CID Mismatch')
  }
  if (nonce !== nonceRecv) {
    throw Error('Nonce mismatch')
  }
  return offers
}
