package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

const subsidy = 12

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

type TXInput struct {
	RefTxID   []byte
	Vout      int
	ScriptSig string
}

type TXOutput struct {
	Value        int
	ScriptPubKey string
}

// Set txn ID based on its hash digest
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var digest [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	digest = sha256.Sum256(encoded.Bytes())
	tx.ID = digest[:]

}

func NewCoinbaseTX(to, script string) *Transaction {
	if script == "" {
		script = fmt.Sprint("Minting %d coins by: %s", subsidy, to)
	}

	txin := TXInput{[]byte{}, -1, script}
	txout := TXOutput{subsidy, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()

	return &tx
}

func (txn *Transaction) IsCoinbase() bool {
	return len(txn.Vin) == 1 &&
		len(txn.Vin[0].RefTxID) == 0 &&
		txn.Vin[0].Vout == -1
}

func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}
