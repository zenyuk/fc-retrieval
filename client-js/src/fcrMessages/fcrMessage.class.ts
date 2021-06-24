const defaultProtocolVersion = 1
const defaultAlternativeProtocolVersion = 1
const protocolSupported = [defaultProtocolVersion, defaultAlternativeProtocolVersion]

export class FCRMessage {
  message_type: number
  protocol_version: number
  protocol_supported: number[]
  message_body: string
  message_signature: string

  constructor(msgType: number, msgBody: string) {
    this.message_type = msgType
    this.message_body = Buffer.from(msgBody).toString('base64')
    this.protocol_version = defaultProtocolVersion
    this.protocol_supported = protocolSupported
    this.message_signature = ''
  }

  verify(pubKey: string): boolean {
    return true
  }
}
