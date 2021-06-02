export interface FCRMessage {
  messageType: number
  protocolVersion: number
  protocolSupported: number[]
  messageBody: string
  signature: string
}
