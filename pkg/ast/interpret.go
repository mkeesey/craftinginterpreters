package ast

import (
	"fmt"

	"github.com/mkeesey/craftinginterpreters/pkg/token"
)

// Visitor[interface{}]
type TreeWalkInterpreter struct {
}

func NewInterpreter() *TreeWalkInterpreter {
	return &TreeWalkInterpreter{}
}

func (p *TreeWalkInterpreter) Interpret(e Expr) (ret interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Error: %v", r)
		}
	}()
	ret = p.evaluate(e)
	return ret, nil
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

func (p *TreeWalkInterpreter) VisitGrouping(e *Grouping) interface{} {
	return p.evaluate(e.Expression)
}

func (p *TreeWalkInterpreter) VisitLiteral(e *Literal) interface{} {
	return e.Value
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

func (p *TreeWalkInterpreter) evaluate(e Expr) interface{} {
	return Visit(e, p)
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
