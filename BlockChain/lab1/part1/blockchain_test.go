package main

import (
	"bytes"
	"testing"
)

func TestBlockchain(t *testing.T) {
	// NewBlockchain
	bc := NewBlockchain()
	if bc == nil {
		t.Error("Got an error while creating a new blockchain")
	}

	// GetGenesisBlock
	gb := bc.GetGenesisBlock()
	if gb == nil {
		t.Error("Got an error while getting the Genesis Block")
	}

	// addBlock
	txs := []*Transaction{
		{[]byte("Send 1 coin to Satoshi")},
	}

	b1 := bc.addBlock(txs)
	if b1 == nil {
		t.Fatal("Got an error while adding a new block")
	}
	if !bytes.Equal(gb.Hash, b1.PrevBlockHash) {
		t.Errorf("Genesis block Hash %x isn't equal to current PrevBlockHash %x", gb.Hash, b1.PrevBlockHash)
	}

	txs2 := []*Transaction{
		{[]byte("Send 3 coins to Satoshi")},
	}
	b2 := bc.addBlock(txs2)
	if b2 == nil {
		t.Fatal("Got an error while adding a new block")
	}
	if !bytes.Equal(b1.Hash, b2.PrevBlockHash) {
		t.Errorf("Previous block Hash %x isn't equal to current PrevBlockHash %x", b1.Hash, b2.PrevBlockHash)
	}

	// CurrentBlock
	b3 := bc.CurrentBlock()
	if b3 == nil {
		t.Fatal("Got an error while getting the last block")
	}
	if !bytes.Equal(b3.Hash, b2.Hash) {
		t.Errorf("Current block Hash %x isn't the expected %x", b3.Hash, b2.Hash)
	}

	// GetBlock
	b4, err := bc.GetBlock(b2.Hash)
	if b4 == nil || err != nil {
		t.Fatal("Got an error while getting the block by hash")
	}
	if !bytes.Equal(b4.Hash, b2.Hash) {
		t.Errorf("Block Hash %x isn't the expected %x", b4.Hash, b2.Hash)
	}
}
