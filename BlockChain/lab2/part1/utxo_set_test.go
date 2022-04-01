package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindSpendableOutputsFromOneOutput(t *testing.T) {
	utxos := getTestExpectedUTXOSet("block0")
	expectedUnspentOutputs := getTestSpendableOutputs(utxos, "rodrigo")
	expectedOut := utxos["e2404638779673c7c3e772e12dc3343e6d38f1d71625419d12a8468522b5cc2d"][0]
	expectedValue := expectedOut.Value

	accumulatedAmount, unspentOutputs := utxos.FindSpendableOutputs("rodrigo", 5)
	assert.Equal(t, expectedValue, accumulatedAmount)
	assert.Equal(t, expectedUnspentOutputs, unspentOutputs)
}

func TestFindSpendableOutputsFromMultipleOutputs(t *testing.T) {
	utxos := getTestExpectedUTXOSet("block2")
	out1 := utxos["4709fe985b8cb59ee67ec3d9ebf968ec82953806332573ba8b2b68d05d9d143a"][1]
	out2 := utxos["ac255688b1df7c2f16fa23cf02f4fe6cb1e500793cc6f9b7d58b58547bfa660c"][0]
	expectedValue := out1.Value + out2.Value
	expectedUnspentOutputs := getTestSpendableOutputs(utxos, "rodrigo")

	accumulatedAmount, unspentOutputs := utxos.FindSpendableOutputs("rodrigo", 5)

	assert.Equal(t, expectedValue, accumulatedAmount)
	assert.Equal(t, expectedUnspentOutputs, unspentOutputs)
}

func TestFindUTXO(t *testing.T) {
	// Rodrigo create a coinbase transaction, receiving 10 "coins"
	utxos := getTestExpectedUTXOSet("block0")

	utxoRodrigo := utxos.FindUTXO("rodrigo")
	assert.Equal(t, []TXOutput{{10, "rodrigo"}}, utxoRodrigo)

	utxoLeander := utxos.FindUTXO("leander")
	assert.Equal(t, []TXOutput(nil), utxoLeander)

	// Rodrigo sent 5 "coins" to Leander
	utxos = getTestExpectedUTXOSet("block1")

	utxoRodrigo = utxos.FindUTXO("rodrigo")
	assert.Equal(t, []TXOutput{{5, "rodrigo"}}, utxoRodrigo)

	utxoLeander = utxos.FindUTXO("leander")
	assert.Equal(t, []TXOutput{{5, "leander"}}, utxoLeander)

	// Rodrigo sent 1 "coin" to Leander and
	// Leander sent 3 "coins" to Rodrigo
	utxos = getTestExpectedUTXOSet("block2")

	utxoRodrigo = utxos.FindUTXO("rodrigo")
	assert.ElementsMatch(t, []TXOutput{
		{4, "rodrigo"},
		{3, "rodrigo"},
	}, utxoRodrigo)
	assert.Equal(t, 2, len(utxoRodrigo))

	utxoLeander = utxos.FindUTXO("leander")
	assert.ElementsMatch(t, []TXOutput{
		{2, "leander"},
		{1, "leander"},
	}, utxoLeander)
	assert.Equal(t, 2, len(utxoLeander))
}

func TestCountUTXOs(t *testing.T) {
	utxos := getTestExpectedUTXOSet("block0")
	assert.Equal(t, 1, utxos.CountUTXOs())

	utxos = getTestExpectedUTXOSet("block1")
	assert.Equal(t, 3, utxos.CountUTXOs())

	utxos = getTestExpectedUTXOSet("block2")
	assert.Equal(t, 5, utxos.CountUTXOs())
}

func TestUpdate(t *testing.T) {
	for k, m := range testUTXOs {
		t.Run(k, func(t *testing.T) {
			utxos := m.utxos
			failMsg := fmt.Sprintf("UTXO update failed for %s", k)
			switch k {
			case "block0":
				utxos.Update(testBlockchainData["block0"].Transactions)
				diff(t, m.expectedUTXOs, utxos, failMsg)
			case "block1":
				utxos.Update(testBlockchainData["block1"].Transactions)
				diff(t, m.expectedUTXOs, utxos, failMsg)
			case "block2":
				utxos.Update(testBlockchainData["block2"].Transactions)
				diff(t, m.expectedUTXOs, utxos, failMsg)
			case "block3":
				utxos.Update(testBlockchainData["block3"].Transactions)
				diff(t, m.expectedUTXOs, utxos, failMsg)
			case "block4":
				utxos.Update(testBlockchainData["block4"].Transactions)
				diff(t, m.expectedUTXOs, utxos, failMsg)
			}
		})
	}
}
