package ast

//go:generate go run ../../bin/genast/genast.go .

type Callable interface {
	Call(interpreter *TreeWalkInterpreter, arguments []interface{}) interface{}
	Arity() int
}

type LoxFunction struct {
	declaration  *Function
	closure      *Environment
	isIntializer bool
}

func NewLoxFunction(declaration *Function, closure *Environment, isInitializer bool) *LoxFunction {
	return &LoxFunction{declaration: declaration, closure: closure, isIntializer: isInitializer}
}

func (l *LoxFunction) Call(interpreter *TreeWalkInterpreter, arguments []interface{}) (ret interface{}) {
	env := WithEnvironment(l.closure)

	for i, param := range l.declaration.Params {
		env.Define(param.Lexeme, arguments[i])
	}

	defer func() {
		if r := recover(); r != nil {
			returnval, ok := r.(returnval)
			if !ok { // some other panic other than our return value
				panic(r)
			}

			if l.isIntializer { // initializers always return 'this'
				var err error
				ret, err = env.GetAt(0, "this")
				if err != nil {
					panic(err)
				}
			} else {
				ret = returnval.Value
			}
		}
	}()
	interpreter.executeBlock(l.declaration.Body, env)
	if l.isIntializer {
		this, err := l.closure.GetAt(0, "this")
		if err != nil {
			panic(err)
		}
		return this
	}
	return nil
}

func (l *LoxFunction) Arity() int {
	return len(l.declaration.Params)
}

func (l *LoxFunction) Bind(instance *LoxInstance) *LoxFunction {
	env := WithEnvironment(l.closure)
	env.Define("this", instance)
	return NewLoxFunction(l.declaration, env, l.isIntializer)
}

func (l *LoxFunction) String() string {
	return "<fn " + l.declaration.Name.Lexeme + ">"
}

type returnval struct {
	Value interface{}
}
