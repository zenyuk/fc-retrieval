import { verifyMessage } from '../fcrcrypto/msg_signing'
import { KeyPair } from '../fcrcrypto/key_pair.class'

const defaultProtocolVersion = 1
const defaultAlternativeProtocolVersion = 1
const defaultProtocolSupported = [defaultProtocolVersion, defaultAlternativeProtocolVersion]

export class FCRMessage {
  message_type: number
  protocol_version: number
  protocol_supported: number[]
  message_body: string
  message_signature: string

  constructor({
    message_type,
    protocol_version = defaultProtocolVersion,
    protocol_supported = defaultProtocolSupported,
    message_body,
    message_signature = '',
  }: any) {
    this.message_type = message_type
    this.message_body = message_body
    this.protocol_version = protocol_version
    this.protocol_supported = protocol_supported
    this.message_signature = message_signature
  }

  verify(pubKey: KeyPair): boolean {
    return verifyMessage(
      pubKey,
      this.message_signature,
      JSON.stringify({
        message_type: this.message_type,
        protocol_version: this.protocol_version,
        protocol_supported: this.protocol_supported,
        message_body: this.message_body,
        message_signature: '',
      }),
    )
  }
}
