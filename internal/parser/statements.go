package parser

import (
	"errors"
	"fmt"

	"github.com/matt-hoiland/glox/internal/ast"
	"github.com/matt-hoiland/glox/internal/loxtype"
	"github.com/matt-hoiland/glox/internal/token"
)

// declaration implements the production:
//
//	declaration -> funDecl
//	             | varDecl
//	             | statement ;
func (p *Parser) declaration() (ast.Stmt, error) {
	var (
		stmt ast.Stmt
		err  error
	)

	defer func() {
		if err != nil {
			p.synchronize()
		}
	}()

	switch {
	case p.match(token.TypeFun):
		stmt, err = p.function(function)
	case p.match(token.TypeVar):
		stmt, err = p.varDeclaration()
	default:
		stmt, err = p.statement()
	}

	return stmt, err
}

type funcKind string

const (
	function funcKind = "function"
	method   funcKind = "method"
)

// function implements the productions:
//
//	funDecl  -> "fun" function ;
//	function -> IDENTIFIER "(" parameters? ")" block ;
func (p *Parser) function(kind funcKind) (ast.Stmt, error) {
	var (
		name   *token.Token
		params []*token.Token
		body   []ast.Stmt
		err    error
	)

	if name, err = p.consume(token.TypeIdentifier, fmt.Errorf("expect %s name", kind)); err != nil {
		return nil, err
	}

	if _, err = p.consume(token.TypeLeftParen, fmt.Errorf("expect '(' after %s name", kind)); err != nil {
		return nil, err
	}

	if !p.check(token.TypeRightParen) {
		for {
			// TODO: check len(params) <= 255

			var param *token.Token
			if param, err = p.consume(token.TypeIdentifier, errors.New("expect parameter name")); err != nil {
				return nil, err
			}
			params = append(params, param)

			if !p.match(token.TypeComma) {
				break
			}
		}
	}

	if _, err = p.consume(token.TypeRightParen, errors.New("expect ')' after parameters")); err != nil {
		return nil, err
	}

	if _, err = p.consume(token.TypeLeftBrace, fmt.Errorf("expect '{' before %s body", kind)); err != nil {
		return nil, err
	}

	if body, err = p.block(); err != nil {
		return nil, err
	}

	return ast.NewFunctionStmt(name, params, body), nil
}

// varDeclaration implements the production:
//
//	varDecl -> "var" IDENTIFIER ( "=" expression )? ";" ;
func (p *Parser) varDeclaration() (ast.Stmt, error) {
	var (
		name        *token.Token
		initializer ast.Expr
		err         error
	)

	if name, err = p.consume(token.TypeIdentifier, ErrNoVariableName); err != nil {
		return nil, err
	}

	if p.match(token.TypeEqual) {
		if initializer, err = p.expression(); err != nil {
			return nil, err
		}
	}

	if _, err = p.consume(token.TypeSemicolon, ErrUnterminatedStatement); err != nil {
		return nil, err
	}

	return ast.NewVarStmt(name, initializer), nil
}

// statement implements the production:
//
//	statement -> exprStmt
//	           | forStmt
//	           | ifStmt
//	           | printStmt
//	           | returnStmt
//	           | whileStmt
//	           | block ;
func (p *Parser) statement() (ast.Stmt, error) {
	switch {
	case p.match(token.TypeFor):
		return p.forStatement()

	case p.match(token.TypeIf):
		return p.ifStatement()

	case p.match(token.TypePrint):
		return p.printStatement()

	case p.match(token.TypeReturn):
		return p.returnStatement()

	case p.match(token.TypeWhile):
		return p.whileStatement()

	case p.match(token.TypeLeftBrace):
		stmts, err := p.block()
		if err != nil {
			return nil, err
		}
		block := ast.NewBlockStmt(stmts)
		return block, nil

	default:
		return p.expressionStatement()
	}
}

// expressionStatement implements the production:
//
//	exprStmt -> expression ";" ;
func (p *Parser) expressionStatement() (ast.Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	if p.replMode && p.peek().Type == token.TypeEOF {
		return ast.NewExpressionStmt(value), nil
	}
	if _, err = p.consume(token.TypeSemicolon, ErrUnterminatedStatement); err != nil {
		return nil, err
	}
	return ast.NewExpressionStmt(value), nil
}

// forStatement implements the production:
//
//	forStmt -> "for" "(" ( varDecl | exprStmt | ";" ) expression? ";" expression? ")" statement ;
func (p *Parser) forStatement() (ast.Stmt, error) {
	var (
		initializer ast.Stmt
		condition   ast.Expr
		increment   ast.Expr
		body        ast.Stmt
		err         error
	)

	if _, err = p.consume(token.TypeLeftParen, ErrMissingOpeningParenthesis); err != nil {
		return nil, err
	}

	switch {
	case p.match(token.TypeSemicolon):
		break
	case p.match(token.TypeVar):
		initializer, err = p.varDeclaration()
	default:
		initializer, err = p.expressionStatement()
	}
	if err != nil {
		return nil, err
	}

	if !p.check(token.TypeSemicolon) {
		if condition, err = p.expression(); err != nil {
			return nil, err
		}
	}
	if _, err = p.consume(token.TypeSemicolon, errors.New("expect ';' after loop condition")); err != nil {
		return nil, err
	}

	if !p.check(token.TypeRightParen) {
		if increment, err = p.expression(); err != nil {
			return nil, err
		}
	}
	if _, err = p.consume(token.TypeRightParen, errors.New("expect ')' after for clauses")); err != nil {
		return nil, err
	}

	if body, err = p.statement(); err != nil {
		return nil, err
	}

	if increment != nil {
		body = ast.NewBlockStmt([]ast.Stmt{body, ast.NewExpressionStmt(increment)})
	}

	if condition == nil {
		condition = ast.NewLiteralExpr(loxtype.Boolean(true))
	}
	body = ast.NewWhileStmt(condition, body)

	if initializer != nil {
		body = ast.NewBlockStmt([]ast.Stmt{initializer, body})
	}

	return body, nil
}

// ifStatement implements the production:
//
//	ifStmt -> "if" "(" expression ")" statement ( "else" statement )? ;
func (p *Parser) ifStatement() (ast.Stmt, error) {
	if _, err := p.consume(token.TypeLeftParen, ErrMissingOpeningParenthesis); err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	if _, err = p.consume(token.TypeRightParen, ErrUnterminatedExpression); err != nil {
		return nil, err
	}

	var thenBranch, elseBranch ast.Stmt
	if thenBranch, err = p.statement(); err != nil {
		return nil, err
	}
	if p.match(token.TypeElse) {
		if elseBranch, err = p.statement(); err != nil {
			return nil, err
		}
	}

	return ast.NewIfStmt(condition, thenBranch, elseBranch), nil
}

// printStatement implements the production:
//
//	printStmt -> "print" expression ";" ;
func (p *Parser) printStatement() (ast.Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err = p.consume(token.TypeSemicolon, ErrUnterminatedStatement); err != nil {
		return nil, err
	}
	return ast.NewPrintStmt(value), nil
}

// returnStatement implements the production:
//
//	returnStmt -> "return" expression? ";" ;
func (p *Parser) returnStatement() (ast.Stmt, error) {
	var (
		keyword *token.Token
		value   ast.Expr
		err     error
	)

	keyword = p.previous()

	if !p.check(token.TypeSemicolon) {
		if value, err = p.expression(); err != nil {
			return nil, err
		}
	}

	if _, err = p.consume(token.TypeSemicolon, errors.New("expect ';' after semicolon")); err != nil {
		return nil, err
	}

	return ast.NewReturnStmt(keyword, value), nil
}

// whileStatement implements the production:
//
//	whileStmt -> "while" "(" expression ")" statement ;
func (p *Parser) whileStatement() (ast.Stmt, error) {
	var (
		condition ast.Expr
		body      ast.Stmt
		err       error
	)

	if _, err = p.consume(token.TypeLeftParen, ErrMissingOpeningParenthesis); err != nil {
		return nil, err
	}
	if condition, err = p.expression(); err != nil {
		return nil, err
	}
	if _, err = p.consume(token.TypeRightParen, ErrUnterminatedExpression); err != nil {
		return nil, err
	}
	if body, err = p.statement(); err != nil {
		return nil, err
	}

	return ast.NewWhileStmt(condition, body), nil
}

// block implements the production:
//
//	block -> "{" declaration* "}" ;
func (p *Parser) block() ([]ast.Stmt, error) {
	var stmts []ast.Stmt

	for !p.check(token.TypeRightBrace) && !p.isAtEnd() {
		decl, err := p.declaration()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, decl)
	}

	if _, err := p.consume(token.TypeRightBrace, ErrUnterminatedBlock); err != nil {
		return nil, err
	}

	return stmts, nil
}
