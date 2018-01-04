package blockchain

import (
	"time"
	"encoding/json"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"github.com/satori/go.uuid"
	"encoding/hex"
)

type Transaction struct {
	Id uuid.UUID
	Sender    string
	Recipient string
	Amount    float64
}

type Block struct {
	Timestamp    time.Time
	Transactions []Transaction
	Proof        uint64
	PreviousHash string
}

type BlockChain struct {
	Chain []Block
	difficulty uint8
}

func newTransAction(sender, recipient string, amount float64) *Transaction {
	return &Transaction{Sender: sender, Recipient: recipient, Amount: amount}
}

func newBlock(trans []Transaction, proof uint64, previousHash string) *Block {
	return &Block{time.Now().UTC(), trans, proof, previousHash}
}

func newBlockChain(difficulty uint8) *BlockChain {
	return &BlockChain{Chain: []Block{Block{Proof:100, PreviousHash: "1"}}, difficulty:difficulty}
}

func (b *Block) addTransaction(t *Transaction) *Transaction {
	t.Id, _ = uuid.NewV1()
	b.Transactions = append(b.Transactions, *t)
	return t
}

func (b *Block) Hash() string {
	js, _ := json.Marshal(b)
	h := sha256.New()
	h.Write(js)
	return hex.EncodeToString(h.Sum(nil))
}

func (bc *BlockChain) MineBlock(node *Node) *Block {
	proof := bc.ProofOfWork()
	t := newTransAction("0", node.NodeId.String(), 1)
	bc.AddTransaction(t)
	trans := make([]Transaction,1)
	trans[0] = *t

	b := newBlock(trans, proof, bc.LastBlock().Hash())
	bc.AddBlock(*b)

	return b
}

func (bc *BlockChain) AddBlock(b Block) {
	bc.Chain = append(bc.Chain, b)
}

func (bc *BlockChain) AddTransaction(t *Transaction) *Transaction {
	return bc.Chain[len(bc.Chain)-1].addTransaction(t)
}

func (bc *BlockChain) validateProof(lastProof, proof uint64) bool {
	h := sha256.New()
	gb := make([]byte, 8)
	binary.LittleEndian.PutUint64(gb, lastProof)
	h.Write(gb)

	lb := make([]byte, 8)
	binary.LittleEndian.PutUint64(lb, proof)
	h.Write(lb)

	guess := h.Sum(nil)

	byteDifficulty := bc.difficulty /8
	bitDifficulty := bc.difficulty % 8

	//log.Printf("Guess: %d, %024x", proof, binary.LittleEndian.Uint64(guess))
	//log.Printf("byte - bit : %d - %d", byteDifficulty, bitDifficulty)
	for i := 0; i < int(byteDifficulty); i++ {
		if guess[i] != 0 {
			return false
		}
	}

	if bitDifficulty == 0 {
		log.Printf("Last byte: %08b", uint8(guess[byteDifficulty-1]))
		log.Printf("Proof: %d", proof)
		return true
	}

    if (guess[byteDifficulty] >> (8 - bitDifficulty)) == 0  {
		log.Printf("Last byte: %08b", uint8(guess[byteDifficulty]))
    	log.Printf("Proof: %d", proof)
		return true
	}
	return false
}

func (bc *BlockChain) ProofOfWork() (proof uint64) {
	proof = 0
	for {
		if bc.validateProof( bc.Chain[len(bc.Chain)-1].Proof, proof) {
			return
		}
		proof += 1
	}
}

func (bc *BlockChain) LastBlock() *Block {
	return &bc.Chain[len(bc.Chain)-1]
}

func (bc *BlockChain) ValidateChain(chain *BlockChain) bool {

	lastBlock := bc.Chain[0]

	for _, block := range chain.Chain[1:] {
		if block.PreviousHash != lastBlock.Hash() {
			return false
		}

		if !chain.validateProof(lastBlock.Proof, block.Proof) {
			return false
		}
		lastBlock = block
	}
	return true
}