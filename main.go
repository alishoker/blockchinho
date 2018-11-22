package main

import (
	"github.com/alishoker/blockchinho/blockchain"
	"fmt"
)

func main() {

	bc := blockchain.NewBlockchain()
	b1:= blockchain.NewBlock(bc.Blocks[len(bc.Blocks)-1].Header,"Transactions 1")
	bc.AddBlock(b1)
	b2:= blockchain.NewBlock(bc.Blocks[len(bc.Blocks)-1].Header,"Transactions 2")
	bc.AddBlock(b2)

	for _, block := range bc.Blocks {
		fmt.Println("==========================")
		fmt.Printf("Header: %x\n", block.Header)
		fmt.Printf("Previous Header: %x\n", block.PreHeader)
		fmt.Printf("Transactions: %s\n", block.Transactions)
	}

}
