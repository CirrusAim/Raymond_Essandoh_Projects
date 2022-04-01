package main

//go:generate abigen --sol ../../contracts/Betting.sol --pkg contract --out ./go-bindings/betting/betting.go

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"reflect"
	"sync"

	contract "betting-cli/go-bindings/betting"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	defaultAddress     = "0x90f8bf6a479f320ead074411a4b0e7944ea8c9c1" //contract owner address
	defaultPKHex       = "4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d" //contract owner PKey
	validOutcomes      = []string{"England", "New Zealand", "Australia", "Pakistan"}
	availableAddresses = []string{
		"0xFFcf8FDEE72ac11b5c542428B35EEF5769C409f0",
		"0x22d491Bde2303f2f43325b2108D26f1eAbA1e32b",
		"0xE11BA2b4D45Eaed5996Cd0823791E0C93114882d",
		"0xd03ea8624C8C5987235048901fB614fDcA89b117",
		"0x95cED938F7991cd0dFcb48F0a06a40FA1aF46EBC",
		"0x3E5e9111Ae8eB78Fe1CC3bb8915d5D461F3Ef9A9",
		"0x28a8746e75304c0780E011BEd21C72cD78cd535E",
		"0xACa94ef8bD5ffEE41947b4585a84BdA5a3d3DA6E",
		"0x1dF62f291b2E969fB0849d99D9Ce41e2F137006e"}

	availablePriKeys = []string{
		"6cbed15c793ce57650b9877cf6fa156fbef513c4e6134f022a85b1ffdd59b2a1",
		"6370fd033278c143179d81c5526140625662b8daa446c22ee2d73db3707e620c",
		"646f1ce2fdad0e6deeeb5c7e8e5543bdde65e86029e2fd9fc169899c440a7913",
		"add53f9a7e588d003326d1cbf9e4a43c061aadd9bc938c843a79e7b4fd2ad743",
		"395df67f0c2d2d9fe1ad08d1bc8b6627011959b79c53d7dd6a3536a33ab8a4fd",
		"e485d098507f54e7733a205420dfddbe58db035fa577fc294ebd14db90767a52",
		"a453611d9419d0e56f499079478fd72c37b251a94bfde4d19872c44cf65386e3",
		"829e924fdf021ba3dbbc4225edfece9aca04b929d6e75613329ca6f1d31c0bb4",
		"b0057716d5917badaf911b193b12b910811c1497b5bada8d7711f758981c3773"}

	commands = []string{
		"Deploy", "Choose Oracle", "Make Bet", "Make Decision", "Withdraw",
		"Restart Betting", "Get All Winners", "List All Gamblers", "Check Winnings",
		"Is Oracle", "List Possible Outcomes", "AvailableAddresses", "Quit"} // 
)

type Client struct {
	scanner         *bufio.Scanner
	backend         *ethclient.Client
	walletAddress   common.Address
	bettingInstance *contract.Betting
}

func connect() *ethclient.Client {
	backend, err := ethclient.Dial("ws://127.0.0.1:7545")
	if err != nil {
		log.Fatal(err)
	}
	return backend
}

func (c Client) getAuth(privateKey *ecdsa.PrivateKey) *bind.TransactOpts {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := c.backend.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := c.backend.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(6721975) // in units
	auth.GasPrice = gasPrice
	return auth
}

func (c Client) deploy(auth *bind.TransactOpts) (common.Address, *types.Transaction, *contract.Betting, error) {
	var validOutcomesBytes [][32]byte = make([][32]byte, 4)
	for i := 0; i < len(validOutcomes); i++ {
		copy(validOutcomesBytes[i][:], []byte(validOutcomes[i]))
	}
	return contract.DeployBetting(auth, c.backend, validOutcomesBytes)
}

func (c Client) getBalance(address common.Address) *big.Int {
	balance, err := c.backend.BalanceAt(context.Background(), address, nil)
	if err != nil {
		log.Fatal(err)
	}
	return balance
}

// subscribe and listen to contract events
func (c Client) listenBettingEvents(address common.Address, quit chan struct{}) {
	errs := make(chan error, 1)
	logs := make(chan types.Log)

	sub, err := c.backend.SubscribeFilterLogs(context.Background(), ethereum.FilterQuery{
		Addresses: []common.Address{address},
	}, logs)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			errs <- err
		case evLog := <-logs:
			abi, err := contract.BettingMetaData.GetAbi()
			if err != nil {
				errs <- err
			}
			eventType, err := abi.EventByID(evLog.Topics[0])
			if err != nil {
				errs <- err
			}
			fmt.Println("Events:")
			switch eventType.Name {
			case "OracleChanged":
				event, err := c.bettingInstance.ParseOracleChanged(evLog)
				if err != nil {
					errs <- err
				}
				fmt.Printf("Oracle changed from %v to : %v\n", event.PreviousOracle, event.NewOracle)
			case "BetMade":
				event, err := c.bettingInstance.ParseBetMade(evLog)
				if err != nil {
					errs <- err
				}
				fmt.Printf("Bet received from %v on outcome : %v of amount: %v\n", event.Gambler, string(event.Outcome[:]), event.Amount)
			case "Withdrawn":
				event, err := c.bettingInstance.ParseWithdrawn(evLog)
				if err != nil {
					errs <- err
				}
				fmt.Printf("%v wei withdrawn by %v \n", event.Amount, event.Gambler)
				// case "Winners":
				// 	event, err := c.bettingInstance.ParseWinners(evLog)
				// 	if err != nil {
				// 		errs <- err
				// 	}
				// 	fmt.Printf("Decision made. Win amount is  : %v. Winners are : %+v\n", event.TotalPrize, event.Wins)
			}
		case err := <-errs:
			log.Fatal(err)
		case <-quit:
			return
		}
	}
}

func isZeroAddress(address common.Address) bool {
	return reflect.DeepEqual(address.Bytes(), common.FromHex("0x0000000000000000000000000000000000000000"))
}

func main() {
	var cmd int
	var wg sync.WaitGroup

	client := &Client{}
	client.backend = connect()
	defer client.backend.Close()

	// Get first account as the deployer/sender account
	senderPK, err := crypto.HexToECDSA(defaultPKHex)
	if err != nil {
		log.Fatal(err)
	}

	addKeyMap := make(map[string]string)
	for i, add := range availableAddresses {
		addKeyMap[add] = availablePriKeys[i]
	}

	client.scanner = bufio.NewScanner(os.Stdin)
	quit := make(chan struct{})

	var oracleAdd string
	for {
		defaultAccount := common.HexToAddress(defaultAddress)
		balance := client.getBalance(defaultAccount)
		fmt.Println("------------------")
		fmt.Printf("Balance of owner account %s : %v\n", defaultAccount.Hex(), balance)
		fmt.Println("------------------\nChoose a command:")
		for i, c := range commands {
			fmt.Printf("(%v) %s\n", i+1, c)
		}
		fmt.Scanln(&cmd)
		fmt.Println("------------------")
		switch cmd {
		case 1:
			if !isZeroAddress(client.walletAddress) {
				fmt.Printf("Contract already deployed at: %v\n", client.walletAddress)
				continue
			}
			wg.Add(1)

			auth := client.getAuth(senderPK)
			address, _, instance, err := client.deploy(auth)
			if err != nil {
				fmt.Printf("An error occur: %v\n", err)
				continue
			}
			client.walletAddress = address
			client.bettingInstance = instance
			fmt.Println("Contract deployed at:", address.Hex())

			// listen to all contract events
			go func() {
				client.listenBettingEvents(address, quit)
				defer wg.Done()
			}()

			fmt.Println("\n-----------------------------------------\n")
		case 2:
			fmt.Println("Enter the oracle address")
			client.scanner.Scan()
			oracleAdd = client.scanner.Text()
			oracle := common.HexToAddress(oracleAdd)
			auth := client.getAuth(senderPK)
			_, err := client.bettingInstance.ChooseOracle(&bind.TransactOpts{
				From:   auth.From,
				Signer: auth.Signer,
			}, oracle)

			if err != nil {
				fmt.Printf("An error occur: %v\n", err)
				oracleAdd = ""
				continue
			}
			fmt.Printf("Oracle chosen successfully\n")
		case 3:
			fmt.Println("Enter the gambler address")
			client.scanner.Scan()
			gambler := client.scanner.Text()

			fmt.Println("Enter the gambler's outcome")
			client.scanner.Scan()
			outcome := client.scanner.Text()

			var outcomeByte [32]byte
			copy(outcomeByte[:], []byte(outcome))
			fmt.Println("Enter the bet amount (in wei):")
			client.scanner.Scan()
			amount, _ := big.NewInt(0).SetString(client.scanner.Text(), 10)

			gamblerPk := addKeyMap[gambler]

			pk, err := crypto.HexToECDSA(gamblerPk)
			if err != nil {
				log.Fatal(err)
			}
			auth := client.getAuth(pk)
			tx, err := client.bettingInstance.MakeBet(&bind.TransactOpts{
				From:   auth.From,
				Signer: auth.Signer,
				Value:  amount,
			}, outcomeByte)
			if err != nil {
				fmt.Printf("An error occur: %v\n", err)
				continue
			}
			fmt.Printf("Transaction 0x%x successfully created\n", tx.Hash())

		case 4:
			// Make Decision
			fmt.Println("Enter winning outcome")
			client.scanner.Scan()
			outcome := client.scanner.Text()

			var outcomeByte [32]byte
			copy(outcomeByte[:], []byte(outcome))

			oraclePk := addKeyMap[oracleAdd]
			pk, err := crypto.HexToECDSA(oraclePk)
			if err != nil {
				log.Fatal(err)
			}
			auth := client.getAuth(pk)
			tx, err := client.bettingInstance.MakeDecision(&bind.TransactOpts{
				From:   auth.From,
				Signer: auth.Signer,
			}, outcomeByte)
			if err != nil {
				fmt.Printf("An error occur: %v\n", err)
				continue
			}
			fmt.Printf("Transaction 0x%x successfully created\n", tx.Hash())
		case 5:
			fmt.Println("Enter the winner's address")
			client.scanner.Scan()
			winner := client.scanner.Text()
			fmt.Println("Enter the amount to withdraw (in wei):")
			client.scanner.Scan()
			amount, _ := big.NewInt(0).SetString(client.scanner.Text(), 10)

			winnerPk := addKeyMap[winner]

			pk, err := crypto.HexToECDSA(winnerPk)
			if err != nil {
				log.Fatal(err)
			}
			auth := client.getAuth(pk)
			tx, err := client.bettingInstance.Withdraw(&bind.TransactOpts{
				From:   auth.From,
				Signer: auth.Signer,
			}, amount)

			if err != nil {
				fmt.Printf("An error occur: %v\n", err)
				continue
			}
			fmt.Printf("Transaction 0x%x successfully created\n", tx.Hash())

		case 6:
			auth := client.getAuth(senderPK)
			tx, err := client.bettingInstance.ContractReset(&bind.TransactOpts{
				From:   auth.From,
				Signer: auth.Signer,
			})

			if err != nil {
				fmt.Printf("An error occur: %v\n", err)
				continue
			}
			fmt.Printf("Transaction 0x%x successfully created\n", tx.Hash())
			fmt.Printf("Contract reset successfully")

		case 7:
			auth := client.getAuth(senderPK)
			winners, err := client.bettingInstance.GetWinners(&bind.CallOpts{
				From: auth.From,
			})

			if err != nil {
				fmt.Printf("An error occur: %v\n", err)
				continue
			}
			fmt.Printf("Winners are : %+v\n", winners)

		case 8:
			auth := client.getAuth(senderPK)
			winners, err := client.bettingInstance.GetGamblers(&bind.CallOpts{
				From: auth.From,
			})

			if err != nil {
				fmt.Printf("An error occur: %v\n", err)
				continue
			}
			fmt.Printf("Gamblers are : %+v\n", winners)
		case 9:
			// Check winnings
			fmt.Println("Enter the gambler address")
			client.scanner.Scan()
			gambler := common.HexToAddress(client.scanner.Text())

			winAmount, err := client.bettingInstance.CheckWinnings(&bind.CallOpts{
				From: gambler,
			})
			if err != nil {
				fmt.Printf("An error occur: %v\n", err)
				continue
			}
			fmt.Printf("Gambler %v has won : %v\n", gambler, winAmount)
		case 10:
			fmt.Println("Enter the oracle address")
			client.scanner.Scan()
			oracle := common.HexToAddress(client.scanner.Text())
			auth := client.getAuth(senderPK)
			isOracle, err := client.bettingInstance.IsOracle(&bind.CallOpts{
				From: auth.From,
			}, oracle)

			if err != nil {
				fmt.Printf("An error occur: %v\n", err)
				continue
			}
			if isOracle {
				fmt.Printf("Given address is the Oracle\n")
			} else {
				fmt.Printf("Given address is not the Oracle\n")
			}
		case 11:
			// List Possible Outcomes

			allOutcomes, err := client.bettingInstance.GetOutcomes(&bind.CallOpts{
				From: common.HexToAddress(defaultAddress),
			})
			if err != nil {
				fmt.Printf("An error occur: %v\n", err)
				continue
			}
			fmt.Printf("Possible outcomes are: ")
			for _, c := range allOutcomes {
				fmt.Printf("'%v'\t", string(c[:]))
			}
			fmt.Println("\n-------------------------------------\n")
		case 12:
			fmt.Println("--------List of available address--------\n")
			for _, add := range availableAddresses {
				fmt.Printf("%s\t", add)
			}
			fmt.Println("\n-----------------------------------------\n")
		case 13:
			quit <- struct{}{}
			wg.Wait()
			return

		}
	}
}
