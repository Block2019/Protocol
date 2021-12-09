package main

import (
	"fmt"
	"strconv"

	"github.com/tensor-programming/golang-blockchain/blockchain"
)

func main() {

	youthchain := blockchain.InitBlockChain()

	youthchain.AddBlock("Send 1 BTC to Ivan")
	youthchain.AddBlock("Send 2 more BTC to Ivan")
	youthchain.AddBlock("Send 3 more BTC to Ivan")

	for _, block := range youthchain.Blocks {
		println(string(block.Data))     // Prints the data stored in each block
		println(string(block.Hash))     // Prints the hash of each block
		println(string(block.PrevHash)) // Prints the hash of the previous block

		// pow := blockchain.NewProof(block)                         // Creates a new proof of work
		// fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate)) // Validates the proof of work
		// println()

		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		println()

	}

	println(len(youthchain.Blocks))

}
