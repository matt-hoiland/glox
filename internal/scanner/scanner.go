package scanner

import (
	"errors"

	ierrors "github.com/matt-hoiland/glox/internal/errors"
	"github.com/matt-hoiland/glox/internal/loxtype"
	"github.com/matt-hoiland/glox/internal/runes"
	"github.com/matt-hoiland/glox/internal/token"
)

var (
	ErrUnexpectedRune     = errors.New("unexpected rune")
	ErrUnterminatedString = errors.New("unterminated string")
)

type Scanner struct {
	source  []runes.Rune
	tokens  []*token.Token
	start   int
	current int
	line    int
}

func New(source string) *Scanner {
	return &Scanner{
		source: []runes.Rune(source),
		line:   1,
	}
}

func (s *Scanner) ScanTokens() ([]*token.Token, error) {
	for !s.isAtEnd() {
		// We are at the beginning of the next lexeme.
		s.start = s.current
		if err := s.scanToken(); err != nil {
			return s.tokens, err
		}
	}

	s.tokens = append(s.tokens, &token.Token{
		Type:    token.TypeEOF,
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

func (s *Scanner) emitIdentifier() *token.Token {
	for s.peek().IsAlphaNumeric() {
		s.advance()
	}

	text := string(s.source[s.start:s.current])
	tokenType, ok := keywords[text]
	if !ok {
		tokenType = token.TypeIdentifier
	}
	return s.emitToken(tokenType)
}

func (s *Scanner) emitNumber() *token.Token {
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
	number := loxtype.ParseNumber(s.source[s.start:s.current])

	return s.emitToken(token.TypeNumber, number)
}

func (s *Scanner) emitString() (*token.Token, error) {
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
	value := loxtype.String(s.source[s.start+1 : s.current-1])
	return s.emitToken(token.TypeString, value), nil
}

func (s *Scanner) emitToken(tokenType token.Type, literal ...token.Literal) *token.Token {
	token := &token.Token{
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

func (s *Scanner) matchTernary(expected runes.Rune, t, f token.Type) token.Type {
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
		r   runes.Rune
		tok *token.Token
		err error
	)

	switch r = s.advance(); r {
	case '(':
		tok = s.emitToken(token.TypeLeftParen)
	case ')':
		tok = s.emitToken(token.TypeRightParen)
	case '{':
		tok = s.emitToken(token.TypeLeftBrace)
	case '}':
		tok = s.emitToken(token.TypeRightBrace)
	case ',':
		tok = s.emitToken(token.TypeComma)
	case '.':
		tok = s.emitToken(token.TypeDot)
	case '-':
		tok = s.emitToken(token.TypeMinus)
	case '+':
		tok = s.emitToken(token.TypePlus)
	case ';':
		tok = s.emitToken(token.TypeSemicolon)
	case '*':
		tok = s.emitToken(token.TypeStar)
	case '!':
		tok = s.emitToken(s.matchTernary('=', token.TypeBangEqual, token.TypeBang))
	case '=':
		tok = s.emitToken(s.matchTernary('=', token.TypeEqualEqual, token.TypeEqual))
	case '<':
		tok = s.emitToken(s.matchTernary('=', token.TypeLessEqual, token.TypeLess))
	case '>':
		tok = s.emitToken(s.matchTernary('=', token.TypeGreaterEqual, token.TypeGreater))
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			tok = s.emitToken(token.TypeSlash)
		}
	case ' ', '\r', '\t':
		// Ignore whitespace.
		break
	case '\n':
		s.line++
	case '"':
		if tok, err = s.emitString(); err != nil {
			return err
		}
	default:
		if r.IsDigit() {
			tok = s.emitNumber()
		} else if r.IsAlpha() {
			tok = s.emitIdentifier()
		} else {
			return &ierrors.Error{Line: s.line, Err: ErrUnexpectedRune}
		}
	}
	if tok != nil {
		s.tokens = append(s.tokens, tok)
	}
	return nil
}
