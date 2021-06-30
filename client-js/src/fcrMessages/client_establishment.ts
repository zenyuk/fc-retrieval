import { NodeID } from '../nodeid/nodeid.interface'
import { FCRMessage } from './fcrMessage.class'
import { FCRMessageType } from './type.enum'

export interface ClientEstablishmentRequest {
  client_id: NodeID
  challenge: string
  ttl: number
}

export interface ClientEstablishmentResponse {
  gateway_id: NodeID
  challenge: string
}

/**
 * Encode a client establishment request message
 *
 * @param {NodeID} client_id
 * @param {string} challenge
 * @param {number} ttl
 * @returns {FCRMessage}
 */
export const encodeClientEstablishmentRequest = (client_id: NodeID, challenge: string, ttl: number): FCRMessage => {
  const body = JSON.stringify({
    client_id: client_id.toString(),
    challenge,
    ttl,
  })

  return new FCRMessage({
    message_type: FCRMessageType.ClientEstablishmentRequestType,
    message_body: Buffer.from(body).toString('base64'),
  })
}

/**
 * Decode a client establishment response message
 *
 * @param {FCRMessage} fcrMsg
 * @returns {ClientEstablishmentResponse}
 */
export const decodeClientEstablishmentResponse = (fcrMsg: FCRMessage): ClientEstablishmentResponse => {
  if (fcrMsg.message_type !== FCRMessageType.ClientEstablishmentResponseType) {
    throw Error('Message type mismatch')
  }
  const msgBody = JSON.parse(Buffer.from(fcrMsg.message_body, 'base64').toString('utf-8'))
  return {
    gateway_id: new NodeID(msgBody.gateway_id),
    challenge: msgBody.challenge,
  }
}
