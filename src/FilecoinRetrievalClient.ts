import { Settings } from './constants/defaults'
import { FCRPaymentMgr } from './fcrPaymentMgr/fcrPaymentMgr.interfaces'
import { GatewaysToUse } from './gateway/gateway.interface'
import { ContentID } from './cid/cid.interface'
import { NodeID } from './nodeid/nodeid.interface'

export interface payResponse {
  paychAddrs: string, 
  voucher: string, 
  topup: boolean, 
  subCIDOffers: 
}

export class FilecoinRetrievalClient {
  settings: Settings
  activeGateways: GatewaysToUse
  gatewaysToUse: GatewaysToUse
  paymentMgrs: FCRPaymentMgr[]

  constructor(settings: Settings) {
    this.settings = settings
    this.activeGateways = {} as GatewaysToUse
    this.gatewaysToUse = {} as GatewaysToUse
    this.paymentMgrs = [] as FCRPaymentMgr[]
  }

  // AddActiveGateways adds one or more gateways to active gateway map.
  // Returns the number of gateways added.
  AddGatewaysToUse(): number {
    return 42
  }

  // AddActiveGateways adds one or more gateways to active gateway map.
  // Returns the number of gateways added.
  AddActiveGateways(gatewayIDs: NodeID[]): number {
    return 42
  }

  // FindOffersStandardDiscoveryV2 finds offer using standard discovery from given gateways
  findOffersStandardDiscoveryV2(cid: ContentID, gatewayID: NodeID, maxOffers: number) {
    const gw = this.activeGateways[gatewayID.id]

    return ['hello']
  }

  pay(gateway: GatewaysToUse) {}
}
