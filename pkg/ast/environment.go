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

func (e *Environment) AssignAt(distance int, tok *token.Token, value interface{}) error {
	env := e.ancestor(distance)
	env.values[tok.Lexeme] = value
	return nil
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

func (e *Environment) GetAt(distance int, name string) (interface{}, error) {
	env := e.ancestor(distance)
	val, ok := env.values[name]
	if ok {
		return val, nil
	}

	panic(fmt.Sprintf("Undefined variable '%s' which was supposed to be a defined local.", name))
}

func (e *Environment) ancestor(distance int) *Environment {
	env := e
	for i := 0; i < distance; i++ {
		env = env.enclosing
	}
	return env
}
