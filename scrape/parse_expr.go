package scrape

import (
	"errors"
	"math/big"

	"github.com/koba-e964/crypto-primality-proof/primality"
)

// Grammar:
// <expr> ::= <pow> | <pow> " * " <expr>
// <pow> ::= <num> | <num> "^" <num>
// <num> = [0-9]+
// This grammar is LL(1) and can be parsed by a recursive descent parser.

func ParseExpr(s string) (*primality.FactoredInt, error) {
	_, result, err := parseExpr(s)
	return result, err
}

func parseExpr(s string) (int, *primality.FactoredInt, error) {
	// <expr> ::= <pow> | <pow> " * " <expr>
	i, value, pow, err := parsePow(s)
	if err != nil {
		return i, nil, err
	}
	if i+3 <= len(s) && s[i:i+3] == " * " {
		i += 3
		j, expr, err := parseExpr(s[i:])
		if err != nil {
			return i + j, nil, err
		}
		return i + j, &primality.FactoredInt{
			Int:           (*primality.BigInt)(value.Mul(value, (*big.Int)(expr.Int))),
			Factorization: append([]primality.FactorEntry{pow}, expr.Factorization...),
		}, nil
	}
	return i, &primality.FactoredInt{
		Int:           (*primality.BigInt)(value),
		Factorization: []primality.FactorEntry{pow},
	}, nil
}

func parsePow(s string) (int, *big.Int, primality.FactorEntry, error) {
	// <pow> ::= <num> | <num> "^" <num>
	i, value, err := parseNum(s)
	if err != nil {
		return i, nil, primality.FactorEntry{}, err
	}
	if i+1 < len(s) && s[i] == '^' {
		i++
		j, exp, err := parseNum(s[i:])
		if err != nil {
			return i + j, nil, primality.FactorEntry{}, err
		}
		if !exp.IsInt64() {
			return i + j, nil, primality.FactorEntry{}, errors.New("exponent is not an int64")
		}
		power := new(big.Int).Exp(value, exp, nil)
		return i + j, power, primality.FactorEntry{Prime: (*primality.BigInt)(value), Exponent: int(exp.Int64())}, nil
	}
	return i, value, primality.FactorEntry{Prime: (*primality.BigInt)(value), Exponent: 1}, nil
}

func parseNum(s string) (int, *big.Int, error) {
	// <num> = [0-9]+
	i := 0
	for i < len(s) && '0' <= s[i] && s[i] <= '9' {
		i++
	}
	if i == 0 {
		return 0, nil, ErrNotANumber
	}
	value, ok := new(big.Int).SetString(s[:i], 10)
	if !ok {
		return i, nil, ErrNotANumber
	}
	return i, value, nil
}
