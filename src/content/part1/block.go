package main

import (
	"time"
	"strconv"
	"bytes"
	"crypto/sha256"
)

// Block 由区块头和交易两部分构成
// Timestamp, PrevBlockHash, Hash 属于区块头（block header）
// Timestamp     : 当前时间戳，也就是区块创建的时间
// PrevBlockHash : 前一个块的哈希
// Hash          : 当前块的哈希
// Data          : 区块实际存储的信息，比特币中也就是交易
type Block struct {
	Timestamp int64
	PrevBlockHash []byte
	Hash []byte
	Data []byte
}

//NewBlock 生成一个新的区块,参数需要prevBlockHash,data
//当前的Hash会根据prevBlockHash和data生成得到
func NewBlock(prevBlockHash []byte, data string) *Block {
	block := &Block{time.Now().Unix(), prevBlockHash, []byte{}, []byte(data)}
	block.setHash()
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