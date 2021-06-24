import { GatewayRegister, ProviderRegister } from './register.class'
import axios from 'axios'

/**
 * Get registered gateways
 * 
 * @param {string} registerUrl 
 * @returns {Promise<GatewayRegister[]>}
 */
export const getGateways = async (registerUrl: string): Promise<GatewayRegister[]> => {
  const response = await axios.get<GatewayRegister[]>(`${registerUrl}/registers/gateway`)
  if (response.status !== 200) {
    throw Error('Get gateways failed')
  }
  return response.data.map(item => new GatewayRegister(item))
}

/**
 * Get a registered gateway
 * 
 * @param {string} registerUrl 
 * @param {string} nodeId 
 * @returns {Promise<GatewayRegister>}
 */
export const getGatewayByID = async (registerUrl: string, nodeId: string): Promise<GatewayRegister> => {
  const response = await axios.get<GatewayRegister>(`${registerUrl}/registers/gateway/${nodeId}`)
  if (response.status !== 200) {
    throw Error('Get gateway by ID failed')
  }
  return new GatewayRegister(response.data)
}

/**
 * Get a registered provider
 * 
 * @param {string} registerUrl 
 * @param {string} nodeId 
 * @returns {Promise<ProviderRegister>}
 */
export const getProviderByID = async (registerUrl: string, nodeId: string): Promise<ProviderRegister> => {
  const response = await axios.get<ProviderRegister>(`${registerUrl}/registers/provider/${nodeId}`)
  if (response.status !== 200) {
    throw Error('Get provider by ID failed')
  }
  return new ProviderRegister(response.data)
}