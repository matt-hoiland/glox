package parser

import (
	"errors"

	"github.com/matt-hoiland/glox/internal/ast"
	ierrors "github.com/matt-hoiland/glox/internal/errors"
	"github.com/matt-hoiland/glox/internal/token"
)

var (
	ErrNoVariableName            = errors.New("expect variable name")
	ErrUnimplemented             = errors.New("unimplemented")
	ErrMissingOpeningParenthesis = errors.New("expect '(' after 'if', 'while', or 'for'")
	ErrUnterminatedExpression    = errors.New("expect ')' after expression")
	ErrUnterminatedStatement     = errors.New("expect ';' after expression")
	ErrUnterminatedBlock         = errors.New("expect '}' after block")
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
		default:
			break
		}

		p.advance()
	}
}
