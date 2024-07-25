// Code generated by genast.go; DO NOT EDIT.
package ast

import (
	"fmt"

	"github.com/mkeesey/craftinginterpreters/pkg/token"
)

type StmtVisitor interface {
	VisitBlock(*Block)
	VisitClass(*Class)
	VisitExpression(*Expression)
	VisitFunction(*Function)
	VisitIf(*If)
	VisitPrint(*Print)
	VisitReturn(*Return)
	VisitStmtVar(*StmtVar)
	VisitWhile(*While)
}

func VisitStmt(stmt Stmt, visitor StmtVisitor) {
	switch n := stmt.(type) {
	case *Block:
		visitor.VisitBlock(n)
	case *Class:
		visitor.VisitClass(n)
	case *Expression:
		visitor.VisitExpression(n)
	case *Function:
		visitor.VisitFunction(n)
	case *If:
		visitor.VisitIf(n)
	case *Print:
		visitor.VisitPrint(n)
	case *Return:
		visitor.VisitReturn(n)
	case *StmtVar:
		visitor.VisitStmtVar(n)
	case *While:
		visitor.VisitWhile(n)
	default:
		panic(fmt.Sprintf("Unknown Stmt type %T", stmt))
	}
}

type Stmt interface {
	stmt()
}

type Block struct {
	Statements []Stmt
}

func (b *Block) stmt() {}

type Class struct {
	Name *token.Token
	Methods []*Function
	Superclass *ExprVar
}

func (b *Class) stmt() {}

type Expression struct {
	Expression Expr
}

func (b *Expression) stmt() {}

type Function struct {
	Name *token.Token
	Params []*token.Token
	Body []Stmt
}

func (b *Function) stmt() {}

type If struct {
	Condition Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (b *If) stmt() {}

type Print struct {
	Expression Expr
}

func (b *Print) stmt() {}

type Return struct {
	Keyword *token.Token
	Value Expr
}

func (b *Return) stmt() {}

type StmtVar struct {
	Name *token.Token
	Initializer Expr
}

func (b *StmtVar) stmt() {}

type While struct {
	Condition Expr
	Body Stmt
}

func (b *While) stmt() {}


