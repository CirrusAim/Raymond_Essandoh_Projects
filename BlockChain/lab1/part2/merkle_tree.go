package main

import (
<<<<<<< HEAD
	"bytes"
	"crypto/sha256"
=======
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
	"fmt"
)

// MerkleTree represents a merkle tree
type MerkleTree struct {
	RootNode *Node
	Leafs    []*Node
}

// Node represents a merkle tree node
type Node struct {
	Parent *Node
	Left   *Node
	Right  *Node
	Hash   []byte
}

const (
	leftNode = iota
	rightNode
)

// MerkleProof represents way to prove element inclusion on the merkle tree
type MerkleProof struct {
	proof [][]byte
	index []int64
}

// NewMerkleTree creates a new Merkle tree from a sequence of data
func NewMerkleTree(data [][]byte) *MerkleTree {
<<<<<<< HEAD

	if len(data) == 0 {
		panic("No merkle tree nodes")
	}
	mt := &MerkleTree{}

	for _, dt := range data {
		mt.Leafs = append(mt.Leafs, NewMerkleNode(nil, nil, dt))
	}

	if len(mt.Leafs) == 1 {
		mt.RootNode = &Node{
			Left:  mt.Leafs[0],
			Right: mt.Leafs[0],
			Hash:  mt.Leafs[0].Hash,
		}
	} else {
		mt.RootNode = buildInternalNodes(mt.Leafs, mt)
	}

	return mt
=======
	// TODO(student)
	return nil
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
}

// NewMerkleNode creates a new Merkle tree node
func NewMerkleNode(left, right *Node, data []byte) *Node {
<<<<<<< HEAD
	var dtHash [32]byte
	if data == nil {
		var chash []byte
		chash = append(chash, left.Hash...)
		chash = append(chash, right.Hash...)
		dtHash = sha256.Sum256(chash)
	} else {
		dtHash = sha256.Sum256(data)
	}

	return &Node{
		Left:  left,
		Right: right,
		Hash:  dtHash[:],
	}
=======
	// TODO(student)
	return &Node{}
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
}

// MerkleRootHash return the hash of the merkle root
func (mt *MerkleTree) MerkleRootHash() []byte {
	return mt.RootNode.Hash
}

// MakeMerkleProof returns a list of hashes and indexes required to
// reconstruct the merkle path of a given hash
//
// @param hash represents the hashed data (e.g. transaction ID) stored on
// the leaf node
// @return the merkle proof (list of intermediate hashes), a list of indexes
// indicating the node location in relation with its parent (using the
// constants: leftNode or rightNode), and a possible error.
func (mt *MerkleTree) MakeMerkleProof(hash []byte) ([][]byte, []int64, error) {
<<<<<<< HEAD
	for _, current := range mt.Leafs {
		if bytes.Equal(current.Hash, hash) {
			currentParent := current.Parent
			index := []int64{}
			merklePath := [][]byte{}
			for currentParent != nil {
				if bytes.Equal(currentParent.Left.Hash, current.Hash) {
					merklePath = append(merklePath, currentParent.Right.Hash)
					index = append(index, rightNode)
				} else {
					merklePath = append(merklePath, currentParent.Left.Hash)
					index = append(index, leftNode)
				}
				current = currentParent
				currentParent = currentParent.Parent
			}
			return merklePath, index, nil
		}
	}
=======
	// TODO(student)
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
	return [][]byte{}, []int64{}, fmt.Errorf("Node %x not found", hash)
}

// VerifyProof verifies that the correct root hash can be retrieved by
// recreating the merkle path for the given hash and merkle proof.
//
// @param rootHash is the hash of the current root of the merkle tree
// @param hash represents the hash of the data (e.g. transaction ID)
// to be verified
// @param mProof is the merkle proof that contains the list of intermediate
// hashes and their location on the tree required to reconstruct
// the merkle path.
func VerifyProof(rootHash []byte, hash []byte, mProof MerkleProof) bool {
<<<<<<< HEAD
	genHash := func(d1 []byte, d2 []byte) [32]byte {
		var hash []byte
		hash = append(hash, d1...)
		hash = append(hash, d2...)
		return sha256.Sum256(hash)
	}

	proof := mProof.proof
	index := mProof.index

	var cummulativeHash []byte

	cummulativeHash = hash

	for idx, proofHash := range proof {
		if index[idx] == rightNode {
			newHash := genHash(cummulativeHash, proofHash)
			cummulativeHash = newHash[:]
		} else {
			newHash := genHash(proofHash, cummulativeHash)
			cummulativeHash = newHash[:]
		}
	}

	return bytes.Equal(cummulativeHash, rootHash)
}

func buildInternalNodes(leafs []*Node, mt *MerkleTree) *Node {
	var nodes []*Node

	for i := 0; i < len(leafs); i += 2 {
		var left, right int = i, i + 1

		if i+1 == len(leafs) {
			right = i
		}
		newNode := NewMerkleNode(leafs[left], leafs[right], nil)

		nodes = append(nodes, newNode)
		leafs[left].Parent = newNode
		leafs[right].Parent = newNode
		if len(leafs) == 2 {
			return newNode
		}
	}
	return buildInternalNodes(nodes, mt)
=======
	// TODO(student)
	return false
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
}
