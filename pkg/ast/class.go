package ast

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
	class *LoxClass
}

func NewLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{class: class}
}

func (l *LoxInstance) String() string {
	return "<instance of " + l.class.name + ">"
}
