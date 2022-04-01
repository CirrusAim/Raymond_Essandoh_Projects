package main

import "bytes"

// TXInput represents a transaction input
type TXInput struct {
	Txid      []byte // The ID of the referenced transaction containing the output used
	OutIdx    int    // The index of the specific output in the transaction. The first output is 0, etc.
	Signature []byte // The signature of this input
	PubKey    []byte // The logic that authorizes the use of this input by satisfying the output's PubKeyHash. In this demo we will be using the raw public key (not hashed)
}

// UsesKey checks whether the address initiated the transaction
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	// Notice that inputs store raw public keys (i.e., not hashed), but the function takes a hashed one as a parameter and compares it with the hashed version of the PubKey on the input
	return bytes.Equal(HashPubKey(in.PubKey), pubKeyHash)
}
