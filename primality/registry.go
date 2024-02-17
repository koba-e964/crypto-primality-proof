package primality

import (
	"errors"
	"fmt"
	"math/big"
)

var ErrMissingDependency = errors.New("missing dependency")

type Registry struct {
	Proofs []Proof `json:"proofs"`
}

// Check checks if the proofs in the registry is correct and self-contained.
func (r *Registry) Check() error {
	seen := map[string]struct{}{}
	for _, proof := range r.Proofs {
		// if the proof is incorrect, there is no way the registry is correct
		if err := proof.Check(); err != nil {
			return err
		}
		seen[(*big.Int)(proof.N).String()] = struct{}{}
	}
	for _, proof := range r.Proofs {
		dep := proof.Dep()
		for _, d := range dep {
			if _, ok := seen[d.String()]; !ok {
				return errors.Join(fmt.Errorf("error in verifying %s (missing dependency: %s)", (*big.Int)(proof.N).String(), d.String()), ErrMissingDependency)
			}
		}
	}
	return nil
}
