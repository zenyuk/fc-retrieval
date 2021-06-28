import { FilecoinRPC } from '../filecoin/filecoin-rpc.class'
import BN from 'bn.js'

const filecoin_signer = require('@zondax/filecoin-signing-tools')

export class LaneState {
  nonce: number
  redeemed: BN
  vouchers: string[]

  constructor() {
    this.nonce = 0
    this.redeemed = new BN(0)
    this.vouchers = new Array()
  }
}

export class ChannelState {
  addr: string
  balance: BN
  redeemed: BN
  laneStates: Map<number, LaneState>

  constructor(addr: string, balance: BN) {
    this.addr = addr
    this.balance = balance
    this.redeemed = new BN(0)
    this.laneStates = new Map<number, LaneState>()
  }
}

export class FCRPaymentMgr {
  recoveredKey: any
  filRPC: any
  header: any

  outboundChs: Map<string, ChannelState>

  // Make sure the private key provided is in hex string
  constructor(privateKeyHex: string, lotusAPIAddr: string, authToken: string) {
    this.recoveredKey = filecoin_signer.keyRecover(Buffer.from(privateKeyHex, 'hex').toString('base64'))
    this.filRPC = new FilecoinRPC({ url: lotusAPIAddr, token: authToken })
    this.header = { Authorization: `Bearer ${authToken}` }

    this.outboundChs = new Map<string, ChannelState>()
  }

  async topup(recipient: string, amount: BN) {
    if (this.outboundChs.has(recipient)) {
      this.hasRecipient(recipient, amount)
    } else {
      this.hasNoRecipient(recipient, amount)
    }
    return 0
  }

  /**
   * Has recipient
   * @param recipient
   * @param amount
   * @returns
   */
  async hasRecipient(recipient: string, amount: BN) {
    const nonceRes = await this.filRPC.getNonce(this.recoveredKey.address)
    const nonce = nonceRes.result
    const topupQuery = {
      from: this.recoveredKey.address,
      to: this.outboundChs.get(recipient)?.addr,
      nonce: nonce,
      value: amount.toString(10),
      gaslimit: 0,
      gasfeecap: '0',
      gaspremium: '0',
      Method: 0,
      Params: '',
    } as any
    const topupRes = await this.filRPC.getGasEstimation(topupQuery)

    if (!topupRes.result) {
      throw new Error('Error in estimating gas cost')
    }
    const topup = topupRes.result
    const signedMessage = JSON.parse(filecoin_signer.transactionSignLotus(topup, this.recoveredKey.private_base64))
    await this.filRPC.sendSignedMessage(signedMessage)
    const cs = this.outboundChs.get(recipient)
    if (cs != undefined) {
      cs.balance = cs.balance.add(amount)
    } else {
      throw new Error('Internal error')
    }
  }

  /**
   * Has no recipient
   * @param recipient
   * @param amount
   * @returns
   */
  async hasNoRecipient(recipient: string, amount: BN) {
    // Need to create a channel
    // Get nonce
    const nonceRes = await this.filRPC.getNonce(this.recoveredKey.address)
    if (!nonceRes.result) {
      throw new Error('Error in getting nonce')
    }
    const nonce = nonceRes.result
    let create_channel = filecoin_signer.createPymtChan(
      this.recoveredKey.address,
      recipient,
      amount.toString(),
      nonce,
      '0',
      '0',
      '0',
    )
    create_channel = await this.filRPC.getGasEstimation(create_channel)
    if ('result' in create_channel) {
      create_channel = create_channel.result
    } else {
      console.log('Error in estimating gas cost')
      return 1
    }
    const signedMessage = JSON.parse(
      filecoin_signer.transactionSignLotus(create_channel, this.recoveredKey.private_base64),
    )
    // Send message
    const res = await this.filRPC.sendSignedMessage(signedMessage)
    if (res.result.ReturnDec.IDAddress == undefined) {
      throw new Error('Error in creating payment channel')
      return 1
    }

    // Success, add a new entry to the out bound channels
    const cAddr = String(res.result.ReturnDec.RobustAddress)
    this.outboundChs.set(recipient, new ChannelState(cAddr, amount))
  }

  pay(recipient: string, lane: number, amount: BN) {
    const cs = this.outboundChs.get(recipient)
    if (cs == null) {
      // Channel not existed
      return { paychAddr: '', voucher: '', topup: true, error: 0 } as any
    }

    // Check if balance is enough
    if (cs.balance.lt(cs.redeemed.add(amount))) {
      // Balance not enough
      return { paychAddr: '', voucher: '', topup: true, error: 0 } as any
    }

    // Balance enough
    let ls = cs.laneStates.get(lane)
    if (ls == null) {
      // Need to create a lane state
      ls = new LaneState()
      cs.laneStates.set(lane, ls)
    }

    // Create a voucher
    let voucher = filecoin_signer.createVoucher(
      cs.addr,
      '0',
      '0',
      ls.redeemed.add(amount).toString(),
      lane.toString(),
      ls.nonce.toString(),
      '0',
    )
    voucher = filecoin_signer.signVoucher(voucher, this.recoveredKey.private_base64)
    if (voucher == null) {
      // Fail to generate a voucher
      return { paychAddr: '', voucher: '', topup: false, error: 1 } as any
    }
    // This fixes a bug from the original library
    voucher = voucher.replace(/\+/g, '-').replace(/\//g, '_').replace(/\=+$/, '')
    // Update lane state and channel state
    ls.nonce++
    ls.redeemed = ls.redeemed.add(amount)
    ls.vouchers.push(voucher)
    cs.redeemed = cs.redeemed.add(amount)
    return { paychAddr: cs.addr, voucher: voucher, topup: false, error: 0 } as any
  }
}
