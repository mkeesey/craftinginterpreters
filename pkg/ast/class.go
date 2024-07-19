package ast

import "github.com/mkeesey/craftinginterpreters/pkg/token"

type LoxClass struct {
	name string
}

func NewLoxClass(name string) *LoxClass {
	return &LoxClass{name: name}
}

func (l *LoxClass) Call(interpreter *TreeWalkInterpreter, arguments []interface{}) interface{} {
	instance := NewLoxInstance(l)
	return instance
}

func (l *LoxClass) Arity() int {
	return 0
}

func (l *LoxClass) String() string {
	return "<class " + l.name + ">"
}

type LoxInstance struct {
	class  *LoxClass
	fields map[string]interface{}
}

func NewLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{class: class, fields: make(map[string]interface{})}
}

func (l *LoxInstance) Get(name *token.Token) interface{} {
	if value, ok := l.fields[name.Lexeme]; ok {
		return value
	}

	panic("Undefined property '" + name.Lexeme + "'.")
}

func (l *LoxInstance) Set(name *token.Token, value interface{}) {
	l.fields[name.Lexeme] = value
}

func (l *LoxInstance) String() string {
	return "<instance of " + l.class.name + ">"
}
