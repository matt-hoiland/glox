// Code generated by tools/generate-ast. DO NOT EDIT.
package ast

import (
	"github.com/matt-hoiland/glox/internal/environment"
	"github.com/matt-hoiland/glox/internal/loxtype"
	"github.com/matt-hoiland/glox/internal/token"
)

type Stmt interface {
	Accept(*environment.Environment, StmtVisitor) (loxtype.Type, error)
}

type StmtVisitor interface {
	VisitBlockStmt(*environment.Environment, *BlockStmt) (loxtype.Type, error)
	VisitExpressionStmt(*environment.Environment, *ExpressionStmt) (loxtype.Type, error)
	VisitIfStmt(*environment.Environment, *IfStmt) (loxtype.Type, error)
	VisitPrintStmt(*environment.Environment, *PrintStmt) (loxtype.Type, error)
	VisitVarStmt(*environment.Environment, *VarStmt) (loxtype.Type, error)
	VisitWhileStmt(*environment.Environment, *WhileStmt) (loxtype.Type, error)
}

type BlockStmt struct {
	Statements []Stmt
}

var _ Stmt = (*BlockStmt)(nil)

func NewBlockStmt(Statements []Stmt) *BlockStmt {
	return &BlockStmt{
		Statements: Statements,
	}
}

func (e *BlockStmt) Accept(env *environment.Environment, visitor StmtVisitor) (loxtype.Type, error) {
	return visitor.VisitBlockStmt(env, e)
}

type ExpressionStmt struct {
	Expression Expr
}

var _ Stmt = (*ExpressionStmt)(nil)

func NewExpressionStmt(Expression Expr) *ExpressionStmt {
	return &ExpressionStmt{
		Expression: Expression,
	}
}

func (e *ExpressionStmt) Accept(env *environment.Environment, visitor StmtVisitor) (loxtype.Type, error) {
	return visitor.VisitExpressionStmt(env, e)
}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

var _ Stmt = (*IfStmt)(nil)

func NewIfStmt(Condition Expr, ThenBranch Stmt, ElseBranch Stmt) *IfStmt {
	return &IfStmt{
		Condition:  Condition,
		ThenBranch: ThenBranch,
		ElseBranch: ElseBranch,
	}
}

func (e *IfStmt) Accept(env *environment.Environment, visitor StmtVisitor) (loxtype.Type, error) {
	return visitor.VisitIfStmt(env, e)
}

type PrintStmt struct {
	Expression Expr
}

var _ Stmt = (*PrintStmt)(nil)

func NewPrintStmt(Expression Expr) *PrintStmt {
	return &PrintStmt{
		Expression: Expression,
	}
}

func (e *PrintStmt) Accept(env *environment.Environment, visitor StmtVisitor) (loxtype.Type, error) {
	return visitor.VisitPrintStmt(env, e)
}

type VarStmt struct {
	Name        *token.Token
	Initializer Expr
}

var _ Stmt = (*VarStmt)(nil)

func NewVarStmt(Name *token.Token, Initializer Expr) *VarStmt {
	return &VarStmt{
		Name:        Name,
		Initializer: Initializer,
	}
}

func (e *VarStmt) Accept(env *environment.Environment, visitor StmtVisitor) (loxtype.Type, error) {
	return visitor.VisitVarStmt(env, e)
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

var _ Stmt = (*WhileStmt)(nil)

func NewWhileStmt(Condition Expr, Body Stmt) *WhileStmt {
	return &WhileStmt{
		Condition: Condition,
		Body:      Body,
	}
}

func (e *WhileStmt) Accept(env *environment.Environment, visitor StmtVisitor) (loxtype.Type, error) {
	return visitor.VisitWhileStmt(env, e)
}
