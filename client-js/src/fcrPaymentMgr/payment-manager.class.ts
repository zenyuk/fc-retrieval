import { Address } from './types'
import BN from 'bn.js'

export interface LaneState {
  nonce: number
  redeemed: BN
  vouchers: string[]
}

export interface ChannelState {
  addr: Address
  balance: BN
  redeemed: BN
  laneStates: Map<number, LaneState>
}

export class FCRPaymentMgr {
  privKey: string
  address: Address
  authToken: string
  lotusAPIAddr: string
  outboundChs: Map<string, ChannelState>
  inboundChs: Map<string, ChannelState>

  constructor(privateKey: string, lotusAPIAddr: string, authToken: string) {
    // TODO
    this.privKey = ''
    this.address = {} as Address
    this.authToken = ''
    this.lotusAPIAddr = ''
    this.outboundChs = {} as Map<string, ChannelState>
    this.inboundChs = {} as Map<string, ChannelState>
  }

  topup(recipient: string, amount: string) {
    // TODO
  }
}
