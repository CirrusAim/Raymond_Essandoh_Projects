package main

import (
	"crypto/sha256"
	"math"
	"math/big"
)

var maxNonce = math.MaxInt64

// TARGETBITS define the mining difficulty
const TARGETBITS = 8

// ProofOfWork represents a block mined with a target difficulty
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// NewProofOfWork builds a ProofOfWork
func NewProofOfWork(block *Block) *ProofOfWork {
	// TODO(student)
	targetVal := new(big.Int).Exp(big.NewInt(2), big.NewInt(256-TARGETBITS), nil)
	return &ProofOfWork{
		block:  block,
		target: targetVal,
	}
}

// setupHeader prepare the header of the block
func (pow *ProofOfWork) setupHeader() []byte {
	// TODO(student)
	var data []byte
	var allTrans [][]byte
	tstamp := IntToHex(pow.block.Timestamp)
	tbit := IntToHex(TARGETBITS)
	for _, tran := range pow.block.Transactions {
		allTrans = append(allTrans, tran.Serialize())
	}
	mt := NewMerkleTree(allTrans)
	mtHash := mt.MerkleRootHash()
	data = append(data, pow.block.PrevBlockHash...)
	data = append(data, mtHash...)
	data = append(data, tstamp...)
	data = append(data, tbit...)

	return data
}

// addNonce adds a nonce to the header
func addNonce(nonce int, header []byte) []byte {
	// TODO(student)
	var data []byte
	data = append(data, header...)
	data = append(data, IntToHex(int64(nonce))...)
	return data
}

// Run performs the proof-of-work
func (pow *ProofOfWork) Run() (int, []byte) {
	// TODO(student)
	header := pow.setupHeader()
	var nonce int
	for nonce = 1; nonce < maxNonce; nonce++ {
		data := addNonce(nonce, header)
		hash := sha256.Sum256(data)
		z := new(big.Int)
		z.SetBytes(hash[:])
		if z.Cmp(pow.target) <= 0 {
			return int(nonce), hash[:]
		}
	}

	return 0, nil
}

// Validate validates block's Proof-Of-Work
// This function just validates if the block header hash
// is less than the target AND equals to the mined block hash.
func (pow *ProofOfWork) Validate() bool {
	// TODO(student)
	z := new(big.Int)
	if pow.block.Hash == nil {
		return false
	}
	z.SetBytes(pow.block.Hash)
	return z.Cmp(pow.target) <= 0
}
