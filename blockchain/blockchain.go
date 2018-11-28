package blockchain

import (
	"fmt"
	"github.com/alishoker/blockchinho/blockchain"
	"github.com/boltdb/bolt"
	"log"
)

const (
	bucketBlocks = "blocks"
	dbFile = "blockchain.db"
)

var keyLastBlock = []byte("l")

type Blockchain struct {

	lastBlockHeader []byte
	db *bolt.DB
}

type BlockchainIterator struct{
	currentHash []byte
	db *bolt.DB
}

func (bc *Blockchain) AddBlock(trans string) {

	var lastHeader []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		buck := tx.Bucket([]byte(bucketBlocks))
		lastHeader = buck.Get(keyLastBlock)

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	newBlock := NewBlock(lastHeader,trans)

    err = db.Update(func(tx *bolt.Tx) error {

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
	return &BlockchainIterator{bc.lastBlockHeader, bc.db}
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

	db, err := bolt.Open(dbFile,0600,nil)

	if err != nil {
		log.Fatal(err)
	}
	//We need to keep db opened
	//defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {

		var bucket bolt.Bucket
		var lastHash []byte
		//b := tx.CreateBucketIfNotExists([]byte(bucketBlocks))

		_, err := bucket.CreateBucketIfNotExists([]byte(bucketBlocks))
		if err != nil {
			log.Panic(err)
		}

		if lastHash = bucket.Get(keyLastBlock); lastHash == nil {

			fmt.Println("Blockchain not existing. Creating a new one...")

			genesis := NewBlock([]byte{}, "The Genesis Block")

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
