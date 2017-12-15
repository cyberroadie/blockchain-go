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
)

var bc *blockchain.BlockChain
var nodeId uuid.UUID

func main() {
	logger := log.New(os.Stdout, "blockchain", log.Lshortfile)
	log.SetOutput(os.Stdout)

	logger.Print("Starting blockchain")
	bc = blockchain.NewBlockChain(16)
	nodeId = uuid.NewV1()

	http.HandleFunc("/transaction", addTransaction)
	http.HandleFunc("/block", mineBlock)
	http.HandleFunc("/chain", getBlockChain)
	http.ListenAndServe(":8080", nil)
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