// Code generated by tools/generate-ast. DO NOT EDIT.
package expr

import (
	"fmt"

	"github.com/matt-hoiland/glox/internal/token"
)

type Expr[R any] interface {
	Accept(Visitor[R]) R
}

type Value interface {
	fmt.Stringer
}

type Visitor[R any] interface {
	VisitBinary(*Binary[R]) R
	VisitGrouping(*Grouping[R]) R
	VisitLiteral(*Literal[R]) R
	VisitUnary(*Unary[R]) R
}

type Binary[R any] struct {
	Left     Expr[R]
	Operator *token.Token
	Right    Expr[R]
}

var _ Expr[any] = (*Binary[any])(nil)

func NewBinary[R any](Left Expr[R], Operator *token.Token, Right Expr[R]) *Binary[R] {
	return &Binary[R]{
		Left:     Left,
		Operator: Operator,
		Right:    Right,
	}
}

func (e *Binary[R]) Accept(visitor Visitor[R]) R {
	return visitor.VisitBinary(e)
}

type Grouping[R any] struct {
	Expression Expr[R]
}

var _ Expr[any] = (*Grouping[any])(nil)

func NewGrouping[R any](Expression Expr[R]) *Grouping[R] {
	return &Grouping[R]{
		Expression: Expression,
	}
}

func (e *Grouping[R]) Accept(visitor Visitor[R]) R {
	return visitor.VisitGrouping(e)
}

type Literal[R any] struct {
	Value Value
}

var _ Expr[any] = (*Literal[any])(nil)

func NewLiteral[R any](Value Value) *Literal[R] {
	return &Literal[R]{
		Value: Value,
	}
}

func (e *Literal[R]) Accept(visitor Visitor[R]) R {
	return visitor.VisitLiteral(e)
}

type Unary[R any] struct {
	Operator *token.Token
	Right    Expr[R]
}

var _ Expr[any] = (*Unary[any])(nil)

func NewUnary[R any](Operator *token.Token, Right Expr[R]) *Unary[R] {
	return &Unary[R]{
		Operator: Operator,
		Right:    Right,
	}
}

func (e *Unary[R]) Accept(visitor Visitor[R]) R {
	return visitor.VisitUnary(e)
}
