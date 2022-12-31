package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const walletFile = "./tmp/wallets.data"

type Wallets struct {
	Wallets map[string]*Wallet
}

func CreateWallets() (*Wallets, error) {
	wallets := Wallets{}

	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFile()

	return &wallets, err
}

func (wallets *Wallets) AddWallet() string {
	wallet := MakeWallet()
	address := fmt.Sprintf("%s", wallet.Address())

	wallets.Wallets[address] = wallet

	return address
}

func (wallets *Wallets) GetAllAddresses() []string {
	var addresses []string

	for address := range wallets.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

func (wallets Wallets) GetWallet(address string) Wallet {
	return *wallets.Wallets[address]
}

func (wallets *Wallets) LoadFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	var ws Wallets

	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&ws)
	if err != nil {
		log.Panic(err)
	}

	wallets.Wallets = ws.Wallets

	return nil
}

func (wallets *Wallets) SaveFile() {
	var content bytes.Buffer

	gob.Register(elliptic.P256()) // It is used to encode the buffer into the file

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(wallets)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644) // 0644 gives read and write perms
	if err != nil {
		log.Panic(err)
	}

}
