import { Settings } from './config/settings.interface'
import { FCRPaymentMgr } from './fcrPaymentMgr/payment-manager.class'
import { GatewaysToUse } from './gateway/gateway.interface'
import { ContentID } from './cid/cid.interface'
import { NodeID } from './nodeid/nodeid.interface'
import { SubCIDOffer } from './cidoffer/subcidoffer.class'
import { getGatewayByID, getProviderByID } from './register/register.class'
import { requestStandardDiscoverOffer } from './clientapi/standard_discover_offer_requester'
import { requestStandardDiscoverV2 } from './clientapi/standard_discover_requester_v2'
import { GatewayRegister } from './register/register.class'
import { decodeGatewayDHTDiscoverResponseV2, requestDHTDiscoverV2 } from './clientapi/find_offers_dht_discovery_v2'

export interface payResponse {
  paychAddrs: string
  voucher: string
  topup: boolean
  subCIDOffers: SubCIDOffer[]
}

export class FilecoinRetrievalClient {
  settings: Settings
  activeGateways: GatewaysToUse
  gatewaysToUse: GatewaysToUse
  paymentMgr: FCRPaymentMgr

  constructor(settings: Settings) {
    this.settings = settings
    this.activeGateways = {} as GatewaysToUse
    this.gatewaysToUse = {} as GatewaysToUse
    this.paymentMgr = {} as FCRPaymentMgr
  }

  // AddActiveGateways adds one or more gateways to active gateway map.
  // Returns the number of gateways added.
  AddGatewaysToUse(): number {
    //
    return 42
  }

  // AddActiveGateways adds one or more gateways to active gateway map.
  // Returns the number of gateways added.
  AddActiveGateways(gatewayIDs: NodeID[]): number {
    //
    return 42
  }

  findOffersDHTDiscoveryV2(
    contentID: ContentID,
    gatewayID: NodeID,
    numDHT: number,
    offersNumberLimit: number,
  ): SubCIDOffer[] {
    const gw = this.activeGateways[gatewayID.id]

    const defaultPaymentLane = 0
    const initialRequestPaymentAmount = numDHT * this.settings.searchPrice
    let payResult = this.paymentMgr.pay(gw.address, defaultPaymentLane, initialRequestPaymentAmount)

    if (payResult.topup) {
      this.paymentMgr.topup(gw.address, this.settings.topUpAmount)

      payResult = this.paymentMgr.pay(gw.address, defaultPaymentLane, initialRequestPaymentAmount)
      if (payResult.topup) {
        // Unable to make payment for initial DHT offers discovery
        return [] as SubCIDOffer[]
      }
    }

    const nonce = 0
    const ttl = 0
    const request = requestDHTDiscoverV2(
      gw,
      contentID,
      nonce,
      ttl,
      numDHT,
      false,
      payResult.paychAddrs,
      payResult.voucher,
    )
    let addedSubOffersCount = 0
    let offersDigestsFromAllGateways: string[][]

    for (let i = 0; i < request.contactedResp.length; i++) {
      const contactedGatewayID = request.contactedGateways[i]
      const resp = request.contactedResp[i]
      const gatewayInfo = getGatewayByID(this.settings.registerURL, contactedGatewayID.id)
      if (!this.validateGatewayInfo(gatewayInfo)) {
        // logging.Error("Gateway register info not valid.")
        continue
      }
      const pubKey = gatewayInfo.getSigningKey()
      if (pubKey == undefined) {
        //logging.Error('Fail to obtain public key.')
        continue
      }
      if (!resp.verify(pubKey)) {
        //logging.Error('Fail to verify sub response.')
        continue
      }

      const decoded = decodeGatewayDHTDiscoverResponseV2(resp)
      if (decoded === undefined) {
        // logging.Error('Fail to decode response')
        continue
      }
      if (!decoded.found) {
        return [] as SubCIDOffer[]
      }
    }

    const offers = [] as SubCIDOffer[]

    return offers
  }

  validateGatewayInfo = (gatewayInfo: GatewayRegister): boolean => {
    return false
  }

  // FindOffersStandardDiscoveryV2 finds offer using standard discovery from given gateways
  findOffersStandardDiscoveryV2(cid: ContentID, gatewayID: NodeID, maxOffers: number) {
    const gw = this.activeGateways[gatewayID.id]

    const payResponse = this.pay(gw)

    if (payResponse.topup == true) {
      this.paymentMgr.topup(gw.nodeID, this.settings.topUpAmount)
      const payResponse = this.pay(gw)
    }

    const offerDigests = requestStandardDiscoverV2(
      gw,
      cid,
      Math.floor(Math.random() * 1000),
      Date.now() + this.settings.establishmentTTL,
      payResponse.paychAddrs,
      payResponse.voucher,
    )

    const offers = requestStandardDiscoverOffer(
      gw,
      cid,
      Math.floor(Math.random() * 1000),
      Date.now(),
      offerDigests,
      payResponse.paychAddrs,
      payResponse.voucher,
    )

    const validOffers = [] as SubCIDOffer[]
    for (const offer of offers) {
      const providerInfo = getProviderByID(this.settings.registerURL, offer.getProviderID())
      const pubKey = providerInfo.signingKey
      if (offer.verify(pubKey) != null) {
        // console.log('Offer signature fail to verify.')
        continue
      }
      if (offer.verifyMerkleProof() != null) {
        // console.log('Merkle proof verification failed.')
        continue
      }

      validOffers.push(offer)
      if (validOffers.length >= maxOffers) {
        break
      }
    }
    return validOffers
  }

  pay(gateway: GatewayRegister): payResponse {
    //
    return {} as payResponse
  }
}
