package blockchain

type Blockchain struct {
	Blocks []*Block
}

func (bc *Blockchain) AddBlock(trans string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
    b := NewBlock(prevBlock.Header,trans)
	bc.Blocks = append(bc.Blocks, b)
}

func NewBlockchain() *Blockchain {
	genesis := NewBlock([]byte{}, "The Genesis Block")
	return &Blockchain{[]*Block{genesis}}
}
