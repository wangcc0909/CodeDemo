package main

import "fmt"

func main() {

	blockChain := NewBlockChain()

	blockChain.AddBlock("Send 1 BTC to Ivan")
	blockChain.AddBlock("Send 2 more BTC to Ivan")

	for _,bc := range blockChain.blocks {
		fmt.Printf("bc.prevHash: %x\n",bc.PrevBlockHash)
		fmt.Printf("bc.Hash: %x\n",bc.Hash)
		fmt.Printf("bc.Data: %s\n",bc.Data)
		fmt.Println()
	}
}
