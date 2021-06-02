// GatewayRegister stores information of a registered gateway
export type GatewayRegisterType = {
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

export const getGatewayByID = (nodeID: string) => {
  return {
    nodeID: '9876543210',
    address: 'f01234',
    networkInfoAdmin: '127.0.0.1:80',
    networkInfoClient: '127.0.0.1:80',
    networkInfoGateway: '127.0.0.1:80',
    networkInfoProvider: '127.0.0.1:80',
    regionCode: 'FR',
    rootSigningKey: '0xABCDE123456789',
    signingKey: '0x987654321EDCBA',
  }
}
