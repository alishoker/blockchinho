package main

import (
	"github.com/alishoker/blockchinho/blockchain"
)

func main() {

	var bc *blockchain.Blockchain
	bc = blockchain.NewBlockchain()
	defer bc.DB.Close()


	cli:= CLI{bc}
	cli.Run()
}
