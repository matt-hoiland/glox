package interpreter

import (
	"fmt"

	"github.com/matt-hoiland/glox/internal/ast"
	"github.com/matt-hoiland/glox/internal/environment"
	"github.com/matt-hoiland/glox/internal/loxtype"
)

type callable interface {
	loxtype.Type
	Arity() int
	Call(*environment.Environment, *Interpreter, []loxtype.Type) (loxtype.Type, error)
}

type function ast.FunctionStmt

var _ callable = (*function)(nil)

func (f *function) Arity() int                        { return len(f.Params) }
func (*function) Equals(loxtype.Type) loxtype.Boolean { return false }
func (*function) IsTruthy() loxtype.Boolean           { return true }
func (f *function) String() string                    { return fmt.Sprintf("<fn: %s>", f.Name.Lexeme) }

func (f *function) Call(parent *environment.Environment, i *Interpreter, args []loxtype.Type) (loxtype.Type, error) {
	env := parent.MakeChild()
	for i, param := range f.Params {
		env.Define(param, args[i])
	}
	if err := i.executeBlock(env, f.Body); err != nil {
		return nil, err
	}
	return nil, nil //nolint:nilnil // TODO: Fix this later.
}

type nativeFunction struct {
	name  string
	arity int
	impl  func(*Interpreter, []loxtype.Type) (loxtype.Type, error)
}

var _ callable = (*nativeFunction)(nil)

func (nf *nativeFunction) Arity() int                       { return nf.arity }
func (*nativeFunction) Equals(loxtype.Type) loxtype.Boolean { return false }
func (*nativeFunction) IsTruthy() loxtype.Boolean           { return true }
func (nf *nativeFunction) String() string                   { return fmt.Sprintf("<native fn: %s>", nf.name) }

func (nf *nativeFunction) Call(_ *environment.Environment, i *Interpreter, args []loxtype.Type) (loxtype.Type, error) {
	return nf.impl(i, args)
}
