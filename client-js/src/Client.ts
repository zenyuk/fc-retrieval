import { defaults, Settings } from './defaults'

export interface Gateway {
  address: string
  networkInfoAdmin: string
  networkInfoClient: string
  networkInfoGateway: string
  networkInfoProvider: string
  nodeId: string
  regionCode: string
  rootSigningKey: string
  sigingKey: string
}

export class CreateSettings {
  settings: Settings

  constructor() {
    this.settings = defaults
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
  gateways: Gateway[]

  constructor(settings: Settings) {
    this.settings = settings
    this.gateways = []
  }

  async findGateways(location: string = '', maxNumToLocate: number = 16) {
    try {
      const response = await fetch(this.settings.defaultRegisterURL + '/registers/gateway')
      this.gateways = await response.json()
      return this.gateways
    } catch (error) {
      throw new Error('Fail to fetch data')
    }
  }

  async connectedGateways() {
    return this.gateways
  }
}
