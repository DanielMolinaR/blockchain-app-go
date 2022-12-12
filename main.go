package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

type Blockchain struct {
	blocks []*Block
}

type Block struct {
	Hash     []byte // hash of the block
	Data     []byte // Data stored in the block
	PrevHash []byte // Previous block's hash in the chain
}

// Create the hash based on the previous hash and the data from the block
func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{}) // Takes 2 dimensional slice of bytes and combine them with an empty slice of bytes
	hash := sha256.Sum256(info)                                // Using SHA256 as placeholder for now
	b.Hash = hash[:]
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash} // create a block with an empty hash
	block.DeriveHash()                                // create the hash of the Block
	return block
}

func (chain *Blockchain) AddBlock(data string) {
	prevBlock := chain.blocks[len(chain.blocks)-1]
	newBlock := CreateBlock(data, prevBlock.Hash)
	chain.blocks = append(chain.blocks, newBlock)
}

// The Genesis block is the first block of a Blockchain
func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

func InitBlockchain() *Blockchain {
	return &Blockchain{[]*Block{Genesis()}}
}

func (chain *Blockchain) DisplayData() {
	for _, block := range chain.blocks {
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data in block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
	}
}

func main() {
	chain := InitBlockchain()

	chain.AddBlock("First block")
	chain.AddBlock("Second block")
	chain.AddBlock("Third block")
	chain.AddBlock("Forth block")

	chain.DisplayData()
}
