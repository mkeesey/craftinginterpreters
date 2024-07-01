package ast

import "fmt"

type Environment struct {
	values map[string]interface{}
}

func NewEnvironment() *Environment {
	return &Environment{
		values: map[string]interface{}{},
	}
}

func (e *Environment) Define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Environment) Get(name string) (interface{}, error) {
	val, ok := e.values[name]
	if !ok {
		return nil, fmt.Errorf("Undefined variable '%s'.", name)
	}
	return val, nil
}
