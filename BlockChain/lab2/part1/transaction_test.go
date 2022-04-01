package main

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestEquals(t *testing.T) {
	tx0ID := testTransactions["tx0"].ID
	for _, test := range []struct {
		name   string
		tx     *Transaction
		result bool
	}{
		{
			name:   "equal",
			tx:     testTransactions["tx0"],
			result: true,
		},
		{
			name:   "not equal",
			tx:     testTransactions["tx1"],
			result: false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.result, test.tx.Equals(tx0ID))
		})
	}
}

func TestSerialize(t *testing.T) {
	for name, tx := range testTransactions {
		t.Run(name, func(t *testing.T) {
			serialized := tx.Serialize()

			dec := gob.NewDecoder(bytes.NewReader(serialized))
			decoded := &Transaction{}
			err := dec.Decode(decoded)
			if err != nil {
				t.Fatalf("error decoding tx: %v", err)
			}

			if diff := cmp.Diff(tx, decoded); diff != "" {
				t.Errorf("wrong serialization: (-want +got)\n%s", diff)
			}
		})
	}
}

func TestHash(t *testing.T) {
	for name, tx := range testTransactions {
		t.Run(name, func(t *testing.T) {
			txhash := tx.Hash()
			if !bytes.Equal(tx.ID, txhash) {
				t.Errorf("wrong tx hash:\n-want: %x\n+got: %x\n", tx.ID, txhash)
			}
		})
	}
}

func TestIsCoinbase(t *testing.T) {
	tx := testTransactions["tx0"]
	assert.True(t, tx.IsCoinbase())

	tx = testTransactions["tx1"]
	assert.False(t, tx.IsCoinbase())
}

func TestNewCoinbaseTX(t *testing.T) {
	// Passing data to the coinbase transaction
	tx := NewCoinbaseTX("leander", "test")
	if tx == nil {
		t.Fatal("NewCoinbaseTX returned nil")
	}
	assert.Equal(t, -1, tx.Vin[0].OutIdx)
	assert.Nil(t, tx.Vin[0].Txid)
	assert.Equal(t, "test", tx.Vin[0].ScriptSig)
	assert.Equal(t, BlockReward, tx.Vout[0].Value)
	assert.Equal(t, "leander", tx.Vout[0].ScriptPubKey)

	// Using default data
	tx = NewCoinbaseTX("leander", "")
	assert.Equal(t, -1, tx.Vin[0].OutIdx)
	assert.Nil(t, tx.Vin[0].Txid)
	assert.Equal(t, "Reward to leander", tx.Vin[0].ScriptSig)
	assert.Equal(t, BlockReward, tx.Vout[0].Value)
	assert.Equal(t, "leander", tx.Vout[0].ScriptPubKey)
}

func TestNewUTXOTransaction(t *testing.T) {
	from := "rodrigo"
	to := "leander"

	// "From" address have 10 (i.e., genesis coinbase)
	// and "to" address have 0
	utxos := UTXOSet{
		"e2404638779673c7c3e772e12dc3343e6d38f1d71625419d12a8468522b5cc2d": {0: testTransactions["tx0"].Vout[0]},
	}

	// Reject if there is not sufficient funds
	tx1, err := NewUTXOTransaction(to, from, 5, utxos)
	assert.Nil(t, tx1)
	assert.Errorf(t, err, "Not enough funds")

	// Accept otherwise
	tx1, err = NewUTXOTransaction(from, to, 5, utxos)
	assert.Nil(t, err)
	diff(t, testTransactions["tx1"], tx1, "incorrect transaction")

	utxos = UTXOSet{
		"bce268225bc12a0015bcc39e91d59f47fd176e64ca42e4f8aecf107fe38f3bfa": {
			0: testTransactions["tx1"].Vout[0],
			1: testTransactions["tx1"].Vout[1],
		},
	}

	tx2, err := NewUTXOTransaction(to, from, 3, utxos)
	assert.Nil(t, err)
	diff(t, testTransactions["tx2"], tx2, "incorrect transaction")

	tx3, err := NewUTXOTransaction(from, to, 1, utxos)
	assert.Nil(t, err)
	diff(t, testTransactions["tx3"], tx3, "incorrect transaction")
}
