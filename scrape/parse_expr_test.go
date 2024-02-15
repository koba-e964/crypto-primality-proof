package scrape

import (
	"math/big"
	"testing"

	"github.com/koba-e964/crypto-primality-proof/primality"
	"github.com/stretchr/testify/assert"
)

func TestParseExpr(t *testing.T) {
	tests := []struct {
		expr     string
		expected *primality.FactoredInt
	}{
		{
			expr: "3^2 * 5",
			expected: &primality.FactoredInt{
				Int: (*primality.BigInt)(big.NewInt(45)),
				Factorization: []primality.FactorEntry{
					{Prime: (*primality.BigInt)(big.NewInt(3)), Exponent: 2},
					{Prime: (*primality.BigInt)(big.NewInt(5)), Exponent: 1},
				},
			},
		},
	}
	for _, test := range tests {
		actual, err := ParseExpr(test.expr)
		if assert.NoError(t, err) {
			assert.Equal(t, test.expected, actual)
		}
	}
}
