package main

import (
	"os"
	"fmt"
	"flag"
	"log"
)

type CLI struct {
	bc *BlockChain
}

const usage = `
usage:
addblock -data BLOCK_DATA  add a block to the blockchain
printchain     print all the blocks of the blockchain
`

func (cli *CLI) PrintUsage() {
	fmt.Println(usage)
}

func (cli *CLI) ValidateArgs() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
		os.Exit(1)
	}
}

func (cli *CLI) Run() {
	cli.ValidateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block Data")
	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.PrintUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			cli.PrintUsage()
			os.Exit(1)
		}
		cli.AddBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}
