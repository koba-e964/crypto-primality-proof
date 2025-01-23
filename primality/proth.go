package primality

import (
	"errors"
	"math/big"
)

var ErrNotProth = errors.New("not Proth number")

// ProveProth tries to prove that a Proth number n is prime.
//
// https://en.wikipedia.org/wiki/Proth%27s_theorem
func ProveProth(n *big.Int) (*Proof, error) {
	if n.Cmp(big.NewInt(1)) <= 0 {
		return nil, ErrNotPrime
	}
	if n.Bit(0) == 0 {
		if n.Cmp(big.NewInt(2)) == 0 {
			return nil, ErrNotProth
		}
		return nil, ErrNotPrime
	}
	a := big.NewInt(1)
	apow := 0
	ncp := big.NewInt(0).Sub(n, big.NewInt(1))
	if ncp.Cmp(big.NewInt(0)) <= 0 {
		panic("n - 1 <= 0 cannot happen")
	}
	for ncp.Bit(0) == 0 {
		ncp.Rsh(ncp, 1)
		a.Lsh(a, 1)
		apow++
	}
	if ncp.Cmp(a) >= 0 {
		return nil, ErrNotProth
	}
	if a == nil {
		return nil, ErrNotPrime
	}
	aFactorization := &FactoredInt{
		Int: (*BigInt)(a),
		Factorization: []FactorEntry{
			{
				Prime:    (*BigInt)(big.NewInt(2)),
				Exponent: apow,
			},
		},
	}
	base := big.NewInt(2)
	for base.Cmp(big.NewInt(100)) < 0 {
		inverses, err := checkGen(n,
			aFactorization,
			base,
		)
		if err != nil {
			base.Add(base, big.NewInt(1))
			continue
		}
		return &Proof{
			N: (*BigInt)(n),
			GeneralizedPocklington: &GeneralizedPocklingtonProof{
				A:        aFactorization,
				Base:     (*BigInt)(base),
				Inverses: inverses,
			},
		}, nil
	}
	return nil, ErrNotPrime
}
