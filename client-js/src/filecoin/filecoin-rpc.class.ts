import axios, { AxiosInstance, AxiosRequestConfig } from 'axios'

export type Address = any

class FilecoinRPC {
  requester: AxiosInstance

  constructor(url: string, token: string) {
    const config: AxiosRequestConfig = {
      baseURL: url,
      headers: { Authorization: `Bearer ${token}` },
    }
    this.requester = axios.create(config)
  }

  /**
   * Get nonce
   * @param address
   * @returns
   */
  async getNonce(address: string) {
    let response = await this.requester.post('', {
      jsonrpc: '2.0',
      method: 'Filecoin.MpoolGetNonce',
      id: 1,
      params: [address],
    })

    return response.data
  }

  /**
   * Send signed message
   * @param signedMessage
   * @returns
   */
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

  /**
   * Get gas estimation
   * @param message
   * @returns
   */
  async getGasEstimation(message: any) {
    let response = await this.requester.post('', {
      jsonrpc: '2.0',
      method: 'Filecoin.GasEstimateMessageGas',
      id: 1,
      params: [message, { MaxFee: '0' }, null],
    })

    return response.data
  }

  /**
   * Read state
   * @param address
   * @returns
   */
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

export { FilecoinRPC }
