package main

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/koba-e964/crypto-primality-proof/primality"
)

func main() {
	// https://safecurves.cr.yp.to/proof/257.html
	cert := primality.Proof{
		N: (*primality.BigInt)(big.NewInt(257)),
		A: &primality.FactoredInt{
			Int: (*primality.BigInt)(big.NewInt(256)),
			Factorization: []primality.FactorEntry{
				{Prime: (*primality.BigInt)(big.NewInt(2)), Exponent: 8},
			},
		},
		Base: (*primality.BigInt)(big.NewInt(3)),
		Inverses: []primality.Inverse{
			{
				Mod:   (*primality.BigInt)(big.NewInt(257)),
				Value: (*primality.BigInt)(big.NewInt(255)),
				Inv:   (*primality.BigInt)(big.NewInt(128)),
			},
		},
	}
	if err := cert.Check(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("proof is correct")
		str, err := json.MarshalIndent(cert, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(str))
		var data primality.Proof
		if err := json.Unmarshal(str, &data); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("unmarshalled proof is correct")
		}
	}
	for i := int64(3); i < 100; i++ {
		cert1, err := primality.Prove(big.NewInt(i))
		if err != nil {
			continue
		}
		str, err := json.MarshalIndent(cert1, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(str))
		if err := cert1.Check(); err != nil {
			panic(err)
		}
	}
}
