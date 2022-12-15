package blockchain

import (
	"fmt"
	"strconv"

	"github.com/dgraph-io/badger"
)

const dbPath = "./tmp/blocks"

type Blockchain struct {
	LastHash []byte
	Databse  *badger.DB
}

func InitBlockchain() *Blockchain {
	var lastHash []byte

	opts := badger.DefaultOptions()
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	HandleError(err)

	// Update lets make Read and Write transactions into the database
	// We are sending an enclosure which takes in a pointer to a badger transaction
	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")
			genesis := genesis()
			fmt.Println("Genesis proved")
			err = txn.Set(genesis.Hash, genesis.Serialize())
			HandleError(err)
			err = txn.Set([]byte("lh"), genesis.Hash)

			lastHash = genesis.Hash

		} else {
			item, err := txn.Get([]byte("lh"))
			HandleError(err)
			err = item.Value(func(val []byte) error {
				lastHash = val
				return nil
			})
		}

		return err
	})

	HandleError(err)

	return &Blockchain{lastHash, db}
}

func (chain *Blockchain) AddBlock(data string) {
	var lastHash []byte

	err := chain.Databse.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		HandleError(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		return err
	})
	HandleError(err)

	newBlock := createBlock(data, lastHash)

	err = chain.Databse.Update(func(txn *badger.Txn) error {
		err = txn.Set(newBlock.Hash, newBlock.Serialize())
		HandleError(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})

	HandleError(err)
}

func (chain *Blockchain) DisplayData() {
	for _, block := range chain.Blocks {
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data in block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
