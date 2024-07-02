package ast

import (
	"fmt"

	"github.com/mkeesey/craftinginterpreters/pkg/token"
)

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

func (e *Environment) Assign(tok *token.Token, value interface{}) error {
	_, ok := e.values[tok.Lexeme]
	if !ok {
		return fmt.Errorf("Undefined variable '%s'.", tok.Lexeme)
	}
	e.values[tok.Lexeme] = value
	return nil
}

func (e *Environment) Get(name string) (interface{}, error) {
	val, ok := e.values[name]
	if !ok {
		return nil, fmt.Errorf("Undefined variable '%s'.", name)
	}
	return val, nil
}
