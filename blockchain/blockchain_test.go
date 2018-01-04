package blockchain

import "testing"

func TestNewBlockChain(t *testing.T) {
	bc := newBlockChain(16)
	if bc.Chain[0].Proof != 100 {
		t.Errorf("Proof of genesis block is unequal to 100: %d", bc.Chain[0].Proof != 100)
	}

	if bc.Chain[0].PreviousHash != "1" {
		t.Errorf("Previous has of genesis block is unequal to '1': %s", bc.Chain[0].PreviousHash != "1")
	}
}

func TestBlockChain_ProofOfWork(t *testing.T) {
	bc := newBlockChain(23)
	proof := bc.ProofOfWork()

	if proof != 1134054 {
		t.Errorf("First proof is unequal to 1134054: %d", proof)
	}
}

func TestValidateProof(t *testing.T) {
	bc := newBlockChain(16)
	if !bc.validateProof( bc.Chain[len(bc.Chain)-1].Proof, 8387) {
		t.Errorf("First Proof based on genesis block not valid")
	}
}

func TestBlockChain_ValidateChain(t *testing.T) {
	bc := newBlockChain(16)

	node := newNode(7070)
	for i := 0; i < 10; i++ {
		bc.MineBlock(node)
	}

	if !bc.ValidateChain(bc) {
		t.Error("block chain should be valid, validating it returned false")
	}

}