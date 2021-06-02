import { defaults, Settings } from './defaults'
import { FCRPaymentMgrInterface } from './data/fcrPaymentMgr.interfaces'
import { GatewaysToUseInterface } from './data/gateway.interface'

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
  paymentMgrs: FCRPaymentMgrInterface[]

  constructor(settings: Settings) {
    this.settings = settings
    this.gatewaysToUse = {} as GatewaysToUseInterface
    this.paymentMgrs = [] as FCRPaymentMgrInterface[]
  }

  // FindOffersStandardDiscoveryV2 finds offer using standard discovery from given gateways
  findOffersStandardDiscoveryV2() {
    return ['hello']
  }
}
