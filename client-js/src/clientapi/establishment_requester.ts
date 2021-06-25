import {
  decodeClientEstablishmentResponse,
  encodeClientEstablishmentRequest,
} from '../fcrMessages/client_establishment'
import { NodeID } from '../nodeid/nodeid.interface'
import { GatewayRegister } from '../register/register.class'
import { sendMessage } from '../request/request'

/**
 * Request a client establishment
 *
 * @param {GatewayRegister} gatewayInfo
 * @param {Buffer} challenge
 * @param {NodeID} clientID
 * @param {number} ttl
 * @returns {Promise<boolean>}
 */
export const requestEstablishment = async (
  gatewayInfo: GatewayRegister,
  challenge: Buffer,
  clientID: NodeID,
  ttl: number,
): Promise<boolean> => {
  if (challenge.length != 32) {
    throw Error('Challenge is not 32 bytes')
  }
  const sentChallenge = challenge.toString('base64')
  const request = encodeClientEstablishmentRequest(clientID, sentChallenge, ttl)

  const response = await sendMessage(`http://${gatewayInfo.networkInfoClient}/v1`, request)

  if (!response.verify(gatewayInfo.getRootSigningKeyPair())) {
    throw Error('Fail to verify response')
  }
  const { gateway_id, challenge: recvChallenge } = decodeClientEstablishmentResponse(response)
  if (gatewayInfo.nodeId.toLowerCase() != gateway_id.toString().toLowerCase()) {
    throw Error('Gateway ID not match')
  }
  if (recvChallenge != sentChallenge) {
    throw Error('Challenge mismatch')
  }
  return true
}
