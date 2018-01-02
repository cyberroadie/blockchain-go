package blockchain

import (
	"net/url"
	"os"
	"log"
	"net/http"
	"encoding/json"
	"encoding/hex"
	"fmt"
)

type server struct {
	difficulty uint8
	RegURL url.URL
	nodeRegistry *NodeRegistry
	blockchain *BlockChain
	serverNode *Node
}

func NewBlockChainServer(difficulty uint, regURL url.URL) *server {
	bc := newBlockChain(uint8(difficulty))
	nr := NewNodeRegistry()
	n := newNode()
	return &server{difficulty:uint8(difficulty), RegURL:regURL, nodeRegistry:nr, blockchain:bc, serverNode:n}
}

func (bcs *server) StartServer() {
	logger := log.New(os.Stdout, "blockchain", log.Lshortfile)
	log.SetOutput(os.Stdout)

	logger.Printf("Starting blockchain with difficulty %d\n", bcs.difficulty)

	http.HandleFunc("/transaction", bcs.addTransaction)
	http.HandleFunc("/block", bcs.mineBlock)
	http.HandleFunc("/chain", bcs.getBlockChain)
	http.HandleFunc("/node/register", bcs.registerNode)
	http.HandleFunc("/node/resolve", resolveConflict)
	http.HandleFunc("/node/registry", bcs.getAllNodes)
	http.ListenAndServe(":8080", nil)

}

func (bcs *server) getAllNodes(writer http.ResponseWriter, request *http.Request) {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "   ")
	encoder.Encode(bcs.nodeRegistry)
}

func resolveConflict(writer http.ResponseWriter, request *http.Request) {

}

func (bcs *server) registerNode(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var n Node
	err := decoder.Decode(&n)
	if err != nil {
		http.Error(writer, "Invalid json", http.StatusBadRequest)
	}
	bcs.nodeRegistry.RegisterNode(n)
}


func (bcs *server) getBlockChain(writer http.ResponseWriter, request *http.Request) {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "    ")
	encoder.Encode(bcs.blockchain)
}

func (bcs *server) mineBlock(writer http.ResponseWriter, request *http.Request) {
	proof := bcs.blockchain.ProofOfWork()
	t := newTransAction("0", bcs.serverNode.nodeId.String(), 1)
	bcs.blockchain.AddTransaction(t)
	trans := make([]Transaction,1)
	trans[0] = *t

	b := newBlock(trans, proof, hex.EncodeToString(bcs.blockchain.LastBlock().Hash()))
	bcs.blockchain.AddBlock(*b)

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "    ")
	encoder.Encode(b)
}

func (bcs *server) addTransaction(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var t Transaction
	err := decoder.Decode(&t)
	if err != nil {
		http.Error(writer, "Invalid json", http.StatusBadRequest)
	}
	bcs.blockchain.AddTransaction(&t)
	fmt.Fprint(writer, "Transaction added")
}
