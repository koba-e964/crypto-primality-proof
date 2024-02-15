package main

import (
	"fmt"

	"github.com/koba-e964/crypto-primality-proof/scrape"
)

func main() {
	url := "https://safecurves.cr.yp.to/primeproofs.html"
	reader, err := scrape.GetContent(url)
	if err != nil {
		panic(err)
	}
	defer reader.Close()
	pps, err := scrape.ParsePrimeProofsPage(reader, "Curve25519")
	if err != nil {
		panic(err)
	}
	fmt.Println(pps.Numbers)
}
