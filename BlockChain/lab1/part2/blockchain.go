package main

<<<<<<< HEAD
import (
	"bytes"
	"errors"
	"time"
)

=======
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
// https://en.bitcoin.it/wiki/File:Jonny1000thetimes.png
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

// Blockchain keeps a sequence of Blocks
type Blockchain struct {
	blocks []*Block
}

// NewBlockchain creates a new blockchain with genesis Block
func NewBlockchain() *Blockchain {
<<<<<<< HEAD
	block := NewGenesisBlock(
		&Transaction{
			Data: []byte(genesisCoinbaseData),
		},
	)

	bc := &Blockchain{
		blocks: []*Block{block},
	}
	return bc
=======
	// TODO(student)
	return nil
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
}

// addBlock saves the block into the blockchain
func (bc *Blockchain) addBlock(transactions []*Transaction) *Block {
<<<<<<< HEAD
	lastBlock := bc.CurrentBlock()
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
		PrevBlockHash: lastBlock.Hash,
	}
	block.SetHash()
	bc.blocks = append(bc.blocks, block)
	return block
=======
	// TODO(student)
	return nil
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
}

// GetGenesisBlock returns the Genesis Block
func (bc *Blockchain) GetGenesisBlock() *Block {
<<<<<<< HEAD
	// Genesis block is always the first block for blockchain created via NewBlockChain
	return bc.blocks[0]
=======
	// TODO(student)
	return nil
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
}

// CurrentBlock returns the last block
func (bc *Blockchain) CurrentBlock() *Block {
<<<<<<< HEAD
	lastBlock := bc.blocks[len(bc.blocks)-1]
	return lastBlock
=======
	// TODO(student)
	return nil
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
}

// GetBlock returns the block of a given hash
func (bc *Blockchain) GetBlock(hash []byte) (*Block, error) {
<<<<<<< HEAD
	for _, block := range bc.blocks {
		if bytes.Equal(block.Hash, hash) {
			return block, nil
		}
	}
	return nil, errors.New("Block not found in blockchain")
=======
	// TODO(student)
	return nil, nil
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
}
