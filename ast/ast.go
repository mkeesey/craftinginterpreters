package ast

//go:generate go run ../bin/genast/genast.go .

// type PrintVisitor struct {
// 	builder strings.Builder
// }

// func (p *PrintVisitor) Print(e Expr) string {
// 	return e.Accept(p)
// }

// func (p *PrintVisitor) VisitBinary(e *Binary) Visitor {
// 	return p
// }

// func (p *PrintVisitor) VisitGrouping(e *Grouping) Visitor {
// 	return p
// }

// func (p *PrintVisitor) VisitLiteral(e *Literal) Visitor {
// 	return p
// }

// func (p *PrintVisitor) VisitUnary(e *Unary) Visitor {
// 	return p
// }

// func (p *PrintVisitor) parenthesize(name string, expr ...Expr) string {
// 	p.builder.Reset()

// 	p.builder.WriteString("(")
// 	p.builder.WriteString(name)
// 	for _, e := range expr {
// 		p.builder.WriteString(" ")
// 		p.builder.WriteString(e.Accept(p).(string))
// 	}
// 	p.builder.WriteString(")")

// 	return p.builder.String()
// }
