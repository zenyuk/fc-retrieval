import { Content, MerkleTree } from './types'
import { FCRMerkleProof } from './proof.class'

export class FCRMerkleTree {
  tree: MerkleTree

  getMerkleRoot(): string {
    // TODO
    return ''
  }

  generateMerkleProof(content: Content): FCRMerkleProof {
    // TODO
    return undefined
  }
}

export const createMerkleTree = (contents: Content[]): FCRMerkleTree => {
  // TODO
  return undefined
}
