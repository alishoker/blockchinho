package blockchain

import (
	"encoding/hex"
	"fmt"
	. "github.com/alishoker/blockchinho/transaction"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

const (
	bucketBlocks = "blocks"
	dbFile       = "blockchain.DB"
)

var keyLastBlock = []byte("l")

type Blockchain struct {
	lastBlockHeader []byte
	DB              *bolt.DB
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (bc *Blockchain) AddBlock(trans []*Transaction) {

	var lastHeader []byte

	err := bc.DB.View(func(tx *bolt.Tx) error {
		buck := tx.Bucket([]byte(bucketBlocks))
		lastHeader = buck.Get(keyLastBlock)

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	newBlock := NewBlock(lastHeader, trans)

	err = bc.DB.Update(func(tx *bolt.Tx) error {

		buck := tx.Bucket([]byte(bucketBlocks))

		err = buck.Put(newBlock.Header, newBlock.Serialize())
		if err != nil {
			log.Fatal(err)
		}

		err = buck.Put(keyLastBlock, newBlock.Header)
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
		block = DeserializeBlock(serialBlock)

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	bci.currentHash = block.PreHeader

	return block

}

func NewBlockchain() *Blockchain {

	if !dbExists() {
		fmt.Println("Blockchain is not existing. Please create one.")
		os.Exit(1)
	}

	var lastHash []byte

	db, err := bolt.Open(dbFile, 0600, nil)

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

	return &Blockchain{lastHash, db}
}

func dbExists() bool {

	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true

}
func CreateBlockchain(address string) *Blockchain {

	if dbExists() {
		fmt.Println("Blockchain already exists. Exiting...")
		os.Exit(1)
	}

	var lastHash []byte

	db, err := bolt.Open(dbFile, 0600, nil)

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

			coinbaseTX := NewCoinbaseTX(address, "The Genesis Block")

			genesis := NewBlock([]byte{}, []*Transaction{coinbaseTX})

			err = bucket.Put(genesis.Header, genesis.Serialize())
			if err != nil {
				log.Panic(err)
			}

			err = bucket.Put(keyLastBlock, genesis.Header)
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

	return &Blockchain{lastHash, db}
}

func (bc *Blockchain) FindUnspentTransactions(address string) []Transaction { //refactor this
	var unspentTXs []Transaction
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				// Check if the output is already spent
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				if out.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			//Gather all inputs that could unlock outputs locked

			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					if in.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(in.RefTxID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}

		if len(block.PreHeader) == 0 {
			break
		}
	}

	return unspentTXs
}

func (bc *Blockchain) FindUTXO(address string) []TXOutput {

	var TXOuts []TXOutput
	unspentTxns := bc.FindUnspentTransactions(address) //get list of unspent txns

	for _, txn := range unspentTxns {
		for _, out := range txn.Vout { //could be many outputs in a Vout
			if out.CanBeUnlockedWith(address) {
				TXOuts = append(TXOuts, out)
			}
		}
	}

	return TXOuts

}

