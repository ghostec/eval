package eval

import (
	"errors"
	"fmt"
)

type expr struct {
	tokens    []Token
	operators operators
}

func (expr *expr) Eval() (Token, error) {
	var opFunc OperatorFunc
	var tokens []Token

	switch len(expr.tokens) {
	case 0:
		return Token{}, errors.New("expr: can't eval zero tokens")
	case 1:
		return expr.tokens[0], nil
	}

	for _, token := range expr.tokens {
		switch token.Kind() {
		case TokenOperator:
			f, ok := expr.operators.Get(token.String())
			if !ok {
				return Token{}, fmt.Errorf("expr: operator %s not found", token.String())
			}

			opFunc = f
		case TokenInvalid:
			return Token{}, errors.New("expr: invalid token")
		default:
			tokens = append(tokens, token)
		}
	}

	return opFunc(tokens...)
}
