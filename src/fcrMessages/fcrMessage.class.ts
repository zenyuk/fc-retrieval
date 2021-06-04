export class FCRMessage {
  messageType: number = 0
  protocolVersion: number = 0
  protocolSupported: number[] = []
  messageBody: string = ''
  signature: string = ''

  verify(pubKey: string): boolean {
    return true
  }
}

export const createFCRMessage = (messageType: number, body: string): FCRMessage => {
  return {} as FCRMessage
}
