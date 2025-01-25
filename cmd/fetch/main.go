package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/koba-e964/crypto-primality-proof/scrape"
)

func main() {
	urlBase := "https://safecurves.cr.yp.to"
	curveNames := []string{
		"Curve25519",
		"NIST P-256",
		"secp256k1",
		"Ed448-Goldilocks",
		"BN(2,254)",
	}
	for _, curveName := range curveNames {
		result, err := scrape.ReadPrimeProofsPage(urlBase, curveName)
		if err != nil {
			panic(err)
		}
		reg, err := result.Translate()
		if err != nil {
			panic(err)
		}
		jsonString, err := json.MarshalIndent(reg, "", "  ")
		if err != nil {
			panic(err)
		}
		jsonString = append(jsonString, '\n')
		os.WriteFile(curveName+".json", jsonString, 0o644)
		log.Println("wrote", curveName+".json")
	}
}
