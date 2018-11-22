package main

import (
	"github.com/alishoker/blockchinho/blockchain"
	"fmt"
)

func main() {

	bc := blockchain.NewBlockchain()
	for i:=1;i<5;i++ {
	bc.AddBlock(fmt.Sprintf("Transactions %d",i))
      }

	for nb, block := range bc.Blocks {
		fmt.Println("========================================")
		fmt.Printf("Block #: %d\n", nb)
		fmt.Printf("Header: %x\n", block.Header)
		fmt.Printf("Previous Header: %x\n", block.PreHeader)
		fmt.Printf("Transactions: %s\n", block.Transactions)
	}

}
