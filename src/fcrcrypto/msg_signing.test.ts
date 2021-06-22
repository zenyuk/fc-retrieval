import { decodePrivateKey, getToBeSigned, signMessage } from './msg_signing';

import { KeyVersion } from './key_version.class';
import { decodeStringAsHexArray } from './encode.helper';

const privKey = '015ed053eab6fdf18c03954373ff7f89089992017d56beb8b05305b19800d6afe0';

describe('Client msg_signing', () => {
  it('getToBeSigned', async () => {
    const out = getToBeSigned({
      MessageType: 1,
      ProtocolVersion: 1,
      GatewayID: '1234567890abcdef01234567890abcdef01234567890abcdef01234567890abcdef0',
      Challenge: 'a4b2345654665646461234567890abcdef01234567890abcdef01234567890abcdef',
    });

    expect(out).toStrictEqual(
      Buffer.from(
        '111234567890abcdef01234567890abcdef01234567890abcdef01234567890abcdef0a4b2345654665646461234567890abcdef01234567890abcdef01234567890abcdef',
      ),
    );
  });

  it('TestDecodePrivKeyWithError', async () => {
    const t = () => {
      decodePrivateKey('abcdefghijklmn');
    };
    expect(t).toThrow(Error);
  });

  it('TestSignEmptyMsg', async () => {
    const keyPair = decodePrivateKey(privKey);

    const sig = signMessage(
      keyPair,
      KeyVersion.InitialKeyVersion,
      /*empty msg*/ { message_type: 0, protocol_version: 0 },
    );

    expect(sig).toStrictEqual(
      decodeStringAsHexArray(
        '0000000149a4c1bab090ce563b003bc017cf255b89d26e1f7170da57e58e58afe51407aa51255338e77a828434cf97f71716b0e0c70a8f5c762fd64380916a8d98d0daf700',
      ),
    );
  });

  it('signMessage not empty', async () => {
    const msg = {
      MessageType: 1,
      ProtocolVersion: 1,
      GatewayID: '1234567890abcdef01234567890abcdef01234567890abcdef01234567890abcdef0',
      Challenge: 'a4b2345654665646461234567890abcdef01234567890abcdef01234567890abcdef',
    };

    const keyPair = decodePrivateKey(privKey);
    const out = signMessage(keyPair, KeyVersion.InitialKeyVersion, msg);

    expect(out).toStrictEqual(
      decodeStringAsHexArray(
        '00000001b29b643d232313afbbad00d6b10e23aa82e09b3183d619046de42cf56d9acc24411f8547fa761b416cc4804539ca859c3b4681b86cf0158a880668514855089000',
      ),
    );
  });

  it('TestSignMsgWithError', async () => {
    const t = () => {
      const keyPair = decodePrivateKey(privKey);

      keyPair.alg.algorithm = 2;

      const out = signMessage(keyPair, KeyVersion.InitialKeyVersion, {});
      expect(out).toEqual('');
    };
    expect(t).toThrow(Error);
  });
});
