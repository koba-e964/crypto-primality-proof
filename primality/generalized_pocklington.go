package primality

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"
)

type GeneralizedPocklingtonProof struct {
	A        *FactoredInt `json:"a,omitempty"` // N = A * B
	Base     *BigInt      `json:"base,omitempty"`
	Inverses []Inverse    `json:"inverses,omitempty"`
}

func (p *GeneralizedPocklingtonProof) Check(N *big.Int) error {
	if err := p.A.Check(); err != nil {
		return errors.Join(fmt.Errorf("invalid A in verifying %s", N.String()), err)
	}
	A := (*big.Int)(p.A.Int)
	NMinus1 := big.NewInt(0).Sub(N, big.NewInt(1))
	B, NModA := big.NewInt(0).DivMod(NMinus1, A, big.NewInt(0))
	if NModA.Cmp(big.NewInt(0)) != 0 {
		return fmt.Errorf("pocklington: N-1 is not divisible by A: not (%s | %s)", A.String(), NMinus1.String())
	}
	if B.Cmp(A) >= 0 {
		return fmt.Errorf("A^2 > N must hold")
	}
	if B.ModInverse(B, A) == nil {
		return fmt.Errorf("pocklington: gcd(A, B) != 1")
	}
	fromInverse := map[string]struct{}{}
	for _, inv := range p.Inverses {
		if err := inv.Check(); err != nil {
			return err
		}
		if (*big.Int)(inv.Mod).Cmp(N) != 0 {
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
		value := big.NewInt(0).Exp((*big.Int)(p.Base), exp, N)
		value.Sub(value, big.NewInt(1))
		value.Mod(value, N)
		fromBase[value.String()] = struct{}{}
	}
	if !reflect.DeepEqual(fromInverse, fromBase) {
		return fmt.Errorf("set of inverses is not correct")
	}
	return nil
}

func (p *GeneralizedPocklingtonProof) Dep() []*big.Int {
	dep := []*big.Int{}
	for _, entry := range p.A.Factorization {
		dep = append(dep, (*big.Int)(entry.Prime))
	}
	return dep
}
