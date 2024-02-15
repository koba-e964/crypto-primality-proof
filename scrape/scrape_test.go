package scrape

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGo(t *testing.T) {
	// https://safecurves.cr.yp.to/proof/7.html
	text := `Primality proof for n = 7:
	Take b = 2.
	
	b^(n-1) mod n = 1.
	
	3 is prime.
	b^((n-1)/3)-1 mod n = 3, which is a unit, inverse 5.
	
	(3) divides n-1.
	
	(3)^2 > n.
	
	n is prime by Pocklington's theorem.
	`
	s, err := Parse(text)
	if assert.NoError(t, err) {
		cert, err := s.Translate()
		if assert.NoError(t, err) {
			err := cert.Check()
			assert.NoError(t, err)
		}
	}
}
