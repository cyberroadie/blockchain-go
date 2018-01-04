package blockchain

import (
	"net/url"
	"os"
	"log"
	"net/http"
	"encoding/json"
	"encoding/hex"
	"bytes"
	"errors"
	"fmt"
	"reflect"
)

type server struct {
	difficulty   uint8
	regURL       url.URL
	nodeRegistry *NodeRegistry
	blockchain   *BlockChain
	serverNode   *Node
	serverPort   int
}

var logger = log.New(os.Stdout, "blockchain", log.Lshortfile)

func NewBlockChainServer(difficulty uint, regURL url.URL, serverPort int) *server {
	bc := newBlockChain(uint8(difficulty))
	nr := NewNodeRegistry()
	n := newNode()
	return &server{difficulty:uint8(difficulty), regURL:regURL, nodeRegistry:nr, blockchain:bc, serverNode:n, serverPort:serverPort}
}

func (bcs *server) StartServer() {

	logger.Printf("starting blockchain with difficulty %d, listening on port %d", bcs.difficulty, bcs.serverPort)

	http.HandleFunc("/transaction", bcs.AddTransaction)
	http.HandleFunc("/block", bcs.MineBlock)
	http.HandleFunc("/chain", bcs.getBlockChain)
	http.HandleFunc("/node/register", bcs.RegisterNode)
	http.HandleFunc("/node/resolve", resolveConflict)
	http.HandleFunc("/node/registry", bcs.GetAllNodes)
	logger.Printf("initialized web server")

	proceed := make(chan bool, 1)
	go func() {
		<- proceed
		err := bcs.registerSelf()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	done := make(chan bool)
	go func() {
		serverAddr := fmt.Sprintf("0.0.0.0:%d", bcs.serverPort)
		err := http.ListenAndServe(serverAddr, nil)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		done <- true
	}()

	v := reflect.ValueOf(http.DefaultServeMux).Elem()
	logger.Printf("started web server with the following routes: %v\n", v.FieldByName("m"))
	proceed <- true
	<-done

}

func (bcs *server) registerSelf() error {
	buf := bytes.NewBufferString("")
	encoder := json.NewEncoder(buf)
	encoder.Encode(bcs.serverNode)

	url := fmt.Sprintf("%s/node/register", bcs.regURL.String())
	res, err := http.Post( url, "application/json", buf)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New(fmt.Sprintf("error registering server %s %s", bcs.regURL.String(), res.Status))
	}

	return nil
}

func (bcs *server) GetAllNodes(writer http.ResponseWriter, request *http.Request) {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "   ")
	encoder.Encode(bcs.nodeRegistry)
}

func resolveConflict(writer http.ResponseWriter, request *http.Request) {

}

func (bcs *server) RegisterNode(writer http.ResponseWriter, request *http.Request) {
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

func (bcs *server) MineBlock(writer http.ResponseWriter, request *http.Request) {
	proof := bcs.blockchain.ProofOfWork()
	t := newTransAction("0", bcs.serverNode.NodeId.String(), 1)
	bcs.blockchain.AddTransaction(t)
	trans := make([]Transaction,1)
	trans[0] = *t

	b := newBlock(trans, proof, hex.EncodeToString(bcs.blockchain.LastBlock().Hash()))
	bcs.blockchain.AddBlock(*b)

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "    ")
	encoder.Encode(b)
}

func (bcs *server) AddTransaction(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var t Transaction
	err := decoder.Decode(&t)
	if err != nil {
		http.Error(writer, "Invalid json", http.StatusBadRequest)
		return
	}

	t = *bcs.blockchain.AddTransaction(&t)
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "    ")
	encoder.Encode(t)
}
