package main

// TXInput represents a transaction input
type TXInput struct {
	Txid      []byte // The ID of the referenced transaction containing the output used
	OutIdx    int    // The index of the specific output in the transaction. The first output is 0, etc.
	ScriptSig string // The logic that authorizes the use of this input by satisfying the output's ScriptPubKey
}

// CanUnlockOutputWith checks whether the address initiated the transaction
func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return unlockingData == in.ScriptSig
}
