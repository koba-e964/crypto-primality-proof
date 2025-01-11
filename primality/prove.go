package primality

import (
	"errors"
	"math/big"
)

var ErrNotPrime = errors.New("not prime")

func findA(n *big.Int) *FactoredInt {
	p := big.NewInt(2)
	rem := big.NewInt(0).Set(n)
	factors := []FactorEntry{}
	for rem.Cmp(big.NewInt(1)) > 0 && !rem.ProbablyPrime(20) {
		e := 0
		for big.NewInt(0).Rem(rem, p).Cmp(big.NewInt(0)) == 0 {
			rem.Div(rem, p)
			e++
		}
		if e > 0 {
			factors = append(factors, FactorEntry{Prime: (*BigInt)(new(big.Int).Set(p)), Exponent: e})
		}
		p.Add(p, big.NewInt(1))
	}
	if rem.Cmp(big.NewInt(1)) > 0 && n.Cmp(new(big.Int).Mul(rem, rem)) < 0 {
		return &FactoredInt{
			Int: (*BigInt)(rem),
			Factorization: []FactorEntry{
				{
					Prime:    (*BigInt)(rem),
					Exponent: 1,
				},
			},
		}
	}
	return &FactoredInt{
		Int:           (*BigInt)(new(big.Int).Div(n, rem)),
		Factorization: factors,
	}
}

func checkGen(n *big.Int, a *FactoredInt, base *big.Int) ([]Inverse, error) {
	nMinus1 := big.NewInt(0).Sub(n, big.NewInt(1))
	invs := []Inverse{}
	seen := map[string]struct{}{}
	for _, entry := range a.Factorization {
		pr := (*big.Int)(entry.Prime)
		exp := big.NewInt(0).Div(nMinus1, pr)
		value := big.NewInt(0).Exp(base, exp, n)
		value.Sub(value, big.NewInt(1))
		value.Mod(value, n)
		inv := new(big.Int).ModInverse(value, n)
		if inv == nil {
			return nil, errors.New("not invertible")
		}
		valueString := value.String()
		if _, ok := seen[valueString]; !ok {
			invs = append(invs, Inverse{
				Mod:   (*BigInt)(n),
				Value: (*BigInt)(value),
				Inv:   (*BigInt)(inv),
			})
		}
	}

	return invs, nil
}

func Prove(n *big.Int) (*Proof, error) {
	if n.Cmp(big.NewInt(2)) == 0 {
		return &Proof{
			N: (*BigInt)(n),
		}, nil
	}
	if !n.ProbablyPrime(20) {
		return nil, ErrNotPrime
	}
	nMinus1 := big.NewInt(0).Sub(n, big.NewInt(1))
	a := findA(nMinus1)
	base := big.NewInt(2)
	for {
		if invs, err := checkGen(n, a, base); err == nil {
			return &Proof{
				N: (*BigInt)(n),
				Proof: &GeneralizedPocklingtonProof{
					A:        a,
					Base:     (*BigInt)(base),
					Inverses: invs,
				},
			}, nil
		}
		base.Add(base, big.NewInt(1))
	}
}
