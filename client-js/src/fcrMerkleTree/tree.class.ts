import { Content } from './types'
const { MerkleTree } = require('merkletreejs')
const SHA256 = require('crypto-js/sha256')
import { MerkleProof } from './proof'

export class FCRMerkleTree {
  tree: any

  /**
   * Create merkle tree
   * @param contents string[]
   */
  constructor(contents: string[]) {
    if (contents.length % 2 != 0) {
      contents.push(contents[contents.length - 1])
    }
    const hashedContents = contents.map(x => SHA256(x))
    this.tree = new MerkleTree(hashedContents, SHA256)
  }

  /**
   * Get merkle root
   * @returns string
   */
  getMerkleRoot(): string {
    return this.tree.getRoot().toString('hex')
  }

  /**
   * Generate merkle proof
   * @param content
   * @returns string
   */
  generateMerkleProof(content: Content): MerkleProof[] {
    return this.tree.getProof(SHA256(content))
  }
}
