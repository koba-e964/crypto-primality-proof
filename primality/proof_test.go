package primality

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFactoredIntCheck(t *testing.T) {
	a := FactoredInt{
		Int: (*BigInt)(big.NewInt(45)),
		Factorization: []FactorEntry{
			{Prime: (*BigInt)(big.NewInt(3)), Exponent: 2},
			{Prime: (*BigInt)(big.NewInt(5)), Exponent: 2}, // should be 1
		},
	}
	assert.Error(t, a.Check())
}

func TestInverseCheck(t *testing.T) {
	inv := Inverse{
		Mod:   (*BigInt)(big.NewInt(7)),
		Value: (*BigInt)(big.NewInt(3)),
		Inv:   (*BigInt)(big.NewInt(3)), // should be 5
	}
	assert.Error(t, inv.Check())
}

func TestProofCheck2(t *testing.T) {
	cert := Proof{
		N: (*BigInt)(big.NewInt(2)),
	}
	assert.NoError(t, cert.Check())
	assert.Len(t, cert.Dep(), 0)
}

func TestProofCheck181(t *testing.T) {
	// https://safecurves.cr.yp.to/proof/181.html
	cert := Proof{
		N: (*BigInt)(big.NewInt(181)),
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
	}

	assert.NoError(t, cert.Check())
	assert.Equal(t, []*big.Int{big.NewInt(3), big.NewInt(5)}, cert.Dep())
	str, err := json.MarshalIndent(cert, "", "  ")
	if assert.NoError(t, err) {
		var data Proof
		assert.NoError(t, json.Unmarshal(str, &data))
	}
}

func TestProofCheck15(t *testing.T) {
	cert := Proof{
		N: (*BigInt)(big.NewInt(15)),
		A: &FactoredInt{
			Int: (*BigInt)(big.NewInt(2)),
			Factorization: []FactorEntry{
				{Prime: (*BigInt)(big.NewInt(2)), Exponent: 1},
			},
		},
		Base: (*BigInt)(big.NewInt(14)),
		Inverses: []Inverse{
			{
				Mod:   (*BigInt)(big.NewInt(15)),
				Value: (*BigInt)(big.NewInt(13)),
				Inv:   (*BigInt)(big.NewInt(7)),
			},
		},
	}
	assert.EqualError(t, cert.Check(), "A^2 > N must hold")
}

func TestProofCheck255(t *testing.T) {
	cert := Proof{
		N: (*BigInt)(big.NewInt(255)),
		A: &FactoredInt{
			Int: (*BigInt)(big.NewInt(8)),
			Factorization: []FactorEntry{
				{Prime: (*BigInt)(big.NewInt(2)), Exponent: 3},
			},
		},
		Base: (*BigInt)(big.NewInt(2)),
		Inverses: []Inverse{
			{
				Mod:   (*BigInt)(big.NewInt(255)),
				Value: (*BigInt)(big.NewInt(1)),
				Inv:   (*BigInt)(big.NewInt(1)),
			},
		},
	}
	assert.EqualError(t, cert.Check(), "pocklington: N-1 is not divisible by A")
}

func TestProofCheck257(t *testing.T) {
	// https://safecurves.cr.yp.to/proof/257.html
	cert := Proof{
		N: (*BigInt)(big.NewInt(257)),
		A: &FactoredInt{
			Int: (*BigInt)(big.NewInt(256)),
			Factorization: []FactorEntry{
				{Prime: (*BigInt)(big.NewInt(2)), Exponent: 8},
			},
		},
		Base: (*BigInt)(big.NewInt(3)),
		Inverses: []Inverse{
			{
				Mod:   (*BigInt)(big.NewInt(257)),
				Value: (*BigInt)(big.NewInt(255)),
				Inv:   (*BigInt)(big.NewInt(128)),
			},
		},
	}
	assert.NoError(t, cert.Check())
	assert.Equal(t, []*big.Int{big.NewInt(2)}, cert.Dep())
	str, err := json.MarshalIndent(cert, "", "  ")
	if assert.NoError(t, err) {
		var data Proof
		assert.NoError(t, json.Unmarshal(str, &data))
	}
}

func TestProofCheckInvalidMod(t *testing.T) {
	cert := Proof{
		N: (*BigInt)(big.NewInt(257)),
		A: &FactoredInt{
			Int: (*BigInt)(big.NewInt(256)),
			Factorization: []FactorEntry{
				{Prime: (*BigInt)(big.NewInt(2)), Exponent: 8},
			},
		},
		Base: (*BigInt)(big.NewInt(3)),
		Inverses: []Inverse{
			{
				Mod:   (*BigInt)(big.NewInt(4)),
				Value: (*BigInt)(big.NewInt(3)),
				Inv:   (*BigInt)(big.NewInt(3)),
			},
		},
	}
	assert.EqualError(t, cert.Check(), "invalid modulus in inverse")
}

func TestProofCheckInverseSetNotCorrect(t *testing.T) {
	cert := Proof{
		N: (*BigInt)(big.NewInt(257)),
		A: &FactoredInt{
			Int: (*BigInt)(big.NewInt(256)),
			Factorization: []FactorEntry{
				{Prime: (*BigInt)(big.NewInt(2)), Exponent: 8},
			},
		},
		Base: (*BigInt)(big.NewInt(3)),
		Inverses: []Inverse{
			{
				Mod:   (*BigInt)(big.NewInt(257)),
				Value: (*BigInt)(big.NewInt(1)),
				Inv:   (*BigInt)(big.NewInt(1)),
			},
		},
	}
	assert.EqualError(t, cert.Check(), "set of inverses is not correct")
}
