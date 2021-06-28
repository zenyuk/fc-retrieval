import axios from 'axios'
import { FCRMessage } from '../fcrMessages/fcrMessage.class'

/**
 * Send a FCR message
 * 
 * @param {string} url 
 * @param {FCRMessage} message 
 * @returns {Promise<FCRMessage>} 
 */
export const sendMessage = async (url: string, message: FCRMessage): Promise<FCRMessage> => {
  const response = await axios.post<FCRMessage>(url, message, {
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json; charset=utf-8',
    }
  })
  if (response.status != 200) {
    throw Error("Send message failed")
  }
  return response.data
}
