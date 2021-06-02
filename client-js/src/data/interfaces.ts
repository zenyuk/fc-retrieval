import { GatewayRegisterType } from '../Register'

export interface GatewayInterface {
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

export interface GatewaysToUseInterface {
  [index: string]: GatewayRegisterType
}
