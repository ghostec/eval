package main

import (
	"github.com/ghostec/eval"
)

func main() {
	vars := eval.Vars(
		eval.Var("player", eval.Vars(
			eval.Var("level", 3),
			eval.Var("class", "mage"),
		)),
	)

	result, err := eval.Expression("player.level < 4").Vars(vars).Eval()
	if err != nil {
		panic(err)
	}

	b, err := result.Bool()
	if err != nil {
		panic(err)
	}

	println("result", b)
}
