package ast

//go:generate go run ../../bin/genast/genast.go .

type Callable interface {
	Call(interpreter *TreeWalkInterpreter, arguments []interface{}) interface{}
	Arity() int
}

type LoxCallable struct {
	declaration *Function
}

func NewLoxCallable(declaration *Function) *LoxCallable {
	return &LoxCallable{declaration: declaration}
}

func (l *LoxCallable) Call(interpreter *TreeWalkInterpreter, arguments []interface{}) interface{} {
	env := WithEnvironment(interpreter.env)

	for i, param := range l.declaration.Params {
		env.Define(param.Lexeme, arguments[i])
	}

	interpreter.executeBlock(l.declaration.Body, env)
	return nil
}

func (l *LoxCallable) Arity() int {
	return len(l.declaration.Params)
}

func (l *LoxCallable) String() string {
	return "<fn " + l.declaration.Name.Lexeme + ">"
}
