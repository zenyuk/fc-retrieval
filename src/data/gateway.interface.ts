export interface GatewayRegister {
  nodeID: string
  address: string
  rootSigningKey: string
  signingKey: string
  regionCode: string
  networkInfoGateway: string
  networkInfoProvider: string
  networkInfoClient: string
  networkInfoAdmin: string
}

export interface GatewaysToUse {
  [index: string]: GatewayRegister
}
