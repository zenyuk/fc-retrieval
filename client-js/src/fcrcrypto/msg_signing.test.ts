import { decodePrivateKey, decodePublicKey, getToBeSigned, signMessage, verifyMessage } from './msg_signing'

import { KeyVersion } from './key_version.class'

const privKey = '015ed053eab6fdf18c03954373ff7f89089992017d56beb8b05305b19800d6afe0'
const pubKey =
  '01047799f37b014564e23578447d718e5c70a786b0e4e58ca25cb2a086b822434594d910b9b8c0fcbfe9f4c2db321e874819e0614be5b57fbb5080accd69adb2eaad'

describe('Client msg_signing', () => {
  it('TestVerifyMsgShortSig', async () => {
    const t = () => {
      const keyPair = decodePublicKey(pubKey)
      const out = verifyMessage(keyPair, '0x12', '{}')
      expect(out).toBeFalsy()
    }
    expect(t).toThrow(Error)
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

  it('TestVerifyMsg true runtime', async () => {
    const keyPair = decodePublicKey(
      '010472e9ab95ed0171cc9f07e9ac0cde6ad23040a97f079ac5702c39867c59149c7c071415ca41c565ef7dda4ebf3cfeb4d52703329e06234720c2e3d25211737ad5',
    )
    const raw = `{"message_type":101,"protocol_version":1,"protocol_supported":[1,1],"message_body":"eyJnYXRld2F5X2lkIjoiMDgwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMCIsImNoYWxsZW5nZSI6Ik1lZitXMGJKR2w3OW8zcjRmNXN1OUdZSGpDd3RtY2pGUmQ2aXNmRk1kMU09In0=","message_signature":""}`

    const signature =
      '0000000106ccdf77b9f655f7f61ca64a219b91891799bbce0373402b2aba763694aed6834b7b369efb39717370d689ac1ac25b45b760cc777653b56fcef6854527e28e2e01'

    const out = verifyMessage(keyPair, signature, raw)
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
