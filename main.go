package main

import (
	"os"

	"github.com/blockchain-app-go/wallet"
)

func main() {
	defer os.Exit(0)

	// cli := cli.CommandLine{}
	// cli.Run()

	wallet := wallet.MakeWallet()
	wallet.Address()
}
