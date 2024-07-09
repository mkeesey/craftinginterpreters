package ast

//go:generate go run ../../bin/genast/genast.go .

type Callable interface {
	Call(interpreter *TreeWalkInterpreter, arguments []interface{}) interface{}
	Arity() int
}
