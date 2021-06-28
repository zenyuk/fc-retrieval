import { Address } from './types'
import { payResponse } from '../FilecoinRetrievalClient'
import BN from 'bn.js'

export interface LaneState {
  nonce: number
  redeemed: string
  vouchers: string[]
}

export interface ChannelState {
  addr: Address
  balance: string
  redeemed: string
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

  topup(recipient: string, amount: BN) {
    // TODO
  }

  pay(address: string, defaultPaymentLane: BN, initialRequestPaymentAmount: BN) {
    return { paychAddrs: '', voucher: '', topup: false } as payResponse
  }
}
