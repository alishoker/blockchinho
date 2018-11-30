package main

import (
	"flag"
	"fmt"
	"github.com/alishoker/blockchinho/blockchain"
	"log"
	"os"
	"strconv"
)

type CLI struct{
	bc *blockchain.Blockchain
}

func newCLI(bc *blockchain.Blockchain) *CLI {
	return &CLI{bc}
}

func (cli *CLI) usage(){
	fmt.Println("A command line client interface (CLI) to interact with Blockchinho blockchain.")
	fmt.Println("Usage:")
	fmt.Println("\taddblock\t trans TRANSACTIONS\t -- add a block to the blockchain")
	fmt.Println("\tprintchain\t\t\t\t -- print all the blockchain")
}

func (cli *CLI) validateArgs(){
	if len(os.Args) < 2 {
		cli.usage()
		os.Exit(1)
	}
}

func (cli *CLI) addBlock(trans string){
	cli.bc.AddBlock(trans)
	fmt.Println("Block added successfully!")
}

func (cli *CLI) printBlockchain() {

	var block *blockchain.Block
	bci := cli.bc.Iterator()

	for{
		block = bci.Next()
		fmt.Println("========================================")
		//fmt.Printf("Block #: %d\n", nb)
		fmt.Printf("Nonce #: %d\n", block.Nonce)
		fmt.Printf("Header: %x\n", block.Header)
		fmt.Printf("Previous Header: %x\n", block.PreHeader)
		fmt.Printf("Transactions: %s\n", block.Transactions)
		pow:=blockchain.NewProofOfWork(block)
		fmt.Printf("Valid PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		//Genesis reached
		if len(block.PreHeader) == 0 {
			break
		}
	}

}

func (cli *CLI) Run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addblock",flag.ExitOnError)
	addBlockTransaction := addBlockCmd.String("trans", "", "Transactions")

	printBlockchainCmd := flag.NewFlagSet("printblockchain",flag.ExitOnError)

	switch os.Args[1] {
	case "addblock":
		err:=addBlockCmd.Parse(os.Args[2:])
		if err!= nil {
			log.Panic(err)
		}
		if addBlockCmd.Parsed() {
			if *addBlockTransaction == "" {
				cli.usage()
				os.Exit(1)
			}
			cli.addBlock(*addBlockTransaction)
		}

	case "printblockchain":
		err := printBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		if printBlockchainCmd.Parsed() {
			cli.printBlockchain()
		}

	default:
		cli.usage()
		os.Exit(1)
	
	}



}