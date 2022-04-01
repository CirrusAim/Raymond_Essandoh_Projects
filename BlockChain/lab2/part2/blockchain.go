package main

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
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
func NewBlockchain(address string) (*Blockchain, error) {
	// TODO(student)
	conibaseTx, err := NewCoinbaseTX(address, "")
	if err != nil {
		return nil, ErrNoValidTx
	}
	block := NewGenesisBlock(time.Now().Unix(),
		conibaseTx,
	)

	block.Transactions[0].ID = block.Transactions[0].Hash()
	block.Mine()
	bc := &Blockchain{
		blocks: []*Block{block},
	}
	return bc, nil
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

// VerifyTransaction verifies transaction input signatures
func (bc *Blockchain) VerifyTransaction(tx *Transaction) bool {
	// TODO(student)
	// Modify the function to get the inputs referred in tx
	// and return false in case of some error (i.e. not found the input).
	// Then call Verify for tx passing those inputs as parameter and return the result.
	// Remember that coinbase transaction doesn't have input or signature. Thus all coinbase tx are valid.
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
			if ip.IsLockedWithKey(HashPubKey(vin.PubKey)) {
				result = true
				break
			}
		}
	}
	return result
}

// FindTransaction finds a transaction by its ID in the whole blockchain
func (bc Blockchain) FindTransaction(ID []byte) (*Transaction, error) {
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
	// TODO(student) -- YOU DON'T NEED TO CHANGE YOUR PREVIOUS METHOD
	utxoSet := make(UTXOSet)
	for _, block := range bc.blocks {
		for _, tran := range block.Transactions {
			mp := make(map[int]TXOutput)
			for idx, out := range tran.Vout {
				mp[idx] = out
			}
			id := hex.EncodeToString(tran.ID)
			utxoSet[id] = mp
		}
	}

	for _, block := range bc.blocks {
		for _, tran := range block.Transactions {
			for _, in := range tran.Vin {
				// if the prevHash is found in there, delete the entry
				// tid := fmt.Sprintf("%x", in.Txid)
				tid := hex.EncodeToString(in.Txid)
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

// GetInputTXsOf returns a map index by the ID,
// of all transactions used as inputs in the given transaction
func (bc *Blockchain) GetInputTXsOf(tx *Transaction) (map[string]*Transaction, error) {
	// TODO(student)
	// Use bc.FindTransaction to search over all transactions
	// in the blockchain and if the referred input into tx exists,
	// if so, get the transaction of this input and add it
	// to a map, where the key is the id of the transaction found
	// and the value is the pointer to transaction itself.
	// To use the id as key in the map, convert it to string
	// using the function: hex.EncodeToString
	// https://golang.org/pkg/encoding/hex/#EncodeToString
	prevTxs := make(map[string]*Transaction)

	for _, vin := range tx.Vin {
		tran, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			return nil, err
		}
		tid := hex.EncodeToString(vin.Txid)
		prevTxs[tid] = tran
	}

	return prevTxs, nil
}

// SignTransaction signs inputs of a Transaction
func (bc *Blockchain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) error {
	// TODO(student)
	// Get the previous transactions referred in the input of tx
	// and call Sign for tx.
	prevTxs, err := bc.GetInputTXsOf(tx)
	if err != nil {
		return err
	}
	err = tx.Sign(privKey, prevTxs)
	return err
}

func (bc Blockchain) String() string {
	var lines []string
	for _, block := range bc.blocks {
		lines = append(lines, fmt.Sprintf("%v", block))
	}
	return strings.Join(lines, "\n")
}
