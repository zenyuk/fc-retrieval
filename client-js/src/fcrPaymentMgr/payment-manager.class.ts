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
    this.laneStates = new Map<number, LaneState>()
  }
}

export class FCRPaymentMgr {
  privateKey: any
  recoveredKey: any
  filRPC: any
  header: any

  outboundChs: Map<string, ChannelState>

  constructor(privateKey: string, lotusAPIAddr: string, authToken: string) {
    this.recoveredKey = filecoin_signer.keyRecover(privateKey)
    this.filRPC = new FilecoinRPC({url: lotusAPIAddr, token: authToken})
    this.header = { 'Authorization': `Bearer ${authToken}` }
    
    this.outboundChs = new Map<string, ChannelState>()
  }

  async topup(recipient: string, amount: BigNumber) {
    if (this.outboundChs.has(recipient)) {
      console.log("I'm here now.")
      // There is an existing channel, TODO. Topup
      var nonce = await this.filRPC.getNonce(this.recoveredKey.address)
      nonce = nonce.result
      var topup = {
        from: this.recoveredKey.address,
        to: this.outboundChs.get(recipient)?.addr,
        nonce: nonce,
        value: amount.toString(10),
        GasLimit: '0',
        GasFeeCap: '0',
        GasPremium: '0',
        Method: 0,
        Params: ""
      } as any
      topup = await this.filRPC.getGasEstimation(topup)
      console.log(topup)
      if ('result' in topup) {
        topup = topup.result
      } else {
        console.log("Error in estimating gas cost")
        return 1
      }
      var signedMessage = JSON.parse(filecoin_signer.transactionSignLotus(topup, this.recoveredKey.private_base64))// Send message
      console.log("I'm here.1")
      console.log(signedMessage)
      var res = await this.filRPC.sendSignedMessage(signedMessage)
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
      var nonce = await this.filRPC.getNonce(this.recoveredKey.address)
      nonce = nonce.result
      let create_channel = filecoin_signer.createPymtChan(this.recoveredKey.address, recipient, amount.toString(), nonce, '0', '0', '0')
      create_channel = await this.filRPC.getGasEstimation(create_channel)
      if ('result' in create_channel) {
        create_channel = create_channel.result
      } else {
        console.log("Error in estimating gas cost")
        return 1
      }
      var signedMessage = JSON.parse(filecoin_signer.transactionSignLotus(create_channel, this.recoveredKey.private_base64))
      // Send message
      var res = await this.filRPC.sendSignedMessage(signedMessage)
      if (res.result.ReturnDec.IDAddress == undefined) {
        // Error in creating payment channel
        console.log("Error in creating payment channel")
        return 1
      }

      // Success, add a new entry to the out bound channels
      var cAddr = String(res.result.ReturnDec.RobustAddress)
      console.log(cAddr)
      this.outboundChs.set(recipient, new ChannelState(cAddr, amount))
    }
    return 0
  }

  pay(recipient: string, lane: number, amount: BigNumber) {
    return { paychAddrs: '', voucher: '', topup: false } as payResponse
  }
}

async function test() {
  var lotusAPIAddr = "http://127.0.0.1:1234/rpc/v0"
  var lotusToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.lYpknouVX5M_BJOEbHZrdcHdxHkfu0ih1W0NCTFlJz0"
  // Private key: d54193a9668ae59befa59498cdee16b78cdc8228d43814442a64588fd1648a29
  // addr: f12yybez3cfe2yb2nsartagpwkk23q5hmmiluqafi
  var base64String = Buffer.from("d54193a9668ae59befa59498cdee16b78cdc8228d43814442a64588fd1648a29", 'hex').toString('base64')

  var mgr = new FCRPaymentMgr(base64String, lotusAPIAddr, lotusToken)
  await mgr.topup("t1cldh4eiwlx47wkjgf2j37piatb5xevjgld4vjua", new BigNumber(1000000))
  console.log(mgr.outboundChs)
  await mgr.topup("t1cldh4eiwlx47wkjgf2j37piatb5xevjgld4vjua", new BigNumber(1000000))
  console.log(mgr.outboundChs)  
}

test()