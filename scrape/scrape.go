package scrape

import (
	"errors"
	"io"
	"math/big"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/koba-e964/crypto-primality-proof/primality"
)

var (
	primalityProof = regexp.MustCompile(`Primality proof for n = (\d+):`)
	take           = regexp.MustCompile(`Take b = (\d+).`)
	divides        = regexp.MustCompile(`\((.+)\) divides n-1.`)
	inv            = regexp.MustCompile(`n = (\d+), which is a unit, inverse (\d+).`)
)

var ErrNotANumber = errors.New("not a number")

type RawPrimeProofs struct {
	Numbers []string
}

type RawProofPage struct {
	N        string
	AExpr    string
	B        string
	Inverses [][2]string
}

func GetContent(url string) (io.ReadCloser, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}

func ParsePrimeProofsPage(reader io.Reader, curveName string) (*RawPrimeProofs, error) {
	document, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	numbers := []string{}
	document.Find("tbody").Each(func(_ int, s *goquery.Selection) {
		if len(s.Children().Nodes) >= 2 {
			s.Find("tr").Each(func(x int, s *goquery.Selection) {
				data := []string{}
				links := []string{}
				s.Find("td").Each(func(x int, s *goquery.Selection) {
					html, err := s.Html()
					if err != nil {
						panic(err)
					}
					data = append(data, html)
					s.Find("a").Each(func(x int, s *goquery.Selection) {
						linkText := s.Text()
						links = append(links, linkText)
					})
				})
				if len(data) >= 1 && strings.Contains(data[0], curveName) {
					numbers = links
					return
				}
			})
		}
	})
	return &RawPrimeProofs{
		Numbers: numbers,
	}, nil
}

func ParseRawProofPage(s string) (*RawProofPage, error) {
	nString := primalityProof.FindStringSubmatch(s)[1]
	bString := take.FindStringSubmatch(s)[1]
	aString := divides.FindStringSubmatch(s)[1]
	invs := inv.FindAllStringSubmatch(s, -1)
	inverses := make([][2]string, len(invs))
	for i, inv := range invs {
		inverses[i] = [2]string{inv[1], inv[2]}
	}
	return &RawProofPage{
		N:        nString,
		AExpr:    aString,
		B:        bString,
		Inverses: inverses,
	}, nil
}

func (r *RawProofPage) Translate() (*primality.Proof, error) {
	var p primality.Proof
	n, ok := new(big.Int).SetString(r.N, 10)
	if !ok {
		return nil, ErrNotANumber
	}
	p.N = (*primality.BigInt)(n)
	a, err := ParseExpr(r.AExpr)
	if err != nil {
		return nil, err
	}
	p.A = a
	base, ok := new(big.Int).SetString(r.B, 10)
	if !ok {
		return nil, ErrNotANumber
	}
	p.Base = (*primality.BigInt)(base)
	inverses := make([]primality.Inverse, len(r.Inverses))
	for i, inv := range r.Inverses {
		value, ok := new(big.Int).SetString(inv[0], 10)
		if !ok {
			return nil, ErrNotANumber
		}
		inv, ok := new(big.Int).SetString(inv[1], 10)
		if !ok {
			return nil, ErrNotANumber
		}
		inverses[i] = primality.Inverse{
			Mod:   (*primality.BigInt)(n),
			Value: (*primality.BigInt)(value),
			Inv:   (*primality.BigInt)(inv),
		}
	}
	p.Inverses = inverses
	return &p, nil
}
