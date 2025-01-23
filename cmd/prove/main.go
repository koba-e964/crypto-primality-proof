package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"slices"

	"github.com/koba-e964/crypto-primality-proof/primality"
)

func main() {
	argv := os.Args
	if len(argv) < 2 {
		panic("missing argument: integer")
	}
	nString := argv[1]
	n := big.NewInt(0)
	if _, ok := n.SetString(nString, 10); !ok {
		panic("invalid integer: " + nString)
	}
	seen := map[string]struct{}{}
	stack := []*big.Int{n}
	registry := primality.Registry{}
	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		nString := n.String()
		if _, ok := seen[nString]; ok {
			continue
		}
		seen[nString] = struct{}{}
		proof, err := primality.Prove(n)
		if err != nil {
			panic(err)
		}
		registry.Proofs = append(registry.Proofs, *proof)
		dep := proof.Dep()
		stack = append(stack, dep...)
	}
	slices.Reverse(registry.Proofs)
	jsonString, err := json.MarshalIndent(registry, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonString))
}
