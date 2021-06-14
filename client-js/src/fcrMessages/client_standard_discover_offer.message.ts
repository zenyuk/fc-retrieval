import { ContentID } from '../cid/cid.interface';
import { SubCIDOffer } from '../cidoffer/subcidoffer.class';
import { FCRMessage } from './fcrMessage.class';
import { FCRMessageType } from './type.enum';

export interface ClientStandardDiscoverOfferRequest {
	pieceCID: ContentID
  nonce: number
	ttl: number
  offerDigests: string[]
	paychAddr: string
	voucher: string
}

export interface ClientStandardDiscoverOfferResponse {
	pieceCID: ContentID
	nonce: number
	found: boolean
	subCIDOffers: SubCIDOffer[]
	fundedPaymentChannel: boolean[]
}

export const encodeClientStandardDiscoverOfferRequest = (
	pieceCID: ContentID,
	nonce: number,
	ttl: number,
	offerDigests: string[],
	paychAddr: string,
	voucher: string,
): FCRMessage => {
	const body = JSON.stringify({
		pieceCID,
		nonce,
		ttl,
		offerDigests,
		paychAddr,
		voucher,
	})
	return new FCRMessage(FCRMessageType.ClientStandardDiscoverOfferRequestType, body)
}

export const decodeClientStandardDiscoverOfferResponse = (fcrMsg: FCRMessage): ClientStandardDiscoverOfferResponse => {
	if (fcrMsg.messageType !== FCRMessageType.ClientStandardDiscoverOfferResponseType) {
    throw Error("Message type mismatch")
	}
	const msg = JSON.parse(fcrMsg.messageBody)
  // TODO: handle errors
	// if err != nil {
	// 	return nil, 0, false, nil, nil, err
	// }
	return {
    pieceCID: msg.pieceCID,
    nonce: msg.nonce,
    found: msg.found,
    subCIDOffers: msg.subCIDOffers,
    fundedPaymentChannel: msg.fundedPaymentChannel,
  }
}