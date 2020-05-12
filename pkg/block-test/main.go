package main

func main() {
	b := NewGenesisBlock([]byte("Hi GenesisBlock"))
	NewBlock([]byte("trans"), b.PrevBlockHash, b.Height+1)
}
