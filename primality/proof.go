package primality

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
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
		return fmt.Errorf("factorization is incorrect")
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
	N        *big.Int     `json:"n"`
	A        *FactoredInt `json:"a"` // N = A * B
	Base     *big.Int     `json:"base"`
	Inverses []Inverse    `json:"inverses"`
}

// Check checks the correctness of the proof per se,
// i.e., it does not check if its dependencies are correct.
func (p *Proof) Check() error {
	if err := p.A.Check(); err != nil {
		return err
	}
	A := (*big.Int)(p.A.Int)
	NMinus1 := big.NewInt(0).Sub(p.N, big.NewInt(1))
	NModA := big.NewInt(0)
	NModA.Mod(NMinus1, A)
	if NModA.Cmp(big.NewInt(0)) != 0 {
		return fmt.Errorf("pocklington: N-1 is not divisible by A")
	}
	fromInverse := map[string]struct{}{}
	for _, inv := range p.Inverses {
		if err := inv.Check(); err != nil {
			return err
		}
		if (*big.Int)(inv.Mod).Cmp(p.N) != 0 {
			return fmt.Errorf("invalid modulus in inverse")
		}
		invString := (*big.Int)(inv.Value).String()
		if _, ok := fromInverse[invString]; ok {
			return fmt.Errorf("duplicate inverse")
		}
		fromInverse[invString] = struct{}{}
	}
	fromBase := map[string]struct{}{}
	for _, entry := range p.A.Factorization {
		pr := (*big.Int)(entry.Prime)
		exp := big.NewInt(0).Div(NMinus1, pr)
		value := big.NewInt(0).Exp(p.Base, exp, p.N)
		value.Sub(value, big.NewInt(1))
		value.Mod(value, p.N)
		fromBase[value.String()] = struct{}{}
	}
	if !reflect.DeepEqual(fromInverse, fromBase) {
		return fmt.Errorf("set of inverses is not correct")
	}
	tmp := big.NewInt(0)
	tmp.Exp(A, big.NewInt(2), nil)
	if tmp.Cmp(p.N) <= 0 {
		return fmt.Errorf("A^2 > N must hold")
	}
	return nil
}

// Dep returns the dependencies of the proof.
func (p *Proof) Dep() []*big.Int {
	deps := []*big.Int{}
	for _, entry := range p.A.Factorization {
		deps = append(deps, (*big.Int)(entry.Prime))
	}
	return deps
}
