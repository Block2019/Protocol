package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/tensor-programming/golang-blockchain/blockchain"
)

type CommandLine struct {
	//blockchain *blockchain.BlockChain
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	//fmt.Println(" add -block BLOCK_DATA - add a block to the chain")
	//fmt.Println(" print - print all the blocks of the chain")
	fmt.Println(" getbalance - get the balance of the wallet")
	//fmt.Println(" createwallet - create a new wallet")
	//fmt.Println(" listaddresses - list all addresses")
	fmt.Println(" send - send amount of coins from FROM address to TO address")
	fmt.Println(" createblockchain - create a new blockchain address")
	fmt.Println(" printchain - print the blockchain")

}

//Validate Arguments given to the program
func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

//Function to add block to the chain
// func (cli *CommandLine) addBlock(data string) {
// 	cli.blockchain.AddBlock(data)
// 	fmt.Println("Added new block Success!")
// }

//FUnction to printout the blocks of the chain
func (cli *CommandLine) printChain() {
	chain := blockchain.ContinueBlockChain("")
	defer chain.Database.Close()
	bci := chain.Iterator() // Create an iterator

	for {
		block := bci.Next() // Get the next block

		fmt.Printf("Prev. hash: %x\n", block.PrevHash)              // Print the previous hash
		fmt.Printf("Hash: %x\n", block.Hash)                        // Print the hash of the block
		pow := blockchain.NewProof(block)                           // Create a new proof of work
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate())) // Validate the proof of work
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

//Create a new blockchain
func (cli *CommandLine) createBlockChain(address string) {

	fmt.Println("Creating a new blockchain...")
	fmt.Println("Starting mining...")
	chain := blockchain.InitBlockChain(address)
	defer chain.Database.Close()
	fmt.Println("Done!")

}

//FUnction to send coins from one address to another
func (cli *CommandLine) send(from, to string, amount int) {
	chain := blockchain.ContinueBlockChain(from)
	defer chain.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})

	fmt.Println("Success!")
}

//Function get balance
func (cli *CommandLine) getBalance(address string) {
	blockchain := blockchain.ContinueBlockChain(address)
	defer blockchain.Database.Close()

	balance := 0
	UTXOs := blockchain.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

//Run Function
func (cli *CommandLine) Run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)

	//addBlockData := addBlockCmd.String("block", "", "Block data")
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	//createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	sendTo := sendCmd.String("to", "", "Send to address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")
	sendFrom := sendCmd.String("from", "", "Send from address")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "send":
		err := sendCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	default:
		cli.printUsage()
		os.Exit(1)
	}

	if createBlockchainCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}

		cli.createBlockChain(*sendFrom)
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

	if getBalanceCmd.Parsed() {
		if *sendFrom == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}

		cli.getBalance(*sendFrom)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

func main() {

	defer os.Exit(0)
	cli := CommandLine{}
	cli.Run()

	// youthchain.AddBlock("Send 1 BTC to Ivan")
	// youthchain.AddBlock("Send 2 more BTC to Ivan")
	// youthchain.AddBlock("Send 3 more BTC to Ivan")

	// for _, block := range youthchain.Blocks {
	// 	println(string(block.Data))     // Prints the data stored in each block
	// 	println(string(block.Hash))     // Prints the hash of each block
	// 	println(string(block.PrevHash)) // Prints the hash of the previous block

	// 	// pow := blockchain.NewProof(block)                         // Creates a new proof of work
	// 	// fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate)) // Validates the proof of work
	// 	// println()

	// 	pow := blockchain.NewProof(block)
	// 	fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
	// 	println()

	// }

	//println(len(youthchain.Blocks))

}
