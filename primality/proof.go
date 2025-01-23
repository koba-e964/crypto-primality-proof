package primality

import (
	"encoding/json"
	"fmt"
	"math/big"
)

// Wrapper type. We need a big.Int to marshal/unmarshal to/from a string
type BigInt big.Int

func (b *BigInt) MarshalJSON() ([]byte, error) {
	return json.Marshal((*big.Int)(b).String())
}

func (b *BigInt) UnmarshalJSON(data []byte) error {
	aux := ""
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if _, ok := (*big.Int)(b).SetString(aux, 10); !ok {
		return fmt.Errorf("invalid int")
	}
	return nil
}

type FactorEntry struct {
	Prime    *BigInt `json:"prime"`
	Exponent int     `json:"exponent"`
}

type FactoredInt struct {
	Int           *BigInt       `json:"int"`
	Factorization []FactorEntry `json:"factorization"`
}

func (f *FactoredInt) Check() error {
	// check if the factorization is correct
	product := big.NewInt(1)
	for _, entry := range f.Factorization {
		prime := (*big.Int)(entry.Prime)
		exponent := entry.Exponent
		product.Mul(product, big.NewInt(0).Exp((*big.Int)(prime), big.NewInt(int64(exponent)), nil))
	}
	if product.Cmp((*big.Int)(f.Int)) != 0 {
		return fmt.Errorf(
			"factorization is incorrect: %s != %s",
			product.String(),
			(*big.Int)(f.Int).String(),
		)
	}
	return nil
}

type Inverse struct {
	Mod   *BigInt `json:"mod"`
	Value *BigInt `json:"value"`
	Inv   *BigInt `json:"inv"`
}

func (i *Inverse) Check() error {
	// check if the inverse is correct
	prod := big.NewInt(0)
	prod.Mul((*big.Int)(i.Value), (*big.Int)(i.Inv))
	prod.Mod(prod, (*big.Int)(i.Mod))
	if prod.Cmp(big.NewInt(1)) != 0 {
		return fmt.Errorf("inverse is incorrect")
	}
	return nil
}

type Proof struct {
	N                      *BigInt                      `json:"n"`
	GeneralizedPocklington *GeneralizedPocklingtonProof `json:"generalized-pocklington,omitempty"`
}

// Check checks the correctness of the proof per se,
// i.e., it does not check if its dependencies are correct.
func (p *Proof) Check() error {
	// if N = 2, N is prime.
	N := (*big.Int)(p.N)
	if N.Cmp(big.NewInt(2)) == 0 {
		return nil
	}
	proved := false
	if p.GeneralizedPocklington != nil {
		if err := p.GeneralizedPocklington.Check(N); err != nil {
			return err
		}
		proved = true
	}
	if !proved {
		return fmt.Errorf("no proof provided")
	}
	return nil
}

// Dep returns the dependencies of the proof.
func (p *Proof) Dep() []*big.Int {
	// if N = 2, N is known to be prime and the proof depends on nothing.
	if (*big.Int)(p.N).Cmp(big.NewInt(2)) == 0 {
		return nil
	}
	deps := []*big.Int{}
	if p.GeneralizedPocklington != nil {
		deps = append(deps, p.GeneralizedPocklington.Dep()...)
	}
	return deps
}
