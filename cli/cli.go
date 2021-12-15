package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/tensor-programming/golang-blockchain/blockchain"
	"github.com/tensor-programming/golang-blockchain/wallet"
)

type CommandLine struct {
	//blockchain *blockchain.BlockChain
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	//fmt.Println(" add -block BLOCK_DATA - add a block to the chain")
	fmt.Println("print - print all the blocks of the chain")
	fmt.Println("getbalance - get the balance of the wallet")
	fmt.Println("createwallet - create a new wallet")
	fmt.Println("listaddresses - list all addresses")
	fmt.Println("send - send amount of coins from FROM address to TO address")
	fmt.Println("createblockchain -address ADDRESS creates a blockchain and sends genesis reward to address")
	fmt.Println("printchain - print the blockchain")

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
			fmt.Println("Only Genesis block Exists")
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

//Get all addresss
func (cli *CommandLine) listAddresses() {
	wallets, err := wallet.CreateWallet()
	if err != nil {
		fmt.Println(err)
	}
	addresses := wallets.GetAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}

	fmt.Printf("Total number of addresses: %d\n", len(addresses))
}

//Create Wallet/Address
func (cli *CommandLine) createWallet() {
	wallets, err := wallet.CreateWallet()
	if err != nil {
		fmt.Println(err)
	}
	address := wallets.AddWallet()
	wallets.SaveToFile()

	fmt.Printf("Your new address: %s\n", address)

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
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
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
		//fmt.Println("Hello World")
		err := createBlockchainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "send":
		err := sendCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "listaddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	default:
		cli.printUsage()
		os.Exit(1)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockChain(*createBlockchainAddress)
	}

	if createWalletCmd.Parsed() {
		cli.createWallet()
	}

	if listAddressesCmd.Parsed() {
		cli.listAddresses()
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
