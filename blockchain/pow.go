package blockchain

import (
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
	"strconv"
)

var maxNonce = math.MaxInt64


const targetBits = 12


type ProofOfWork struct {
	block *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork{

	target:=big.NewInt(1)
	target=target.Lsh(target,uint(256-targetBits))

	return &ProofOfWork{b, target}
	}

func (pow *ProofOfWork) Mine() ( int, []byte){

	//var header [32]byte
	var (checkHeader big.Int
	 header [32]byte
	 nonce  = 0
	 data []byte
	)

	preData:=bytes.Join([][]byte{
		pow.block.PreHeader,
		[]byte(strconv.FormatInt(int64(pow.block.TimeStamp), 16)),
		pow.block.TransactionsDigest(),
		[]byte(strconv.FormatInt(int64(targetBits), 16))},
		[]byte{})

	for ; nonce < maxNonce; nonce++ {
		data=bytes.Join([][]byte{preData,
			[]byte(strconv.FormatInt(int64(nonce), 16))},
			[]byte{})
		header = sha256.Sum256(data)

		//fmt.Printf("nonce: %d, Hash: %x\n", nonce, header)
		checkHeader.SetBytes(header[:])

		if checkHeader.Cmp(pow.target) == -1 {
			break
		}
		}


	return nonce, header[:]

}

func (pow *ProofOfWork) Validate() bool {

	var header [32]byte
	var checkHeader big.Int

	data:=bytes.Join([][]byte{
		pow.block.PreHeader,
		[]byte(strconv.FormatInt(int64(pow.block.TimeStamp), 16)),
		pow.block.TransactionsDigest(),
		[]byte(strconv.FormatInt(int64(targetBits), 16))},
		[]byte{})

	data=bytes.Join([][]byte{data,
		[]byte(strconv.FormatInt(int64(pow.block.Nonce), 16))},
		[]byte{})

		header = sha256.Sum256(data)
		checkHeader.SetBytes(header[:])

		return checkHeader.Cmp(pow.target) == -1
}