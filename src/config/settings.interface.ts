import { NodeID } from '../nodeid/nodeid.interface'

export interface Settings {
  establishmentTTL: number
  registerURL: string
  client: NodeID
  logLevel: string
  logTarget: string
  logServiceName: string
  blockchainPrivateKey: any
  retrievalPrivateKey: any
  retrievalPrivateKeyVer: any
  walletPrivateKey: string
  lotusAP: string
  lotusAuthToken: string
  searchPrice: number
  offerPrice: number
  topUpAmount: string
}
