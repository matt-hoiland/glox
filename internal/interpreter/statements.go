package interpreter

import (
	"fmt"

	"github.com/matt-hoiland/glox/internal/ast"
	"github.com/matt-hoiland/glox/internal/environment"
	"github.com/matt-hoiland/glox/internal/loxtype"
)

var _ ast.StmtVisitor = (*Interpreter)(nil)

func (i *Interpreter) execute(env *environment.Environment, s ast.Stmt) error {
	if _, err := s.Accept(env, i); err != nil {
		return err
	}
	return nil
}

func (i *Interpreter) VisitBlockStmt(env *environment.Environment, s *ast.BlockStmt) (loxtype.Type, error) {
	if err := i.executeBlock(env.MakeChild(), s.Statements); err != nil {
		return nil, err
	}
	return nil, nil //nolint:nilnil // TODO: Emit final type?
}

func (i *Interpreter) executeBlock(env *environment.Environment, stmts []ast.Stmt) error {
	for _, stmt := range stmts {
		if err := i.execute(env, stmt); err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) VisitExpressionStmt(env *environment.Environment, s *ast.ExpressionStmt) (loxtype.Type, error) {
	_, err := i.evaluate(env, s.Expression)
	if err != nil {
		return nil, err
	}
	return nil, nil //nolint:nilnil // TODO: Emit final type?
}

func (i *Interpreter) VisitFunctionStmt(env *environment.Environment, s *ast.FunctionStmt) (loxtype.Type, error) {
	env.Define(s.Name, newFunction(env, s))
	return nil, nil //nolint:nilnil // TODO: Emit final type?
}

func (i *Interpreter) VisitIfStmt(env *environment.Environment, s *ast.IfStmt) (loxtype.Type, error) {
	cond, err := i.evaluate(env, s.Condition)
	if err != nil {
		return nil, err
	}

	if cond.IsTruthy() {
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

func (i *Interpreter) VisitWhileStmt(env *environment.Environment, s *ast.WhileStmt) (loxtype.Type, error) {
	var (
		cond loxtype.Type
		err  error
	)

	for {
		if cond, err = i.evaluate(env, s.Condition); err != nil {
			return nil, err
		}
		if !cond.IsTruthy() {
			break
		}
		if err = i.execute(env, s.Body); err != nil {
			return nil, err
		}
	}

	return nil, nil //nolint:nilnil // TODO: Emit final type?
}
