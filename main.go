package main

import (
	"blockchain/blockchain"
	"log"
	"os"
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"encoding/hex"
	"flag"
)

var bc *blockchain.BlockChain
var nr *blockchain.NodeRegistry
var nodeId uuid.UUID
var nodes []blockchain.Node

func main() {
	logger := log.New(os.Stdout, "blockchain", log.Lshortfile)
	log.SetOutput(os.Stdout)

	var difficulty = flag.Uint("difficulty", 16, "set difficulty block chain (max 256)")
	if *difficulty > 256 || *difficulty < 1 {
		logger.Fatalf("Difficulty of blockchain has to be greater than 0 and less than 256, " +
			                 "given input is %d", difficulty)
	}

	nodeId = uuid.NewV1()

	logger.Printf("Starting blockchain with difficulty %d and node ID %s", difficulty, nodeId.String())
	nodes = make([]blockchain.Node, 10)
	bc = blockchain.NewBlockChain(uint8(*difficulty))
	nr = blockchain.NewNodeRegistry()

	http.HandleFunc("/transaction", addTransaction)
	http.HandleFunc("/block", mineBlock)
	http.HandleFunc("/chain", getBlockChain)
	http.HandleFunc("/node/register", registerNode)
	http.HandleFunc("/node/resolve", resolveConflict)
	http.HandleFunc("/node/registry", getAllNodes)
	http.ListenAndServe(":8080", nil)
}

func getAllNodes(writer http.ResponseWriter, request *http.Request) {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "   ")
	encoder.Encode(nr)
}

func resolveConflict(writer http.ResponseWriter, request *http.Request) {

}

func registerNode(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var n blockchain.Node
	err := decoder.Decode(&n)
	if err != nil {
		http.Error(writer, "Invalid json", http.StatusBadRequest)
	}
	nr.RegisterNode(n)
}


func getBlockChain(writer http.ResponseWriter, request *http.Request) {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "    ")
	encoder.Encode(bc)
}

func mineBlock(writer http.ResponseWriter, request *http.Request) {
	proof := bc.ProofOfWork()
	t := blockchain.NewTransAction("0", nodeId.String(), 1)
	bc.AddTransaction(t)
	trans := make([]blockchain.Transaction,1)
	trans[0] = *t

	b := blockchain.NewBlock(trans, proof, hex.EncodeToString(bc.LastBlock().Hash()))
	bc.AddBlock(*b)

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "    ")
	encoder.Encode(b)
}

func addTransaction(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var t blockchain.Transaction
	err := decoder.Decode(&t)
	if err != nil {
		http.Error(writer, "Invalid json", http.StatusBadRequest)
	}
	bc.AddTransaction(&t)
	fmt.Fprint(writer, "Transaction added")
}