package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

var maxNonce=math.MaxInt64
const TARGET_BITS  = 16

type ProofOfWork struct {
	block *Block
	target *big.Int
}

func NewProofOfWork(b *Block) * ProofOfWork{
	target := big.NewInt(1);
	target.Lsh(target, uint(256-TARGET_BITS))
	return & ProofOfWork{b, target}
}


func (p *ProofOfWork)prepareData(nonce int) []byte{
	data := bytes.Join([][]byte{
		p.block.PrevBlockHash,
		p.block.Data,
		IntToHex(int64(TARGET_BITS)),
		IntToHex(int64(p.block.Nonce)),
	}, []byte{})

	return data
}

// Run performs a proof-of-work
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining a new block")
	for nonce < maxNonce {
		hash = sha256.Sum256(pow.prepareData(nonce))
		if math.Remainder(float64(nonce), 100000) == 0 {
			fmt.Printf("\r%x", hash)
		}
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 {
			break
		}
		nonce++

	}
	fmt.Print("\n\n")
	return nonce, hash[:]
}