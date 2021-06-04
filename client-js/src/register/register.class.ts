export class Register {}

export class GatewayRegister {
  nodeID: string = ''
  address: string = ''
  rootSigningKey: string = ''
  signingKey: string = ''
  regionCode: string = ''
  networkInfoGateway: string = ''
  networkInfoProvider: string = ''
  networkInfoClient: string = ''
  networkInfoAdmin: string = ''

  getSigningKey(): string {
    return ''
  }
}

export interface ProviderRegister {
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

export const getProviderByID = (registerURL: string, nodeID: string): ProviderRegister => {
  //
  return {} as ProviderRegister
}

export const getGatewayByID = (registerURL: string, nodeID: string): GatewayRegister => {
  //
  return {} as GatewayRegister
}
