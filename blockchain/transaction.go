package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

type TxOutput struct {
	Value  int
	PubKey string
}

type TxInput struct {
	ID  []byte
	Out int
	Sig string
}

func (tx *Transaction) SetId() {
	var encoded bytes.Buffer
	var hash [32]byte // Hash based on the bytes that represent our tx

	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	HandleError(err)

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]

}

func CoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}

	txIn := TxInput{[]byte{}, -1, data} // Since is not referecing to any Output the ID is empty and the OUT int -1
	txOut := TxOutput{100, to}

	tx := Transaction{nil, []TxInput{txIn}, []TxOutput{txOut}}
	tx.SetId()

	return &tx
}

func NewTransaction(from, to string, amount int, chain *Blockchain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	accumulated, validOutputs := chain.FindSPendableOutputs(from, amount)

	if accumulated < amount {
		log.Panic("Error: Not enough funds")
	}

	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		HandleError(err)

		for _, out := range outs {
			input := TxInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{amount, to})

	if accumulated > amount {
		outputs = append(outputs, TxOutput{accumulated - amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetId()

	return &tx
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == 0
}

func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data
}

func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}
