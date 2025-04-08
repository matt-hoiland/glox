package parser

import (
	"errors"

	"github.com/matt-hoiland/glox/internal/ast"
	ierrors "github.com/matt-hoiland/glox/internal/errors"
	"github.com/matt-hoiland/glox/internal/loxtype"
	"github.com/matt-hoiland/glox/internal/token"
)

var (
	ErrNoVariableName         = errors.New("expect variable name")
	ErrUnimplemented          = errors.New("unimplemented")
	ErrUnterminatedExpression = errors.New("expect ')' after expression")
	ErrUnterminatedStatement  = errors.New("expect ';' after expression")
	ErrUnterminatedBlock      = errors.New("expect '}' after block")
)

type Parser struct {
	Tokens  []*token.Token
	Current int
}

func New(tokens []*token.Token) *Parser {
	return &Parser{
		Tokens:  tokens,
		Current: 0,
	}
}

func (p *Parser) Parse() ([]ast.Stmt, error) {
	var statements []ast.Stmt
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}
	return statements, nil
}

//----------------------------------------------------------------------------
// Token Stream management methods.
//----------------------------------------------------------------------------

// advance consumes the current token and returns it.
// This is similar to how [scanner.Scanner]'s corresponding method crawled through characters.
func (p *Parser) advance() *token.Token {
	if !p.isAtEnd() {
		p.Current++
	}
	return p.previous()
}

// check returns true if the current token is of the given type.
// Unlike [Parser.match], it never consumes the token, it only looks at it.
func (p *Parser) check(tokenType token.Type) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokenType
}

func (p *Parser) consume(tokenType token.Type, err error) (*token.Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}

	return nil, ierrors.New(p.peek(), err)
}

// isAtEnd checks if we’ve run out of tokens to parse.
func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.TypeEOF
}

// match checks to see if the current token has string of the given types.
// If so, it consumes the token and returns true.
// Otherwise, it returns false and leaves the current token alone.
func (p *Parser) match(tokenTypes ...token.Type) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

// peek returns the current token we have yet to consume.
func (p *Parser) peek() *token.Token {
	return p.Tokens[p.Current]
}

// previous returns the most recently consumed token.
func (p *Parser) previous() *token.Token {
	return p.Tokens[p.Current-1]
}

//nolint:unused // This will be used in a future chapter.
func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == token.TypeSemicolon {
			return
		}

		switch p.peek().Type {
		case token.TypeClass,
			token.TypeFun,
			token.TypeVar,
			token.TypeFor,
			token.TypeIf,
			token.TypeWhile,
			token.TypePrint,
			token.TypeReturn:
			return
		}

		p.advance()
	}
}

//----------------------------------------------------------------------------
// Grammar production methods.
//----------------------------------------------------------------------------

// declaration implements the production:
//
//	declaration -> varDecl
//	             | statement ;
func (p *Parser) declaration() (stmt ast.Stmt, err error) {
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

	if _, err := p.consume(token.TypeSemicolon, ErrUnterminatedStatement); err != nil {
		return nil, err
	}

	return ast.NewVarStmt(name, initializer), nil
}

// statement implements the production:
//
//	statement -> exprStmt
//	           | printStmt
//	           | block ;
func (p *Parser) statement() (ast.Stmt, error) {
	if p.match(token.TypePrint) {
		return p.printStatement()
	}
	if p.match(token.TypeLeftBrace) {
		stmts, err := p.block()
		if err != nil {
			return nil, err
		}
		block := ast.NewBlockStmt(stmts)
		return block, nil
	}
	return p.expressionStatement()
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

// expression implements the production:
//
//	expression -> assignment ;
func (p *Parser) expression() (ast.Expr, error) {
	return p.assignment()
}

// assignment implements the production:
//
//	assignment -> IDENTIFIER "=" assignment
//	            | equality ;
func (p *Parser) assignment() (ast.Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	if p.match(token.TypeEqual) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
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
		right, err := p.comparison()
		if err != nil {
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
		right, err := p.term()
		if err != nil {
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
		right, err := p.factor()
		if err != nil {
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
		right, err := p.factor()
		if err != nil {
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
		if _, err := p.consume(token.TypeRightParen, ErrUnterminatedExpression); err != nil {
			return nil, err
		}
		return ast.NewGroupingExpr(expression), nil
	}

	return nil, ErrUnimplemented
}
