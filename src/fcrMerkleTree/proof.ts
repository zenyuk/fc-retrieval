const { MerkleTree } = require('merkletreejs')

export interface MerkleProof {
  position: string
  data: {
    type: string
    data: number[]
  }
}

export const verifyContent = (merkleProof: MerkleProof[], leaf: string, root: string): boolean => {
  const tree = new MerkleTree([])
  return tree.verify(merkleProof, leaf, root)
}
