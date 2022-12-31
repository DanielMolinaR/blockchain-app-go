package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/dgraph-io/badger"
)

const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"           // This file is used to check if the blockchain database already exists
	genesisData = "First Transasction from Genesis" // Arbitray data
)

type Blockchain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockchainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func DbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func OpenDatabase() *badger.DB {
	opts := badger.DefaultOptions("")
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	HandleError(err)
	return db
}

func InitBlockchain(adress string) *Blockchain {
	if DbExists() {
		fmt.Println("Blockchain already exists")
		runtime.Goexit()
	}

	var lastHash []byte

	db := OpenDatabase()

	// Update lets make Read and Write transactions into the database
	// We are sending an enclosure which takes in a pointer to a badger transaction
	err := db.Update(func(txn *badger.Txn) error {
		cbtx := CoinbaseTx(adress, genesisData)
		genesis := genesis(cbtx)
		fmt.Println("Genesis created")
		err := txn.Set(genesis.Hash, genesis.Serialize())
		HandleError(err)
		err = txn.Set([]byte("lh"), genesis.Hash)

		lastHash = genesis.Hash

		return err
	})

	HandleError(err)

	return &Blockchain{lastHash, db}
}

func ContinueBlockchain(address string) *Blockchain {
	if DbExists() == false {
		fmt.Println("Blockchain already exists")
		runtime.Goexit()
	}

	var lastHash []byte

	db := OpenDatabase()

	err := db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		HandleError(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		return err
	})

	HandleError(err)

	return &Blockchain{lastHash, db}

}

func (chain *Blockchain) AddBlock(transactions []*Transaction) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		HandleError(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		return err
	})
	HandleError(err)

	newBlock := createBlock(transactions, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err = txn.Set(newBlock.Hash, newBlock.Serialize())
		HandleError(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})

	HandleError(err)
}

func (chain *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{chain.LastHash, chain.Database}
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

/*
	This method finds the unused transactions assign to an user
 	unused transactions are transactions taht have output which are not referenced
	other inputs. If an output hasn't been used means that those transactions
	still exits for a certain user so by counting all the used transaction that are
	assigned to a certain user we can find how many tokens are assigned to that user.const

	FindUnspentTransactions and FindUnspentTxO are made to find the balance of the user account
*/
func (chain *Blockchain) FindUnspentTransactions(pubKeyHash []byte) []Transaction {
	var unspentTxs []Transaction

	spentTxos := make(map[string][]int) // map where the keys are string and the values are a slice of ints

	iter := chain.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs: // label used to break out the for loop of Outputs and not from the other two for loops
			for outIdx, out := range tx.Outputs {
				if spentTxos[txID] != nil {
					for _, spentOut := range spentTxos[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.IsLockedWithKey(pubKeyHash) {
					unspentTxs = append(unspentTxs, *tx)
				}

			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Inputs {
					if in.UsesKey(pubKeyHash) { // If we can unlock it means that the transaction is related to the user
						inTxID := hex.EncodeToString(in.ID)
						spentTxos[inTxID] = append(spentTxos[inTxID], in.Out)
					}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return unspentTxs
}

func (chain *Blockchain) FindUnspentTxO(pubKeyHash []byte) []TxOutput {
	var UnspentTxOs []TxOutput
	unspentTransactions := chain.FindUnspentTransactions(pubKeyHash)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.IsLockedWithKey(pubKeyHash) {
				UnspentTxOs = append(UnspentTxOs, out)
			}
		}
	}

	return UnspentTxOs
}

/*
	Enables the creation of normal transactions which are not coinbase
*/
func (chain *Blockchain) FindSPendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := chain.FindUnspentTransactions(pubKeyHash)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
			if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOuts
}

func (chain *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
	iter := chain.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("Transaction does not exit")
}

func (chain *Blockchain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTxs := chain.GetPreviousTransactions(tx)

	tx.Sign(privKey, prevTxs)
}

func (chain *Blockchain) VerifyTransaction(tx *Transaction) bool {
	prevTxs := chain.GetPreviousTransactions(tx)

	return tx.Verify(prevTxs)
}

func (chain *Blockchain) GetPreviousTransactions(tx *Transaction) map[string]Transaction {
	prevTxs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		prevTx, err := chain.FindTransaction(in.ID)
		HandleError(err)
		prevTxs[hex.EncodeToString(prevTx.ID)] = prevTx
	}

	return prevTxs
}
