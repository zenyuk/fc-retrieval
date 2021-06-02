import { defaults, Settings } from './constants/defaults'
import { FCRPaymentMgr } from './interfaces/fcrPaymentMgr.interfaces'
import { GatewaysToUse } from './interfaces/gateway.interface'

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

export class FilecoinRetrievalClient {
  settings: Settings
  gatewaysToUse: GatewaysToUse
  paymentMgrs: FCRPaymentMgr[]

  constructor(settings: Settings) {
    this.settings = settings
    this.gatewaysToUse = {} as GatewaysToUse
    this.paymentMgrs = [] as FCRPaymentMgr[]
  }

  // FindOffersStandardDiscoveryV2 finds offer using standard discovery from given gateways
  findOffersStandardDiscoveryV2() {
    return ['hello']
  }
}
