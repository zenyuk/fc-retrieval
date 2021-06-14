import axios from 'axios'
import qs = require('qs')

import { FCRMessage } from '../fcrMessages/fcrMessage.class'
import { GatewayRegister } from '../register/register.class'

// getGateways retrives the gateways list
export const getGateways = async (url: string): Promise<GatewayRegister[]> => {
  try {
    const response = await axios.get(url)
    return response.data as GatewayRegister[]
  } catch (error) {
    throw error
  }
}

// SendMessage send a message and get a response message
export const sendMessage = async (url: string, message: FCRMessage): Promise<FCRMessage> => {
  try {
    const response = await axios.post(url, qs.stringify(message))
    return response.data as FCRMessage
  } catch (error) {
    throw error
  }
}
