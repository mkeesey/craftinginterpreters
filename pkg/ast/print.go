package ast

import (
	"fmt"
	"strings"
)

type PrintVisitor struct {
}

func (p *PrintVisitor) Print(e Expr) string {
	return VisitExpr(e, p)
}

func (p *PrintVisitor) VisitBinary(e *Binary) string {
	return p.parenthesize(e.Operator.Lexeme, e.Left, e.Right)
}

func (p *PrintVisitor) VisitGrouping(e *Grouping) string {
	return p.parenthesize("group", e.Expression)
}

func (p *PrintVisitor) VisitLiteral(e *Literal) string {
	if e.Value == nil {
		return "nil"
	}
	return fmt.Sprint(e.Value)
}

func (p *PrintVisitor) VisitUnary(e *Unary) string {
	return p.parenthesize(e.Operator.Lexeme, e.Right)
}

func (p *PrintVisitor) parenthesize(name string, expr ...Expr) string {
	builder := strings.Builder{}

	builder.WriteString("(")
	builder.WriteString(name)
	for _, e := range expr {
		builder.WriteString(" ")
		builder.WriteString(VisitExpr(e, p))
	}
	builder.WriteString(")")

	return builder.String()
}
