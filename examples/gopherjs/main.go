package main

import (
	"github.com/ghostec/eval"
	"github.com/goplusjs/gopherjs/js"
)

func main() {
	eli := Eligibility{}
	js.Global.Set("eligibility", eli.Eval)
}

type Eligibility struct{}

func (e Eligibility) Eval(expression string, vars map[string]interface{}) bool {
	result, err := eval.Expression(expression).Vars(vars).Eval()
	if err != nil {
		panic(err)
	}

	b, err := result.Bool()
	if err != nil {
		panic(err)
	}

	return b
}
