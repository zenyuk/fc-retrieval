import { FCRMerkleTree } from './tree.class'
const SHA256 = require('crypto-js/sha256')

// const { MerkleTree } = require('merkletreejs')/

const mockedContents = [
  '0000000000000000000000000000000000000000000000000000000000000001',
  '0000000000000000000000000000000000000000000000000000000000000002',
  '0000000000000000000000000000000000000000000000000000000000000003',
  '0000000000000000000000000000000000000000000000000000000000000004',
  '0000000000000000000000000000000000000000000000000000000000000005',
]

describe('MerkteTree', () => {
  describe('when create a tree', () => {
    it('gets the merkle root', async () => {
      const merkleTree = new FCRMerkleTree(mockedContents)
      const merkleRoot = await merkleTree.getMerkleRoot()
      expect(merkleRoot).toEqual('934b70c92a21b9a793e221d2f9d95ac9f147fa9eee29ceedf818dd80327e07c1')
    })
    it('gets the proof', async () => {
      const merkleTree = new FCRMerkleTree(mockedContents)
      const merkleRoot = await merkleTree.getMerkleRoot()
      const merkleProof = await merkleTree.generateMerkleProof(mockedContents[0])

      expect(merkleRoot).toEqual('934b70c92a21b9a793e221d2f9d95ac9f147fa9eee29ceedf818dd80327e07c1')
      expect(JSON.stringify(merkleProof[0])).toMatch(
        `{"position":"right","data":{"type":"Buffer","data":[104,228,215,93,35,227,27,210,63,67,38,111,118,149,103,178,88,141,159,61,48,16,250,124,172,174,126,253,111,92,224,173]}}`,
      )
      expect(JSON.stringify(merkleProof[1])).toMatch(
        `{"position":"right","data":{"type":"Buffer","data":[244,184,34,142,90,224,230,102,75,200,67,36,20,150,240,100,216,192,187,16,56,34,197,77,140,123,213,116,148,88,179,104]}}`,
      )
      expect(JSON.stringify(merkleProof[2])).toMatch(
        `{"position":"right","data":{"type":"Buffer","data":[207,197,22,125,201,94,15,228,38,220,78,53,3,45,216,143,22,253,208,247,29,221,233,54,43,29,171,35,152,101,202,92]}}`,
      )
    })
  })
  describe('when verify with good tree proof and root', () => {
    it('succeed', async () => {
      const merkleTree = new FCRMerkleTree(mockedContents)

      const root = await merkleTree.getMerkleRoot()
      const leaf = SHA256(mockedContents[0])
      const merkleProof = await merkleTree.generateMerkleProof(mockedContents[0])

      const isVerified = merkleTree.tree.verify(merkleProof, leaf, root)
      expect(isVerified).toEqual(true)
    })
  })
  describe('when verify with bad proof', () => {
    it('fails', async () => {
      const merkleTree = new FCRMerkleTree(mockedContents)

      const root = await merkleTree.getMerkleRoot()
      const leaf = SHA256(mockedContents[0])
      const badMerkleProof = await merkleTree.generateMerkleProof(mockedContents[1])

      const isVerified = merkleTree.tree.verify(badMerkleProof, leaf, root)
      expect(isVerified).toEqual(false)
    })
  })
})
