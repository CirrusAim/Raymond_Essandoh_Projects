package main

import "fmt"

// UTXOSet represents a set of UTXO as an in-memory cache
// The key of the most external map is the transaction ID
// (encoded as string) that contains these outputs
// {map of transaction ID -> {map of TXOutput Index -> TXOutput}}
type UTXOSet map[string]map[int]TXOutput

// FindSpendableOutputs finds and returns unspent outputs in the UTXO Set
// to reference in inputs and the current accumulated balance
func (u UTXOSet) FindSpendableOutputs(unlockingData string, amount int) (int, map[string][]int) {
	// TODO(student)
	var accumulatedBal int
	unspentOutputs := make(map[string][]int)

	for txID, outputs := range u {
		for outIdx, out := range outputs {
			// if out.ScriptPubKey == unlockingData {
			if out.CanBeUnlockedWith(unlockingData) {
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				accumulatedBal += out.Value
			}
		}
	}

	return accumulatedBal, unspentOutputs
}

// FindUTXO finds all UTXO in the UTXO Set for a given unlockingData key (e.g., address)
// This function ignores the index of each output and returns
// a list of all outputs in the UTXO Set that can be unlocked by the user
func (u UTXOSet) FindUTXO(unlockingData string) []TXOutput {
	var UTXO []TXOutput
	// TODO(student)
	// Search for UTXO that unlockingData can unlock
	for _, outputs := range u {
		for _, out := range outputs {
			// if out.ScriptPubKey == unlockingData {
			if out.CanBeUnlockedWith(unlockingData) {
				UTXO = append(UTXO, out)
			}
		}
	}
	return UTXO
}

// CountUTXOs returns the number of transactions outputs in the UTXO set
func (u UTXOSet) CountUTXOs() int {
	var UXTOCnt int
	for _, outputs := range u {
		UXTOCnt += len(outputs)
	}
	return UXTOCnt
}

// Update updates the UTXO Set with the new set of transactions
func (u UTXOSet) Update(transactions []*Transaction) {
	// TODO(student)
	// Iterate over the transactions  and update
	// the current UTXOSet with the new
	// transactions.
	//

	for _, tx := range transactions {

		newTxOutMap := make(map[int]TXOutput)

		// Check if the transaction has valid scripts

		for _, txVin := range tx.Vin {
			// prev trans id

			if txVin.Txid != nil {
				// not a genesis block
				prevTxId := fmt.Sprintf("%x", txVin.Txid)

				prevOutBlock := u[prevTxId][txVin.OutIdx]
				if !prevOutBlock.CanBeUnlockedWith(txVin.ScriptSig) {
					// not authorized
					continue
				}

				var bal int
				for _, currTxVout := range tx.Vout {
					bal += currTxVout.Value
				}

				// if balance
				if bal >= prevOutBlock.Value {
					delete(u[prevTxId], txVin.OutIdx)
					if len(u[prevTxId]) == 0 {
						delete(u, prevTxId)
					}
				}
			}
		}

		currTxID := fmt.Sprintf("%x", tx.ID)

		for outIdx, txVout := range tx.Vout {
			newTxOutMap[outIdx] = txVout
		}

		u[currTxID] = newTxOutMap
	}
}
