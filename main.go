package main

import (
	"context"
	"fmt"
	"log"

	token "event-listen-program/contracts"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// Create an instance of the Ethereum client
	// client, err := ethclient.Dial("https://testnet.cronoslabs.com/v1/8ab55f6930dcf5912f5df9c551e07a0dE")

	// TESTNET
	// mode := "TESTNET"
	// wss := "YOUR TESTNET WSS"
	// tokenAddr := "YOUR TESTNET ERC20 TOKEN ADDR"
	// client, err := ethclient.Dial(wss)
	// if err != nil {
	// 	fmt.Println("client dial err....")
	// 	log.Fatal(err)
	// }

	// MAINNET
	mode := "MAINNET"
	wss := "YOUR MAINNET WSS"
	tokenAddr := "YOUR MAINNET ERC20 TOKEN ADDR"
	client, err := ethclient.Dial(wss)
	if err != nil {
		fmt.Println("client dial err....")
		log.Fatal(err)
	}

	// Load the contract ABI and address
	contractAddress := common.HexToAddress(tokenAddr)
	instance, err := token.NewToken(contractAddress, client)
	if err != nil {
		fmt.Println("create instance err....")
		log.Fatal(err)
	}

	// Create a new filter query
	query := &bind.WatchOpts{
		Start:   nil,
		Context: context.Background(),
	}

	// Subscribe to the Transfer event
	eventChan := make(chan *token.TokenTransfer)
	sub, err := instance.WatchTransfer(query, eventChan, nil, nil)
	if err != nil {
		fmt.Println("subscribe event err....")
		log.Fatal(err)
	}

	fmt.Println("Running....")

	// Wait for events
	for {
		select {
		case err := <-sub.Err():
			fmt.Println("sub event listening loop err....")
			log.Fatal(err)
		case event := <-eventChan:
			fmt.Printf("%s : Transfer event received: from=%s, to=%s, value=%v\n", mode,
				event.From.String(), event.To.String(), event.Value)
		}
	}
}

// // VER 2
// package main

// import (
// 	"context"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"strings"

// 	token "event-listen-program/contracts"

// 	"github.com/ethereum/go-ethereum"
// 	"github.com/ethereum/go-ethereum/accounts/abi"
// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/core/types"
// 	"github.com/ethereum/go-ethereum/ethclient"
// )

// func main() {
// 	// MAINNET
// 	// Create an instance of the Ethereum client
// 	client, err := ethclient.Dial("YOUR MAINNET WSS")
// 	contractAddress := common.HexToAddress("YOUR MAINNET ERC20 TOKEN ADDR")

// 	// TESTNET
// 	// Create an instance of the Ethereum client
// 	client, err := ethclient.Dial("YOUR TESTNET WSS")
// 	contractAddress := common.HexToAddress("YOUR TESTNET ERC20 TOKEN ADDR")

// 	if err != nil {
// 		fmt.Println("Err in ClientDial...")
// 		log.Fatal(err)
// 	}

// 	// Load the contract ABI and address
// 	abiBytes, err := ioutil.ReadFile("ERC20.abi")
// 	if err != nil {
// 		fmt.Println("Err in AbiBytes...")
// 		log.Fatal(err)
// 	}
// 	contractAbi, err := abi.JSON(strings.NewReader(string(abiBytes)))
// 	if err != nil {
// 		fmt.Println("Err in contactsABI...")
// 		log.Fatal(err)
// 	}

// 	// Create a channel to receive the event logs
// 	logsCh := make(chan types.Log)

// 	// Create a filter query
// 	query := ethereum.FilterQuery{
// 		Addresses: []common.Address{contractAddress},
// 		Topics:    [][]common.Hash{{contractAbi.Events["Transfer"].ID}},
// 	}

// 	// Subscribe to the event logs
// 	sub, err := client.SubscribeFilterLogs(context.Background(), query, logsCh)
// 	if err != nil {
// 		fmt.Println("Subscribe Filter Logs err...")
// 		log.Fatal(err)
// 	}
// 	defer sub.Unsubscribe()

// 	fmt.Println("Running...")

// 	// Parse the logs and print the event data
// 	event := &token.TokenTransfer{}
// 	for {
// 		select {
// 		case err := <-sub.Err():
// 			log.Fatal(err)
// 		case log := <-logsCh:
// 			err := contractAbi.UnpackIntoInterface(event, "Transfer", log.Data)
// 			if err != nil {
// 				fmt.Println("Err on Event logs channel")
// 			}
// 			event.From = common.HexToAddress(log.Topics[1].Hex())
// 			event.To = common.HexToAddress(log.Topics[2].Hex())
// 			fmt.Printf("Transfer event received: from=%s, to=%s, value=%v\n",
// 				event.From.String(), event.To.String(), event.Value)
// 		}
// 	}

// }
