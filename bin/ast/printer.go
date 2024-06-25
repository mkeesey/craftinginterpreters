package main

import (
	"fmt"

	"github.com/mkeesey/craftinginterpreters/ast"
	"github.com/mkeesey/craftinginterpreters/token"
)

func main() {
	expr := &ast.Binary{
		Left: &ast.Unary{
			Operator: token.NewToken(token.MINUS, "-", nil, 1),
			Right:    &ast.Literal{Value: 123},
		},
		Operator: token.NewToken(token.STAR, "*", nil, 1),
		Right: &ast.Grouping{
			Expression: &ast.Literal{Value: 45.67},
		},
	}

	printer := ast.PrintVisitor{}
	fmt.Println(printer.Print(expr))
}
