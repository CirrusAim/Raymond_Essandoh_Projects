package main

// TXOutput represents a transaction output
type TXOutput struct {
	Value        int    // The transaction value
	ScriptPubKey string // The conditions to claim this output
}

// CanBeUnlockedWith checks if the output can be unlocked with the provided data
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return unlockingData == out.ScriptPubKey
}
