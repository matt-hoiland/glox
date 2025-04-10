package ast

import (
	"strings"

	"github.com/matt-hoiland/glox/internal/environment"
	"github.com/matt-hoiland/glox/internal/loxtype"
)

func Print(e Stmt) string {
	s, _ := Printer{}.Print(e)
	return s.String()
}

type Printer struct{}

var _ ExprVisitor = Printer{}
var _ StmtVisitor = Printer{}

func (ap Printer) Print(e Stmt) (loxtype.Type, error) {
	return e.Accept(nil, ap)
}

func (ap Printer) parenthesize(env *environment.Environment, name string, expressions ...Expr) (loxtype.Type, error) {
	var builder strings.Builder
	builder.WriteRune('(')
	builder.WriteString(name)
	for _, e := range expressions {
		builder.WriteRune(' ')
		s, _ := e.Accept(env, ap)
		builder.WriteString(s.String())
	}
	builder.WriteRune(')')
	return loxtype.String(builder.String()), nil
}

func (ap Printer) VisitAssignExpr(*environment.Environment, *AssignExpr) (loxtype.Type, error) {
	panic("unimplemented")
}

func (ap Printer) VisitBinaryExpr(env *environment.Environment, e *BinaryExpr) (loxtype.Type, error) {
	return ap.parenthesize(env, e.Operator.Lexeme, e.Left, e.Right)
}

func (ap Printer) VisitGroupingExpr(env *environment.Environment, e *GroupingExpr) (loxtype.Type, error) {
	return ap.parenthesize(env, "group", e.Expression)
}

func (ap Printer) VisitLiteralExpr(_ *environment.Environment, e *LiteralExpr) (loxtype.Type, error) {
	if e.Value == nil {
		return loxtype.Nil{}, nil
	}
	return e.Value, nil
}

func (Printer) VisitLogicalExpr(*environment.Environment, *LogicalExpr) (loxtype.Type, error) {
	panic("unimplemented")
}

func (ap Printer) VisitUnaryExpr(env *environment.Environment, e *UnaryExpr) (loxtype.Type, error) {
	return ap.parenthesize(env, e.Operator.Lexeme, e.Right)
}

func (ap Printer) VisitVariableExpr(*environment.Environment, *VariableExpr) (loxtype.Type, error) {
	panic("unimplemented")
}

func (ap Printer) VisitBlockStmt(*environment.Environment, *BlockStmt) (loxtype.Type, error) {
	panic("unimplemented")
}

func (ap Printer) VisitExpressionStmt(env *environment.Environment, s *ExpressionStmt) (loxtype.Type, error) {
	value, _ := s.Expression.Accept(env, ap)
	return loxtype.String(value.String() + ";"), nil
}

func (ap Printer) VisitIfStmt(*environment.Environment, *IfStmt) (loxtype.Type, error) {
	panic("unimplemented")
}

func (ap Printer) VisitPrintStmt(env *environment.Environment, s *PrintStmt) (loxtype.Type, error) {
	value, _ := s.Expression.Accept(env, ap)
	return loxtype.String("print " + value.String() + ";"), nil
}

func (ap Printer) VisitVarStmt(*environment.Environment, *VarStmt) (loxtype.Type, error) {
	panic("unimplemented")
}

func (Printer) VisitWhileStmt(*environment.Environment, *WhileStmt) (loxtype.Type, error) {
	panic("unimplemented")
}
