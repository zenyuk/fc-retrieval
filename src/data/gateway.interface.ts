export interface GatewayRegisterInterface {
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

export interface GatewaysToUseInterface {
  [index: string]: GatewayRegisterInterface
}
