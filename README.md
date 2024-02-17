# crypto-primality-proof
This tool reads from https://safecurves.cr.yp.to/primeproofs.html and verifies the displayed proofs are valid.

# How to use from external packages

```go
package main

import (
	"encoding/json"
	"os"

	"github.com/koba-e964/crypto-primality-proof/scrape"
)

func main() {
	urlBase := "https://safecurves.cr.yp.to"
	curveName := "Curve25519"
	result, err := scrape.ReadPrimeProofsPage(urlBase, curveName)
	if err != nil {
		panic(err)
	}
	reg, err := result.Translate()
	if err != nil {
		panic(err)
	}
	// Checks if the proofs are valid and self-contained as a whole
	if err := reg.Check(); err != nil {
		panic(err)
	}
}
```
