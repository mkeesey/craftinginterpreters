package ast

import (
	"fmt"

	"github.com/mkeesey/craftinginterpreters/pkg/token"
)

type Environment struct {
	enclosing *Environment
	values    map[string]interface{}
}

func NewEnvironment() *Environment {
	return &Environment{
		enclosing: nil,
		values:    map[string]interface{}{},
	}
}

func WithEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		values:    map[string]interface{}{},
	}
}

func (e *Environment) Define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Environment) Assign(tok *token.Token, value interface{}) error {
	_, ok := e.values[tok.Lexeme]
	if ok {
		e.values[tok.Lexeme] = value
		return nil
	}

	if e.enclosing != nil {
		return e.enclosing.Assign(tok, value)
	}

	return fmt.Errorf("Undefined variable '%s'.", tok.Lexeme)
}

func (e *Environment) Get(name string) (interface{}, error) {
	val, ok := e.values[name]
	if ok {
		return val, nil
	}

	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}

	return nil, fmt.Errorf("Undefined variable '%s'.", name)
}
