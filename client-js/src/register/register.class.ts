import { NodeID } from '../nodeid/nodeid.interface'

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

export class ProviderRegister {
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

export const getGatewayByID = (registerURL: string, nodeID: NodeID): ProviderRegister => {
  //
  return {} as ProviderRegister
}

export const getProviderByID = (registerURL: string, nodeID: NodeID): ProviderRegister => {
  //
  return {} as ProviderRegister
}

export const validateProviderInfo = (register: ProviderRegister): boolean => {
  //
  return {} as boolean
}

export const validateGatewayInfo = (register: GatewayRegister): boolean => {
  //
  return {} as boolean
}
