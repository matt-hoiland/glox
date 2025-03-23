package scanner

import (
	"errors"

	ierrors "github.com/matt-hoiland/glox/internal/errors"
	"github.com/matt-hoiland/glox/internal/literal"
	"github.com/matt-hoiland/glox/internal/runes"
	"github.com/matt-hoiland/glox/internal/scanner/tokentype"
)

var (
	ErrUnexpectedRune     = errors.New("unexpected rune")
	ErrUnterminatedString = errors.New("unterminated string")
)

type Scanner struct {
	source  []runes.Rune
	tokens  []*Token
	start   int
	current int
	line    int
}

func New(source string) *Scanner {
	return &Scanner{
		source: []runes.Rune(source),
	}
}

func (s *Scanner) ScanTokens() ([]*Token, error) {
	for !s.isAtEnd() {
		// We are at the beginning of the next lexeme.
		s.start = s.current
		if err := s.scanToken(); err != nil {
			return s.tokens, err
		}
	}

	s.tokens = append(s.tokens, &Token{
		Type:    tokentype.EOF,
		Lexeme:  "",
		Literal: nil,
		Line:    s.line,
	})
	return s.tokens, nil
}

func (s *Scanner) advance() runes.Rune {
	r := s.source[s.current]
	s.current++
	return r
}

func (s *Scanner) emitIdentifier() *Token {
	for s.peek().IsAlphaNumeric() {
		s.advance()
	}

	text := string(s.source[s.start:s.current])
	tokenType, ok := keywords[text]
	if !ok {
		tokenType = tokentype.Identifier
	}
	return s.emitToken(tokenType)
}

func (s *Scanner) emitNumber() *Token {
	for s.peek().IsDigit() {
		s.advance()
	}

	// Look for the fractional part.
	if s.peek() == '.' && s.peekNext().IsDigit() {
		// Consume the '.'
		s.advance()

		for s.peek().IsDigit() {
			s.advance()
		}
	}
	number := literal.ParseNumber(s.source[s.start:s.current])

	return s.emitToken(tokentype.Number, number)
}

func (s *Scanner) emitString() (*Token, error) {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		return nil, &ierrors.Error{Line: s.line, Err: ErrUnterminatedString}
	}

	s.advance()
	value := literal.String(s.source[s.start+1 : s.current-1])
	return s.emitToken(tokentype.String, value), nil
}

func (s *Scanner) emitToken(tokenType tokentype.TokenType, literal ...Literal) *Token {
	token := &Token{
		Type:   tokenType,
		Lexeme: string(s.source[s.start:s.current]),
		Line:   s.line,
	}
	if len(literal) > 0 {
		token.Literal = literal[0]
	}
	return token
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) match(expected runes.Rune) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) matchTernary(expected runes.Rune, t, f tokentype.TokenType) tokentype.TokenType {
	if s.match(expected) {
		return t
	}
	return f
}

func (s *Scanner) peek() runes.Rune {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() runes.Rune {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

func (s *Scanner) scanToken() error {
	var (
		r     runes.Rune
		token *Token
		err   error
	)

	switch r = s.advance(); r {
	case '(':
		token = s.emitToken(tokentype.LeftParen)
	case ')':
		token = s.emitToken(tokentype.RightParen)
	case '{':
		token = s.emitToken(tokentype.LeftBrace)
	case '}':
		token = s.emitToken(tokentype.RightBrace)
	case ',':
		token = s.emitToken(tokentype.Comma)
	case '.':
		token = s.emitToken(tokentype.Dot)
	case '-':
		token = s.emitToken(tokentype.Minus)
	case '+':
		token = s.emitToken(tokentype.Plus)
	case ';':
		token = s.emitToken(tokentype.Semicolon)
	case '*':
		token = s.emitToken(tokentype.Star)
	case '!':
		token = s.emitToken(s.matchTernary('=', tokentype.BangEqual, tokentype.Bang))
	case '=':
		token = s.emitToken(s.matchTernary('=', tokentype.EqualEqual, tokentype.Equal))
	case '<':
		token = s.emitToken(s.matchTernary('=', tokentype.LessEqual, tokentype.Less))
	case '>':
		token = s.emitToken(s.matchTernary('=', tokentype.GreaterEqual, tokentype.Greater))
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			token = s.emitToken(tokentype.Slash)
		}
	case ' ', '\r', '\t':
		// Ignore whitespace.
		break
	case '\n':
		s.line++
	case '"':
		if token, err = s.emitString(); err != nil {
			return err
		}
	default:
		if r.IsDigit() {
			token = s.emitNumber()
		} else if r.IsAlpha() {
			token = s.emitIdentifier()
		} else {
			return &ierrors.Error{Line: s.line, Err: ErrUnexpectedRune}
		}
	}
	if token != nil {
		s.tokens = append(s.tokens, token)
	}
	return nil
}
