package scrape

import (
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Takes about 30 seconds
func TestReadPrimeProofsPage(t *testing.T) {
	urlBase := "https://safecurves.cr.yp.to"
	result, err := ReadPrimeProofsPage(urlBase, "Curve25519")
	if assert.NoError(t, err) {
		assert.Len(t, result.Numbers, 151)
		reg, err := result.Translate()
		if assert.NoError(t, err) {
			err := reg.Check()
			assert.NoError(t, err)
		}
	}
}

func TestParsePrimeProofsPage(t *testing.T) {
	// Modified version of https://safecurves.cr.yp.to/primeproofs.html
	text := `
	<html>
	<body>
	<table border>
	<tr>
	<th><p>Curve</p></th>
	<th><p>Relevant proven primes</p></th>
	</tr>
	<tr>
	<td><p>
	TestCurve
	</p></td>
	<td><p><font size=1>
	<a href=proof/2.html>2</a>
	<a href=proof/3.html>3</a>
	<a href=proof/5.html>5</a>
	</font></p></td>
	</tr>
	</table>
	</body>
	</html>
		`
	s, err := ParsePrimeProofsPage(strings.NewReader(text), "TestCurve")
	if assert.NoError(t, err) {
		assert.Len(t, s.Numbers, 3)
		assert.Equal(t, []string{"2", "3", "5"}, s.Numbers)
	}
}

func TestParseRawProofPage3(t *testing.T) {
	text := `Primality proof for n = 3:
	Take b = 2.
	
	b^(n-1) mod n = 1.
	
	2 is prime.
	b^((n-1)/2)-1 mod n = 1, which is a unit, inverse 1.
	
	(2) divides n-1.
	
	(2)^2 > n.
	
	n is prime by Pocklington's theorem.`

	s, err := ParseRawProofPage(text)
	if assert.NoError(t, err) {
		cert, err := s.Translate()
		assert.Equal(t, "3", (*big.Int)(cert.N).String())
		if assert.NoError(t, err) {
			err := cert.Check()
			assert.NoError(t, err)
		}
	}
}

func TestParseRawProofPage7(t *testing.T) {
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
	s, err := ParseRawProofPage(text)
	if assert.NoError(t, err) {
		cert, err := s.Translate()
		assert.Equal(t, "7", (*big.Int)(cert.N).String())
		if assert.NoError(t, err) {
			err := cert.Check()
			assert.NoError(t, err)
		}
	}
}

func TestParseRawProofPage2(t *testing.T) {
	// https://safecurves.cr.yp.to/proof/2.html
	text := "2 is prime."
	s, err := ParseRawProofPage(text)
	if assert.NoError(t, err) {
		cert, err := s.Translate()
		if assert.NoError(t, err) {
			err := cert.Check()
			assert.NoError(t, err)
		}
	}
}
