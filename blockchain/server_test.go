package blockchain

import (
	"testing"
	"net/url"
	"net/http/httptest"
	"net/http"
	"strings"
	"log"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"github.com/satori/go.uuid"
)

var bcs = NewBlockChainServer(2, url.URL{Scheme: "http", Host: "localhost:8080"})

func TestServer_AddTransaction_Invalid_JSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(bcs.AddTransaction))
	defer ts.Close()

	res, err := http.Post(ts.URL, "application/json", strings.NewReader("not json"))

	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != 400 {
		log.Fatal("expected error with non json input")
	}

	msg, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", msg)

}

func TestServer_AddTransaction(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(bcs.AddTransaction))
	defer ts.Close()

	tr := Transaction{Sender:"me", Recipient:"you", Amount: 1}
	buf := bytes.NewBufferString("")
	encoder := json.NewEncoder(buf)
	encoder.Encode(tr)

	res, err := http.Post(ts.URL, "application/json", strings.NewReader(buf.String()))

	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != 200 {
		log.Fatal("unexpected HTTP code $d", res.StatusCode)
	}

	msg, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", msg)

}

func TestServer_MineBlock(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(bcs.MineBlock))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		log.Fatalf("unexpected HTTP status code: %d", res.StatusCode)
	}

	var block Block
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&block)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("block: %v", block)

}

func TestServer_RegisterNode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(bcs.RegisterNode))
	defer ts.Close()

	n := Node{uuid.NewV1(), url.URL{Scheme:"http", Host:"localhost:5000"}}
	buf := bytes.NewBufferString("")
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(n)

	res, err := http.Post(ts.URL, "application/json", strings.NewReader(buf.String()))
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != 200 {
		log.Fatal("unexpected HTTP code $d", res.StatusCode)
	}
}
