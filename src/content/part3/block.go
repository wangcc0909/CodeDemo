package main

import (
	"time"
	"strconv"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

// Nonce 在对工作量证明进行验证时用到
type Block struct {
	Timestamp int64
	PrevBlockHash []byte
	Hash []byte
	Data []byte
	Nonce int
}

//将block序列化一个字节数组
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		fmt.Println(err)
	}
	return result.Bytes()
}

func DeSerializeBlock(b []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewBuffer(b))
	err := decoder.Decode(&block)
	if err != nil {
		fmt.Println(err)
	}
	return &block
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