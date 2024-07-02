// Code generated by genast.go; DO NOT EDIT.
package ast

import (
	"fmt"

	"github.com/mkeesey/craftinginterpreters/pkg/token"
)

type ExprVisitor[T any] interface {
	VisitAssign(*Assign) T
	VisitBinary(*Binary) T
	VisitGrouping(*Grouping) T
	VisitLiteral(*Literal) T
	VisitLogical(*Logical) T
	VisitUnary(*Unary) T
	VisitExprVar(*ExprVar) T
}

func VisitExpr[T any](expr Expr, visitor ExprVisitor[T]) T {
	switch n := expr.(type) {
	case *Assign:
		return visitor.VisitAssign(n)
	case *Binary:
		return visitor.VisitBinary(n)
	case *Grouping:
		return visitor.VisitGrouping(n)
	case *Literal:
		return visitor.VisitLiteral(n)
	case *Logical:
		return visitor.VisitLogical(n)
	case *Unary:
		return visitor.VisitUnary(n)
	case *ExprVar:
		return visitor.VisitExprVar(n)
	default:
		panic(fmt.Sprintf("Unknown Expr type %T", expr))
	}
}

type Expr interface {
	expr()
}

type Assign struct {
	Name *token.Token
	Value Expr
}

func (b *Assign) expr() {}

type Binary struct {
	Left Expr
	Operator *token.Token
	Right Expr
}

func (b *Binary) expr() {}

type Grouping struct {
	Expression Expr
}

func (b *Grouping) expr() {}

type Literal struct {
	Value any
}

func (b *Literal) expr() {}

type Logical struct {
	Left Expr
	Operator *token.Token
	Right Expr
}

func (b *Logical) expr() {}

type Unary struct {
	Operator *token.Token
	Right Expr
}

func (b *Unary) expr() {}

type ExprVar struct {
	Name *token.Token
}

func (b *ExprVar) expr() {}


