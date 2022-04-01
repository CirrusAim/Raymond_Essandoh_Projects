package main

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrTxNotFound    = errors.New("transaction not found")
	ErrNoValidTx     = errors.New("there is no valid transaction")
	ErrBlockNotFound = errors.New("block not found")
	ErrInvalidBlock  = errors.New("block is not valid")
)

// Blockchain keeps a sequence of Blocks
type Blockchain struct {
	blocks []*Block
}

// NewBlockchain creates a new blockchain with genesis Block
func NewBlockchain(address string) *Blockchain {
	// TODO(student)
	block := NewGenesisBlock(time.Now().Unix(),
		NewCoinbaseTX(address, ""),
	)

	block.Transactions[0].ID = block.Transactions[0].Hash()
	block.Mine()
	bc := &Blockchain{
		blocks: []*Block{block},
	}
	return bc
}

// addBlock saves the block into the blockchain
func (bc *Blockchain) addBlock(block *Block) error {
	// TODO(student) -- make sure you only add valid blocks!
	ok := bc.ValidateBlock(block)
	if !ok {
		return ErrInvalidBlock
	}

	bc.blocks = append(bc.blocks, block)
	return nil
}

// GetGenesisBlock returns the Genesis Block
func (bc Blockchain) GetGenesisBlock() *Block {
	return bc.blocks[0]
}

// CurrentBlock returns the last block
func (bc Blockchain) CurrentBlock() *Block {
	return bc.blocks[len(bc.blocks)-1]
}

// GetBlock returns the block of a given hash
func (bc Blockchain) GetBlock(hash []byte) (*Block, error) {
	for _, block := range bc.blocks {
		if bytes.Equal(block.Hash, hash) {
			return block, nil
		}
	}
	return nil, ErrBlockNotFound
}

// ValidateBlock validates the block before adding it to the blockchain
func (bc *Blockchain) ValidateBlock(block *Block) bool {
	// TODO(student) -- a valid block cannot be nil and must contain txs.
	// Also, it should has the result of a valid PoW.
	if block == nil || len(block.Transactions) == 0 {
		return false
	}

	// check if it has coinbase transaction
	tx := block.Transactions[0]

	if !tx.IsCoinbase() {
		return false
	}

	pow := NewProofOfWork(block)
	if !pow.Validate() || pow.block.Nonce <= 1 {
		return false
	}
	return true
}

// MineBlock mines a new block with the provided transactions
func (bc *Blockchain) MineBlock(transactions []*Transaction) (*Block, error) {
	// TODO(student)
	// 1) Verify the existence of transactions inputs and discard invalid transactions that make reference to unknown inputs
	// 2) Add a block if there is a list of valid transactions
	prevBlock := bc.CurrentBlock()
	newBlock := NewBlock(time.Now().Unix(), transactions, prevBlock.Hash)
	if newBlock == nil || len(newBlock.Transactions) == 0 {
		return nil, ErrNoValidTx
	}

	// verify each transaction
	for _, tx := range transactions {
		if !bc.VerifyTransaction(tx) {
			return nil, ErrNoValidTx
		}
	}

	newBlock.Mine()
	if bc.ValidateBlock(newBlock) {
		bc.addBlock(newBlock)
		return newBlock, nil
	}

	return nil, ErrNoValidTx
}

// VerifyTransaction verifies if referred inputs exist
func (bc Blockchain) VerifyTransaction(tx *Transaction) bool {
	// TODO(student)
	// Check if all inputs of a given transaction refer to a existent transaction made previously
	// if not, you should return false!
	// TIP: remember that Coinbase transaction doesn't have input. Thus all coinbase tx are valid
	result := false
	for _, vin := range tx.Vin {
		// Check if coinbase transaction
		if vin.OutIdx == -1 {
			result = true
			break
		}

		prevTx, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			return false
		}
		for _, ip := range prevTx.Vout {
			if ip.ScriptPubKey == vin.ScriptSig {
				result = true
				break
			}
		}
	}
	return result
}

// FindTransaction finds a transaction by its ID in the whole blockchain
func (bc Blockchain) FindTransaction(ID []byte) (*Transaction, error) {
	// TODO(student)
	// TIP: the chain is made of what?
	for _, block := range bc.blocks {
		for _, tran := range block.Transactions {
			if bytes.Equal(tran.ID, ID) {
				return tran, nil
			}
		}
	}
	return nil, ErrTxNotFound
}

// FindUTXOSet finds and returns all unspent transaction outputs
func (bc Blockchain) FindUTXOSet() UTXOSet {
	// TODO(student)
	// 1) Search in the blockchain for unspent transactions outputs
	// 2) Ignore an already spent output
	// TIP: what determines that an output was spent?
	utxoSet := make(UTXOSet)
	for _, block := range bc.blocks {
		for _, tran := range block.Transactions {
			mp := make(map[int]TXOutput)
			for idx, out := range tran.Vout {
				mp[idx] = out
			}
			id := fmt.Sprintf("%x", tran.ID)
			utxoSet[id] = mp
		}
	}

	for _, block := range bc.blocks {
		for _, tran := range block.Transactions {
			for _, in := range tran.Vin {
				// if the prevHash is found in there, delete the entry
				tid := fmt.Sprintf("%x", in.Txid)
				if _, ok := utxoSet[tid]; ok {
					delete(utxoSet[tid], in.OutIdx)
					if len(utxoSet[tid]) == 0 {
						delete(utxoSet, tid)
					}
				}
			}
		}
	}

	return utxoSet
}

func (bc Blockchain) String() string {
	var lines []string
	for _, block := range bc.blocks {
		lines = append(lines, fmt.Sprintf("%v", block))
	}
	return strings.Join(lines, "\n")
}
