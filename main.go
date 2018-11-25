package main

import (
	"fmt"
	"github.com/alishoker/blockchinho/blockchain"
	"strconv"
)

func main() {

	bc := blockchain.NewBlockchain()
	for i := 1; i < 2; i++ {
		bc.AddBlock(fmt.Sprintf("Transactions %d", i))

	}

	for nb, block := range bc.Blocks {
		fmt.Println("========================================")
		fmt.Printf("Block #: %d\n", nb)
		fmt.Printf("Nonce #: %d\n", block.Nonce)
		fmt.Printf("Header: %x\n", block.Header)
		fmt.Printf("Previous Header: %x\n", block.PreHeader)
		fmt.Printf("Transactions: %s\n", block.Transactions)
		pow:=blockchain.NewProofOfWork(block)
		fmt.Printf("Valid PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

	}

}
