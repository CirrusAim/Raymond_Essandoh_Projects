package main

import (
<<<<<<< HEAD
	"crypto/sha256"
=======
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
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
<<<<<<< HEAD
	targetVal := new(big.Int).Exp(big.NewInt(2), big.NewInt(256-TARGETBITS), nil)
	return &ProofOfWork{
		block:  block,
		target: targetVal,
	}
=======
	return &ProofOfWork{}
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
}

// setupHeader prepare the header of the block
func (pow *ProofOfWork) setupHeader() []byte {
	// TODO(student)
<<<<<<< HEAD
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
=======
	return nil
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
}

// addNonce adds a nonce to the header
func addNonce(nonce int, header []byte) []byte {
	// TODO(student)
<<<<<<< HEAD
	var data []byte
	data = append(data, header...)
	data = append(data, IntToHex(int64(nonce))...)
	return data
=======
	return nil
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
}

// Run performs the proof-of-work
func (pow *ProofOfWork) Run() (int, []byte) {
	// TODO(student)
<<<<<<< HEAD
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

=======
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
	return 0, nil
}

// Validate validates block's Proof-Of-Work
// This function just validates if the block header hash
// is less than the target AND equals to the mined block hash.
<<<<<<< HEAD

//NEW ISSUE TO LOOK AT AS THAT WILL FORM THE BASIS FOR THE NEXT WORK TO IMPROVE BELOW
// header := pow.setupHeader()
// then you reconstruct the hash using this header and the pow.block.Nonce
// if the resconstructed hash you compare it with the pow.block.Hash and see if they match and if the value of the reconstructed hash is
// smaller than the target.

func (pow *ProofOfWork) Validate() bool {
	// TODO(student)
	z := new(big.Int)
	if pow.block.Hash == nil {
		return false
	}
	z.SetBytes(pow.block.Hash)
	return z.Cmp(pow.target) <= 0
=======
func (pow *ProofOfWork) Validate() bool {
	// TODO(student)
	return false
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
}
