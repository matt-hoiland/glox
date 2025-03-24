package interpreter

import (
	"errors"
	"fmt"

	"github.com/matt-hoiland/glox/internal/expr"
	"github.com/matt-hoiland/glox/internal/loxtype"
	"github.com/matt-hoiland/glox/internal/token"
)

var (
	ErrUnimplemented = errors.New("unimplemented")

	ErrType           = errors.New("type-error")
	ErrNonBooleanType = fmt.Errorf("non-boolean %w", ErrType)
	ErrNonNumericType = fmt.Errorf("non-numeric %w", ErrType)
	ErrNonStringType  = fmt.Errorf("non-string %w", ErrType)
)

type Interpreter struct{}

var _ expr.Visitor = (*Interpreter)(nil)

func New() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) Evaluate(e expr.Expr) (loxtype.Type, error) {
	return i.evaluate(e)
}

func (i *Interpreter) evaluate(e expr.Expr) (loxtype.Type, error) {
	return e.Accept(i)
}

func (i *Interpreter) isEqual(a, b loxtype.Type) loxtype.Boolean {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	ae, ok := a.(loxtype.Equalser)
	if !ok {
		return false
	}
	return ae.Equals(b)
}

func (i *Interpreter) isTruthy(value loxtype.Type) loxtype.Boolean {
	if value == nil {
		return false
	}
	if _, ok := value.(loxtype.Nil); ok {
		return false
	}
	if b, ok := value.(loxtype.Boolean); ok {
		return b
	}
	return true
}

func (i *Interpreter) VisitBinary(e *expr.Binary) (loxtype.Type, error) {
	left, err := i.evaluate(e.Left)
	if err != nil {
		return nil, fmt.Errorf("could not evaluate left operand of binary expression: %w", err)
	}
	right, err := i.evaluate(e.Right)
	if err != nil {
		return nil, fmt.Errorf("could not evaluate right operand of binary expression: %w", err)
	}

	switch e.Operator.Type {
	case token.TypeBangEqual:
		return !i.isEqual(left, right), nil
	case token.TypeEqualEqual:
		return i.isEqual(left, right), nil
	case token.TypeGreater:
		a, b, err := convertBoth[loxtype.Number](left, right, ErrNonNumericType)
		if err != nil {
			return nil, fmt.Errorf("greater expression: %w", err)
		}
		return a.Greater(b), nil
	case token.TypeGreaterEqual:
		a, b, err := convertBoth[loxtype.Number](left, right, ErrNonNumericType)
		if err != nil {
			return nil, fmt.Errorf("greater-equal expression: %w", err)
		}
		return a.GreaterEqual(b), nil
	case token.TypeLess:
		a, b, err := convertBoth[loxtype.Number](left, right, ErrNonNumericType)
		if err != nil {
			return nil, fmt.Errorf("less expression: %w", err)
		}
		return a.Less(b), nil
	case token.TypeLessEqual:
		a, b, err := convertBoth[loxtype.Number](left, right, ErrNonNumericType)
		if err != nil {
			return nil, fmt.Errorf("less-equal expression: %w", err)
		}
		return a.LessEqual(b), nil
	case token.TypeMinus:
		a, b, err := convertBoth[loxtype.Number](left, right, ErrNonNumericType)
		if err != nil {
			return nil, fmt.Errorf("minus expression: %w", err)
		}
		return a.Subtract(b), nil
	case token.TypeSlash:
		a, b, err := convertBoth[loxtype.Number](left, right, ErrNonNumericType)
		if err != nil {
			return nil, fmt.Errorf("slash expression: %w", err)
		}
		return a.Divide(b), nil
	case token.TypeStar:
		a, b, err := convertBoth[loxtype.Number](left, right, ErrNonNumericType)
		if err != nil {
			return nil, fmt.Errorf("star expression: %w", err)
		}
		return a.Multiply(b), nil
	case token.TypePlus:
		na, nb, nErr := convertBoth[loxtype.Number](left, right, ErrNonNumericType)
		sa, sb, sErr := convertBoth[loxtype.String](left, right, ErrNonStringType)
		if nErr != nil && sErr != nil {
			return nil, fmt.Errorf("operands to plus expression must be either string or numeric: %w", ErrType)
		}
		if nErr != nil {
			return sa.Add(sb), nil
		}
		return na.Add(nb), nil
	}

	return nil, ErrUnimplemented
}

func (i *Interpreter) VisitGrouping(e *expr.Grouping) (loxtype.Type, error) {
	return i.evaluate(e.Expression)
}

func (i *Interpreter) VisitLiteral(e *expr.Literal) (loxtype.Type, error) {
	return e.Value, nil
}

func (i *Interpreter) VisitUnary(e *expr.Unary) (loxtype.Type, error) {
	right, err := i.evaluate(e.Right)
	if err != nil {
		return nil, fmt.Errorf("could not evaluate operand of unary expression: %w", err)
	}

	switch e.Operator.Type {
	case token.TypeBang:
		return i.isTruthy(right).Negate(), nil
	case token.TypeMinus:
		n, ok := right.(loxtype.Number)
		if !ok {
			return nil, fmt.Errorf("cannot apply minus operator: %w", ErrNonNumericType)
		}
		return n.Negate(), nil
	}

	return nil, ErrUnimplemented
}

func convertBoth[T any](a, b loxtype.Type, err error) (T, T, error) {
	at, aok := a.(T)
	bt, bok := b.(T)
	if !aok {
		return at, bt, fmt.Errorf("type-error: left-hand operand: %w", err)
	}
	if !bok {
		return at, bt, fmt.Errorf("type-error: right-hand operand: %w", err)
	}
	return at, bt, nil
}
