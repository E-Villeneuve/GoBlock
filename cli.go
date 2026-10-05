package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type CLI struct {
	bc *Blockchain
}

func (cli *CLI) validate() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  addblock -data BLOCK_DATA - add a block to the blockchain")
		fmt.Println("  printchain - print all the blocks of the blockchain")
		os.Exit(1)
	}
}

func (cli *CLI) printChain() {
	bci := cli.bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)

		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))

		fmt.Println()

		if len(block.PrevBlockHash) == 0 { // Genesis block has no previous hash, so loop needs to break
			break
		}

	}
}

func (cli *CLI) Run() {
	cli.validate()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	addBlockData := addBlockCmd.String("data", "", "Block data")

	switch os.Args[1] {
	case "addblock":
		addBlockCmd.Parse(os.Args[2:])
	case "printchain":
		printChainCmd.Parse(os.Args[2:])
	default:
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.bc.AddBlock(*addBlockData)

	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}
