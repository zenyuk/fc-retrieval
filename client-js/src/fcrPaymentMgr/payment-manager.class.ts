import { FilecoinRPC } from './types'
import { payResponse } from '../FilecoinRetrievalClient'
import { BigNumber } from "bignumber.js";

const filecoin_signer = require('@zondax/filecoin-signing-tools')

export class LaneState {
  nonce: number
  redeemed: BigNumber
  vouchers: string[]

  constructor() {
    this.nonce = 0
    this.redeemed = new BigNumber(0)
    this.vouchers = new Array()
  }
}

export class ChannelState {
  addr: string
  balance: BigNumber
  redeemed: BigNumber
  laneStates: Map<number, LaneState>

  constructor(addr: string, balance: BigNumber) {
    this.addr = addr
    this.balance = balance
    this.redeemed = new BigNumber(0)
    this.laneStates = {} as Map<number, LaneState>
  }
}

export class FCRPaymentMgr {
  recoveredKey: any
  filRPC: any
  header: any

  outboundChs: Map<string, ChannelState>

  constructor(privateKey: string, lotusAPIAddr: string, authToken: string) {
    this.recoveredKey = filecoin_signer.keyRecover(privateKey)
    this.filRPC = new FilecoinRPC({url: lotusAPIAddr, token: authToken})
    this.header = { 'Authorization': `Bearer ${authToken}` }
    
    this.outboundChs = {} as Map<string, ChannelState>
  }

  topup(recipient: string, amount: BigNumber) {
    if (recipient in this.outboundChs) {
      // There is an existing channel, TODO. Topup
      var nonce = this.filRPC.getNonce(this.recoveredKey.address)
      nonce = nonce.result
      var topup = {
        from: this.recoveredKey.address,
        to: recipient,
        nonce: nonce,
        value: amount.toString(10),
        GasLimit: '0',
        GasFeeCap: '0',
        GasPremium: '0',
        Method: 0,
        Params: ""
      } as any
      topup = this.filRPC.getGasEstimation(topup)
      if ('result' in topup) {
        topup = topup.result
      } else {
        console.log("Error in estimating gas cost")
        return 1
      }
      var signedMessage = JSON.parse(filecoin_signer.transactionSignLotus(topup, this.recoveredKey.privateKey))// Send message
      var res = this.filRPC.sendSignedMessage(signedMessage)
      // TODO
      console.log(res)
      var cs = this.outboundChs.get(recipient)
      if (cs != undefined) {
        cs.balance = cs.balance.plus(amount)
      } else {
        console.log("Internal error")
        return 1
      }
      return 0
    } else {
      // Need to create a channel
      // Get nonce
      var nonce = this.filRPC.getNonce(this.recoveredKey.address)
      nonce = nonce.result
      let create_channel = filecoin_signer.createPymtChan(this.recoveredKey.address, recipient, amount.toString(), nonce, '0', '0', '0')
      create_channel = this.filRPC.getGasEstimation(create_channel)
      if ('result' in create_channel) {
        create_channel = create_channel.result
      } else {
        console.log("Error in estimating gas cost")
        return 1
      }
      var signedMessage = JSON.parse(filecoin_signer.transactionSignLotus(create_channel, this.recoveredKey.privateKey))
      // Send message
      var res = this.filRPC.sendSignedMessage(signedMessage)
      if (res.result.ReturnDec.IDAddress == undefined) {
        // Error in creating payment channel
        console.log("Error in creating payment channel")
        return 1
      }
      // Success, add a new entry to the out bound channels
      var cAddr = String(res.result.ReturnDec.IDAddress)
      this.outboundChs.set(cAddr, new ChannelState(cAddr, amount))
      console.log(cAddr)
    }
    return 0
  }

  pay(recipient: string, lane: number, amount: BigNumber) {
    return { paychAddrs: '', voucher: '', topup: false } as payResponse
  }
}