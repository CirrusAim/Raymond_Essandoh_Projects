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
