package main

//BlockChain 是一个指针数组
type BlockChain struct {
	blocks []*Block
}

//AddBlock 向链中添加一个块
//data 在实际中就是交易
func (bc *BlockChain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks) - 1]
	newBlock := NewBlock(prevBlock.Hash,data)
	bc.blocks = append(bc.blocks, newBlock)
}

//创建一个有创世纪块的区块链
func NewBlockChain() *BlockChain {
	return &BlockChain{[]*Block{NewGenesisBlock()}}
}
