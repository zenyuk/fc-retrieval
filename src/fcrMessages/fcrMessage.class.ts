const defaultProtocolVersion = 1
const defaultAlternativeProtocolVersion = 1
const protocolSupported = [defaultProtocolVersion, defaultAlternativeProtocolVersion]

export class FCRMessage {
  messageType: number
  protocolVersion: number
  protocolSupported: number[]
  messageBody: object
  signature: string

  constructor(msgType: number, msgBody: object) {
    this.messageType = msgType
    this.messageBody = msgBody
    this.protocolVersion = defaultProtocolVersion
    this.protocolSupported = protocolSupported
    this.signature = ''
  }

  verify(pubKey: string): boolean {
    return true
  }
}
