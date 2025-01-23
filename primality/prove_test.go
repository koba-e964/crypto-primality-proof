package primality

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProve(t *testing.T) {
	for i := int64(2); i < 100; i++ {
		cert1, err := Prove(big.NewInt(i))
		if err != nil {
			continue
		}
		assert.NoError(t, cert1.Check())
	}
}

func TestProveLarge(t *testing.T) {
	nums := []string{
		"3221225473",
		"221360928884514619393",
	}
	for _, n := range nums {
		bigInt, ok := big.NewInt(0).SetString(n, 10)
		if !ok {
			t.Fatalf("failed to parse %s", n)
		}
		cert1, err := Prove(bigInt)
		assert.NoError(t, err)
		assert.NoError(t, cert1.Check())
	}
}

func BenchmarkProve100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for i := int64(2); i < 100; i++ {
			cert1, err := Prove(big.NewInt(i))
			if err != nil {
				continue
			}
			assert.NoError(b, cert1.Check())
		}
	}
}
