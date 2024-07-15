package ast

import (
	"github.com/mkeesey/craftinginterpreters/pkg/failure"
	"github.com/mkeesey/craftinginterpreters/pkg/token"
)

type Resolver struct {
	reporter    *failure.Reporter
	interpreter *TreeWalkInterpreter
	scopes      []map[string]bool
}

func NewResolver(interpreter *TreeWalkInterpreter, reporter *failure.Reporter) *Resolver {
	return &Resolver{interpreter: interpreter, reporter: reporter}
}

func (r *Resolver) VisitAssign(a *Assign) interface{} {
	r.resolveExpr(a.Value)
	r.resolveLocal(a, a.Name)
	return nil
}

func (r *Resolver) VisitBinary(bin *Binary) interface{} {
	r.resolveExpr(bin.Left)
	r.resolveExpr(bin.Right)
	return nil
}

func (r *Resolver) VisitCall(call *Call) interface{} {
	r.resolveExpr(call.Callee)
	for _, arg := range call.Arguments {
		r.resolveExpr(arg)
	}
	return nil
}

func (r *Resolver) VisitGrouping(g *Grouping) interface{} {
	r.resolveExpr(g.Expression)
	return nil
}

func (r *Resolver) VisitLiteral(lit *Literal) interface{} {
	return nil
}

func (r *Resolver) VisitLogical(logical *Logical) interface{} {
	r.resolveExpr(logical.Left)
	r.resolveExpr(logical.Right)
	return nil
}

func (r *Resolver) VisitUnary(unary *Unary) interface{} {
	r.resolveExpr(unary.Right)
	return nil
}

func (r *Resolver) VisitExprVar(expr *ExprVar) interface{} {
	scope, ok := r.peekScope()
	if ok && scope[expr.Name.Lexeme] == false {
		r.reporter.Error(expr.Name.Line, "Cannot read local variable in its own initializer")
	}

	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) VisitBlock(b *Block) {
	r.beginScope()
	r.Resolve(b.Statements)
	r.endScope()
}

func (r *Resolver) VisitExpression(exp *Expression) {
	r.resolveExpr(exp.Expression)
}

func (r *Resolver) VisitFunction(fun *Function) {
	r.declare(fun.Name)
	r.define(fun.Name)
	r.resolveFunction(fun)
}

func (r *Resolver) VisitIf(i *If) {
	r.resolveExpr(i.Condition)
	r.resolveStmt(i.ThenBranch)
	if i.ElseBranch != nil {
		r.resolveStmt(i.ElseBranch)
	}
}

func (r *Resolver) VisitPrint(p *Print) {
	r.resolveExpr(p.Expression)
}

func (r *Resolver) VisitReturn(ret *Return) {
	if ret.Value != nil {
		r.resolveExpr(ret.Value)
	}
}

func (r *Resolver) VisitStmtVar(s *StmtVar) {
	r.declare(s.Name)
	if s.Initializer != nil {
		r.resolveExpr(s.Initializer)
	}
	r.define(s.Name)
}

func (r *Resolver) VisitWhile(while *While) {
	r.resolveExpr(while.Condition)
	r.resolveStmt(while.Body)
}

func (r *Resolver) Resolve(statements []Stmt) {
	for _, stmt := range statements {
		r.resolveStmt(stmt)
	}
}

func (r *Resolver) resolveStmt(stmt Stmt) {
	VisitStmt(stmt, r)
}

func (r *Resolver) resolveExpr(expr Expr) {
	VisitExpr(expr, r)
}

func (r *Resolver) resolveFunction(fun *Function) {
	r.beginScope()
	for _, param := range fun.Params {
		r.declare(param)
		r.define(param)
	}

	r.Resolve(fun.Body)
	r.endScope()
}

func (r *Resolver) resolveLocal(expr Expr, name *token.Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.Lexeme]; ok {
			r.interpreter.resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) beginScope() {
	//TODO reuse map instead?
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	lastIdx := len(r.scopes) - 1
	r.scopes[lastIdx] = nil
	r.scopes = r.scopes[:lastIdx]
}

func (r *Resolver) declare(name *token.Token) {
	scope, ok := r.peekScope()
	if !ok {
		return
	}
	scope[name.Lexeme] = false
}

func (r *Resolver) define(name *token.Token) {
	scope, ok := r.peekScope()
	if !ok {
		return
	}
	scope[name.Lexeme] = true
}

func (r *Resolver) peekScope() (map[string]bool, bool) {
	if len(r.scopes) == 0 {
		return nil, false
	}
	return r.scopes[len(r.scopes)-1], true
}
