package parser

import (
	"errors"

	ierrors "github.com/matt-hoiland/glox/internal/errors"
	"github.com/matt-hoiland/glox/internal/expr"
	"github.com/matt-hoiland/glox/internal/literal"
	"github.com/matt-hoiland/glox/internal/token"
)

var (
	ErrUnterminatedExpression = errors.New("expect ')' after expression")
	ErrUnimplemented          = errors.New("unimplemented")
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

func (p *Parser) Parse() (expr.Expr[string], error) {
	return p.expression()
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

	return nil, p.error(p.peek(), err)
}

func (p *Parser) error(token *token.Token, err error) error {
	return &ierrors.Error{
		Line: token.Line,
		Err:  err,
	}
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

// expression implements the production:
//
//	expression -> equality ;
func (p *Parser) expression() (expr.Expr[string], error) {
	return p.equality()
}

// equality implements the production:
//
//	equality -> comparison ( ( "!=" | "==" ) comparison )* ;
func (p *Parser) equality() (expr.Expr[string], error) {
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
		left = expr.NewBinary(left, operator, right)
	}

	return left, nil
}

// comparison implements the production:
//
//	comparison -> term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
func (p *Parser) comparison() (expr.Expr[string], error) {
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
		left = expr.NewBinary(left, operator, right)
	}

	return left, nil
}

// term implements the production:
//
//	term -> factor ( ( "-" | "+" ) factor )* ;
func (p *Parser) term() (expr.Expr[string], error) {
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
		left = expr.NewBinary(left, operator, right)
	}

	return left, nil
}

// factor implements the production:
//
//	factor -> unary ( ( "/" | "*" ) unary )* ;
func (p *Parser) factor() (expr.Expr[string], error) {
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
		left = expr.NewBinary(left, operator, right)
	}

	return left, nil
}

// unary implements the production:
//
//	 unary -> ( "!" | "-" ) unary
//		    | primary ;
func (p *Parser) unary() (expr.Expr[string], error) {
	if p.match(token.TypeBang, token.TypeMinus) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return expr.NewUnary(operator, right), nil
	}

	return p.primary()
}

// primary implements the production:
//
//	 primary -> NUMBER | STRING | "true" | "false" | "nil"
//		      | "(" expression ")" ;
func (p *Parser) primary() (expr.Expr[string], error) {
	if p.match(token.TypeFalse) {
		return expr.NewLiteral[string](literal.Boolean(false)), nil
	}
	if p.match(token.TypeTrue) {
		return expr.NewLiteral[string](literal.Boolean(true)), nil
	}
	if p.match(token.TypeNil) {
		return expr.NewLiteral[string](literal.Nil{}), nil
	}
	if p.match(token.TypeNumber, token.TypeString) {
		return expr.NewLiteral[string](p.previous().Literal), nil
	}
	if p.match(token.TypeLeftParen) {
		expression, err := p.expression()
		if err != nil {
			return nil, err
		}
		if _, err := p.consume(token.TypeRightParen, ErrUnterminatedExpression); err != nil {
			return nil, err
		}
		return expr.NewGrouping(expression), nil
	}

	return nil, ErrUnimplemented
}
