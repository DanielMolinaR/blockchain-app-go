package main

import (
	"os"

	"github.com/blockchain-app-go/cli"
)

func main() {
	defer os.Exit(0)

	cli := cli.CommandLine{}
	cli.Run()
}
