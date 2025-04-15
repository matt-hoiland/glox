package environment

import (
	ierrors "github.com/matt-hoiland/glox/internal/errors"
	"github.com/matt-hoiland/glox/internal/loxtype"
	"github.com/matt-hoiland/glox/internal/token"
)

type Environment struct {
	enclosing *Environment
	values    map[string]loxtype.Type
}

func New() *Environment {
	return &Environment{
		values: map[string]loxtype.Type{},
	}
}

func (e *Environment) ancestor(distance int) *Environment {
	env := e
	for range distance {
		env = env.enclosing
	}
	return env
}

func (e *Environment) Assign(name *token.Token, value loxtype.Type) error {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return nil
	}
	if e.enclosing != nil {
		return e.enclosing.Assign(name, value)
	}
	return ierrors.New(name, newUndefinedVariableError(name))
}

func (e *Environment) AssignAt(distance int, name *token.Token, value loxtype.Type) error {
	return e.ancestor(distance).Assign(name, value)
}

func (e *Environment) Define(name *token.Token, value loxtype.Type) {
	e.values[name.Lexeme] = value
}

func (e *Environment) Get(name *token.Token) (loxtype.Type, error) {
	value, ok := e.values[name.Lexeme]
	if ok {
		return value, nil
	}
	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}
	return nil, ierrors.New(name, newUndefinedVariableError(name))
}

func (e *Environment) GetAt(name *token.Token, distance int) (loxtype.Type, error) {
	return e.ancestor(distance).Get(name)
}

func (e *Environment) MakeChild() *Environment {
	child := New()
	child.enclosing = e
	return child
}
