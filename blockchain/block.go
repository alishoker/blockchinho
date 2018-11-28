package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"time"
)

type Block struct {
	Header       []byte
	PreHeader    []byte
	TimeStamp    int64
	Transactions []byte
	Nonce		int
}

func NewBlock(preHeader []byte, transactions string) *Block {
	block := &Block{Header: []byte{},
		PreHeader:    preHeader,
		TimeStamp:    time.Now().Unix(),
		Transactions: []byte(transactions),
		Nonce:0}
	block.setHeader()
	fmt.Printf("NewBlock: Header: %x\n", block)
	return block
}

func (block *Block) setHeader() {
	pow := NewProofOfWork(block)
	nonce, header:=pow.Mine()

	block.Nonce=nonce
	block.Header=header[:]
	//fmt.Printf("setHeader: Header: %x\n", block)
}

/*
func (block *Block) setHeader() {
	hash := sha256.Sum256(bytes.Join([][]byte{block.PreHeader,
		[]byte(strconv.FormatInt(block.TimeStamp, 10)),
		block.Transactions}, []byte{}))
	block.Header = hash[:]
}
*/

//Serialization gob

func (b *Block) Serialize() []byte {
	var out bytes.Buffer
	marshal := gob.NewEncoder(&out)

	if err:=marshal.Encode(b); err!=nil {
		log.Panic(err)
	}

	return out.Bytes()
}

func DeserializeBlock(in []byte) *Block{

	var block Block

	marshal := gob.NewDecoder( bytes.NewReader(in))

	if err:=marshal.Decode(&block); err!=nil {
		log.Panic(err)
	}

	return &block
}