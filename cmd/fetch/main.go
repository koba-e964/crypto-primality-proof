package main

import (
	"encoding/json"
	"os"

	"github.com/koba-e964/crypto-primality-proof/scrape"
)

func main() {
	urlBase := "https://safecurves.cr.yp.to"
	result, err := scrape.ReadPrimeProofsPage(urlBase, "Curve25519")
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
	os.WriteFile("Curve25519.json", jsonString, 0o644)
}
