package main

import (
	"fmt"
	"strconv"
)

func (cli *CLI) AddBlock(data string) {
	cli.bc.AddBlock(data)
	fmt.Println("SUCCESS")
}

func (cli *CLI) printChain() {
	bci := cli.bc.Iterator()
	for {
		block := bci.Next()
		fmt.Printf("Prev Hash:%x\n",block.PrevBlockHash)
		fmt.Printf("Data:%s\n",block.Data)
		fmt.Printf("Hash:%x\n",block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW:%s\n",strconv.FormatBool(pow.validate()))
		fmt.Println()
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
