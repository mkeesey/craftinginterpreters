package ast

import (
	"github.com/mkeesey/craftinginterpreters/pkg/failure"
	"github.com/mkeesey/craftinginterpreters/pkg/token"
)

type LoxClass struct {
	name       string
	superclass *LoxClass
	methods    map[string]*LoxFunction
}

func NewLoxClass(name string, superclass *LoxClass, methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{name: name, superclass: superclass, methods: methods}
}

func (l *LoxClass) Call(interpreter *TreeWalkInterpreter, arguments []interface{}) interface{} {
	instance := NewLoxInstance(l)
	initializer := l.findMethod("init")
	if initializer != nil {
		initializer.Bind(instance).Call(interpreter, arguments)
	}

	return instance
}

func (l *LoxClass) Arity() int {
	initializer := l.findMethod("init")
	if initializer == nil {
		return 0
	}

	return initializer.Arity()
}

func (l *LoxClass) String() string {
	return l.name
}

func (l *LoxClass) findMethod(name string) *LoxFunction {
	if method, ok := l.methods[name]; ok {
		return method
	}

	if l.superclass != nil {
		return l.superclass.findMethod(name)
	}

	return nil
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

	if method := l.class.findMethod(name.Lexeme); method != nil {
		return method.Bind(l)
	}

	panic(failure.RuntimeError{Token: name, Message: "Undefined property '" + name.Lexeme + "'."})
}

func (l *LoxInstance) Set(name *token.Token, value interface{}) {
	l.fields[name.Lexeme] = value
}

func (l *LoxInstance) String() string {
	return l.class.name + " instance"
}
