import axios from 'axios';
import { BigNumber } from "bignumber.js";

export type Address = any;

class FilecoinRPC {
    requester: any
    constructor(args: { url: any; token: any }) {
      if (!('url' in args && 'token' in args)) {
        throw new Error('FilecoinRPC required an `url` and a `token` to communicate with the node.')
      }
  
      this.requester = axios.create({
        baseURL: args.url,
        headers: { Authorization: `Bearer ${args.token}` },
      })
    }
  
    async getNonce(address: any) {
      let response = await this.requester.post('', {
        jsonrpc: '2.0',
        method: 'Filecoin.MpoolGetNonce',
        id: 1,
        params: [address],
      })
  
      return response.data
    }
  
    async sendSignedMessage(signedMessage: any) {
      let response = await this.requester.post('', {
        jsonrpc: '2.0',
        method: 'Filecoin.MpoolPush',
        id: 1,
        params: [signedMessage],
      })
  
      if ('error' in response.data) {
        throw new Error(response.data.error.message)
      }
  
      let cid = response.data.result
  
      response = await this.requester.post('', {
        jsonrpc: '2.0',
        method: 'Filecoin.StateWaitMsg',
        id: 1,
        params: [cid, null],
      })
  
      return response.data
    }
  
    async getGasEstimation(message: any) {
      let response = await this.requester.post('', {
        jsonrpc: '2.0',
        method: 'Filecoin.GasEstimateMessageGas',
        id: 1,
        params: [message, { MaxFee: '0' }, null],
      })
  
      return response.data
    }
  
    async readState(address: any) {
      let response = await this.requester.post('', {
        jsonrpc: '2.0',
        method: 'Filecoin.StateReadState',
        id: 1,
        params: [address, null],
      })
  
      return response.data
    }
  }

  export declare class Message {
    Version?: number;
    To: string;
    From: string;
    Nonce: number;
    Value: BigNumber;
    GasLimit: number;
    GasFeeCap: BigNumber;
    GasPremium: BigNumber;
    Method: number;
    Params: string;
}

export {FilecoinRPC}