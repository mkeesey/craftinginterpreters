package ast

import (
	"github.com/mkeesey/craftinginterpreters/pkg/failure"
	"github.com/mkeesey/craftinginterpreters/pkg/token"
)

type functionType int

const (
	funcTypeNone functionType = iota
	funcTypeFunction
	funcTypeInitializer
	funcTypeMethod
)

type classType int

const (
	classTypeNone classType = iota
	classTypeClass
	classTypeSubclass
)

type Resolver struct {
	reporter      *failure.Reporter
	interpreter   *TreeWalkInterpreter
	scopes        []map[string]bool
	currFuncType  functionType
	currClassType classType
}

func NewResolver(interpreter *TreeWalkInterpreter, reporter *failure.Reporter) *Resolver {
	return &Resolver{interpreter: interpreter,
		reporter:      reporter,
		currFuncType:  funcTypeNone,
		currClassType: classTypeNone,
	}
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

func (r *Resolver) VisitGet(get *Get) interface{} {
	r.resolveExpr(get.Object)
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

func (r *Resolver) VisitSet(set *Set) interface{} {
	r.resolveExpr(set.Value)
	r.resolveExpr(set.Object)
	return nil
}

func (r *Resolver) VisitSuper(super *Super) interface{} {
	if r.currClassType == classTypeNone {
		r.reporter.TokenError(super.Keyword, "Cannot use 'super' outside of a class.")
	} else if r.currClassType != classTypeSubclass {
		r.reporter.TokenError(super.Keyword, "Cannot use 'super' in a class with no superclass.")
	}

	r.resolveLocal(super, super.Keyword)
	return nil
}

func (r *Resolver) VisitThis(this *This) interface{} {
	if r.currClassType == classTypeNone {
		r.reporter.TokenError(this.Keyword, "Cannot use 'this' outside of a class.")
		return nil
	}
	r.resolveLocal(this, this.Keyword)
	return nil
}

func (r *Resolver) VisitUnary(unary *Unary) interface{} {
	r.resolveExpr(unary.Right)
	return nil
}

func (r *Resolver) VisitExprVar(expr *ExprVar) interface{} {
	scope, ok := r.peekScope()
	if ok {
		defined, declared := scope[expr.Name.Lexeme]
		if declared && !defined {
			r.reporter.Report(expr.Name.Line, expr.Name.Lexeme, "Cannot read local variable in its own initializer")
		}
	}

	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) VisitBlock(b *Block) {
	r.beginScope()
	r.Resolve(b.Statements)
	r.endScope()
}

func (r *Resolver) VisitClass(class *Class) {
	priorClassType := r.currClassType
	r.currClassType = classTypeClass
	defer func() {
		r.currClassType = priorClassType
	}()

	r.declare(class.Name)
	r.define(class.Name)
	if class.Superclass != nil {
		r.currClassType = classTypeSubclass
		if class.Name.Lexeme == class.Superclass.Name.Lexeme {
			r.reporter.TokenError(class.Name, "A class cannot inherit from itself.")
		}

		r.resolveExpr(class.Superclass)
		r.beginScope()
		defer r.endScope()
		scope, _ := r.peekScope()
		scope["super"] = true
	}

	r.beginScope()
	defer r.endScope()
	scope, _ := r.peekScope()
	scope["this"] = true

	for _, method := range class.Methods {
		declaration := funcTypeMethod
		if method.Name.Lexeme == "init" {
			declaration = funcTypeInitializer
		}
		r.resolveFunction(method, declaration)
	}
}

func (r *Resolver) VisitExpression(exp *Expression) {
	r.resolveExpr(exp.Expression)
}

func (r *Resolver) VisitFunction(fun *Function) {
	r.declare(fun.Name)
	r.define(fun.Name)
	r.resolveFunction(fun, funcTypeFunction)
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
	if r.currFuncType == funcTypeNone {
		r.reporter.TokenError(ret.Keyword, "Can't return from top-level code.")
	}
	if ret.Value != nil {
		if r.currFuncType == funcTypeInitializer {
			r.reporter.TokenError(ret.Keyword, "Can't return a value from an initializer.")
		}

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

func (r *Resolver) resolveFunction(fun *Function, funcType functionType) {
	enclosingType := r.currFuncType
	r.currFuncType = funcType
	r.beginScope()
	for _, param := range fun.Params {
		r.declare(param)
		r.define(param)
	}

	r.Resolve(fun.Body)
	r.endScope()
	r.currFuncType = enclosingType
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
	if _, alreadyDeclared := scope[name.Lexeme]; alreadyDeclared {
		r.reporter.Report(name.Line, name.Lexeme, "Variable with this name already declared in this scope")
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
