//nolint:nilnil // It's the only way for now.
package interpreter

import (
	"errors"

	"github.com/matt-hoiland/glox/internal/ast"
	"github.com/matt-hoiland/glox/internal/environment"
	ierrors "github.com/matt-hoiland/glox/internal/errors"
	"github.com/matt-hoiland/glox/internal/loxtype"
	"github.com/matt-hoiland/glox/internal/token"
)

type resolver struct {
	i               *Interpreter
	scopes          []map[string]bool
	currentFunction functionType
}

var (
	_ ast.StmtVisitor = (*resolver)(nil)
	_ ast.ExprVisitor = (*resolver)(nil)
)

func newResolver(i *Interpreter) *resolver {
	return &resolver{
		i: i,
		scopes: []map[string]bool{
			{}, // global scope
		},
		currentFunction: none,
	}
}

func (r *resolver) beginScope() {
	r.scopes = append(r.scopes, map[string]bool{})
}

func (r *resolver) currentScope() map[string]bool {
	return r.scopes[len(r.scopes)-1]
}

func (r *resolver) declare(name *token.Token) error {
	if len(r.scopes) == 0 {
		return ierrors.New(name, errors.New("no scope"))
	}

	if _, declared := r.currentScope()[name.Lexeme]; declared {
		return ierrors.New(name, errors.New("redeclaration of scoped variable"))
	}
	r.currentScope()[name.Lexeme] = false
	return nil
}

func (r *resolver) define(name *token.Token) {
	if len(r.scopes) == 0 {
		return
	}

	r.currentScope()[name.Lexeme] = true
}

func (r *resolver) endScope() {
	r.scopes = r.scopes[0 : len(r.scopes)-1]
}

func (r *resolver) resolveFunction(s *ast.FunctionStmt, ft functionType) error {
	enclosingFunction := r.currentFunction
	r.currentFunction = ft

	r.beginScope()
	for _, param := range s.Params {
		if err := r.declare(param); err != nil {
			return err
		}
		r.define(param)
	}
	if err := r.resolveStmts(s.Body); err != nil {
		return err
	}
	r.endScope()

	r.currentFunction = enclosingFunction
	return nil
}

func (r *resolver) resolveLocal(e ast.Expr, name *token.Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.Lexeme]; ok {
			r.i.resolve(e, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *resolver) resolveStmts(stmts []ast.Stmt) error {
	for _, stmt := range stmts {
		if err := r.resolveStmt(stmt); err != nil {
			return err
		}
	}
	return nil
}

func (r *resolver) resolveStmt(s ast.Stmt) error {
	if _, err := s.Accept(nil, r); err != nil {
		return err
	}
	return nil
}

func (r *resolver) resolveExpr(e ast.Expr) error {
	if _, err := e.Accept(nil, r); err != nil {
		return err
	}
	return nil
}

func (r *resolver) VisitBlockStmt(_ *environment.Environment, s *ast.BlockStmt) (loxtype.Type, error) {
	r.beginScope()
	if err := r.resolveStmts(s.Statements); err != nil {
		return nil, err
	}
	r.endScope()
	return nil, nil
}

func (r *resolver) VisitExpressionStmt(_ *environment.Environment, s *ast.ExpressionStmt) (loxtype.Type, error) {
	if err := r.resolveExpr(s.Expression); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *resolver) VisitFunctionStmt(_ *environment.Environment, s *ast.FunctionStmt) (loxtype.Type, error) {
	if err := r.declare(s.Name); err != nil {
		return nil, err
	}
	r.define(s.Name)

	if err := r.resolveFunction(s, function); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *resolver) VisitIfStmt(_ *environment.Environment, s *ast.IfStmt) (loxtype.Type, error) {
	if err := r.resolveExpr(s.Condition); err != nil {
		return nil, err
	}
	if err := r.resolveStmt(s.ThenBranch); err != nil {
		return nil, err
	}
	if s.ElseBranch != nil {
		if err := r.resolveStmt(s.ElseBranch); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (r *resolver) VisitPrintStmt(_ *environment.Environment, s *ast.PrintStmt) (loxtype.Type, error) {
	if err := r.resolveExpr(s.Expression); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *resolver) VisitReturnStmt(_ *environment.Environment, s *ast.ReturnStmt) (loxtype.Type, error) {
	if r.currentFunction == none {
		return nil, ierrors.New(s.Keyword, errors.New("can't return from top-level code"))
	}
	if err := r.resolveExpr(s.Value); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *resolver) VisitVarStmt(_ *environment.Environment, s *ast.VarStmt) (loxtype.Type, error) {
	if err := r.declare(s.Name); err != nil {
		return nil, err
	}
	if s.Initializer != nil {
		if err := r.resolveExpr(s.Initializer); err != nil {
			return nil, err
		}
	}
	r.define(s.Name)
	return nil, nil
}

func (r *resolver) VisitWhileStmt(_ *environment.Environment, s *ast.WhileStmt) (loxtype.Type, error) {
	if err := r.resolveExpr(s.Condition); err != nil {
		return nil, err
	}
	if err := r.resolveStmt(s.Body); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *resolver) VisitAssignExpr(_ *environment.Environment, e *ast.AssignExpr) (loxtype.Type, error) {
	if err := r.resolveExpr(e.Value); err != nil {
		return nil, err
	}
	r.resolveLocal(e, e.Name)
	return nil, nil
}

func (r *resolver) VisitBinaryExpr(_ *environment.Environment, e *ast.BinaryExpr) (loxtype.Type, error) {
	if err := r.resolveExpr(e.Left); err != nil {
		return nil, err
	}
	if err := r.resolveExpr(e.Right); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *resolver) VisitCallExpr(_ *environment.Environment, e *ast.CallExpr) (loxtype.Type, error) {
	if err := r.resolveExpr(e.Callee); err != nil {
		return nil, err
	}

	for _, expr := range e.Arguments {
		if err := r.resolveExpr(expr); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *resolver) VisitGroupingExpr(_ *environment.Environment, e *ast.GroupingExpr) (loxtype.Type, error) {
	if err := r.resolveExpr(e.Expression); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *resolver) VisitLiteralExpr(*environment.Environment, *ast.LiteralExpr) (loxtype.Type, error) {
	return nil, nil
}

func (r *resolver) VisitLogicalExpr(_ *environment.Environment, e *ast.LogicalExpr) (loxtype.Type, error) {
	if err := r.resolveExpr(e.Left); err != nil {
		return nil, err
	}
	if err := r.resolveExpr(e.Right); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *resolver) VisitUnaryExpr(_ *environment.Environment, e *ast.UnaryExpr) (loxtype.Type, error) {
	if err := r.resolveExpr(e.Right); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *resolver) VisitVariableExpr(_ *environment.Environment, e *ast.VariableExpr) (loxtype.Type, error) {
	if defined, ok := r.currentScope()[e.Name.Lexeme]; ok && !defined {
		return nil, ierrors.New(e.Name, errors.New("can't read local variable in its own initializer"))
	}

	r.resolveLocal(e, e.Name)
	return nil, nil
}
