package ast

import (
	"time"
)

type TimeCallable struct {
}

func (t *TimeCallable) Call(interpreter *TreeWalkInterpreter, arguments []interface{}) interface{} {
	return float64(time.Now().Unix())
}

func (t *TimeCallable) Arity() int {
	return 0
}

func (t *TimeCallable) String() string {
	return "<native fn>"
}
