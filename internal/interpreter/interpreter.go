package interpreter

import (
	"errors"
	"fmt"
	"io"

	"github.com/matt-hoiland/glox/internal/ast"
	"github.com/matt-hoiland/glox/internal/environment"
	"github.com/matt-hoiland/glox/internal/loxtype"
	"github.com/matt-hoiland/glox/internal/parser"
	"github.com/matt-hoiland/glox/internal/scanner"
	"github.com/matt-hoiland/glox/internal/token"
)

var (
	ErrUnimplemented = errors.New("unimplemented")

	ErrType           = errors.New("type-error")
	ErrNonBooleanType = fmt.Errorf("non-boolean %w", ErrType)
	ErrNonNumericType = fmt.Errorf("non-numeric %w", ErrType)
	ErrNonStringType  = fmt.Errorf("non-string %w", ErrType)
)

type Interpreter struct {
	w io.Writer
}

var (
	_ ast.ExprVisitor = (*Interpreter)(nil)
	_ ast.StmtVisitor = (*Interpreter)(nil)
)

func New(w io.Writer) *Interpreter {
	return &Interpreter{w}
}

func (i *Interpreter) Run(env *environment.Environment, code string) error {
	var (
		tokens []*token.Token
		stmts  []ast.Stmt
		err    error
	)

	if tokens, err = scanner.New(code).ScanTokens(); err != nil {
		return err
	}

	if stmts, err = parser.New(tokens).Parse(); err != nil {
		return err
	}

	if err = i.Interpret(env, stmts); err != nil {
		return err
	}
	return nil
}

func (i *Interpreter) Interpret(env *environment.Environment, stmts []ast.Stmt) error {
	for _, s := range stmts {
		if err := i.execute(env, s); err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) execute(env *environment.Environment, s ast.Stmt) error {
	if _, err := s.Accept(env, i); err != nil {
		return err
	}
	return nil
}

func (i *Interpreter) executeBlock(env *environment.Environment, stmts []ast.Stmt) error {
	for _, stmt := range stmts {
		if err := i.execute(env, stmt); err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) evaluate(env *environment.Environment, e ast.Expr) (loxtype.Type, error) {
	return e.Accept(env, i)
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

func (i *Interpreter) VisitBlockStmt(env *environment.Environment, s *ast.BlockStmt) (loxtype.Type, error) {
	if err := i.executeBlock(env.Enclose(), s.Statements); err != nil {
		return nil, err
	}
	return nil, nil //nolint:nilnil // TODO: Emit final type?
}

func (i *Interpreter) VisitExpressionStmt(env *environment.Environment, s *ast.ExpressionStmt) (loxtype.Type, error) {
	_, err := i.evaluate(env, s.Expression)
	if err != nil {
		return nil, err
	}
	return nil, nil //nolint:nilnil // TODO: Emit final type?
}

func (i *Interpreter) VisitIfStmt(env *environment.Environment, s *ast.IfStmt) (loxtype.Type, error) {
	cond, err := i.evaluate(env, s.Condition)
	if err != nil {
		return nil, err
	}

	if i.isTruthy(cond) {
		err = i.execute(env, s.ThenBranch)
	} else if s.ElseBranch != nil {
		err = i.execute(env, s.ElseBranch)
	}

	if err != nil {
		return nil, err
	}
	return nil, nil //nolint:nilnil // TODO: Emit final type?
}

func (i *Interpreter) VisitPrintStmt(env *environment.Environment, s *ast.PrintStmt) (loxtype.Type, error) {
	value, err := i.evaluate(env, s.Expression)
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(i.w, value)
	return nil, nil //nolint:nilnil // TODO: Emit final type?
}

func (i *Interpreter) VisitVarStmt(env *environment.Environment, s *ast.VarStmt) (loxtype.Type, error) {
	var (
		value loxtype.Type = loxtype.Nil{}
		err   error
	)
	if s.Initializer != nil {
		if value, err = i.evaluate(env, s.Initializer); err != nil {
			return nil, err
		}
	}

	env.Define(s.Name, value)
	return nil, nil //nolint:nilnil // TODO: Emit final type?
}

func (*Interpreter) VisitWhileStmt(*environment.Environment, *ast.WhileStmt) (loxtype.Type, error) {
	panic("unimplemented")
}

func (i *Interpreter) VisitAssignExpr(env *environment.Environment, e *ast.AssignExpr) (loxtype.Type, error) {
	value, err := i.evaluate(env, e.Value)
	if err != nil {
		return nil, err
	}
	if err = env.Assign(e.Name, value); err != nil {
		return nil, err
	}
	return value, nil
}

func (i *Interpreter) VisitBinaryExpr(env *environment.Environment, e *ast.BinaryExpr) (loxtype.Type, error) {
	left, err := i.evaluate(env, e.Left)
	if err != nil {
		return nil, fmt.Errorf("could not evaluate left operand of binary expression: %w", err)
	}
	right, err := i.evaluate(env, e.Right)
	if err != nil {
		return nil, fmt.Errorf("could not evaluate right operand of binary expression: %w", err)
	}

	switch e.Operator.Type {
	case token.TypeBangEqual:
		return !i.isEqual(left, right), nil
	case token.TypeEqualEqual:
		return i.isEqual(left, right), nil
	case token.TypeGreater:
		return i.biopGreater(left, right)
	case token.TypeGreaterEqual:
		return i.biopGreaterEqual(left, right)
	case token.TypeLess:
		return i.biopLess(left, right)
	case token.TypeLessEqual:
		return i.biopLessEqual(left, right)
	case token.TypeMinus:
		return i.biopMinus(left, right)
	case token.TypeSlash:
		return i.biopSlash(left, right)
	case token.TypeStar:
		return i.biopStar(left, right)
	case token.TypePlus:
		return i.biopPlus(left, right)
	default:
		return nil, ErrUnimplemented
	}
}

func (*Interpreter) biopGreater(left, right loxtype.Type) (loxtype.Type, error) {
	a, b, ok := convertBoth[loxtype.Number](left, right)
	if !ok {
		return nil, fmt.Errorf("greater expression: %w", ErrNonNumericType)
	}
	return a.Greater(b), nil
}

func (*Interpreter) biopGreaterEqual(left loxtype.Type, right loxtype.Type) (loxtype.Type, error) {
	a, b, ok := convertBoth[loxtype.Number](left, right)
	if !ok {
		return nil, fmt.Errorf("greater-equal expression: %w", ErrNonNumericType)
	}
	return a.GreaterEqual(b), nil
}

func (*Interpreter) biopLess(left loxtype.Type, right loxtype.Type) (loxtype.Type, error) {
	a, b, ok := convertBoth[loxtype.Number](left, right)
	if !ok {
		return nil, fmt.Errorf("less expression: %w", ErrNonNumericType)
	}
	return a.Less(b), nil
}

func (*Interpreter) biopLessEqual(left loxtype.Type, right loxtype.Type) (loxtype.Type, error) {
	a, b, ok := convertBoth[loxtype.Number](left, right)
	if !ok {
		return nil, fmt.Errorf("less-equal expression: %w", ErrNonNumericType)
	}
	return a.LessEqual(b), nil
}

func (*Interpreter) biopMinus(left loxtype.Type, right loxtype.Type) (loxtype.Type, error) {
	a, b, ok := convertBoth[loxtype.Number](left, right)
	if !ok {
		return nil, fmt.Errorf("minus expression: %w", ErrNonNumericType)
	}
	return a.Subtract(b), nil
}

func (*Interpreter) biopSlash(left loxtype.Type, right loxtype.Type) (loxtype.Type, error) {
	a, b, ok := convertBoth[loxtype.Number](left, right)
	if !ok {
		return nil, fmt.Errorf("slash expression: %w", ErrNonNumericType)
	}
	return a.Divide(b), nil
}

func (*Interpreter) biopStar(left loxtype.Type, right loxtype.Type) (loxtype.Type, error) {
	a, b, ok := convertBoth[loxtype.Number](left, right)
	if !ok {
		return nil, fmt.Errorf("star expression: %w", ErrNonNumericType)
	}
	return a.Multiply(b), nil
}

func (*Interpreter) biopPlus(left loxtype.Type, right loxtype.Type) (loxtype.Type, error) {
	na, nb, nok := convertBoth[loxtype.Number](left, right)
	sa, sb, sok := convertBoth[loxtype.String](left, right)
	if !nok && !sok {
		return nil, fmt.Errorf("operands to plus expression must be either string or numeric: %w", ErrType)
	}
	if !nok {
		return sa.Add(sb), nil
	}
	return na.Add(nb), nil
}

func (i *Interpreter) VisitGroupingExpr(env *environment.Environment, e *ast.GroupingExpr) (loxtype.Type, error) {
	return i.evaluate(env, e.Expression)
}

func (i *Interpreter) VisitLiteralExpr(_ *environment.Environment, e *ast.LiteralExpr) (loxtype.Type, error) {
	return e.Value, nil
}

func (i *Interpreter) VisitLogicalExpr(env *environment.Environment, e *ast.LogicalExpr) (loxtype.Type, error) {
	left, err := i.evaluate(env, e.Left)
	if err != nil {
		return nil, err
	}

	if e.Operator.Type == token.TypeOr {
		if i.isTruthy(left) {
			return left, nil
		}
	} else {
		if !i.isTruthy(left) {
			return left, nil
		}
	}

	return i.evaluate(env, e.Right)
}

func (i *Interpreter) VisitUnaryExpr(env *environment.Environment, e *ast.UnaryExpr) (loxtype.Type, error) {
	right, err := i.evaluate(env, e.Right)
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
	default:
		return nil, ErrUnimplemented
	}
}

func (i *Interpreter) VisitVariableExpr(env *environment.Environment, e *ast.VariableExpr) (loxtype.Type, error) {
	return env.Get(e.Name)
}

func convertBoth[T any](a, b loxtype.Type) (T, T, bool) {
	at, aok := a.(T)
	bt, bok := b.(T)
	if !aok {
		return at, bt, false
	}
	if !bok {
		return at, bt, false
	}
	return at, bt, true
}
