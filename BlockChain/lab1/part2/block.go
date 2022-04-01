package main

import (
<<<<<<< HEAD
	"crypto/sha256"
=======
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
	"time"
)

// Block keeps block information
type Block struct {
	Timestamp     int64          // the block creation timestamp
	Transactions  []*Transaction // The block transactions
	PrevBlockHash []byte         // the hash of the previous block
	Hash          []byte         // the hash of the block
}

// NewBlock creates and returns Block
func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}}
	block.SetHash()
	return block
}

// NewGenesisBlock creates and returns genesis Block
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

// SetHash calculates and sets the block hash
func (b *Block) SetHash() {
	var headers []byte
<<<<<<< HEAD
	btTimestamp := IntToHex(b.Timestamp)
	headers = append(headers, []byte(b.PrevBlockHash)...)
	headers = append(headers, b.HashTransactions()...)
	headers = append(headers, []byte(btTimestamp)...)
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
=======
	// TODO(student)
	// Use the function IntToHex in utils.go to converts the timestamp to a byte array. In the first part of the lab we just used strconv for simplicity.
	//  - b.HashTransactions() need to be used here when combining the block header data.
	//  - You should set the block hash to be the hash of the header, so the line below should be changed.
	b.Hash = headers[:]
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
}

// HashTransactions returns a hash of the transactions in the block
func (b *Block) HashTransactions() []byte {
<<<<<<< HEAD
	// var txHash [32]byte
	var allTrans [][]byte
	for _, tran := range b.Transactions {
		allTrans = append(allTrans, tran.Data)
	}
	mt := NewMerkleTree(allTrans)
	// txHash = sha256.Sum256(allTrans)
	return mt.MerkleRootHash()
=======
	var merkleRoot [32]byte
	// TODO(student)
	// You should return the merkle root hash
	return merkleRoot[:]
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
}
