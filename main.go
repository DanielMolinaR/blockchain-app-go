package main

import (
	"github.com/blockchain-app-go/blockchain"
)

func main() {
	chain := blockchain.InitBlockchain()

	chain.AddBlock("First block after Genesis")
	chain.AddBlock("Second block after Genesis")
	chain.AddBlock("Third block after Genesis")
	chain.AddBlock("Forth block after Genesis")

	chain.DisplayData()
}
