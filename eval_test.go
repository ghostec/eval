package eval_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/ghostec/eval"
)

func TestAll(t *testing.T) {
	testCases := []struct {
		expression string
		result     string
	}{
		{"3 < 4", "true"},
		{"5 < 4", "false"},
		{"(3 < 4)", "true"},
		{"(3 < 4) && (5 < 4)", "false"},
		{"(2 < 4) && (3 < 4)", "true"},
		{"((2 < 4) && (1 < 2)) && (3 < 4)", "true"},
		{"3 == 4", "false"},
		{"25 == 25", "true"},
		{"25+137", strconv.Itoa(25 + 137)},
		{fmt.Sprintf("(25+ 137) == %d", 25+137), "true"},
		// TODO: precedence, e.g: 25 + 137 == 162 should work without parens
	}

	for _, tc := range testCases {
		t.Run(tc.expression, func(t *testing.T) {
			result, err := eval.Expression(tc.expression).Eval()
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}
			if result.String() != tc.result {
				t.Errorf(`unexpected result for expression "%s". Got: %s, Expected: %s`, tc.expression, result, tc.result)
			}
		})
	}
}
