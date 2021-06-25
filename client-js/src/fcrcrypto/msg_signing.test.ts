import { decodePrivateKey, decodePublicKey, getToBeSigned, signMessage, verifyMessage } from './msg_signing'

import { KeyVersion } from './key_version.class'
const privKey = '015ed053eab6fdf18c03954373ff7f89089992017d56beb8b05305b19800d6afe0'
const pubKey =
  '01047799f37b014564e23578447d718e5c70a786b0e4e58ca25cb2a086b822434594d910b9b8c0fcbfe9f4c2db321e874819e0614be5b57fbb5080accd69adb2eaad'

describe('Client msg_signing', () => {
  it('TestVerifyMsgShortSig', async () => {
    const keyPair = decodePublicKey(pubKey)
    const out = verifyMessage(keyPair, '0x12', '{}')
    expect(out).toBeFalsy()
  })

  it('TestVerifyMsg false', async () => {
    const keyPair = decodePublicKey(pubKey)
    const raw = JSON.stringify({
      message_type: 1,
      protocol_version: 1,
      gateway_id: '1234567890abcdef01234567890abcdef01234567890abcdef01234567890abcdef0',
      challenge: 'a4b2345654665646461234567890abcdef01234567890abcdef01234567890abcdef',
    })

    const out = verifyMessage(
      keyPair,
      '00000001b29b643d232313afbbad00d6b10e73aa82e09b3183d619046de42cf56d9acc24411f8547fa761b416cc4804539ca859c3b4681b86cf0158a880668514855089000',
      raw,
    )
    expect(out).toEqual(false)
  })

  it('TestVerifyMsg true', async () => {
    const keyPair = decodePublicKey(pubKey)
    const raw = JSON.stringify({
      message_type: 1,
      protocol_version: 1,
      gateway_id: '1234567890abcdef01234567890abcdef01234567890abcdef01234567890abcdef0',
      challenge: 'a4b2345654665646461234567890abcdef01234567890abcdef01234567890abcdef',
    })

    const out = verifyMessage(
      keyPair,
      '000000017989a43a3545120d9e9134b592653b582b230f80c6eb18aa18847e0ba47f7518261cac7232388ea18a977718e5f0f81e78c89fc60f7132bed878ac3fccf7063d00',
      raw,
    )
    expect(out).toEqual(true)
  })
  it('getToBeSigned', async () => {
    const out = getToBeSigned({
      MessageType: 1,
      ProtocolVersion: 1,
      GatewayID: '1234567890abcdef01234567890abcdef01234567890abcdef01234567890abcdef0',
      Challenge: 'a4b2345654665646461234567890abcdef01234567890abcdef01234567890abcdef',
    })

    expect(out).toStrictEqual(
      Buffer.from(
        '111234567890abcdef01234567890abcdef01234567890abcdef01234567890abcdef0a4b2345654665646461234567890abcdef01234567890abcdef01234567890abcdef',
      ),
    )
  })

  it('TestDecodePrivKeyWithError', async () => {
    const t = () => {
      decodePrivateKey('abcdefghijklmn')
    }
    expect(t).toThrow(Error)
  })

  it('TestSignEmptyMsg', async () => {
    const keyPair = decodePrivateKey(privKey)

    const sig = signMessage(
      keyPair,
      KeyVersion.InitialKeyVersion,
      /*empty msg*/ { message_type: 0, protocol_version: 0 },
    )

    expect(sig).toStrictEqual(
      '0000000149a4c1bab090ce563b003bc017cf255b89d26e1f7170da57e58e58afe51407aa51255338e77a828434cf97f71716b0e0c70a8f5c762fd64380916a8d98d0daf700',
    )
  })

  it('signMessage not empty', async () => {
    const msg = {
      MessageType: 1,
      ProtocolVersion: 1,
      GatewayID: '1234567890abcdef01234567890abcdef01234567890abcdef01234567890abcdef0',
      Challenge: 'a4b2345654665646461234567890abcdef01234567890abcdef01234567890abcdef',
    }

    const keyPair = decodePrivateKey(privKey)
    const out = signMessage(keyPair, KeyVersion.InitialKeyVersion, msg)

    expect(out).toStrictEqual(
      '00000001b29b643d232313afbbad00d6b10e23aa82e09b3183d619046de42cf56d9acc24411f8547fa761b416cc4804539ca859c3b4681b86cf0158a880668514855089000',
    )
  })

  it('TestSignMsgWithError', async () => {
    const t = () => {
      const keyPair = decodePrivateKey(privKey)

      keyPair.alg.algorithm = 2

      const out = signMessage(keyPair, KeyVersion.InitialKeyVersion, {})
      expect(out).toEqual('')
    }
    expect(t).toThrow(Error)
  })
})
