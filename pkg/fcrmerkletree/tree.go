package fcrmerkletree

import (
	"encoding/hex"

	"github.com/cbergoon/merkletree"
)

// FCRMerkleTrie is used to store
type FCRMerkleTrie struct {
	tree *merkletree.MerkleTree
}

// CreateMerkleTrie creates a merkle tree from a list of cids
func CreateMerkleTrie(contents []merkletree.Content) (*FCRMerkleTrie, error) {
	tree, err := merkletree.NewTree(contents)
	if err != nil {
		return nil, err
	}
	return &FCRMerkleTrie{tree: tree}, nil
}

// GetMerkleRoot returns the merkle root of the tree
func (mt *FCRMerkleTrie) GetMerkleRoot() string {
	return hex.EncodeToString(mt.tree.MerkleRoot())
}

// GenerateMerkleProof gets the merkle proof for a given cid
func (mt *FCRMerkleTrie) GenerateMerkleProof(content merkletree.Content) (*FCRMerkleProof, error) {
	path, index, err := mt.tree.GetMerklePath(content)
	if err != nil {
		return nil, err
	}
	return &FCRMerkleProof{path: path, index: index}, nil
}
