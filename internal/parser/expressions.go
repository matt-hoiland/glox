package parser

import (
	"github.com/matt-hoiland/glox/internal/ast"
	ierrors "github.com/matt-hoiland/glox/internal/errors"
	"github.com/matt-hoiland/glox/internal/loxtype"
	"github.com/matt-hoiland/glox/internal/token"
)

// expression implements the production:
//
//	expression -> assignment ;
func (p *Parser) expression() (ast.Expr, error) {
	return p.assignment()
}

// assignment implements the production:
//
//	assignment -> IDENTIFIER "=" assignment
//	            | logic_or ;
func (p *Parser) assignment() (ast.Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(token.TypeEqual) {
		equals := p.previous()
		var value ast.Expr
		if value, err = p.assignment(); err != nil {
			return nil, err
		}

		if varExpr, ok := expr.(*ast.VariableExpr); ok {
			name := varExpr.Name
			return ast.NewAssignExpr(name, value), nil
		}

		return nil, ierrors.New(equals, nil)
	}
	return expr, nil
}

func (p *Parser) or() (ast.Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(token.TypeOr) {
		operator := p.previous()
		var right ast.Expr
		if right, err = p.and(); err != nil {
			return nil, err
		}
		expr = ast.NewLogicalExpr(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) and() (ast.Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(token.TypeAnd) {
		operator := p.previous()
		var right ast.Expr
		if right, err = p.equality(); err != nil {
			return nil, err
		}
		expr = ast.NewLogicalExpr(expr, operator, right)
	}

	return expr, nil
}

// equality implements the production:
//
//	equality -> comparison ( ( "!=" | "==" ) comparison )* ;
func (p *Parser) equality() (ast.Expr, error) {
	left, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(token.TypeBangEqual, token.TypeEqualEqual) {
		operator := p.previous()
		var right ast.Expr
		if right, err = p.comparison(); err != nil {
			return nil, err
		}
		left = ast.NewBinaryExpr(left, operator, right)
	}

	return left, nil
}

// comparison implements the production:
//
//	comparison -> term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
func (p *Parser) comparison() (ast.Expr, error) {
	left, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(token.TypeGreater, token.TypeGreaterEqual, token.TypeLess, token.TypeLessEqual) {
		operator := p.previous()
		var right ast.Expr
		if right, err = p.term(); err != nil {
			return nil, err
		}
		left = ast.NewBinaryExpr(left, operator, right)
	}

	return left, nil
}

// term implements the production:
//
//	term -> factor ( ( "-" | "+" ) factor )* ;
func (p *Parser) term() (ast.Expr, error) {
	left, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(token.TypeMinus, token.TypePlus) {
		operator := p.previous()
		var right ast.Expr
		if right, err = p.factor(); err != nil {
			return nil, err
		}
		left = ast.NewBinaryExpr(left, operator, right)
	}

	return left, nil
}

// factor implements the production:
//
//	factor -> unary ( ( "/" | "*" ) unary )* ;
func (p *Parser) factor() (ast.Expr, error) {
	left, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(token.TypeSlash, token.TypeStar) {
		operator := p.previous()
		var right ast.Expr
		if right, err = p.factor(); err != nil {
			return nil, err
		}
		left = ast.NewBinaryExpr(left, operator, right)
	}

	return left, nil
}

// unary implements the production:
//
//	unary -> ( "!" | "-" ) unary
//	       | primary ;
func (p *Parser) unary() (ast.Expr, error) {
	if p.match(token.TypeBang, token.TypeMinus) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return ast.NewUnaryExpr(operator, right), nil
	}

	return p.primary()
}

// primary implements the production:
//
//	primary -> "true" | "false" | "nil"
//	         | NUMBER | STRING
//	         | "(" expression ")"
//	         | IDENTIFIER ;
func (p *Parser) primary() (ast.Expr, error) {
	if p.match(token.TypeFalse) {
		return ast.NewLiteralExpr(loxtype.Boolean(false)), nil
	}
	if p.match(token.TypeTrue) {
		return ast.NewLiteralExpr(loxtype.Boolean(true)), nil
	}
	if p.match(token.TypeNil) {
		return ast.NewLiteralExpr(loxtype.Nil{}), nil
	}
	if p.match(token.TypeNumber, token.TypeString) {
		return ast.NewLiteralExpr(p.previous().Literal), nil
	}
	if p.match(token.TypeIdentifier) {
		return ast.NewVariableExpr(p.previous()), nil
	}
	if p.match(token.TypeLeftParen) {
		expression, err := p.expression()
		if err != nil {
			return nil, err
		}
		if _, err = p.consume(token.TypeRightParen, ErrUnterminatedExpression); err != nil {
			return nil, err
		}
		return ast.NewGroupingExpr(expression), nil
	}

	return nil, ErrUnimplemented
}
