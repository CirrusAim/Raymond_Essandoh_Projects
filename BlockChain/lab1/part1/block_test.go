package main

import (
	"bytes"
	"testing"
)

// Fixed block timestamp
const BlockTime int64 = 1563897484

// Genesis hash with timestamp BlockTime and data: "Genesis data info"
// 556b087be95a9918ca21c2d25e8fcbe2a484299bc621dff9402d94088d1070e3
var GenesisHash = []byte{85, 107, 8, 123, 233, 90, 153, 24, 202, 33, 194, 210, 94, 143, 203, 226, 164, 132, 41, 155, 198, 33, 223, 249, 64, 45, 148, 8, 141, 16, 112, 227}

// Set of test transactions
var txs = []*Transaction{
	{[]byte("tx 1")},
	{[]byte("tx 2")},
}

func TestGenesisBlock(t *testing.T) {
	// Genesis block
	tx := Transaction{[]byte("Genesis data info")}
	gb := NewGenesisBlock(&tx)
	if !bytes.Equal([]byte{}, gb.PrevBlockHash) {
		t.Error("Genesis block shouldn't has PrevBlockHash")
	}

	if !bytes.Equal(gb.Transactions[0].Data, []byte("Genesis data info")) {
		t.Error("Genesis data should be stored in the block")
	}
}

func TestBlockHashTransactions(t *testing.T) {
	// Hash of txs: c4b94bd52ffc22fbe19b7fe85ef8322cb006ef26c807e4c08cd768c6d6127425
	txsHash := []byte{196, 185, 75, 213, 47, 252, 34, 251, 225, 155, 127, 232, 94, 248, 50, 44, 176, 6, 239, 38, 200, 7, 228, 192, 140, 215, 104, 198, 214, 18, 116, 37}

	b := &Block{BlockTime, txs, GenesisHash, []byte{}}

	if !bytes.Equal(txsHash, b.HashTransactions()) {
		t.Errorf("The block hash %x isn't equal to %x", b.HashTransactions(), txsHash)
	}
}

func TestNewBlock(t *testing.T) {
	// SetHash
	// Hash of headers containing: {GenesisHash, txsHash, BlockTime}
	// 9549f8795def9bca28a3e1111ccd0facae64d8c6be38fbf147d2bcbae7bbd213
	headerHash := []byte{149, 73, 248, 121, 93, 239, 155, 202, 40, 163, 225, 17, 28, 205, 15, 172, 174, 100, 216, 198, 190, 56, 251, 241, 71, 210, 188, 186, 231, 187, 210, 19}

	b1 := &Block{BlockTime, txs, GenesisHash, []byte{}}
	b1.SetHash()

	if !bytes.Equal(headerHash, b1.Hash) {
		t.Errorf("The block hash %x isn't equal to %x", b1.Hash, headerHash)
	}

	// NewBlock
	b2 := NewBlock(txs, GenesisHash)
	b3 := &Block{b2.Timestamp, txs, GenesisHash, []byte{}}
	b3.SetHash()

	if bytes.Equal(b2.Hash, []byte{}) {
		t.Error("The block hash should have a valid value")
	}

	if !bytes.Equal(b2.Hash, b3.Hash) {
		t.Errorf("The block hash %x isn't equal to %x", b3.Hash, b2.Hash)
	}
}
