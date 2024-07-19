package ast

import (
	"fmt"

	"github.com/mkeesey/craftinginterpreters/pkg/failure"
	"github.com/mkeesey/craftinginterpreters/pkg/token"
)

// ExprVisitor[interface{}]
// StmtVisitor
type TreeWalkInterpreter struct {
	env       *Environment
	globalEnv *Environment
	locals    map[Expr]int
	reporter  *failure.Reporter
}

func NewInterpreter(reporter *failure.Reporter) *TreeWalkInterpreter {
	globalEnv := NewEnvironment()

	globalEnv.Define("clock", &TimeCallable{})

	return &TreeWalkInterpreter{
		globalEnv: globalEnv,
		env:       globalEnv,
		locals:    make(map[Expr]int),
		reporter:  reporter,
	}
}

func (p *TreeWalkInterpreter) Interpret(statements []Stmt) {
	defer func() {
		if r := recover(); r != nil {
			p.reporter.RuntimeError(r)
		}
	}()
	for _, stmt := range statements {
		VisitStmt(stmt, p)
	}
}

func (p *TreeWalkInterpreter) VisitAssign(e *Assign) interface{} {
	value := p.evaluate(e.Value)

	if distance, ok := p.locals[e]; ok {
		err := p.env.AssignAt(distance, e.Name, value)
		if err != nil {
			panic(err)
		}
	} else {
		err := p.globalEnv.Assign(e.Name, value)
		if err != nil {
			panic(err)
		}
	}

	return value
}

func (p *TreeWalkInterpreter) VisitBinary(e *Binary) interface{} {
	left := p.evaluate(e.Left)
	right := p.evaluate(e.Right)

	switch e.Operator.Type {
	case token.PLUS:
		leftFloat, isLeftFloat := left.(float64)
		rightFloat, isRightFloat := right.(float64)
		if isLeftFloat && isRightFloat {
			return leftFloat + rightFloat
		}

		leftStr, isLeftStr := left.(string)
		rightStr, isRightStr := right.(string)
		if isLeftStr && isRightStr {
			return leftStr + rightStr
		}

		panic(fmt.Sprintf("Cannot add operands %v %v", left, right))
	case token.MINUS:
		return requireFloat64(left) - requireFloat64(right)
	case token.SLASH:
		return requireFloat64(left) / requireFloat64(right)
	case token.STAR:
		return requireFloat64(left) * requireFloat64(right)
	case token.GREATER:
		return requireFloat64(left) > requireFloat64(right)
	case token.GREATER_EQUAL:
		return requireFloat64(left) >= requireFloat64(right)
	case token.LESS:
		return requireFloat64(left) < requireFloat64(right)
	case token.LESS_EQUAL:
		return requireFloat64(left) <= requireFloat64(right)
	case token.EQUAL_EQUAL:
		// TODO - ensure this matches lox requirements
		return left == right
	case token.BANG_EQUAL:
		return left != right
	}

	panic(fmt.Sprintf("unknown operator type %s", e.Operator.Type))
}

func (p *TreeWalkInterpreter) VisitCall(e *Call) interface{} {
	callee := p.evaluate(e.Callee)

	args := make([]interface{}, 0)
	for _, arg := range e.Arguments {
		args = append(args, p.evaluate(arg))
	}

	function, ok := callee.(Callable)
	if !ok {
		panic("Can only call functions and classes")
	}
	if len(args) != function.Arity() {
		panic(fmt.Sprintf("Expected %d arguments but got %d", function.Arity(), len(args)))
	}

	ret := function.Call(p, args)

	return ret
}

func (p *TreeWalkInterpreter) VisitGet(e *Get) interface{} {
	obj := p.evaluate(e.Object)
	if instance, ok := obj.(*LoxInstance); ok {
		return instance.Get(e.Name)
	}
	panic("Only instances have properties")
}

func (p *TreeWalkInterpreter) VisitGrouping(e *Grouping) interface{} {
	return p.evaluate(e.Expression)
}

func (p *TreeWalkInterpreter) VisitLiteral(e *Literal) interface{} {
	return e.Value
}

func (p *TreeWalkInterpreter) VisitLogical(e *Logical) interface{} {
	left := p.evaluate(e.Left)

	if e.Operator.Type == token.OR {
		if isTruthy(left) {
			return left
		}
	} else {
		if !isTruthy(left) {
			return left
		}
	}

	return p.evaluate(e.Right)
}

func (p *TreeWalkInterpreter) VisitSet(e *Set) interface{} {
	obj := p.evaluate(e.Object)
	if instance, ok := obj.(*LoxInstance); ok {
		value := p.evaluate(e.Value)
		instance.Set(e.Name, value)
		return value
	}
	panic("Only instances have fields")
}

func (p *TreeWalkInterpreter) VisitUnary(e *Unary) interface{} {
	right := p.evaluate(e.Right)

	switch e.Operator.Type {
	case token.BANG:
		val := isTruthy(right)
		return !val
	case token.MINUS:
		val := requireFloat64(right)
		return -val
	}

	panic(fmt.Sprintf("unknown operator type %s", e.Operator.Type))
}

func (p *TreeWalkInterpreter) VisitExprVar(e *ExprVar) interface{} {
	val, err := p.lookupVariable(e.Name, e)
	if err != nil {
		panic(err)
	}
	return val
}

func (p *TreeWalkInterpreter) lookupVariable(name *token.Token, expr Expr) (interface{}, error) {
	distance, ok := p.locals[expr]
	if ok {
		return p.env.GetAt(distance, name.Lexeme)
	}
	return p.globalEnv.Get(name.Lexeme)
}

func (p *TreeWalkInterpreter) VisitBlock(e *Block) {
	env := WithEnvironment(p.env)
	p.executeBlock(e.Statements, env)
}

func (p *TreeWalkInterpreter) executeBlock(stmts []Stmt, env *Environment) {
	previous := p.env
	defer func() {
		p.env = previous
	}()

	p.env = env

	for _, stmt := range stmts {
		VisitStmt(stmt, p)
	}
}

func (p *TreeWalkInterpreter) VisitClass(class *Class) {
	p.env.Define(class.Name.Lexeme, nil)
	loxClass := NewLoxClass(class.Name.Lexeme)
	p.env.Assign(class.Name, loxClass)
}

func (p *TreeWalkInterpreter) resolve(expr Expr, depth int) {
	p.locals[expr] = depth
}

func (p *TreeWalkInterpreter) VisitExpression(e *Expression) {
	p.evaluate(e.Expression)
}

func (p *TreeWalkInterpreter) VisitIf(e *If) {
	condition := p.evaluate(e.Condition)
	if isTruthy(condition) {
		VisitStmt(e.ThenBranch, p)
	} else if e.ElseBranch != nil {
		VisitStmt(e.ElseBranch, p)
	}
}

func (p *TreeWalkInterpreter) VisitPrint(e *Print) {
	val := p.evaluate(e.Expression)
	fmt.Println(val)
}

func (p *TreeWalkInterpreter) VisitReturn(e *Return) {
	var value interface{}
	if e.Value != nil {
		value = p.evaluate(e.Value)
	}
	panic(returnval{value})
}

func (p *TreeWalkInterpreter) VisitStmtVar(e *StmtVar) {
	var value interface{}
	if e.Initializer != nil {
		value = p.evaluate(e.Initializer)
	}

	p.env.Define(e.Name.Lexeme, value)
}

func (p *TreeWalkInterpreter) VisitFunction(e *Function) {
	function := NewLoxCallable(e, p.env)
	p.env.Define(e.Name.Lexeme, function)
}

func (p *TreeWalkInterpreter) VisitWhile(e *While) {
	for isTruthy(p.evaluate(e.Condition)) {
		VisitStmt(e.Body, p)
	}
}

func (p *TreeWalkInterpreter) evaluate(e Expr) interface{} {
	return VisitExpr(e, p)
}

func isTruthy(value interface{}) bool {
	if value == nil {
		return false
	}
	casted, ok := value.(bool)
	if !ok {
		return true // not a bool, but considered truthy
	}
	return casted
}

func requireFloat64(value interface{}) float64 {
	val, ok := value.(float64)
	if !ok {
		panic(fmt.Sprintf("cannot cast %v to double", value))
	}
	return val
}
