const defaultProtocolVersion = 1
const defaultAlternativeProtocolVersion = 1
const protocolSupported = [defaultProtocolVersion, defaultAlternativeProtocolVersion]

export class FCRMessage {
  messageType: number
  protocolVersion: number
  protocolSupported: number[]
  messageBody: string
  signature: string

  constructor(msgType: number, msgBody: string) {
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
