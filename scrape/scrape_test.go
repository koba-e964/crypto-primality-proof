package scrape

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestParseRawProofPage(t *testing.T) {
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
		if assert.NoError(t, err) {
			err := cert.Check()
			assert.NoError(t, err)
		}
	}
}
