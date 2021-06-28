import { FCRMerkleTree } from './tree.class'
const SHA256 = require('crypto-js/sha256')
import { verifyContent } from './proof'

// const { MerkleTree } = require('merkletreejs')/

const mockedContents = [
  '0000000000000000000000000000000000000000000000000000000000000001',
  '0000000000000000000000000000000000000000000000000000000000000002',
  '0000000000000000000000000000000000000000000000000000000000000003',
  '0000000000000000000000000000000000000000000000000000000000000004',
  '0000000000000000000000000000000000000000000000000000000000000005',
]

describe('Proof', () => {
  describe('when verify with good proof leaf and root', () => {
    it('succeed', async () => {
      const merkleTree = new FCRMerkleTree(mockedContents)
      const root = await merkleTree.getMerkleRoot()
      const leaf = SHA256(mockedContents[0])
      const merkleProof = await merkleTree.generateMerkleProof(mockedContents[0])

      const isVerified = verifyContent(merkleProof, leaf, root)
      expect(isVerified).toEqual(true)
    })
  })
  describe('when verify with bad proof', () => {
    it('fails', async () => {
      const merkleTree = new FCRMerkleTree(mockedContents)
      const root = await merkleTree.getMerkleRoot()
      const leaf = SHA256(mockedContents[0])
      const badMerkleProof = await merkleTree.generateMerkleProof(mockedContents[1])

      const isVerified = verifyContent(badMerkleProof, leaf, root)
      expect(isVerified).toEqual(false)
    })
  })
})
