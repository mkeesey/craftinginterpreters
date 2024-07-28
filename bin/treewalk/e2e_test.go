package main

import (
	"strings"
	"testing"

	"github.com/mkeesey/craftinginterpreters/pkg/failure"
	"github.com/mkeesey/craftinginterpreters/pkg/parser"
	"github.com/mkeesey/craftinginterpreters/pkg/scanner"
	"github.com/stretchr/testify/require"
)

func TestScanParse(t *testing.T) {
	type testcase struct {
		input    string
		expected string
	}

	testcases := []testcase{
		// {
		// 	input:    `"waffles" == "tacos"`,
		// 	expected: `(== waffles tacos)`,
		// },
		// {
		// 	input:    `2342 + 23423 * 23`,
		// 	expected: `(+ 2342 (* 23423 23))`,
		// },
		// {
		// 	input:    `(2342 + 23423) * 23`,
		// 	expected: `(* (group (+ 2342 23423)) 23)`,
		// },
	}

	for _, testcase := range testcases {
		t.Run(testcase.input, func(t *testing.T) {
			reporter := &failure.Reporter{}
			scan := scanner.NewScanner(strings.NewReader(testcase.input), reporter)
			tokens := scan.ScanTokens()
			require.False(t, reporter.HasFailed())

			parser := parser.NewParser(tokens, reporter)
			expr, err := parser.Parse()
			require.NoError(t, err)
			require.False(t, reporter.HasFailed())

			expr = expr //TODO
			// visitor := ast.PrintVisitor{}
			// output := visitor.Print(expr)
			// require.Equal(t, testcase.expected, output)
		})
	}
}
