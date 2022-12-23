package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

type Block struct {
	Hash         []byte         // hash of the block
	Transactions []*Transaction // Data stored in the block
	PrevHash     []byte         // Previous block's hash in the chain
	Nonce        int
}

// Method to provide an unique representation to all the transactions from the block combined
func (block *Block) HashTransaction() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range block.Transactions {
		txHashes = append(txHashes, tx.ID)
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

func createBlock(txs []*Transaction, prevHash []byte) *Block {
	block := &Block{[]byte{}, txs, prevHash, 0} // create a block with an empty hash
	pow := NewProof(block)

	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// The Genesis block is the first block of a Blockchain
func genesis(coinbase *Transaction) *Block {
	return createBlock([]*Transaction{coinbase}, []byte{})
}

func (block *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(block)

	HandleError(err)

	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	HandleError(err)
	return &block
}

func HandleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
