package main

import (
	"blockchain/blockchain"
	"os"
	"fmt"
	"flag"
	"net/url"
)

var bc *blockchain.BlockChain
var nr *blockchain.NodeRegistry
var nodes []blockchain.Node

func main() {

	var difficulty = flag.Uint("dif", 16, "set difficulty block chain (max 256)")
	var regDomain = flag.String("reg", "localhost:8080", "domain where to get the node registry from")

	flag.Parse()

	if *difficulty > 256 || *difficulty < 1 {
		fmt.Errorf("difficulty of blockchain has to be greater than 0 and less than 256, " +
			                 "given input is %d", difficulty)
		os.Exit(1)
	}

	u, err := url.Parse(fmt.Sprintf("http://%s", regDomain))
	if err != nil {
		fmt.Errorf("invalid registry url: %s", fmt.Sprintf("http://%s", regDomain))
		os.Exit(1)
	}

	bcs := blockchain.NewBlockChainServer(*difficulty, *u)
	bcs.StartServer()

}

