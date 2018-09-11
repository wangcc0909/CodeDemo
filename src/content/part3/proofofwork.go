package main

import (
	"math/big"
	"fmt"
	"math"
	"bytes"
	"crypto/sha256"
	"github.com/boltdb/bolt"
	"log"
)

//难度值  这里表示哈希的前24位必须是0
const targetBits = 24

const MaxNonce = math.MaxInt64

//每个块的工作量都必须要有证明,所以有个指向Block的指针
//target是目标 我们最终要找的哈希必须小于目标
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

//target = 1 左移 256 - targetBits 位
func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	//左移多少位
	target.Lsh(target, uint(256-targetBits))
	pow := &ProofOfWork{block: block, target: target}
	return pow
}

//工作量证明用到的数据有: prevBlockHash,data,timestamp,targetBits,nonce
func (pow *ProofOfWork) PrepareData(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.block.PrevBlockHash,
		pow.block.Data,
		IntToHex(pow.block.Timestamp),
		IntToHex(int64(targetBits)),
		IntToHex(int64(nonce)),
	}, []byte{})
	return data
}

//工作量证明的核心就是寻找有效的hash
func (pow *ProofOfWork) Run() (int,[]byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\" \n", pow.block.Data)
	for nonce < MaxNonce {
		data := pow.PrepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:]) //将hash值视为一个大端在前的无符号整数
		if hashInt.Cmp(pow.target) == -1 {  //两个大整数比较  相等返回-1
			fmt.Printf("\r%x",hash)
			break
		}else {
			nonce++
		}
	}
	fmt.Print("\n\n")
	return nonce,hash[:]
}

//验证工作量 只有hash小于目标就是有效工作量
func (pow *ProofOfWork) validate() bool {
	var hashInt big.Int
	data := pow.PrepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	valid := hashInt.Cmp(pow.target) == -1
	return valid
}

type BlockChainIterator struct {
	currentHash []byte
	db *bolt.DB
}

func (bc *BlockChain) Iterator() *BlockChainIterator {
	bci := &BlockChainIterator{currentHash:bc.tip,db:bc.db}
	return bci
}

//返回链中的下一个块
func (i *BlockChainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		encoder := bucket.Get(i.currentHash)
		block = DeSerialezeBlock(encoder)
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.PrevBlockHash
	return block
}

