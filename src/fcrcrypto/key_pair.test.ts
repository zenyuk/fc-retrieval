import { retrievalV1Hash } from './key_pair.class';

describe('Crypto dependencies test', () => {
  it('retrievalV1Hash', async () => {
    const out = retrievalV1Hash(Buffer.from('00'));

    expect(out.length).toEqual(32);
    expect(out[0]).toStrictEqual(203);
    expect(out[1]).toStrictEqual(198);
    expect(out[31]).toStrictEqual(66);
  });
});
