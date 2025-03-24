package token

import (
	"fmt"
)

type Literal interface {
	fmt.Stringer
}

type Token struct {
	Type    Type
	Lexeme  string
	Literal Literal
	Line    int
}

func NewToken(tokenType Type, lexeme string, literal Literal, line int) *Token {
	return &Token{
		Type:    tokenType,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

func (t *Token) String() string {
	if t.Literal != nil {
		return fmt.Sprintf("%s %s '%s'", t.Type, t.Lexeme, t.Literal)
	}
	return fmt.Sprintf("%s %s", t.Type, t.Lexeme)
}
