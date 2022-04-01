package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Fixed block timestamp
const TestBlockTime int64 = 1563897484

// Format error message. Based on:
// https://cs.opensource.google/go/go/+/refs/tags/go1.17:src/testing/testing.go;l=537
func errorf(file string, line int, s string) string {
	if line == 0 {
		line = 1
	}
	buf := new(strings.Builder)
	// Every line is indented at least 4 spaces.
	buf.WriteString("    ")
	fmt.Fprintf(buf, "%s:%d: ", filepath.Base(file), line)
	lines := strings.Split(s, "\n")
	if l := len(lines); l > 1 && lines[l-1] == "" {
		lines = lines[:l-1]
	}
	for i, line := range lines {
		if i > 0 {
			// Second and subsequent lines are indented an additional 4 spaces.
			buf.WriteString("\n        ")
		}
		buf.WriteString(line)
	}
	buf.WriteByte('\n')
	return buf.String()
}

func diff(t *testing.T, want interface{}, got interface{}, message string) {
	if diff := cmp.Diff(want, got); diff != "" {
		_, file, line, _ := runtime.Caller(1)
		fmt.Print(errorf(file, line, fmt.Sprintf("%s: (-want +got)\n%s", message, diff)))
		t.Fail()
	}
}

// Transactions example flow:
// tx0: genesis coinbase tx - Rodrigo received 10 coins
// tx1: Rodrigo sent 5 coins to Leander and get 5 coins as remainder
// tx2: Using tx1 output, Leander sent 3 of his coins to Rodrigo
// and get 2 coins as remainder
// tx3: Using tx1 output, Rodrigo sent 1 coin to Leander and
// get 4 coins as remainder
// tx4: Using tx2 output, Rodrigo sent 2 coins to Leander and
// get 1 coin as remainder
// tx5: Using tx3 and tx4 outputs, Leander sent 3 "coins" to Rodrigo
var testTransactions = map[string]*Transaction{
	"tx0": { // genesis transaction
		ID: Hex2Bytes("e2404638779673c7c3e772e12dc3343e6d38f1d71625419d12a8468522b5cc2d"),
		Vin: []TXInput{
			{Txid: nil, OutIdx: -1, ScriptSig: GenesisCoinbaseData},
		},
		Vout: []TXOutput{
			{Value: BlockReward, ScriptPubKey: "rodrigo"},
		},
	},
	"tx1": {
		ID: Hex2Bytes("bce268225bc12a0015bcc39e91d59f47fd176e64ca42e4f8aecf107fe38f3bfa"),
		Vin: []TXInput{
			{
				Txid:      Hex2Bytes("e2404638779673c7c3e772e12dc3343e6d38f1d71625419d12a8468522b5cc2d"),
				OutIdx:    0,
				ScriptSig: "rodrigo",
			},
		},
		Vout: []TXOutput{
			{Value: 5, ScriptPubKey: "leander"},
			{Value: 5, ScriptPubKey: "rodrigo"},
		},
	},
	"tx2": {
		ID: Hex2Bytes("ac255688b1df7c2f16fa23cf02f4fe6cb1e500793cc6f9b7d58b58547bfa660c"),
		Vin: []TXInput{
			{
				Txid:      Hex2Bytes("bce268225bc12a0015bcc39e91d59f47fd176e64ca42e4f8aecf107fe38f3bfa"),
				OutIdx:    0,
				ScriptSig: "leander",
			},
		},
		Vout: []TXOutput{
			{Value: 3, ScriptPubKey: "rodrigo"},
			{Value: 2, ScriptPubKey: "leander"},
		},
	},
	"tx3": {
		ID: Hex2Bytes("4709fe985b8cb59ee67ec3d9ebf968ec82953806332573ba8b2b68d05d9d143a"),
		Vin: []TXInput{
			{
				Txid:      Hex2Bytes("bce268225bc12a0015bcc39e91d59f47fd176e64ca42e4f8aecf107fe38f3bfa"),
				OutIdx:    1,
				ScriptSig: "rodrigo",
			},
		},
		Vout: []TXOutput{
			{Value: 1, ScriptPubKey: "leander"},
			{Value: 4, ScriptPubKey: "rodrigo"},
		},
	},
	"tx4": {
		ID: Hex2Bytes("536532f5e3a5043e2de1be480761602b8218bf7abf374ece8794e0ca0d5b072b"),
		Vin: []TXInput{
			{
				Txid:      Hex2Bytes("ac255688b1df7c2f16fa23cf02f4fe6cb1e500793cc6f9b7d58b58547bfa660c"),
				OutIdx:    0,
				ScriptSig: "rodrigo",
			},
		},
		Vout: []TXOutput{
			{Value: 2, ScriptPubKey: "leander"},
			{Value: 1, ScriptPubKey: "rodrigo"},
		},
	},
	"tx5": {
		ID: Hex2Bytes("33481b10feed4a93eab2233c68cf119281ac23ed268f1389e076b426ba8b412a:"),
		Vin: []TXInput{
			{
				Txid:      Hex2Bytes("536532f5e3a5043e2de1be480761602b8218bf7abf374ece8794e0ca0d5b072b"),
				OutIdx:    0,
				ScriptSig: "leander",
			},
			{
				Txid:      Hex2Bytes("4709fe985b8cb59ee67ec3d9ebf968ec82953806332573ba8b2b68d05d9d143a"),
				OutIdx:    0,
				ScriptSig: "leander",
			},
		},
		Vout: []TXOutput{
			{Value: 3, ScriptPubKey: "rodrigo"},
		},
	},
}

func getTestInputsTX(tx string) []TXInput {
	return testTransactions[tx].Vin
}

func newMockCoinbaseTX(to, data, txID string) *Transaction {
	tx := &Transaction{
		ID: Hex2Bytes(txID),
		Vin: []TXInput{
			{
				Txid:      nil,
				OutIdx:    -1,
				ScriptSig: data,
			},
		},
		Vout: []TXOutput{
			{
				Value:        BlockReward,
				ScriptPubKey: to,
			},
		},
	}
	return tx
}

var minerCoinbaseTx = map[string]*Transaction{
	"tx1": newMockCoinbaseTX("miner", "1", "9a57281b774df79677d0266d3e740a3133aad5956c5a4020d9e37dc49755e469"),
	"tx2": newMockCoinbaseTX("miner", "2", "c4e66ff1bd12cf01914f6baa83d68a3ef33e4e3f5838762d4cde5cbd75c06de5"),
	"tx3": newMockCoinbaseTX("miner", "3", "fb8be9697a0397f97d83c4300ccc7ac9757742e7fb765a708a2a4c67252cf904"),
	"tx4": newMockCoinbaseTX("miner", "4", "ccbacaced08ce19dbef6f69b1eeb41bd6d7739d4d17eeffe6cb1cd612d45adee"),
}

var testBlockchainData = map[string]*Block{
	"block0": { // genesis block
		Timestamp: TestBlockTime,
		Transactions: []*Transaction{
			testTransactions["tx0"],
		},
		PrevBlockHash: nil,
		Hash:          Hex2Bytes("00ebe2e15ef92721e3940c99b2f068b52e527ecdf549abddf8d162f646ced2c4"),
		Nonce:         352,
	},
	"block1": {
		Timestamp: TestBlockTime,
		Transactions: []*Transaction{
			minerCoinbaseTx["tx1"],
			testTransactions["tx1"],
		},
		PrevBlockHash: Hex2Bytes("00ebe2e15ef92721e3940c99b2f068b52e527ecdf549abddf8d162f646ced2c4"),
		Hash:          Hex2Bytes("0076483ed21f560b5418f90ac263bd58b94acf0ee643694714245a882a4d3b17"),
		Nonce:         270,
	},
	"block2": {
		Timestamp: TestBlockTime,
		Transactions: []*Transaction{
			minerCoinbaseTx["tx2"],
			testTransactions["tx3"],
			testTransactions["tx2"],
		},
		PrevBlockHash: Hex2Bytes("0076483ed21f560b5418f90ac263bd58b94acf0ee643694714245a882a4d3b17"),
		Hash:          Hex2Bytes("006f8f853d2518b039c837a0eafddb80a93b2492d20b376b4a7b11d835fa7b69"),
		Nonce:         168,
	},
	"block3": {
		Timestamp: TestBlockTime,
		Transactions: []*Transaction{
			minerCoinbaseTx["tx3"],
			testTransactions["tx4"],
		},
		PrevBlockHash: Hex2Bytes("006f8f853d2518b039c837a0eafddb80a93b2492d20b376b4a7b11d835fa7b69"),
		Hash:          Hex2Bytes("004219baa29e7de53fa3a72eadadacea1068ccd6e0a57db75b7f9773a11e4316"),
		Nonce:         65,
	},
	"block4": {
		Timestamp: TestBlockTime,
		Transactions: []*Transaction{
			minerCoinbaseTx["tx4"],
			testTransactions["tx5"],
		},
		PrevBlockHash: Hex2Bytes("004219baa29e7de53fa3a72eadadacea1068ccd6e0a57db75b7f9773a11e4316"),
		Hash:          Hex2Bytes("00d52a95c152b4be2895f555fdb1eb335ff7e5d9ba9afd63ceef1efffc9b2563"),
		Nonce:         1783,
	},
}

var testUTXOs = map[string]struct {
	utxos         UTXOSet
	expectedUTXOs UTXOSet
}{
	"block0": { // (0 input -> 1 output, generating "coins"): nil -> tx0
		utxos: UTXOSet{},
		expectedUTXOs: UTXOSet{
			"e2404638779673c7c3e772e12dc3343e6d38f1d71625419d12a8468522b5cc2d": {0: testTransactions["tx0"].Vout[0]},
		},
	},
	"block1": { // (1 input -> 2 outputs, splitting one input): tx0 -> tx1
		utxos: UTXOSet{
			"e2404638779673c7c3e772e12dc3343e6d38f1d71625419d12a8468522b5cc2d": {0: testTransactions["tx0"].Vout[0]},
		},
		expectedUTXOs: UTXOSet{
			"bce268225bc12a0015bcc39e91d59f47fd176e64ca42e4f8aecf107fe38f3bfa": {
				0: testTransactions["tx1"].Vout[0],
				1: testTransactions["tx1"].Vout[1],
			},
			"9a57281b774df79677d0266d3e740a3133aad5956c5a4020d9e37dc49755e469": {
				0: {Value: 10, ScriptPubKey: "miner"},
			},
		},
	},
	"block2": { // (1 input -> 2 output, with multiple txs): tx1 -> tx2,tx3
		utxos: UTXOSet{
			"bce268225bc12a0015bcc39e91d59f47fd176e64ca42e4f8aecf107fe38f3bfa": {
				0: testTransactions["tx1"].Vout[0],
				1: testTransactions["tx1"].Vout[1],
			},
		},
		expectedUTXOs: UTXOSet{
			"4709fe985b8cb59ee67ec3d9ebf968ec82953806332573ba8b2b68d05d9d143a": {
				0: testTransactions["tx3"].Vout[0],
				1: testTransactions["tx3"].Vout[1],
			},
			"ac255688b1df7c2f16fa23cf02f4fe6cb1e500793cc6f9b7d58b58547bfa660c": {
				0: testTransactions["tx2"].Vout[0],
				1: testTransactions["tx2"].Vout[1],
			},
			"c4e66ff1bd12cf01914f6baa83d68a3ef33e4e3f5838762d4cde5cbd75c06de5": {
				0: {Value: 10, ScriptPubKey: "miner"},
			},
		},
	},
	// NOTE: Some outputs were ignored from the initial utxos to reduce combination of possible input sources
	"block3": { // (1 input -> 2 outputs): tx2 -> tx4
		utxos: UTXOSet{
			// tx3 was intentionally ignored
			"ac255688b1df7c2f16fa23cf02f4fe6cb1e500793cc6f9b7d58b58547bfa660c": {
				0: testTransactions["tx2"].Vout[0],
				1: testTransactions["tx2"].Vout[1],
			},
		},
		expectedUTXOs: UTXOSet{
			"ac255688b1df7c2f16fa23cf02f4fe6cb1e500793cc6f9b7d58b58547bfa660c": {1: testTransactions["tx2"].Vout[1]},
			"536532f5e3a5043e2de1be480761602b8218bf7abf374ece8794e0ca0d5b072b": {
				0: testTransactions["tx4"].Vout[0],
				1: testTransactions["tx4"].Vout[1],
			},
			"fb8be9697a0397f97d83c4300ccc7ac9757742e7fb765a708a2a4c67252cf904": {
				0: {Value: 10, ScriptPubKey: "miner"},
			},
		},
	},
	// NOTE: Some outputs were ignored from the initial utxos to reduce combination of possible input sources
	"block4": { // (2 inputs -> 1 output): tx3,tx4 -> tx5
		utxos: UTXOSet{
			"4709fe985b8cb59ee67ec3d9ebf968ec82953806332573ba8b2b68d05d9d143a": {
				0: testTransactions["tx3"].Vout[0],
				1: testTransactions["tx3"].Vout[1],
			},
			"536532f5e3a5043e2de1be480761602b8218bf7abf374ece8794e0ca0d5b072b": {
				0: testTransactions["tx4"].Vout[0],
				1: testTransactions["tx4"].Vout[1],
			},
		},
		expectedUTXOs: UTXOSet{
			"4709fe985b8cb59ee67ec3d9ebf968ec82953806332573ba8b2b68d05d9d143a": {1: testTransactions["tx3"].Vout[1]},
			"536532f5e3a5043e2de1be480761602b8218bf7abf374ece8794e0ca0d5b072b": {1: testTransactions["tx4"].Vout[1]},
			"33481b10feed4a93eab2233c68cf119281ac23ed268f1389e076b426ba8b412a": {0: testTransactions["tx5"].Vout[0]},
			"ccbacaced08ce19dbef6f69b1eeb41bd6d7739d4d17eeffe6cb1cd612d45adee": {
				0: {Value: 10, ScriptPubKey: "miner"},
			},
		},
	},
}

func getTestExpectedUTXOSet(block string) UTXOSet {
	return testUTXOs[block].expectedUTXOs
}

func getTestSpendableOutputs(utxos UTXOSet, unlockingData string) map[string][]int {
	unspentOutputs := make(map[string][]int)

	for txID, outputs := range utxos {
		for outIdx, out := range outputs {
			if out.ScriptPubKey == unlockingData {
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
			}
		}
	}
	return unspentOutputs
}
