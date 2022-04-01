package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
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
func NewCoinbaseTX(to, data string) (*Transaction, error) {
	if len(data) == 0 {
		b := make([]byte, 10)
		if _, err := rand.Read(b); err != nil {
			data = "Random"
		}
		data = hex.EncodeToString(b)
	}
	vin := TXInput{
		Txid:   nil,
		OutIdx: -1,
		PubKey: []byte(data),
	}

	vout := TXOutput{
		Value: BlockReward,
	}
	vout.Lock(to)
	tx := &Transaction{Vin: []TXInput{vin}, Vout: []TXOutput{vout}}
	tx.ID = tx.Hash()
	return tx, nil
}

// NewUTXOTransaction creates a new UTXO transaction
// NOTE: The returned tx is NOT signed!
func NewUTXOTransaction(pubKey []byte, to string, amount int, utxos UTXOSet) (*Transaction, error) {
	// TODO(student)
	// Modify your function to use the address instead of just strings
	// And also sign the new transaction before return
	spendableAmt, unspentOutputs := utxos.FindSpendableOutputs(HashPubKey(pubKey), amount)

	if spendableAmt < amount {
		return nil, ErrNoFunds
	}

	var Vin []TXInput
	var Vout []TXOutput
	for prevTxId, outInfo := range unspentOutputs {
		for _, outIdx := range outInfo {
			vin := TXInput{
				Txid:   Hex2Bytes(prevTxId),
				OutIdx: outIdx,
				PubKey: pubKey,
			}

			Vin = append(Vin, vin)
		}
	}

	//  create output
	vout := TXOutput{
		Value: amount,
	}
	vout.Lock(to)

	Vout = append(Vout, vout)
	if spendableAmt > amount {
		change := TXOutput{
			Value: spendableAmt - amount,
		}
		change.PubKeyHash = HashPubKey(pubKey)
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

// TrimmedCopy creates a trimmed copy of Transaction to be used in signing
func (tx Transaction) TrimmedCopy() Transaction {
	copyTx := new(Transaction)

	buf := new(bytes.Buffer)
	gob.NewEncoder(buf).Encode(tx)
	gob.NewDecoder(buf).Decode(copyTx)

	for idx := range copyTx.Vin {
		copyTx.Vin[idx].Signature = nil
		copyTx.Vin[idx].PubKey = nil
	}

	return *copyTx
}

// Sign signs each input of a Transaction
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]*Transaction) error {
	// TODO(student)
	// 1) coinbase transactions are not signed.
	// 2) Throw a Panic in case of any prevTXs (used inputs) didn't exists
	// Take a look on the tests to see the expected error message
	// 3) Create a copy of the transaction to be signed
	// 4) Sign all the previous TXInputs of the transaction tx using the
	// copy as the payload (serialized) to be signed in the ecdsa.Sig
	// (https://golang.org/pkg/crypto/ecdsa/#Sign)
	// Make sure that each input of the copy to be signed
	// have the correct PubKeyHash of each output in the prevTXs
	// Store the signature as a concatenation of R and S fields

	// Coinbase transactions are not signed.
	if tx.IsCoinbase() {
		return nil
	}

	for _, vin := range tx.Vin {
		// tid := fmt.Sprintf("%x", vin.Txid)
		tid := hex.EncodeToString(vin.Txid)
		if _, ok := prevTXs[tid]; !ok {
			return (ErrTxInputNotFound)
		}
	}

	copyTx := tx.TrimmedCopy()

	for idx, vin := range copyTx.Vin {
		// tid := fmt.Sprintf("%x", vin.Txid)
		tid := hex.EncodeToString(vin.Txid)
		prevOut := prevTXs[tid]
		for _, out := range prevOut.Vout {
			if bytes.Equal(vin.Txid, prevOut.ID) {
				copyTx.Vin[idx].PubKey = out.PubKeyHash
				break
			}
		}
		r, s, _ := ecdsa.Sign(rand.Reader, &privKey, vin.PubKey)
		var concat []byte
		concat = append(concat, r.Bytes()...)
		concat = append(concat, s.Bytes()...)

		tx.Vin[idx].Signature = concat
		copyTx.Vin[idx].PubKey = nil
	}

	return nil
}

// Verify verifies signatures of Transaction inputs
func (tx Transaction) Verify(prevTXs map[string]*Transaction) bool {
	// TODO(student)
	// 1) coinbase transactions are not signed.
	// 2) Throw a Panic in case of any prevTXs (used inputs) didn't exists
	// Take a look on the tests to see the expected error message
	// 3) Create the same copy of the transaction that was signed
	// and get the curve used for sign: P256
	// 4) Doing the opposite operation of the signing, perform the
	// verification of the signature, by recovering the R and S byte fields
	// of the Signature and the X and Y fields of the PubKey from
	// the inputs of tx. Verify the signature of each input using the
	// ecdsa.Verify function (https://golang.org/pkg/crypto/ecdsa/#Verify)
	// Note that to use this function you need to reconstruct the
	// ecdsa.PublicKey. Also notice that the ecdsa.Verify function receive
	// a byte array, you the transaction copy need to be serialized.
	// return true if all inputs have valid signature,
	// and false if any of them have an invalid signature.

	// Coinbase transactions need not to be verified.
	if tx.IsCoinbase() {
		return true
	}

	for _, vin := range tx.Vin {
		// tid := fmt.Sprintf("%x", vin.Txid)
		tid := hex.EncodeToString(vin.Txid)
		if _, ok := prevTXs[tid]; !ok {
			return false
		}
	}

	copyTx := tx.TrimmedCopy()
	curve := elliptic.P256()

	for idx, vin := range copyTx.Vin {
		// tid := fmt.Sprintf("%x", vin.Txid)
		tid := hex.EncodeToString(vin.Txid)
		prevOut := prevTXs[tid]
		for _, out := range prevOut.Vout {
			if bytes.Equal(vin.Txid, prevOut.ID) {
				copyTx.Vin[idx].PubKey = out.PubKeyHash
				break
			}
		}

		r := new(big.Int)
		s := new(big.Int)
		x := new(big.Int)
		y := new(big.Int)

		sign := tx.Vin[idx].Signature
		r.SetBytes(sign[:len(sign)/2])
		s.SetBytes(sign[len(sign)/2:])

		pubKey := tx.Vin[idx].PubKey
		x.SetBytes(pubKey[0 : len(pubKey)/2])
		y.SetBytes(pubKey[len(pubKey)/2:])

		pub := ecdsa.PublicKey{Curve: curve, X: x, Y: y}
		hash := copyTx.Serialize()

		status := ecdsa.Verify(&pub, hash, r, s)

		if !status {
			return false
		}
		copyTx.Vin[idx].PubKey = nil

	}

	return true
}

// String returns a human-readable representation of a transaction
func (tx Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("\n--- Transaction %x :", tx.ID))

	for i, input := range tx.Vin {
		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.Txid))
		lines = append(lines, fmt.Sprintf("       OutIdx:    %d", input.OutIdx))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("       PubKey: %x", input.PubKey))
	}

	for i, output := range tx.Vout {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Value))
		lines = append(lines, fmt.Sprintf("       PubKeyHash: %x", output.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}
