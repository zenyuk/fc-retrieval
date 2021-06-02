import axios from 'axios'

import { defaults, Settings } from './defaults'
import { FCRPaymentMgrType } from './data/types'
import { GatewaysToUseInterface } from './data/interfaces'
import { GatewayRegisterType, getGatewayByID } from './Register'

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

export class Client {
  settings: Settings
  gatewaysToUse: GatewaysToUseInterface
  paymentMgrs: FCRPaymentMgrType[]

  constructor(settings: Settings) {
    this.settings = settings
    this.gatewaysToUse = {} as GatewaysToUseInterface
    this.paymentMgrs = [] as FCRPaymentMgrType[]
  }

  // FindGateways find gateways located near to the specified location. Use AddGateways
  // to use these gateways.
  async findGateways(location: string = '', maxNumToLocate: number = 16) {
    const url = this.settings.defaultRegisterURL + '/registers/gateway'
    try {
      const response = await axios.get(url)
      const gateways = response.data as GatewayRegisterType[]
      return gateways
    } catch (error) {
      throw new Error('Fail to fetch data: ' + url)
    }
  }

  // AddGatewaysToUseInterface adds one or more gateways to use.
  addGatewaysToUseInterface(gwNodeIDs: string[]) {
    const numAdded: number = 0
    for (const nodeID of gwNodeIDs) {
      const gateway = getGatewayByID(nodeID)
      this.gatewaysToUse[nodeID] = gateway
    }
  }

  // FindOffersStandardDiscoveryV2 finds offer using standard discovery from given gateways
  findOffersStandardDiscoveryV2() {
    return ['hello']
  }
}
