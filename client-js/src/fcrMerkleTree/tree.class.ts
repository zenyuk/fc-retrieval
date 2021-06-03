import {Content, MerkleTree} from './types'
import { FCRMerkleProof } from './proof.class'

export class FCRMerkleTree {
	tree: MerkleTree

  constructor(contents: Content[]) {
    // TODO
  }

  getMerkleRoot(): string {
    // TODO
    return ''
  }

  generateMerkleProof(content: Content): FCRMerkleProof {
    // TODO
    return undefined
  }
}
