package eval

import (
	"errors"
	"fmt"
)

type OperatorFunc func(tokens ...Token) (Token, error)

var BaseOperators = Operators(OpAdd, OpLT, OpEq, OpAnd)

var (
	OpAdd = Operator("+", func(tokens ...Token) (Token, error) {
		if len(tokens) != 2 {
			return Token{}, errors.New("operator: + requires exactly 2 operands")
		}

		a, err := tokens[0].Int64()
		if err != nil {
			return Token{}, errors.New("operator: + operands must be int64")
		}

		b, err := tokens[1].Int64()
		if err != nil {
			return Token{}, errors.New("operator: + operands must be int64")
		}

		return Token{value: fmt.Sprintf("%v", a+b), kind: TokenVar}, nil
	})
	OpLT = Operator("<", func(tokens ...Token) (Token, error) {
		if len(tokens) != 2 {
			return Token{}, errors.New("operator: < requires exactly 2 operands")
		}

		if !tokens[0].Kind().IsNumeric() || !tokens[1].Kind().IsNumeric() {
			return Token{}, errors.New("operator: < requires numeric operands")
		}

		a, err := tokens[0].Int64()
		if err != nil {
			return Token{}, errors.New("operator: < operands must be int64")
		}

		b, err := tokens[1].Int64()
		if err != nil {
			return Token{}, errors.New("operator: < operands must be int64")
		}

		return Token{value: fmt.Sprintf("%v", a < b), kind: TokenVar}, nil
	})
	OpEq = Operator("==", func(tokens ...Token) (Token, error) {
		if len(tokens) != 2 {
			return Token{}, errors.New("operator: == requires exactly 2 operands")
		}

		a, b := tokens[0].String(), tokens[1].String()
		return Token{value: fmt.Sprintf("%v", a == b), kind: TokenVar}, nil
	})
	OpAnd = Operator("&&", func(tokens ...Token) (Token, error) {
		if len(tokens) != 2 {
			return Token{}, errors.New("operator: && requires exactly 2 operands")
		}

		a, err := tokens[0].Bool()
		if err != nil {
			return Token{}, errors.New("operator: && operands must be boolean")
		}

		b, err := tokens[1].Bool()
		if err != nil {
			return Token{}, errors.New("operator: && operands must be boolean")
		}

		return Token{value: fmt.Sprintf("%v", a && b), kind: TokenVar}, nil
	})
)

type operators map[string]OperatorFunc

func Operators(operators ...operator) operators {
	ret := map[string]OperatorFunc{}

	for _, operator := range operators {
		ret[operator[0].(string)] = operator[1].(OperatorFunc)
	}

	return ret
}

type operator [2]interface{}

func Operator(operator string, opFunc OperatorFunc) operator {
	return [2]interface{}{operator, opFunc}
}

func (o operators) Get(operator string) (OperatorFunc, bool) {
	f, ok := o[operator]
	return f, ok
}

func (o operators) Set(operator string, opFunc OperatorFunc) {
	o[operator] = opFunc
}
