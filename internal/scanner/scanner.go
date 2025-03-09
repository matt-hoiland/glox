package scanner

import (
	"github.com/matt-hoiland/glox/internal/scanner/literal"
	"github.com/matt-hoiland/glox/internal/scanner/tokentype"
)

type Scanner struct {
	source  []rune
	tokens  []*Token
	start   int
	current int
	line    int
}

func New(source string) *Scanner {
	return &Scanner{
		source: []rune(source),
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

func (s *Scanner) addToken(tokenType tokentype.TokenType, literal ...Literal) {
	token := &Token{
		Type:   tokenType,
		Lexeme: string(s.source[s.start:s.current]),
		Line:   s.line,
	}
	if len(literal) > 0 {
		token.Literal = literal[0]
	}
	s.tokens = append(s.tokens, token)
}

func (s *Scanner) advance() rune {
	r := s.source[s.current]
	s.current++
	return r
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := string(s.source[s.start:s.current])
	tokenType, ok := keywords[text]
	if !ok {
		tokenType = tokentype.Identifier
	}
	s.addToken(tokenType)
}

func (s *Scanner) isAlpha(r rune) bool {
	return (r >= 'A' && r <= 'Z') ||
		(r >= 'a' && r <= 'z') ||
		r == '_'
}

func (s Scanner) isAlphaNumeric(r rune) bool {
	return s.isAlpha(r) || s.isDigit(r)
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) matchTernary(expected rune, t, f tokentype.TokenType) tokentype.TokenType {
	if s.match(expected) {
		return t
	}
	return f
}

func (s *Scanner) number() error {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	// Look for the fractional part.
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		// Consume the '.'
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}
	}
	number, err := literal.ParseNumber(s.source[s.start:s.current])
	if err != nil {
		return &Error{Line: s.line, Err: err}
	}
	s.addToken(tokentype.Number, number)
	return nil
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

func (s *Scanner) scanToken() error {
	r := s.advance()
	switch r {
	case '(':
		s.addToken(tokentype.LeftParen)
	case ')':
		s.addToken(tokentype.RightParen)
	case '{':
		s.addToken(tokentype.LeftBrace)
	case '}':
		s.addToken(tokentype.RightBrace)
	case ',':
		s.addToken(tokentype.Comma)
	case '.':
		s.addToken(tokentype.Dot)
	case '-':
		s.addToken(tokentype.Minus)
	case '+':
		s.addToken(tokentype.Plus)
	case ';':
		s.addToken(tokentype.Semicolon)
	case '*':
		s.addToken(tokentype.Star)
	case '!':
		s.addToken(s.matchTernary('=', tokentype.BangEqual, tokentype.Bang))
	case '=':
		s.addToken(s.matchTernary('=', tokentype.EqualEqual, tokentype.Equal))
	case '<':
		s.addToken(s.matchTernary('=', tokentype.LessEqual, tokentype.Less))
	case '>':
		s.addToken(s.matchTernary('=', tokentype.GreaterEqual, tokentype.Greater))
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(tokentype.Slash)
		}
	case ' ', '\r', '\t':
		// Ignore whitespace.
		break
	case '\n':
		s.line++
	case '"':
		if err := s.string(); err != nil {
			return err
		}
	default:
		if s.isDigit(r) {
			if err := s.number(); err != nil {
				return err
			}
		} else if s.isAlpha(r) {
			s.identifier()
		} else {
			return &Error{Line: s.line, Err: ErrUnexpectedRune}
		}
	}

	return nil
}

func (s *Scanner) string() error {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		return &Error{Line: s.line, Err: ErrUnterminatedString}
	}

	s.advance()
	value := literal.String(s.source[s.start+1 : s.current-1])
	s.addToken(tokentype.String, value)
	return nil
}
