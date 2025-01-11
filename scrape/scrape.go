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
	twoProof       = regexp.MustCompile(`^2 is prime.`)
	take           = regexp.MustCompile(`Take b = (\d+).`)
	divides        = regexp.MustCompile(`\((.+)\) divides n-1.`)
	inv            = regexp.MustCompile(`n = (\d+), which is a unit, inverse (\d+).`)
)

var ErrNotANumber = errors.New("not a number")

type RawPrimeProofs struct {
	CurveName string
	Numbers   []string
	Subpages  []RawProofPage
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
	return res.Body, nil
}

// ReadPrimeProofsPage reads the prime proofs page and returns proofs for the specified curve.
// urlBase should be like "https://safecurves.cr.yp.to"
func ReadPrimeProofsPage(urlBase string, curveName string) (*RawPrimeProofs, error) {
	reader, err := GetContent(urlBase + "/primeproofs.html")
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	proofs, err := ParsePrimeProofsPage(reader, curveName)
	if err != nil {
		return nil, err
	}
	// read subpages
	for _, number := range proofs.Numbers {
		subUrl := urlBase + "/proof/" + number + ".html"
		subReader, err := GetContent(subUrl)
		if err != nil {
			return nil, err
		}
		buf := new(strings.Builder)
		if _, err := io.Copy(buf, subReader); err != nil {
			return nil, err
		}
		subReader.Close()
		subProof, err := ParseRawProofPage(buf.String())
		if err != nil {
			return nil, err
		}
		proofs.Subpages = append(proofs.Subpages, *subProof)
	}

	return proofs, nil
}

// ParsePrimeProofsPage parses the prime proofs page and returns numbers used for the specified curve.
// Subpages are not read.
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

func (r *RawPrimeProofs) Translate() (*primality.Registry, error) {
	reg := primality.Registry{
		Proofs: make([]primality.Proof, 0, len(r.Numbers)),
	}
	for _, subpage := range r.Subpages {
		proof, err := subpage.Translate()
		if err != nil {
			return nil, err
		}
		reg.Proofs = append(reg.Proofs, *proof)
	}
	return &reg, nil
}

func ParseRawProofPage(s string) (*RawProofPage, error) {
	if twoProof.MatchString(s) {
		return &RawProofPage{
			N: "2",
		}, nil
	}
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
	if r.N == "2" {
		p.N = (*primality.BigInt)(big.NewInt(2))
		return &p, nil
	}
	n, ok := new(big.Int).SetString(r.N, 10)
	if !ok {
		return nil, ErrNotANumber
	}
	p.N = (*primality.BigInt)(n)
	a, err := ParseExpr(r.AExpr)
	if err != nil {
		return nil, err
	}
	var innerProof primality.GeneralizedPocklingtonProof
	innerProof.A = a
	base, ok := new(big.Int).SetString(r.B, 10)
	if !ok {
		return nil, ErrNotANumber
	}
	innerProof.Base = (*primality.BigInt)(base)
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
	innerProof.Inverses = inverses
	p.Proof = &innerProof
	return &p, nil
}
