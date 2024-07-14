package ast

//go:generate go run ../../bin/genast/genast.go .

type Callable interface {
	Call(interpreter *TreeWalkInterpreter, arguments []interface{}) interface{}
	Arity() int
}

type LoxCallable struct {
	declaration *Function
	closure     *Environment
}

func NewLoxCallable(declaration *Function, closure *Environment) *LoxCallable {
	return &LoxCallable{declaration: declaration, closure: closure}
}

func (l *LoxCallable) Call(interpreter *TreeWalkInterpreter, arguments []interface{}) (ret interface{}) {
	env := WithEnvironment(l.closure)

	for i, param := range l.declaration.Params {
		env.Define(param.Lexeme, arguments[i])
	}

	defer func() {
		if r := recover(); r != nil {
			returnval, ok := r.(returnval)
			if !ok {
				panic(r)
			}
			ret = returnval.Value
		}
	}()
	interpreter.executeBlock(l.declaration.Body, env)
	return nil
}

func (l *LoxCallable) Arity() int {
	return len(l.declaration.Params)
}

func (l *LoxCallable) String() string {
	return "<fn " + l.declaration.Name.Lexeme + ">"
}

type returnval struct {
	Value interface{}
}
