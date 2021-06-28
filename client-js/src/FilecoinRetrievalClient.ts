import { Settings } from './config/settings.config'
import { FCRPaymentMgr } from './fcrPaymentMgr/payment-manager.class'
import { ContentID } from './cid/cid.interface'
import { NodeID } from './nodeid/nodeid.interface'
import { SubCIDOffer } from './cidoffer/subcidoffer.class'
import { getGatewayByID, getProviderByID } from './register/register.service'
import { requestStandardDiscoverOffer } from './clientapi/standard_discover_offer_requester'
import { requestStandardDiscoverV2 } from './clientapi/standard_discover_requester_v2'
import { GatewayRegister } from './register/register.class'
import { requestDHTOfferAck } from './clientapi/dht_offer_ack_requester'
import { decodeProviderPublishDHTOfferResponse } from './fcrMessages/provider_publish_dht_offer'
import { decodeGatewayDHTDiscoverResponseV2, requestDHTDiscoverV2 } from './clientapi/find_offers_dht_discovery_v2'
import { requestDHTOfferDiscover } from './clientapi/request_dht_offer_discover'
import BN from 'bn.js'
import crypto from 'crypto'
import { requestEstablishment } from './clientapi/establishment_requester'
import { verifyAnyMessage } from './fcrcrypto/msg_signing'

export interface payResponse {
  paychAddrs: string
  voucher: string
  topup: boolean
  subCIDOffers: SubCIDOffer[]
}

export class FilecoinRetrievalClient {
  settings: Settings
  activeGateways: Map<string, GatewayRegister>
  gatewaysToUse: Map<string, GatewayRegister>
  paymentMgr: FCRPaymentMgr

  constructor(settings: Settings) {
    this.settings = Object.assign({}, settings)
    this.activeGateways = new Map()
    this.gatewaysToUse = new Map()
    this.paymentMgr = {} as FCRPaymentMgr
  }

  /**
   * Add one or more gateways to gateways to use map
   * Returns: the number of gateways added
   *
   * @param {NodeID[]} gatewayIDs
   * @returns {Promise<number>}
   */
  async addGatewaysToUse(gatewayIDs: NodeID[]): Promise<number> {
    let numAdded = 0
    for (const gatewayID of gatewayIDs) {
      const _gatewayID = gatewayID.toString()
      if (this.gatewaysToUse.has(_gatewayID)) {
        continue
      }
      try {
        const gateway = await getGatewayByID(this.settings.registerURL, _gatewayID)
        gateway.validateInfo()
        this.gatewaysToUse.set(_gatewayID, gateway)
        numAdded++
      } catch (e) {
        console.error(`Add gateways to use failed for gatewayID=${_gatewayID}:`, e)
        continue
      }
    }
    return numAdded
  }

  /**
   * Add one or more gateways to active gateways map
   * Returns: the number of gateways added
   *
   * @param {NodeID[]} gatewayIDs
   * @returns {Promise<number>}
   */
  async addActiveGateways(gatewayIDs: NodeID[]): Promise<number> {
    let numAdded = 0
    for (const gatewayID of gatewayIDs) {
      const _gatewayID = gatewayID.toString()
      if (this.activeGateways.has(_gatewayID)) {
        continue
      }
      const gatewayInfo = this.gatewaysToUse.get(_gatewayID)
      if (!gatewayInfo) {
        console.log(`gatewayID=${_gatewayID} does not exist in gateways to use. Consider add the gateway first`)
        continue
      }
      try {
        const challenge = crypto.randomBytes(32)
        const ttl = new Date().getTime() + this.settings.establishmentTTL
        const done = await requestEstablishment(gatewayInfo, challenge, this.settings.clientId, ttl)
        if (!done) {
          console.log(`Error in initial establishment: gatewayID=${_gatewayID}`)
          continue
        }
        this.activeGateways.set(_gatewayID, gatewayInfo)
        numAdded++
      } catch (e) {
        console.log(`Add active gateways failed for gatewayID=${_gatewayID}:`, e)
        continue
      }
    }
    return numAdded
  }

  async findOffersDHTDiscoveryV2(
    contentID: ContentID,
    gatewayID: NodeID,
    numDHT: number,
    offersNumberLimit: number,
  ): Promise<Map<string, SubCIDOffer[]>> {
    const offersMap = new Map<string, SubCIDOffer[]>()

    const gw = this.activeGateways.get(gatewayID.id)
    if (!gw) {
      return offersMap
    }

    const defaultPaymentLane = new BN(0)
    const initialRequestPaymentAmount = new BN(numDHT).mul(this.settings.searchPrice)
    let payResponse = this.paymentMgr.pay(gw.address, defaultPaymentLane, initialRequestPaymentAmount)

    if (payResponse.topup) {
      this.paymentMgr.topup(gw.address, this.settings.topUpAmount)

      payResponse = this.paymentMgr.pay(gw.address, defaultPaymentLane, initialRequestPaymentAmount)
      if (payResponse.topup) {
        // Unable to make payment for initial DHT offers discovery
        return offersMap
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
      payResponse.paychAddrs,
      payResponse.voucher,
    )
    let addedSubOffersCount = 0
    const offersDigestsFromAllGateways: string[][] = [[]]

    for (let i = 0; i < request.contactedResp.length; i++) {
      const contactedGatewayID = request.contactedGateways[i]
      const resp = request.contactedResp[i]
      const gatewayInfo = await getGatewayByID(this.settings.registerURL, contactedGatewayID.toString())
      if (!this.validateGatewayInfo(gatewayInfo)) {
        // logging.Error("Gateway register info not valid.")
        continue
      }
      const pubKey = gatewayInfo.getSigningKeyPair()
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
        return offersMap
      }
      // comply with given offers number limit
      if (addedSubOffersCount + decoded.subCidOffersDigest.length > offersNumberLimit) {
        offersDigestsFromAllGateways.push(decoded.subCidOffersDigest.slice(0, offersNumberLimit - addedSubOffersCount))
      } else {
        offersDigestsFromAllGateways.push(decoded.subCidOffersDigest)
      }
      addedSubOffersCount += decoded.subCidOffersDigest.length
      if (addedSubOffersCount >= offersNumberLimit) {
        break
      }
    }
    /*
    let unit = 0;
    for (let entry in offersDigestsFromAllGateways) {
      unit += entry.length;
    }
     const offerRequestPaymentAmount = this.settings.offerPrice.mul(new BN(unit))
    */
    payResponse = this.paymentMgr.pay(gw.address, defaultPaymentLane, initialRequestPaymentAmount)

    if (payResponse.topup) {
      this.paymentMgr.topup(gw.address, this.settings.topUpAmount)

      payResponse = this.paymentMgr.pay(gw.address, defaultPaymentLane, initialRequestPaymentAmount)
      if (payResponse.topup) {
        // Unable to make payment for initial DHT offers discovery
        return offersMap
      }
    }

    const gatewaySubOffers = requestDHTOfferDiscover(
      gw,
      request.contactedGateways,
      contentID,
      nonce,
      offersDigestsFromAllGateways,
      payResponse.paychAddrs,
      payResponse.voucher,
    )

    for (const entry of gatewaySubOffers) {
      offersMap.set(entry.gatewayID.id, entry.subOffers)
    }

    return offersMap
  }

  validateGatewayInfo = (gatewayInfo: GatewayRegister): boolean => {
    return false
  }

  async FindDHTOfferAck(contentID: ContentID, gatewayID: NodeID, providerID: NodeID): Promise<boolean> {
    const provider = await getProviderByID(this.settings.registerURL, providerID.id)

    const pvalidation = provider.validateInfo()
    if (!pvalidation) {
      throw new Error('Invalid register info')
    }

    const dhtOfferAckResponse = await requestDHTOfferAck(provider, contentID, gatewayID)
    if (!dhtOfferAckResponse.found) {
      return false
    }

    const gateway = await getGatewayByID(this.settings.registerURL, gatewayID.toString())

    const gvalidation = gateway.validateInfo()
    if (!gvalidation) {
      throw new Error('Invalid register info')
    }

    const gwPubKey = gateway.getSigningKeyPair()
    const pvdPubKey = provider.getSigningKeyPair()

    dhtOfferAckResponse.offerRequest.verify(pvdPubKey)

    // const offers = decodeProviderPublishDHTOfferRequest(dhtOfferAckResponse.offerRequest);
    const found = false
    // for (const offer in offers) {
    //   // ?
    // }
    if (!found) {
      throw new Error('Initial request does not contain the given cid')
    }

    const verified = dhtOfferAckResponse.offerResponse.verify(pvdPubKey)
    if (!verified) {
      throw new Error('Error in verifying the ack')
    }

    const dhtOfferResponse = decodeProviderPublishDHTOfferResponse(dhtOfferAckResponse.offerResponse)
    verifyAnyMessage(gwPubKey, dhtOfferResponse.signature, dhtOfferAckResponse.offerRequest)

    return true
  }

  // FindOffersStandardDiscoveryV2 finds offer using standard discovery from given gateways
  async findOffersStandardDiscoveryV2(cid: ContentID, gatewayID: NodeID, maxOffers: number) {
    const gw = this.activeGateways.get(gatewayID.id)
    if (!gw) {
      return
    }

    let payResponse = this.paymentMgr.pay(gw.address, new BN(0), this.settings.searchPrice)

    if (payResponse.topup == true) {
      this.paymentMgr.topup(gw.nodeId, this.settings.topUpAmount)
      payResponse = this.paymentMgr.pay(gw.address, new BN(0), this.settings.searchPrice)
    }

    const offerDigests = await requestStandardDiscoverV2(
      gw,
      cid,
      Math.floor(Math.random() * 1000),
      Date.now() + this.settings.establishmentTTL,
      payResponse.paychAddrs,
      payResponse.voucher,
    )

    const offers = await requestStandardDiscoverOffer(
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
      const providerInfo = await getProviderByID(this.settings.registerURL, offer.getProviderID())
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
}
