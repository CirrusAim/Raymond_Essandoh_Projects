package main

import (
<<<<<<< HEAD
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"
=======
	"errors"
	"fmt"
	"strings"
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
)

// Blockchain keeps a sequence of Blocks
type Blockchain struct {
	blocks []*Block
}

// NewBlockchain creates a new blockchain with genesis Block
func NewBlockchain() *Blockchain {
	// TODO(student)
<<<<<<< HEAD
	block := NewGenesisBlock(time.Now().Unix(),
		&Transaction{
			Data: []byte(GenesisCoinbaseData),
		},
	)

	block.Transactions[0].ID = block.Transactions[0].Hash()
	block.Mine()
	bc := &Blockchain{
		blocks: []*Block{block},
	}
	return bc
=======
	return nil
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
}

// addBlock saves a valid block into the blockchain
func (bc *Blockchain) addBlock(block *Block) error {
	// TODO(student) -- make sure you only add valid blocks!
<<<<<<< HEAD
	ok := bc.ValidateBlock(block)
	if !ok {
		return errors.New("unexpected error adding block %x" + string(block.Hash))
	}

	bc.blocks = append(bc.blocks, block)
=======
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
	return nil
}

// GetGenesisBlock returns the Genesis Block
func (bc Blockchain) GetGenesisBlock() *Block {
	// TODO(student) -- YOU DON'T NEED TO CHANGE YOUR PREVIOUS METHOD
<<<<<<< HEAD
	return bc.blocks[0]
=======
	return nil
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
}

// CurrentBlock returns the last block
func (bc Blockchain) CurrentBlock() *Block {
	// TODO(student) -- YOU DON'T NEED TO CHANGE YOUR PREVIOUS METHOD
<<<<<<< HEAD
	return bc.blocks[len(bc.blocks)-1]
=======
	return nil
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
}

// GetBlock returns the block of a given hash
func (bc Blockchain) GetBlock(hash []byte) (*Block, error) {
	// TODO(student) -- YOU DON'T NEED TO CHANGE YOUR PREVIOUS METHOD
<<<<<<< HEAD
	for _, block := range bc.blocks {
		if bytes.Equal(block.Hash, hash) {
			return block, nil
		}
	}
	return nil, errors.New("Block not found in blockchain")
=======
	return nil, nil
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
}

// ValidateBlock validates the block before adding it to the blockchain
func (bc *Blockchain) ValidateBlock(block *Block) bool {
	// TODO(student) -- a valid block cannot be nil and must contain txs.
	// Also, it should has the result of a valid PoW.
<<<<<<< HEAD
	if block == nil || len(block.Transactions) == 0 {
		return false
	}
	pow := NewProofOfWork(block)
	if !pow.Validate() || pow.block.Nonce <= 1 {
		return false
	}
	return true
=======
	return false
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
}

// MineBlock mines a new block with the provided transactions and
// adds the block into blockchain.
func (bc *Blockchain) MineBlock(transactions []*Transaction) (*Block, error) {
	// TODO(student) -- you should mine a block within the given txs and add it to the blockchain.
<<<<<<< HEAD
	prevBlock := bc.CurrentBlock()
	newBlock := NewBlock(time.Now().Unix(), transactions, prevBlock.Hash)
	if newBlock == nil || len(newBlock.Transactions) == 0 {
		return nil, errors.New("Block is empty")
	}
	newBlock.Mine()
	if bc.ValidateBlock(newBlock) {
		bc.addBlock(newBlock)
		return newBlock, nil
	}

	return nil, errors.New("unable to Mine a block")
=======
	return nil, nil
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
}

// FindTransaction finds a transaction by its ID in the whole blockchain
func (bc Blockchain) FindTransaction(ID []byte) (*Transaction, error) {
	// TODO(student) -- what is the easiest way to find a transaction in the whole blockchain?
<<<<<<< HEAD
	for _, block := range bc.blocks {
		for _, tran := range block.Transactions {
			if bytes.Equal(tran.ID, ID) {
				return tran, nil
			}
		}
	}
=======
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
	return nil, errors.New("Transaction not found in any block")
}

func (bc Blockchain) String() string {
	var lines []string
	for _, block := range bc.blocks {
		lines = append(lines, fmt.Sprintf("%v", block))
	}
	return strings.Join(lines, "\n")
}
