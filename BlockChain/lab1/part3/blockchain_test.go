package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newMockBlockchain() *Blockchain {
	return &Blockchain{[]*Block{testBlockchainData["block0"]}}
}

func addMockBlock(bc *Blockchain, newBlock *Block) {
	bc.blocks = append(bc.blocks, newBlock)
}

func TestBlockchain(t *testing.T) {
	bc := NewBlockchain()
	assert.NotNil(t, bc, "Blockchain is nil")
	if bc != nil {
		assert.Equal(t, 1, len(bc.blocks))
	}
}

func TestGetGenesisBlock(t *testing.T) {
	bc := newMockBlockchain()

	gb := bc.GetGenesisBlock()
	assert.NotNil(t, gb, "Genesis block is nil")
	if gb != nil {
		assert.Nil(t, gb.PrevBlockHash, "Genesis block shouldn't has PrevBlockHash")

		// Genesis block should contains a genesis transaction
		if len(gb.Transactions) > 0 {
			tx := gb.Transactions[0]
			assert.Equal(t, 1, len(gb.Transactions))
			assert.Equal(t, testTransactions["tx0"], tx)
		} else {
			t.Errorf("No transactions found on the Genesis block")
		}
	}
}

func TestAddBlock(t *testing.T) {
	bc := newMockBlockchain()
	assert.Equal(t, 1, len(bc.blocks))

	block := testBlockchainData["block1"]

	err := bc.addBlock(block)
	assert.Nil(t, err, "unexpected error adding block %x", block.Hash)
	assert.Equal(t, 2, len(bc.blocks))
}

func TestCurrentBlock(t *testing.T) {
	bc := newMockBlockchain()

	b := bc.CurrentBlock()
	if b == nil {
		t.Fatal("CurrentBlock returned nil")
	}
	expectedBlock := bc.blocks[0]
	assert.Equalf(t, expectedBlock.Hash, b.Hash, "Current block Hash: %x isn't the expected: %x", b.Hash, expectedBlock.Hash)

	addMockBlock(bc, testBlockchainData["block1"])

	b = bc.CurrentBlock()
	expectedBlock = bc.blocks[1]
	assert.Equalf(t, expectedBlock.Hash, b.Hash, "Current block Hash: %x isn't the expected: %x", b.Hash, expectedBlock.Hash)
}

func TestGetBlock(t *testing.T) {
	bc := newMockBlockchain()

	b, err := bc.GetBlock(bc.blocks[0].Hash)
	assert.NotNil(t, b, "GetBlock returned nil block")
	assert.Nil(t, err, "unexpected error getting block")

	if b != nil {
		assert.Equalf(t, bc.blocks[0].Hash, b.Hash, "Block Hash: %x isn't the expected: %x", b.Hash, bc.blocks[0].Hash)
	}
}

func TestMineBlockWithoutTx(t *testing.T) {
	bc := newMockBlockchain()
	b, err := bc.MineBlock([]*Transaction{})
	assert.Error(t, err, "there are no valid transactions to be mined")
	assert.Nil(t, b)
}

func TestMineBlock(t *testing.T) {
	bc := newMockBlockchain()
	tx := testTransactions["tx1"]

	b, err := bc.MineBlock([]*Transaction{tx})
	assert.NotNil(t, b, "MineBlock returned nil")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(bc.blocks))

	if b != nil {
		gb := bc.blocks[0]

		assert.Equalf(t, gb.Hash, b.PrevBlockHash, "Genesis block Hash: %x isn't equal to current PrevBlockHash: %x", gb.Hash, b.PrevBlockHash)

		minedBlock, err := bc.GetBlock(b.Hash)
		assert.Nil(t, err)
		assert.Equal(t, b, minedBlock)

		txMinedBlock := bc.blocks[1].Transactions[0]
		assert.NotNil(t, txMinedBlock)
		assert.Equal(t, tx.ID, txMinedBlock.ID)
	}
}

func TestValidateBlockWithInvalidHash(t *testing.T) {
	bc := newMockBlockchain()
	block := &Block{
		Timestamp:     TestBlockTime,
		Transactions:  []*Transaction{testTransactions["tx0"]},
		PrevBlockHash: nil,
		Hash:          Hex2Bytes("73d40a0510b6327d0fbcd4a2baf6e7a70f2de174ad2c84538a7b09320e9db3f2"),
		Nonce:         164,
	}
	assert.False(t, bc.ValidateBlock(block))
}

func TestValidateBlockWithInvalidNonce(t *testing.T) {
	bc := newMockBlockchain()
	block := &Block{
		Timestamp:     TestBlockTime,
		Transactions:  []*Transaction{testTransactions["tx0"]},
		PrevBlockHash: nil,
		Hash:          Hex2Bytes("00b8075f4a34f54c1cf0c7f6ec9605a52161ee21e974abb4fa8a39ab7553049a"),
		Nonce:         1,
	}
	assert.False(t, bc.ValidateBlock(block))
}

func TestValidateBlock(t *testing.T) {
	bc := newMockBlockchain()

	for _, b := range []struct {
		name  string
		block *Block
		valid bool
	}{
		{
			name:  "genesis",
			block: testBlockchainData["block0"],
			valid: true,
		},
		{
			name: "mined",
			block: &Block{
				Timestamp:     TestBlockTime,
				Transactions:  []*Transaction{testTransactions["tx1"]},
				PrevBlockHash: testBlockchainData["block0"].Hash,
				Hash:          Hex2Bytes("00940171f20a13b9fd2cdf2c5866023c9ba876cf219951c853905bbff18af962"),
				Nonce:         711,
			},
			valid: true,
		},
		{
			name:  "nil",
			block: nil,
			valid: false,
		},
		{
			name:  "no transactions",
			block: NewBlock(TestBlockTime, []*Transaction{}, nil),
			valid: false,
		},
	} {
		t.Run(b.name, func(t *testing.T) {
			assert.Equal(t, b.valid, bc.ValidateBlock(b.block))
		})
	}
}

func TestFindTransactionSuccess(t *testing.T) {
	bc := newMockBlockchain()

	// Find genesis transaction
	tx, err := bc.FindTransaction(Hex2Bytes("30f2e93d7c139e7766fb80b3cb0150e0e764946bb7e4d7d54d69b53f0b1a6af1"))
	assert.Nil(t, err)
	assert.NotNil(t, tx)
}

func TestFindTransactionFailure(t *testing.T) {
	bc := newMockBlockchain()

	notFoundTx, err := bc.FindTransaction(Hex2Bytes("non-existentID"))
	assert.Error(t, err, "Transaction not found")
	assert.Nil(t, notFoundTx)
}
