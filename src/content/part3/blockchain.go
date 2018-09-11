package main

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
)

const dbFile = "blockChain.db"
const blockBucket = "blocks"

// tip 这个词本身有事物尖端或尾部的意思，这里指的是存储最后一个块的哈希
// 在链的末端可能出现短暂分叉的情况，所以选择 tip 其实也就是选择了哪条链
// db 存储数据库连接
type BlockChain struct {
	tip []byte
	db *bolt.DB
}

// 加入区块时，需要将区块持久化到数据库中
func (bc *BlockChain) AddBlock(data string) {
	var lastHash []byte
	// 首先获取最后一个块的哈希用于生成新块的哈希
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(lastHash,data)
	bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		err = b.Put(newBlock.Hash,newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = b.Put([]byte("l"),newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}
		bc.tip = newBlock.Hash
		return nil
	})
}

//创建一个有创世纪块的区块链
func NewBlockChain() *BlockChain {
	var tip []byte
	//打开一个数据库文件
	db,err := bolt.Open(dbFile,0600,nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		// 如果数据库中不存在区块链就创建一个，否则直接读取最后一个块的哈希
		if bucket == nil {
			fmt.Println("not existing blockchain found,creating a new one ...")
			genesis := NewGenesisBlock()
			b,err := tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				log.Panic(err)
			}

			err = b.Put(genesis.Hash,genesis.Serialize())
			if err != nil {
				log.Panic(err)
			}
			err = b.Put([]byte("l"),genesis.Hash)
			if err != nil {
				log.Panic(err)
			}
			tip = genesis.Hash
		}else {
			tip = bucket.Get([]byte("l"))
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
	bc := &BlockChain{tip:tip,db:db}

	return bc
}
