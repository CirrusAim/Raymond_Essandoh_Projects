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
	bc := NewBlockchain("rodrigo")
	if bc == nil {
		t.Fatal("Blockchain is nil")
	}
	assert.Equal(t, 1, len(bc.blocks))
}

func TestGetGenesisBlock(t *testing.T) {
	bc := newMockBlockchain()

	// GetGenesisBlock
	gb := bc.GetGenesisBlock()
	if gb == nil {
		t.Fatal("Genesis block is nil")
	}
	assert.Nil(t, gb.PrevBlockHash, "Genesis block shouldn't has PrevBlockHash")

	// Genesis block should contains a genesis transaction
	if len(gb.Transactions) > 0 {
		coinbaseTx := gb.Transactions[0]
		assert.Equal(t, 1, len(gb.Transactions))
		assert.Equal(t, -1, coinbaseTx.Vin[0].OutIdx)
		assert.Nil(t, coinbaseTx.Vin[0].Txid)
		assert.Equal(t, GenesisCoinbaseData, coinbaseTx.Vin[0].ScriptSig)
		assert.Equal(t, BlockReward, coinbaseTx.Vout[0].Value)
		assert.Equal(t, "rodrigo", coinbaseTx.Vout[0].ScriptPubKey)
	} else {
		t.Errorf("No transactions found on the Genesis block")
	}
}

func TestAddBlock(t *testing.T) {
	bc := newMockBlockchain()
	assert.Equal(t, 1, len(bc.blocks))

	b1 := testBlockchainData["block1"]
	err := bc.addBlock(b1)
	assert.Nil(t, err, "unexpected error adding block %x", b1.Hash)
	assert.Equal(t, 2, len(bc.blocks))

	gb := bc.blocks[0]
	assert.Equalf(t, gb.Hash, b1.PrevBlockHash, "Genesis block Hash: %x isn't equal to current PrevBlockHash: %x", gb.Hash, b1.PrevBlockHash)

	b2 := testBlockchainData["block2"]
	err = bc.addBlock(b2)
	assert.Nil(t, err, "unexpected error adding block %x", b2.Hash)
	assert.Equal(t, 3, len(bc.blocks))
	assert.Equalf(t, b1.Hash, b2.PrevBlockHash, "Previous block Hash: %x isn't equal to the expected: %x", b2.PrevBlockHash, b1.Hash)
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
	assert.Nil(t, err)
	if b == nil {
		t.Fatal("GetBlock returned nil")
	}

	assert.Equalf(t, bc.blocks[0].Hash, b.Hash, "Block Hash: %x isn't the expected: %x", b.Hash, bc.blocks[0].Hash)
}

func TestMineBlockWithInvalidTxInput(t *testing.T) {
	bc := newMockBlockchain()

	// Ignore transaction that refer to non-existent transaction input
	invalidTx := &Transaction{
		ID: Hex2Bytes("bce268225bc12a0015bcc39e91d59f47fd176e64ca42e4f8aecf107fe38f3bfa"),
		Vin: []TXInput{
			{
				Txid:      Hex2Bytes("non-existentID"),
				OutIdx:    0,
				ScriptSig: "rodrigo",
			},
		},
		Vout: []TXOutput{
			{Value: 5, ScriptPubKey: "leander"},
			{Value: 5, ScriptPubKey: "rodrigo"},
		},
	}

	b, err := bc.MineBlock([]*Transaction{invalidTx})
	assert.Error(t, err, "there are no valid transactions to be mined")
	assert.Nil(t, b)
}

func TestMineBlock(t *testing.T) {
	bc := newMockBlockchain()

	b, err := bc.MineBlock([]*Transaction{
		minerCoinbaseTx["tx1"],
		testTransactions["tx1"],
	})
	assert.Nil(t, err)
	if b == nil {
		t.Fatal("MineBlock returned nil")
	}
	assert.Equal(t, 2, len(bc.blocks))

	gb := bc.blocks[0]
	assert.Equalf(t, gb.Hash, b.PrevBlockHash, "Genesis block Hash: %x isn't equal to current PrevBlockHash: %x", gb.Hash, b.PrevBlockHash)

	minedBlock, err := bc.GetBlock(b.Hash)
	assert.Equal(t, b, minedBlock)
	if minedBlock == nil {
		t.Fatal("GetBlock returned nil")
	} else {
		txMinedBlock := bc.blocks[1].Transactions[1] // second tx in block1
		assert.NotNil(t, txMinedBlock)
		assert.Equal(t, testTransactions["tx1"].ID, txMinedBlock.ID)
	}
}

func TestVerifyTransaction(t *testing.T) {
	bc := newMockBlockchain()
	assert.True(t, bc.VerifyTransaction(testTransactions["tx0"]))
	assert.True(t, bc.VerifyTransaction(testTransactions["tx1"]))
}

func TestVerifyTransactionInvalidTxInput(t *testing.T) {
	bc := newMockBlockchain()

	tx := &Transaction{
		ID: Hex2Bytes("bce268225bc12a0015bcc39e91d59f47fd176e64ca42e4f8aecf107fe38f3bfa"),
		Vin: []TXInput{
			{
				Txid:      Hex2Bytes("non-existentID"),
				OutIdx:    0,
				ScriptSig: "rodrigo",
			},
		},
		Vout: []TXOutput{
			{Value: 5, ScriptPubKey: "leander"},
			{Value: 5, ScriptPubKey: "rodrigo"},
		},
	}
	assert.False(t, bc.VerifyTransaction(tx))
}

func TestValidateBlock(t *testing.T) {
	bc := newMockBlockchain()

	for _, b := range []struct {
		name  string
		block *Block
		valid bool
	}{
		{
			name:  "valid genesis",
			block: testBlockchainData["block0"],
			valid: true,
		},
		{
			name: "valid mined",
			block: &Block{
				Timestamp: TestBlockTime,
				Transactions: []*Transaction{
					minerCoinbaseTx["tx1"],
					testTransactions["tx1"],
				},
				PrevBlockHash: testBlockchainData["block0"].Hash,
				Hash:          Hex2Bytes("0076483ed21f560b5418f90ac263bd58b94acf0ee643694714245a882a4d3b17"),
				Nonce:         270,
			},
			valid: true,
		},
		{
			name: "invalid block hash",
			block: &Block{
				Timestamp:     TestBlockTime,
				Transactions:  []*Transaction{testTransactions["tx0"]},
				PrevBlockHash: nil,
				Hash:          Hex2Bytes("73d40a0510b6327d0fbcd4a2baf6e7a70f2de174ad2c84538a7b09320e9db3f2"),
				Nonce:         164,
			},
			valid: false,
		},
		{
			name: "invalid block nonce",
			block: &Block{
				Timestamp:     TestBlockTime,
				Transactions:  []*Transaction{testTransactions["tx0"]},
				PrevBlockHash: nil,
				Hash:          Hex2Bytes("00b8075f4a34f54c1cf0c7f6ec9605a52161ee21e974abb4fa8a39ab7553049a"),
				Nonce:         1,
			},
			valid: false,
		},
		{
			name: "missing coinbase",
			block: &Block{
				Timestamp: TestBlockTime,
				// missing coinbase transaction
				Transactions: []*Transaction{
					testTransactions["tx3"],
				},
				PrevBlockHash: testBlockchainData["block1"].Hash,
				Hash:          Hex2Bytes("001b92bf4f15fccc72d5f3be56c430507f83179014da38d9289dbdc03c790c3f"),
				Nonce:         153,
			},
			valid: false,
		},
		{
			name: "wrong coinbase order",
			block: &Block{
				Timestamp: TestBlockTime,
				// wrong coinbase order; Coinbase must be the first transaction in a block!
				Transactions: []*Transaction{
					testTransactions["tx1"],
					minerCoinbaseTx["tx1"],
				},
				PrevBlockHash: testBlockchainData["block0"].Hash,
				Hash:          Hex2Bytes("005209ca671422e9295965054ce9940f3ecbf1e15823d7fe5b1ce144ad1cc28f"),
				Nonce:         410,
			},
			valid: false,
		},
		{
			name:  "nil block",
			block: nil,
			valid: false,
		},
		{
			name:  "empty transaction list",
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
	tx0, err := bc.FindTransaction(testTransactions["tx0"].ID)
	assert.Nil(t, err)
	diff(t, testTransactions["tx0"], tx0, "incorrect transaction found")

	addMockBlock(bc, testBlockchainData["block1"])

	tx1, err := bc.FindTransaction(testTransactions["tx1"].ID)
	assert.Nil(t, err)
	diff(t, testTransactions["tx1"], tx1, "incorrect transaction found")
}

func TestFindTransactionFailure(t *testing.T) {
	bc := newMockBlockchain()

	notFoundTx, err := bc.FindTransaction(Hex2Bytes("non-existentID"))
	assert.Error(t, err, "Transaction not found")
	assert.Nil(t, notFoundTx)
}

func TestFindUTXOSet(t *testing.T) {
	bc := newMockBlockchain()
	expectedUTXOs := getTestExpectedUTXOSet("block0")

	utxos := bc.FindUTXOSet()
	diff(t, expectedUTXOs, utxos, "incorrect UTXO Set")

	addMockBlock(bc, testBlockchainData["block1"])
	expectedUTXOs = getTestExpectedUTXOSet("block1")

	utxos = bc.FindUTXOSet()
	diff(t, expectedUTXOs, utxos, "incorrect UTXO Set")
}
