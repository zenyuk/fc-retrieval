import axios from 'axios'

import { defaults, Settings } from './defaults'
import { FCRPaymentMgr } from './data/types'
import { GatewayRegister, getGatewayByID } from './Register'

export interface Gateway {
  nodeID: string
  address: string
  networkInfoAdmin: string
  networkInfoClient: string
  networkInfoGateway: string
  networkInfoProvider: string
  regionCode: string
  rootSigningKey: string
  signingKey: string
}

export class CreateSettings {
  settings: Settings

  constructor() {
    this.settings = defaults as Settings
  }

  setDefaultEstablishmentTTL(ttl: number) {
    this.settings.defaultEstablishmentTTL = ttl
  }

  setDefaultLogLevel(logLevel: string) {
    this.settings.defaultLogLevel = logLevel
  }

  setDefaultLogTarget(logTarget: string) {
    this.settings.defaultLogTarget = logTarget
  }

  setDefaultLogServiceName(serviceName: string) {
    this.settings.defaultLogServiceName = serviceName
  }

  setDefaultRegisterURL(url: string) {
    this.settings.defaultRegisterURL = url
  }

  build(): Settings {
    return this.settings
  }
}

export interface GatewaysToUse {
  [index: string]: GatewayRegister
}

export class Client {
  settings: Settings
  gatewaysToUse: GatewaysToUse
  paymentMgrLock: FCRPaymentMgr[]

  constructor(settings: Settings) {
    this.settings = settings
    this.gatewaysToUse = {} as GatewaysToUse
    this.paymentMgrLock = [] as FCRPaymentMgr[]
  }

  paymentMgr() {}

  // FindGateways find gateways located near to the specified location. Use AddGateways
  // to use these gateways.
  async findGateways(location: string = '', maxNumToLocate: number = 16) {
    const url = this.settings.defaultRegisterURL + '/registers/gateway'
    try {
      const response = await axios.get(url)
      const gateways = response.data as GatewayRegister[]
      return gateways
    } catch (error) {
      throw new Error('Fail to fetch data: ' + url)
    }
  }

  // AddGatewaysToUse adds one or more gateways to use.
  addGatewaysToUse(gwNodeIDs: string[]) {
    const numAdded: number = 0
    for (const nodeID of gwNodeIDs) {
      const gateway = getGatewayByID(nodeID)
      this.gatewaysToUse[nodeID] = gateway
    }
  }

  // RemoveGatewaysToUse removes one or more gateways from the list of Gateways to use.
  // This also removes the gateway from gateways in active map.
  removeGatewaysToUse() {}

  // RemoveAllGatewaysToUse removes all gateways from the list of Gateways.
  // This also cleared all gateways in active
  removeAllGatewaysToUse() {}

  // GetGatewaysToUse returns the list of gateways to use.
  getGatewaysToUse() {
    return this.gatewaysToUse
  }

  // AddActiveGateways adds one or more gateways to active gateway map.
  // Returns the number of gateways added.
  addActiveGateways() {}

  // RemoveActiveGateways removes one or more gateways from the list of Gateways in active.
  removeActiveGateways() {}

  // RemoveAllActiveGateways removes all gateways from the list of Gateways in active.
  removeAllActiveGateways() {}

  // GetActiveGateways returns the list of gateways that are active.
  getActiveGateways() {}

  // FindOffersStandardDiscovery finds offer using standard discovery from given gateways
  findOffersStandardDiscovery() {}

  // FindOffersDHTDiscovery finds offer using dht discovery from given gateways
  findOffersDHTDiscovery() {}

  // FindDHTOfferAck finds offer ack for a cid, gateway pair
  findDHTOfferAck() {}

  // FindOffersStandardDiscoveryV2 finds offer using standard discovery from given gateways
  findOffersStandardDiscoveryV2() {}

  private pay() {}
}
