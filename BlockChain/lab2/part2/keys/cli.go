package main

import (
	"crypto/ecdsa"
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

const (
	CREATE_BLCKCHN = "create-blockchain"
	ADD_TRAN       = "add-transaction"
	DEMO_TRAN      = "demo-transaction"
	GET_BALANCE    = "get-balance"
	MINE_BLCK      = "mine-block"
	PRINT_CHAIN    = "print-chain"
	PRINT_BLCK     = "print-block"
	PRINT_TRAN     = "print-transaction"
	EXIT           = "exit"
)

func getBalance(address string, utxo UTXOSet) int {
	pubKeyHash := GetPubKeyHashFromAddress(address)
	amt, _ := utxo.FindSpendableOutputs(pubKeyHash, 0)
	return amt
}

func main() {
	peppers := []string{CREATE_BLCKCHN, DEMO_TRAN, GET_BALANCE, PRINT_CHAIN, EXIT}

	templatesSelect := &promptui.SelectTemplates{
		Label:    "{{ . | green }}",
		Active:   "\U000027a4 {{ . | cyan }}",
		Inactive: "  {{ . | cyan }}",
		Selected: " {{ . | cyan }}",
	}

	searcher := func(input string, index int) bool {
		pepper := peppers
		name := strings.Replace(strings.ToLower(pepper[index]), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}
	var blockchain *Blockchain
	var block *Block
	var trans []*Transaction
	var utxo UTXOSet
	aPri, aPub := newKeyPair()
	aAdd := GetStringAddress(GetAddress(aPub))

	bPri, bPub := newKeyPair()
	bAdd := GetStringAddress(GetAddress(bPub))

	cPri, cPub := newKeyPair()
	cAdd := GetStringAddress(GetAddress(cPub))
	for {

		promptSelect := promptui.Select{
			Label:     "Welcome to Blockchain. Please select options",
			Items:     peppers,
			Templates: templatesSelect,
			Size:      len(peppers),
			Searcher:  searcher,
		}

		_, result, err := promptSelect.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		switch result {
		case CREATE_BLCKCHN:
			// create blockchain with 1HrwWkjdwQuhaHSco9H7u7SVsmo4aeDZBX address. which belongs to
			blockchain, _ = NewBlockchain(aAdd)
			utxo = blockchain.FindUTXOSet()
			fmt.Println("Created blockchain")

		case DEMO_TRAN:
			// Demo tran will first create a coinbase trans, then will create a transaction from coinbase address to a new address

			demoBlockAdd := func(from []byte, fromadd, to string, amount int, priKey ecdsa.PrivateKey) {
				trans = []*Transaction{}
				tran, err := NewCoinbaseTX(aAdd, fromadd+to)
				if err != nil {
					fmt.Printf("Unable to create coinbase transaction for block: Err : %+v\n", err)
					return
				}
				trans = append(trans, tran)

				tran, err = NewUTXOTransaction(from, to, amount, utxo)

				if err != nil {
					fmt.Printf("Unable to create transaction from : %s  to %s: Err : %+v\n", fromadd, to, err)
					return
				}
				blockchain.SignTransaction(tran, priKey)
				trans = append(trans, tran)
				utxo.Update(trans)

				block, err = blockchain.MineBlock(trans)
				if err != nil {
					fmt.Printf("Unable to Mine block : Err : %+v\n", err)
				}

				fmt.Println(trans)
			}

			fmt.Printf("Address of a is : %v\n", aAdd)
			fmt.Printf("Address of b is : %v\n", bAdd)
			fmt.Printf("Address of c is : %v\n", cAdd)

			// From A->B
			demoBlockAdd(aPub, aAdd, bAdd, 10, aPri)

			// Should fail,C has insufficient funds so will not be added to blockchain
			demoBlockAdd(cPub, cAdd, bAdd, 3, cPri)

			// From B->C
			demoBlockAdd(bPub, bAdd, cAdd, 5, bPri)

			aBlnc := getBalance(aAdd, utxo)
			bBlnc := getBalance(bAdd, utxo)
			cBlnc := getBalance(cAdd, utxo)
			fmt.Printf("Balance of %s is : %v\n", aAdd, aBlnc)
			fmt.Printf("Balance of %s is : %v\n", bAdd, bBlnc)
			fmt.Printf("Balance of %s is : %v\n", cAdd, cBlnc)

		case GET_BALANCE:
			var addr string
			fmt.Println("Enter address ->")
			fmt.Scanln(&addr)
			senderBlnc := getBalance(addr, utxo)
			fmt.Printf("Balance of address %s is : %v\n", addr, senderBlnc)

		case MINE_BLCK:
			var err error
			block, err = blockchain.MineBlock(trans)
			if err != nil {
				fmt.Println("Error occurred while mining. Skipping mining.Error : " + err.Error())

			}
			fmt.Println("Block mined successfully")
			trans = []*Transaction{}

		case PRINT_CHAIN:
			fmt.Println("Printing blockchain")
			fmt.Println(blockchain)

		case PRINT_BLCK:
			fmt.Println("Printing block")
			fmt.Println(block)

		case PRINT_TRAN:
			fmt.Println("Printing transactions")
			fmt.Println(trans)
		case EXIT:
			fmt.Println("selected : ", EXIT)
			os.Exit(0)
		}
	}
}
