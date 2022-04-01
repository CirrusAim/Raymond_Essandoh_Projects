package main

<<<<<<< HEAD
import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

=======
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
// Transaction represents a simple transaction
type Transaction struct {
	ID   []byte // The hash of the serialized Data
	Data []byte
}

func NewTransaction(data []byte) *Transaction {
	// TODO(student)
<<<<<<< HEAD
	tx := &Transaction{Data: data}
	tx.ID = tx.Hash()
	return tx
=======
	return &Transaction{}
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
}

// Equals checks if the given transaction ID matches the ID of tx
func (tx Transaction) Equals(ID []byte) bool {
<<<<<<< HEAD
	return bytes.Equal(tx.ID, ID)
=======
	// TODO(student)
	return false
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
}

// Serialize returns a serialized Transaction
func (tx Transaction) Serialize() []byte {
	// TODO(student)
	// This function should encode all fields of the Transaction struct, using the gob encoder
	// Note: This includes the tx.ID!
	// TIP: https://golang.org/pkg/encoding/gob/
<<<<<<< HEAD
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(tx)
	if err != nil {
		return nil
	}
	return buffer.Bytes()
=======
	return nil
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
}

// Hash returns the hash of the Transaction
func (tx *Transaction) Hash() []byte {
	// TODO(student)
	// This function should hash the serialized representation of a transaction but it MUST
	// ignore the ID (set it to nil), since the ID is the hash of the tx itself (if exists).
	// You may need to make a copy of the object, otherwise it will change the original pointer.
<<<<<<< HEAD
	tx1 := Transaction{Data: tx.Data, ID: nil}
	data := tx1.Serialize()
	hash := sha256.Sum256(data)
	return hash[:]
=======
	return nil
>>>>>>> 5c06b006eddeb9b9814aaa40bb36fc4ef6af0707
}
