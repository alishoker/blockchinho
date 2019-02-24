package main

import (
	"flag"
	"fmt"
	. "github.com/alishoker/blockchinho/blockchain"
	"log"
	"os"
	"strconv"
	"text/tabwriter"
)

type CLI struct{}

func (cli *CLI) usage(){

	//format output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, '.', tabwriter.AlignRight|tabwriter.Debug)

	fmt.Println("A command line client interface (CLI) to interact with Blockchinho blockchain.")
	fmt.Println("Usage:")
	fmt.Println("\tcreateblockchain -address 'ADDRESS'\t -- create a new blockchain with rewards to ADDRESS")
	fmt.Println("\taddblock -trans 'TRANSACTIONS'\t -- add a block to the blockchain")
	fmt.Println("\tgetbalance -address 'ADDRESS'\t -- get the current balance of ADDRESS")
	fmt.Println("\tprintchain\t\t\t -- print all the blockchain")
	w.Flush()
}

func (cli *CLI) validateArgs(){
	if len(os.Args) < 2 {
		cli.usage()
		os.Exit(1)
	}
}

/*
func (cli *CLI) addBlock(trans string){
	cli.bc.AddBlock(trans)
	fmt.Println("Block added successfully!")
}
*/
func (cli *CLI) CreateBlockChain(address string){
	bc:=CreateBlockchain(address)
	bc.DB.Close()

	fmt.Println("Blockchain created successfully!")
}

func (cli *CLI) ImportBlockChain(){
	bc:=NewBlockchain()
	defer bc.DB.Close()

	fmt.Println("Blockchain imported successfully!")
}



func (cli *CLI) printBlockchain() {

	bc := NewBlockchain()
	defer bc.DB.Close()

	var block *Block
	bci := bc.Iterator()
	defer bc.DB.Close()

	for{
		block = bci.Next()
		fmt.Println("========================================")
		//fmt.Printf("Block #: %d\n", nb)
		fmt.Printf("Nonce #: %d\n", block.Nonce)
		fmt.Printf("Header: %x\n", block.Header)
		fmt.Printf("Previous Header: %x\n", block.PreHeader)
		fmt.Printf("Transactions: %s\n", block.Transactions)
		pow:=NewProofOfWork(block)
		fmt.Printf("Valid PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		//Genesis reached
		if len(block.PreHeader) == 0 {
			break
		}
	}

}

func (cli *CLI) getBalance(address string) {
	bc := NewBlockchain()
	defer bc.DB.Close()

	balance := 0
	UTXOs := bc.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Current balance of '%s': %d\n", address, balance)
}


func (cli *CLI) Run() {
	cli.validateArgs()

	createBlockchainCmd :=flag.NewFlagSet("createblockchain",flag.ExitOnError)
	createBlockchainAddress :=createBlockchainCmd.String("address","","The address to send reward to")

	addBlockCmd := flag.NewFlagSet("addblock",flag.ExitOnError)
	addBlockTransaction := addBlockCmd.String("trans", "", "Transactions of a block")

	getBalanceCmd := flag.NewFlagSet("getbalance",flag.ExitOnError)
	getBalanceAddress := getBalanceCmd.String("address", "", "Address that holds the balance")

	printBlockchainCmd := flag.NewFlagSet("printblockchain",flag.ExitOnError)

	switch os.Args[1] {
	case "addblock":
		err:=addBlockCmd.Parse(os.Args[2:])
		if err!= nil {
			log.Panic(err)
		}
		if addBlockCmd.Parsed() {
			if *addBlockTransaction == "" {
				addBlockCmd.Usage()
				os.Exit(1)
			}
			//FIXME
			//cli.addBlock(*addBlockTransaction)
		}

	case "printblockchain":
		err := printBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		if printBlockchainCmd.Parsed() {
			cli.printBlockchain()
		}
	case "createblockchain":
		err:=createBlockchainCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)
		}
		if createBlockchainCmd.Parsed(){
			if *createBlockchainAddress==""{
				createBlockchainCmd.Usage()
				os.Exit(1)
			}
			cli.CreateBlockChain(*createBlockchainAddress)

		}
	case "getbalance":
		err:=getBalanceCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)
		}
		if getBalanceCmd.Parsed(){
			if *getBalanceAddress==""{
				getBalanceCmd.Usage()
				os.Exit(1)
			}
			cli.getBalance(*getBalanceAddress)

		}
	default:
		cli.usage()
		os.Exit(1)
	
	}



}