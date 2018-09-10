package main

import (
	"time"
	"strconv"
	"bytes"
	"crypto/sha256"
)

// Nonce 在对工作量证明进行验证时用到
type Block struct {
	Timestamp int64
	PrevBlockHash []byte
	Hash []byte
	Data []byte
	Nonce int
}

// 创建新块时需要运行工作量证明找到有效哈希
func NewBlock(prevBlockHash []byte, data string) *Block {
	block := &Block{time.Now().Unix(), prevBlockHash, []byte{}, []byte(data),0}
	pow := NewProofOfWork(block)
	nonce,hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

//setHash 设置当前块的hash
//Hash = sha256(prevBlockHash + data + timestamp)
func (b *Block) setHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp,10))
	headers := bytes.Join([][]byte{b.PrevBlockHash,b.Data,timestamp},[]byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

//NewGenesisBlock 生成创世块
func NewGenesisBlock() *Block {
	return NewBlock([]byte{},"Genesis Block")
}