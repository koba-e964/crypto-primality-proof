package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/koba-e964/crypto-primality-proof/primality"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		panic("no arguments")
	}
	failed := false
	for _, filename := range args {
		dat, err := os.ReadFile(filename)
		if err != nil {
			failed = true
			log.Print(fmt.Errorf("failed to read %s: %w", filename, err))
			continue
		}
		var reg primality.Registry
		if err := json.Unmarshal(dat, &reg); err != nil {
			failed = true
			log.Print(fmt.Errorf("failed to unmarshal %s: %w", filename, err))
			continue
		}
		if err := reg.Check(); err != nil {
			failed = true
			log.Print(fmt.Errorf("failed to verify %s: %w", filename, err))
			continue
		}
	}
	if failed {
		os.Exit(1)
	}
}
