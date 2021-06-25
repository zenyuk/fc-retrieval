const blake2b = require('blake2b')

// RetrievalV1Hash message digests some data using the algorithm used by version one of the
// Filecoin retrieval protocol.
export const retrievalV1Hash = (data: Uint8Array): Uint8Array => {
  return blake2b(32).update(data).digest('binary')
}
