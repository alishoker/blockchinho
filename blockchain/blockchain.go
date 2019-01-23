package blockchain

import (
	"fmt"
	"github.com/alishoker/blockchinho/transaction"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

const (
	bucketBlocks = "blocks"
	dbFile = "blockchain.DB"
)

var keyLastBlock = []byte("l")

type Blockchain struct {

	lastBlockHeader []byte
	DB              *bolt.DB
}

type BlockchainIterator struct{
	currentHash []byte
	db *bolt.DB
}

func (bc *Blockchain) AddBlock(trans []*transaction.Transaction) {

	var lastHeader []byte

	err := bc.DB.View(func(tx *bolt.Tx) error {
		buck := tx.Bucket([]byte(bucketBlocks))
		lastHeader = buck.Get(keyLastBlock)

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	newBlock := NewBlock(lastHeader,trans)

    err = bc.DB.Update(func(tx *bolt.Tx) error {

		buck:= tx.Bucket([]byte(bucketBlocks))

		err = buck.Put(newBlock.Header,newBlock.Serialize())
		if err != nil {
			log.Fatal(err)
		}

		err = buck.Put(keyLastBlock,newBlock.Header)
		if err != nil {
			log.Fatal(err)
		}

		bc.lastBlockHeader = newBlock.Header

		return nil

	})

	if err != nil {
		log.Fatal(err)
	}

}


func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.lastBlockHeader, bc.DB}
}

func (bci *BlockchainIterator) Next() *Block {
	var block *Block

	err := bci.db.View(func(tx *bolt.Tx) error {
		buck := tx.Bucket([]byte(bucketBlocks))
		serialBlock := buck.Get(bci.currentHash)
		block= DeserializeBlock(serialBlock)

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	bci.currentHash = block.PreHeader

	return block

}

func NewBlockchain() *Blockchain {

	if dbExists(){
		fmt.Println("Blockchain is not existing. Please create one.")
		os.Exit(1)
	}

	var lastHash []byte

	db, err := bolt.Open(dbFile,0600,nil)

	if err != nil {
		log.Fatal(err)
	}
	//We need to keep DB opened, close it at the caller of this method
	//defer DB.Close()

	err = db.Update(func(tx *bolt.Tx) error {

		var bucket *bolt.Bucket
		var err error

		//b := tx.CreateBucketIfNotExists([]byte(bucketBlocks))

		bucket, err = tx.CreateBucketIfNotExists([]byte(bucketBlocks))
		if err != nil {
			log.Panic("Blockchain Update:", err)
		}

		if lastHash = bucket.Get(keyLastBlock); lastHash == nil {
			fmt.Println("Blockchain is not existing. Please create one.")
			log.Fatal(err)
		}

		return nil //closure

	})

	if err != nil {
		log.Panic(err)
	}


	return &Blockchain{lastHash,db}
}

func dbExists() bool{


	if _,err:=os.Stat(dbFile);os.IsNotExist(err){
		return false
	}
	return true

}
func CreateBlockchain(address string) *Blockchain {


	if dbExists(){
		fmt.Println("Blockchain already exists. Exiting...")
		os.Exit(1)
	}


	var lastHash []byte

	db, err := bolt.Open(dbFile,0600,nil)

	if err != nil {
		log.Fatal(err)
	}
	//We need to keep DB opened, close it at the caller of this method
	//defer DB.Close()

	err = db.Update(func(tx *bolt.Tx) error {

		var bucket *bolt.Bucket
		var err error

		//b := tx.CreateBucketIfNotExists([]byte(bucketBlocks))

		bucket, err = tx.CreateBucketIfNotExists([]byte(bucketBlocks))
		if err != nil {
			log.Panic("Creating Blockchain bucket failed:", err)
		}

		if lastHash = bucket.Get(keyLastBlock); lastHash == nil {

			fmt.Println("Creating a brand new Blockchain.")

			coinbaseTX:=transaction.NewCoinbaseTX(address,"The Genesis Block")

			genesis := NewBlock([]byte{}, []*transaction.Transaction{coinbaseTX})

			err = bucket.Put(genesis.Header, genesis.Serialize())
			if err != nil {
				log.Panic(err)
			}

			err = bucket.Put(keyLastBlock,genesis.Header)
			if err != nil {
				log.Panic(err)
			}

			lastHash = genesis.Header

		}
		return nil //closure

	})

	if err != nil {
		log.Panic(err)
	}


	return &Blockchain{lastHash,db}
}
