const { MerkleTree } = require('merkletreejs')
const SHA256 = require('crypto-js/sha256')

export interface MerkleProof {
  position: string
  data: {
    type: string
    data: number[]
  }
}

export const verifyContent = (merkleProof: MerkleProof[], leaf: string, root: string): boolean => {
  const tree = new MerkleTree([], SHA256)
  return tree.verify(merkleProof, leaf, root)
}
