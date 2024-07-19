package ast

import "github.com/mkeesey/craftinginterpreters/pkg/token"

type LoxClass struct {
	name    string
	methods map[string]*LoxCallable
}

func NewLoxClass(name string, methods map[string]*LoxCallable) *LoxClass {
	return &LoxClass{name: name, methods: methods}
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

	if method := l.findMethod(name.Lexeme); method != nil {
		return method
	}

	panic("Undefined property '" + name.Lexeme + "'.")
}

func (l *LoxInstance) Set(name *token.Token, value interface{}) {
	l.fields[name.Lexeme] = value
}

func (l *LoxInstance) String() string {
	return "<instance of " + l.class.name + ">"
}

func (l *LoxInstance) findMethod(name string) *LoxCallable {
	if method, ok := l.class.methods[name]; ok {
		return method
	}

	return nil
}
