package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

const dbPath = "./tmp/blocks"

type Blockchain struct {
	LastHash []byte
	Databse  *badger.DB
}

type BlockchainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func InitBlockchain() *Blockchain {
	var lastHash []byte

	opts := badger.DefaultOptions("")
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

func (chain *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{chain.LastHash, chain.Databse}
}

func (iter *BlockchainIterator) Next() *Block {
	var block *Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		HandleError(err)
		var encodedBlock []byte
		err = item.Value(func(val []byte) error {
			encodedBlock = val
			return nil
		})

		block = Deserialize(encodedBlock)

		return err
	})
	HandleError(err)

	iter.CurrentHash = block.PrevHash

	return block
}
