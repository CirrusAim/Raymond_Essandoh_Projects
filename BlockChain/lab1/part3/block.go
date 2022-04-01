package main

import (
<<<<<<< HEAD
	"bytes"
	"errors"
=======
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
	"fmt"
	"strings"
	"time"
)

// Block keeps block information
type Block struct {
	Timestamp     int64          // the block creation timestamp
	Transactions  []*Transaction // The block transactions
	PrevBlockHash []byte         // the hash of the previous block
	Hash          []byte         // the hash of the block
	Nonce         int            // the nonce of the block
}

// NewBlock creates and returns a non-mined Block
func NewBlock(timestamp int64, transactions []*Transaction, prevBlockHash []byte) *Block {
<<<<<<< HEAD
=======
	// TODO(student) -- Remove SetHash. Mine should be called in `blockchain.MineBlock`
	// after create a new block and before add to the blockchain.
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
	return &Block{timestamp, transactions, prevBlockHash, nil, 0}
}

// NewGenesisBlock creates and returns genesis Block
func NewGenesisBlock(timestamp int64, tx *Transaction) *Block {
	return NewBlock(timestamp, []*Transaction{tx}, nil)
}

// Mine calculates and sets the block hash and nonce.
func (b *Block) Mine() {
	// TODO(student) -- create a new PoW and set the Hash and Nonce of the mined block
<<<<<<< HEAD
	pow := NewProofOfWork(b)
	nonce, hash := pow.Run()
	b.Hash = hash
	b.Nonce = nonce
=======
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
}

// HashTransactions returns a hash of the transactions in the block
func (b *Block) HashTransactions() []byte {
	// TODO(student)
	// This function should iterate over all txs in a block,
	// serialize them and compute the merkle root of it.
	// It returns the merkle root hash.
<<<<<<< HEAD
	var allTrans [][]byte
	for _, tran := range b.Transactions {
		allTrans = append(allTrans, tran.Serialize())
	}
	mt := NewMerkleTree(allTrans)
	return mt.MerkleRootHash()
=======
	return []byte{}
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
}

// FindTransaction finds a transaction by its ID
func (b *Block) FindTransaction(ID []byte) (*Transaction, error) {
<<<<<<< HEAD
	for _, tran := range b.Transactions {
		if bytes.Equal(ID, tran.ID) {
			return tran, nil
		}
	}
	return nil, errors.New("Transaction not found")
=======
	// TODO(student) -- what is the easiest way to find a transaction in a block?
	return nil, nil
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
}

func (b *Block) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("============ Block %x ============", b.Hash))
	lines = append(lines, fmt.Sprintf("Prev. hash: %x", b.PrevBlockHash))
	lines = append(lines, fmt.Sprintf("Timestamp: %v", time.Unix(b.Timestamp, 0)))
	lines = append(lines, fmt.Sprintf("Nonce: %d", b.Nonce))
	lines = append(lines, fmt.Sprintf("Transactions:"))
	for i, tx := range b.Transactions {
		lines = append(lines, fmt.Sprintf("%d: %x", i, tx.ID))
	}
	return strings.Join(lines, "\n")
}
