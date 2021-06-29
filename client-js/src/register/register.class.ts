import { KeyPair } from '../fcrcrypto/key_pair.class'
import { NodeID } from '../nodeid/nodeid.interface'
import { decodePublicKey } from '../fcrcrypto/msg_signing'

export abstract class Register {
  nodeId: string
  address: string
  regionCode: string
  signingKey: string
  rootSigningKey: string

  constructor({ nodeId, address, regionCode, signingKey, rootSigningKey }: any) {
    this.nodeId = nodeId
    this.address = address
    this.regionCode = regionCode
    this.signingKey = signingKey
    this.rootSigningKey = rootSigningKey
  }

  /**
   * Get node ID
   *
   * @returns {NodeID}
   */
  getNodeID(): NodeID {
    return new NodeID(this.nodeId)
  }

  /**
   * Get signing key pair
   *
   * @returns {KeyPair}
   */
  getRootSigningKeyPair(): KeyPair {
    return {} as KeyPair
  }

  /**
   * Get root signing key pair
   *
   * @returns {KeyPair}
   */
  getSigningKeyPair(): KeyPair {
    return decodePublicKey(this.signingKey)
  }

  /**
   * Validate registration information
   *
   * @returns {boolean}
   */
  public validateInfo(): boolean {
    // for (const property of Object.getOwnPropertyNames(this)) {
    //   const value = (this as any)[property]
    //   if (!value || value.trim().length === 0) {
    //     throw Error(`Registration issue: ${property} not set`)
    //   }
    // }
    if (!this.getSigningKeyPair()) {
      throw Error(`Registration issue: Root Signing Public Key error`)
    }
    return true
  }
}

export class ProviderRegister extends Register {
  networkInfoGateway: string
  networkInfoClient: string
  networkInfoAdmin: string

  constructor({
    nodeId,
    address,
    regionCode,
    signingKey,
    rootSigningKey,
    networkInfoGateway,
    networkInfoClient,
    networkInfoAdmin,
  }: any) {
    super({
      nodeId,
      address,
      regionCode,
      signingKey,
      rootSigningKey,
    })
    this.networkInfoGateway = networkInfoGateway
    this.networkInfoClient = networkInfoClient
    this.networkInfoAdmin = networkInfoAdmin
  }
}

export class GatewayRegister extends Register {
  networkInfoGateway: string
  networkInfoProvider: string
  networkInfoClient: string
  networkInfoAdmin: string

  constructor({
    nodeId,
    address,
    regionCode,
    signingKey,
    rootSigningKey,
    networkInfoGateway,
    networkInfoProvider,
    networkInfoClient,
    networkInfoAdmin,
  }: any) {
    super({
      nodeId,
      address,
      regionCode,
      signingKey,
      rootSigningKey,
    })
    this.networkInfoGateway = networkInfoGateway
    this.networkInfoProvider = networkInfoProvider
    this.networkInfoClient = networkInfoClient
    this.networkInfoAdmin = networkInfoAdmin
  }
}
