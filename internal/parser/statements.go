package parser

import (
	"github.com/matt-hoiland/glox/internal/ast"
	"github.com/matt-hoiland/glox/internal/token"
)

// declaration implements the production:
//
//	declaration -> varDecl
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

	if p.match(token.TypeVar) {
		stmt, err = p.varDeclaration()
		return stmt, err
	}
	stmt, err = p.statement()
	return stmt, err
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
//	           | ifStmt
//	           | printStmt
//	           | block ;
func (p *Parser) statement() (ast.Stmt, error) {
	switch {
	case p.match(token.TypeIf):
		return p.ifStatement()

	case p.match(token.TypePrint):
		return p.printStatement()

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

// expressionStatement implements the production:
//
//	exprStmt -> expression ";" ;
func (p *Parser) expressionStatement() (ast.Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err = p.consume(token.TypeSemicolon, ErrUnterminatedStatement); err != nil {
		return nil, err
	}
	return ast.NewExpressionStmt(value), nil
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
