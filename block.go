package blockchain

import (
	"bytes"
	"crypto.sha256"
	"strconv"
	"time"
)

type Block struct {
	Header       []byte
	PreHeader    []byte
	TimeStamp    int64
	Transactions []byte
}

func newBlock(preHeader []byte, transactions string) *Block {
	block := &Block{Header: []byte{},
		PreHeader:    preHeader,
		TimeStamp:    time.Now().Unix(),
		Transactions: []byte(transactions)}
	block.setHeader()
	return block
}

func (block *Block) setHeader() {
	block.Header = sha256.Sum256(bytes.Join(block.PreHeader,
		[]byte(strconv.FormatInt(block.TimeStamp, 10)),
		block.Transactions))
}
