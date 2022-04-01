package main

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func newMockHeader(prevBlockHash []byte, merkleRoot []byte) []byte {
	return bytes.Join(
		[][]byte{
			prevBlockHash,
			merkleRoot,
			IntToHex(TestBlockTime),
			IntToHex(TARGETBITS),
		},
		[]byte{},
	)
}

// TARGETBITS == 8 => target difficulty of 2^248
// Hexadecimal: 100000000000000000000000000000000000000000000000000000000000000
// Big Int: 452312848583266388373324160190187140051835877600158453279131187530910662656
var testTargetDifficulty, _ = new(big.Int).SetString("452312848583266388373324160190187140051835877600158453279131187530910662656", 10)

func TestNewProofOfWork(t *testing.T) {
	b := &Block{
		Timestamp:    TestBlockTime,
		Transactions: []*Transaction{testTransactions["tx0"]},
	}

	pow := NewProofOfWork(b)

	assert.Equal(t, testTargetDifficulty, pow.target)
	if diff := cmp.Diff(b, pow.block); diff != "" {
		t.Errorf("Wrong block mined: (-want +got)\n%s", diff)
	}
}

func TestSetupHeader(t *testing.T) {
	pow := &ProofOfWork{
		block: &Block{
			Timestamp:    TestBlockTime,
			Transactions: []*Transaction{testTransactions["tx0"]},
		},
		target: testTargetDifficulty,
	}
	header := pow.setupHeader()

	expectedHeader := newMockHeader(nil, Hex2Bytes("f1bfaf87ec0117cd2d90bfb6d039f8ad022c9383f39304f29bf0fb1c5ada2b7b"))
	assert.Equalf(t, expectedHeader, header, "The current block header: %x isn't equal to the expected %x\n", header, expectedHeader)
}

func TestAddNonce(t *testing.T) {
	header := newMockHeader(nil, Hex2Bytes("f1bfaf87ec0117cd2d90bfb6d039f8ad022c9383f39304f29bf0fb1c5ada2b7b"))
	expectedHeader := Hex2Bytes("f1bfaf87ec0117cd2d90bfb6d039f8ad022c9383f39304f29bf0fb1c5ada2b7b000000005d372e8c00000000000000080000000000000009")

	if diff := cmp.Diff(expectedHeader, addNonce(9, header)); diff != "" {
		t.Errorf("AddNonce failed: (-want +got)\n%s", diff)
	}
}

func TestRun(t *testing.T) {
	for k, block := range testBlockchainData {
		t.Run(k, func(t *testing.T) {
			b := &Block{
				Timestamp:     TestBlockTime,
				Transactions:  block.Transactions,
				PrevBlockHash: block.PrevBlockHash,
			}
			pow := &ProofOfWork{b, testTargetDifficulty}
			nonce, hash := pow.Run()
			if diff := cmp.Diff(testBlockchainData[k].Nonce, nonce); diff != "" {
				t.Errorf("diff failed for %q: (-want +got)\n%s", k, diff)
			}
			if diff := cmp.Diff(testBlockchainData[k].Hash, hash); diff != "" {
				t.Errorf("diff failed for %q: (-want +got)\n%s", k, diff)
			}
		})
	}
}

func TestValidatePoW(t *testing.T) {
	for k, block := range testBlockchainData {
		t.Run(k, func(t *testing.T) {
			pow := &ProofOfWork{block, testTargetDifficulty}
			assert.True(t, pow.Validate())
		})
	}
}
