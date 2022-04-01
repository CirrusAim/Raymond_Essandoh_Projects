package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNoFunds         = errors.New("not enough funds")
	ErrTxInputNotFound = errors.New("transaction input not found")
)

// Transaction represents a Bitcoin transaction
type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

// NewCoinbaseTX creates a new coinbase transaction
func NewCoinbaseTX(to, data string) *Transaction {
	if len(data) == 0 {
		data = fmt.Sprintf("Reward to %s", to)
	}
	vin := TXInput{
		Txid:      nil,
		OutIdx:    -1,
		ScriptSig: data,
	}

	vout := TXOutput{
		Value:        BlockReward,
		ScriptPubKey: to,
	}
	tx := &Transaction{Vin: []TXInput{vin}, Vout: []TXOutput{vout}}
	tx.ID = tx.Hash()
	return tx
}

// NewUTXOTransaction creates a new UTXO transaction
func NewUTXOTransaction(from, to string, amount int, utxos UTXOSet) (*Transaction, error) {
	// TODO(student)
	// 1) Find valid spendable outputs and the current balance of the sender
	// 2) The sender has sufficient funds? If not return the error:
	// "Not enough funds"
	// 3) Build a list of inputs based on the current valid outputs
	// 4) Build a list of new outputs, creating a "change" output if necessary
	// 5) Create a new transaction with the input and output list.

	spendableAmt, unspentOutputs := utxos.FindSpendableOutputs(from, amount)

	if spendableAmt < amount {
		return nil, ErrNoFunds
	}

	var Vin []TXInput
	var Vout []TXOutput
	for prevTxId, outInfo := range unspentOutputs {
		for _, outIdx := range outInfo {
			vin := TXInput{
				Txid:      Hex2Bytes(prevTxId),
				OutIdx:    outIdx,
				ScriptSig: utxos[prevTxId][outIdx].ScriptPubKey,
			}
			Vin = append(Vin, vin)
		}

	}

	//  create output
	vout := TXOutput{
		Value:        amount,
		ScriptPubKey: to,
	}
	Vout = append(Vout, vout)
	if spendableAmt > amount {
		change := TXOutput{
			Value:        spendableAmt - amount,
			ScriptPubKey: from,
		}
		Vout = append(Vout, change)
	}

	tx := &Transaction{Vin: Vin, Vout: Vout}
	tx.ID = tx.Hash()
	return tx, nil
}

// IsCoinbase checks whether the transaction is coinbase
func (tx Transaction) IsCoinbase() bool {
	return tx.Vin[0].OutIdx == -1
}

// Equals checks if the given transaction ID matches the ID of tx
func (tx Transaction) Equals(ID []byte) bool {
	return bytes.Equal(tx.ID, ID)
}

// Serialize returns a serialized Transaction
func (tx Transaction) Serialize() []byte {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(tx)
	if err != nil {
		return nil
	}
	return buffer.Bytes()
}

// Hash returns the hash of the Transaction
func (tx *Transaction) Hash() []byte {
	tx1 := Transaction{ID: []byte{}, Vin: tx.Vin, Vout: tx.Vout}
	data := tx1.Serialize()
	hash := sha256.Sum256(data)
	return hash[:]
}

// String returns a human-readable representation of a transaction
func (tx Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x :", tx.ID))

	for i, input := range tx.Vin {
		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.Txid))
		lines = append(lines, fmt.Sprintf("       OutIdx:    %d", input.OutIdx))
		lines = append(lines, fmt.Sprintf("       ScriptSig: %s", input.ScriptSig))
	}

	for i, output := range tx.Vout {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Value))
		lines = append(lines, fmt.Sprintf("       ScriptPubKey: %s", output.ScriptPubKey))
	}

	return strings.Join(lines, "\n")
}
