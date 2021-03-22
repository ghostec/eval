package eval

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func Expression(expression string) *eval {
	return &eval{expression: expression}
}

func (eval *eval) Eval() (Token, error) {
	var tokens []Token
	var err error

	ops := BaseOperators

	parsed := Parse(eval.expression, eval.vars, ops)

	for t := Tokenize(parsed).Operators(ops); t.HasNext(); {
		token := t.Next()

		if token.Kind() == TokenFmt {
			// if inside (...) => subexpr
			if token.String() == "(" {
				open := 1
				subexpr := ""

				for ; ; t.HasNext() {
					token = t.Next()
					str := token.String()

					switch str {
					case "(":
						open += 1
					case ")":
						open -= 1
					}

					if open == 0 {
						break
					}

					subexpr += str
				}

				token, err = Expression(subexpr).Eval()
				if err != nil {
					return Token{}, err
				}
			} else {
				continue
			}
		}

		tokens = append(tokens, token)
	}

	expression := expr{tokens: tokens, operators: ops}
	return expression.Eval()
}

type eval struct {
	expression string
	vars       vars
}

func (eval *eval) Vars(vars vars) *eval {
	eval.vars = vars
	return eval
}

func Parse(expression string, vars vars, operators operators) (ret string) {
	for t := Tokenize(expression).Operators(operators); t.HasNext(); {
		token := t.Next()

		var str string
		switch token.Kind() {
		case TokenVar:
			value, ok := vars.Get(token.String())
			if !ok {
				panic("var not found")
			}
			str = fmt.Sprintf("%v", value)
		default:
			str = token.String()
		}

		ret += str
	}

	return
}

type tokenizer struct {
	src       string
	pos       int
	operators operators
	fmt       Set
}

func Tokenize(src string) *tokenizer {
	return &tokenizer{src: src, fmt: NewSet(' ', '(', ')')}
}

func (t *tokenizer) Operators(operators operators) *tokenizer {
	t.operators = operators
	return t
}

func (t *tokenizer) Next() (ret Token) {
	if t.pos >= len(t.src) {
		return Token{}
	}

	cur := rune(t.src[t.pos])
	if t.fmt.Contains(cur) {
		ret = Token{value: string(cur), kind: TokenFmt}
		t.pos += 1
		return
	}

	var str string
	var maybeOperator bool

ParseRunes:
	for i := t.pos; i < len(t.src); i++ {
		rn := rune(t.src[i])

		isAlphanumeric := unicode.IsDigit(rn) || unicode.IsLetter(rn) || rn == '.'
		isFmt := t.fmt.Contains(rn)

		switch {
		case isFmt, isAlphanumeric && maybeOperator:
			// could be parsing either operand or operator but
			// 1. found formatting token
			// OR
			// 2. found another operand/operator
			break ParseRunes
		case !isAlphanumeric && len(str) == 0:
			// just started parsing a token and found a potential operator
			maybeOperator = true
		case !isAlphanumeric && !maybeOperator:
			// could be parsing either operand or operator but
			// 1. was already parsing the token -- e.g: len(str) > 0
			// AND
			// 2. found a rune that couldn't be part of this token
			break ParseRunes
		}

		t.pos += 1
		str += string(rn)
	}

	kind := TokenVar
	if _, err := strconv.ParseInt(str, 10, 64); err == nil {
		kind = TokenInt
	}
	if _, ok := t.operators.Get(str); ok {
		kind = TokenOperator
	}

	return Token{value: str, kind: kind}
}

func (t *tokenizer) HasNext() bool {
	return t.pos < len(t.src)
}

type Token struct {
	value string
	kind  tokenKind
}

func (t Token) String() string {
	return t.value
}

func (t Token) Kind() tokenKind {
	return t.kind
}

func (t Token) Int64() (int64, error) {
	return strconv.ParseInt(t.value, 10, 64)
}

func (t Token) Bool() (bool, error) {
	switch t.value {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, errors.New("token: not boolean")
	}
}

type tokenKind int

const (
	TokenInvalid tokenKind = iota
	TokenOperator
	TokenFmt
	TokenInt
	TokenVar
)

func (k tokenKind) IsNumeric() bool {
	switch k {
	case TokenInt:
		return true
	default:
		return false
	}
}

type vars map[string]interface{}

func (v vars) Get(key string) (interface{}, bool) {
	value := interface{}(v)

	for _, part := range strings.Split(key, ".") {
		switch val := value.(type) {
		case vars:
			value = val[part]
		case map[string]interface{}:
			value = val[part]
		default:
			return nil, false
		}
	}

	return value, true
}

func (v vars) Set(key string, value interface{}) {
	v[key] = value
}

func Var(key string, value interface{}) varsVar {
	return varsVar{key, value}
}

type varsVar [2]interface{}

func (v varsVar) Key() string {
	return v[0].(string)
}

func (v varsVar) Value() interface{} {
	return v[1]
}

func Vars(vs ...varsVar) vars {
	ret := vars(map[string]interface{}{})

	for _, v := range vs {
		ret.Set(v.Key(), v.Value())
	}

	return ret
}
