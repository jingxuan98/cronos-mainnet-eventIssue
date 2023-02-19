// VER 1
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
	// TESTNET
	// mode := "TESTNET"
	// Create an instance of the Ethereum client
	// client, err := ethclient.Dial("YOUR TESTNET WSS")
	// tokenAddr := "YOUR TESTNET ERC20 TOKEN ADDR"

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

// // VERSION 2

// package main

// import (
// 	"context"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"math/big"
// 	"strings"

// 	token "event-listen-program/contracts"

// 	"github.com/ethereum/go-ethereum"
// 	"github.com/ethereum/go-ethereum/accounts/abi"
// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/ethclient"
// )

// func main() {
// 	/// TESTNET
// 	// // Create an instance of the Ethereum client
// 	// client, err := ethclient.Dial("YOUR TESTNET WSS")
// 	// // Load the contract ABI and address
// 	// contractAddress := common.HexToAddress("YOUR TESTNET ERC20 TOKEN ADDR")
// 	// if err != nil {
// 	// 	fmt.Println("Err in ClientDial...")
// 	// 	log.Fatal(err)
// 	// }

// 	// MAINNET
// 	// Create an instance of the Ethereum client
// 	client, err := ethclient.Dial("YOUR MAINNET WSS")
// 	// Load the contract ABI and address
// 	contractAddress := common.HexToAddress("YOUR MAINNET ERC20 TOKEN ADDR")
// 	if err != nil {
// 		fmt.Println("Err in ClientDial...")
// 		log.Fatal(err)
// 	}

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

// 	// Create a filter query
// 	query := ethereum.FilterQuery{
// 		FromBlock: big.NewInt(7039191),
// 		ToBlock:   big.NewInt(7039191),
// 		Addresses: []common.Address{contractAddress},
// 		Topics:    [][]common.Hash{{contractAbi.Events["Transfer"].ID}},
// 	}

// 	// Retrieve logs that match the filter query
// 	logs, err := client.FilterLogs(context.Background(), query)
// 	if err != nil {
// 		fmt.Println("Err in FilterLogs...")
// 		log.Fatal(err)
// 	}

// 	fmt.Println("Running...")
// 	// Parse the logs and print the event data
// 	event := &token.TokenTransfer{}
// 	for _, log := range logs {
// 		err := contractAbi.UnpackIntoInterface(event, "Transfer", log.Data)
// 		if err != nil {
// 			fmt.Println("Err in logs loop", err)
// 		}
// 		event.From = common.HexToAddress(log.Topics[1].Hex())
// 		event.To = common.HexToAddress(log.Topics[2].Hex())
// 		fmt.Printf("Transfer event received: from=%s, to=%s, value=%v\n",
// 			event.From.String(), event.To.String(), event.Value)
// 	}
// }
