package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

// take the data from the block

// create a counter (nonce) which starts at 0

// create a hash of the data plus the counter

// check the hash to see if it meets a set of requirements

// Requierements:
// The First few bytes must contain 0s

const Difficulty = 18

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProof(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))

	pow := &ProofOfWork{block, target}

	return pow
}

func (pow *ProofOfWork) InitData(nonce int) []byte {
	// Takes 2 dimensional slice of bytes and combine them with an empty slice of bytes
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.HashTransaction(),
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)

	return data
}

// Create the hash based on the previous hash, the data and nonce from the block and the difficulty
func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0

	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data) // Using SHA256 as placeholder for now

		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 { // if we have already exceeded the target breaks
			break
		} else {
			nonce++
		}

	}
	fmt.Println()

	return nonce, hash[:]
}

// We´ll use the nonce retrieved from Run() to derive the hash which met the target we wanted
// and we´ll run the cycle one more time to show that the hash is valid or not
func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1

}

func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
