import { ContentID } from '../cid/cid.interface'
import { NodeID } from '../nodeid/nodeid.interface'
import { CIDOffer } from '../cidoffer/cidoffer.class'
import { FCRMessage } from './fcrMessage.class'

// // providerPublishDHTOfferRequest is the request from provider to gateway to publish dht offer
// type providerPublishDHTOfferRequest struct {
// 	ProviderID nodeid.NodeID       `json:"provider_id"`
// 	Nonce      int64               `json:"nonce"`
// 	NumOffers  int64               `json:"num_of_offers"`
// 	Offers     []cidoffer.CIDOffer `json:"single_offers"`
// }

export interface ProviderPublishDHTOfferRequest {
  providerID: NodeID // provider id
  nonce: number // nonce
  offers: CIDOffer[] // offers
}

export interface ProviderPublishDHTOfferResponse {
  nonce: number
  signature: string
}

// EncodeProviderPublishDHTOfferRequest is used to get the FCRMessage of providerPublishDHTOfferRequest
export const encodeProviderPublishDHTOfferRequest = (
  providerID: NodeID,
  nonce: number,
  offers: CIDOffer[],
): FCRMessage => {
  return {} as FCRMessage
}

// DecodeProviderPublishDHTOfferRequest is used to get the fields from FCRMessage of providerPublishDHTOfferRequest
export const decodeProviderPublishDHTOfferRequest = (fcrMsg: FCRMessage): ProviderPublishDHTOfferRequest => {
  return {} as ProviderPublishDHTOfferRequest
}

// EncodeProviderPublishDHTOfferRequest is used to get the FCRMessage of providerPublishDHTOfferRequest
export const encodeProviderPublishDHTOfferResponse = (nonce: number, signature: string): FCRMessage => {
  return {} as FCRMessage
}

// DecodeProviderPublishDHTOfferRequest is used to get the fields from FCRMessage of providerPublishDHTOfferRequest
export const decodeProviderPublishDHTOfferResponse = (fcrMsg: FCRMessage): ProviderPublishDHTOfferResponse => {
  return {} as ProviderPublishDHTOfferResponse
}
