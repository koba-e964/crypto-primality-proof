package primality

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistryCheckSuccess(t *testing.T) {
	// https://safecurves.cr.yp.to/proof/3.html
	cert3 := Proof{
		N: (*BigInt)(big.NewInt(3)),
		Proof: &GeneralizedPocklingtonProof{
			A: &FactoredInt{
				Int: (*BigInt)(big.NewInt(2)),
				Factorization: []FactorEntry{
					{Prime: (*BigInt)(big.NewInt(2)), Exponent: 1},
				},
			},
			Base: (*BigInt)(big.NewInt(2)),
			Inverses: []Inverse{
				{
					Mod:   (*BigInt)(big.NewInt(3)),
					Value: (*BigInt)(big.NewInt(1)),
					Inv:   (*BigInt)(big.NewInt(1)),
				},
			},
		},
	}
	cert2 := Proof{
		N: (*BigInt)(big.NewInt(2)),
	}
	registry := &Registry{
		Proofs: []Proof{cert3, cert2}, // order doesn't matter
	}
	assert.NoError(t, registry.Check())
}

func TestRegistryCheckInsufficient(t *testing.T) {
	// https://safecurves.cr.yp.to/proof/181.html
	cert := Proof{
		N: (*BigInt)(big.NewInt(181)),
		Proof: &GeneralizedPocklingtonProof{
			A: &FactoredInt{
				Int: (*BigInt)(big.NewInt(45)),
				Factorization: []FactorEntry{
					{Prime: (*BigInt)(big.NewInt(3)), Exponent: 2},
					{Prime: (*BigInt)(big.NewInt(5)), Exponent: 1},
				},
			},
			Base: (*BigInt)(big.NewInt(2)),
			Inverses: []Inverse{
				{
					Mod:   (*BigInt)(big.NewInt(181)),
					Value: (*BigInt)(big.NewInt(47)),
					Inv:   (*BigInt)(big.NewInt(104)),
				},
				{
					Mod:   (*BigInt)(big.NewInt(181)),
					Value: (*BigInt)(big.NewInt(58)),
					Inv:   (*BigInt)(big.NewInt(103)),
				},
			},
		},
	}
	registry := &Registry{
		Proofs: []Proof{cert},
	}
	assert.Contains(t, registry.Check().Error(), ErrMissingDependency.Error())
}
