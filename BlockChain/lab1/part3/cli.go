package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

const (
	CREATE_BLCKCHN = "create-blockchain"
	ADD_TRAN       = "add-transaction"
	MINE_BLCK      = "mine-block"
	PRINT_CHAIN    = "print-chain"
	PRINT_BLCK     = "print-block"
	PRINT_TRAN     = "print-transaction"
	EXIT           = "exit"
)

func main() {
	peppers := []string{CREATE_BLCKCHN, ADD_TRAN, MINE_BLCK, PRINT_CHAIN, PRINT_BLCK, PRINT_TRAN, EXIT}

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
			blockchain = NewBlockchain()
			fmt.Println("Created blockchain")
		case ADD_TRAN:
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("Enter transaction -> ")
			text, _ := reader.ReadString('\n')
			tran := NewTransaction([]byte(text))
			trans = append(trans, tran)
			fmt.Println("transaction  added successfully")

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
			var lines []string
			for _, tx := range trans {
				lines = append(lines, fmt.Sprintf("%s: %x", string(tx.Data), tx.ID))
			}
			strings.Join(lines, "\n")
			fmt.Println(lines)
		case EXIT:
			fmt.Println("selected : ", EXIT)
			os.Exit(0)
		}
	}
}
