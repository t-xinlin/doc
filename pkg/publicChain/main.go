package main

import (
	"fmt"
	"math/big"
)

const tarBits = 24

func main() {
	tar := big.NewInt(1)
	fmt.Printf("big int: %v\n", tar)
	tar.Lsh(tar, uint(256-tarBits))
	fmt.Printf("big int: %v\n", tar)
}
