package token

import (
	"fmt"

	"github.com/matt-hoiland/glox/internal/loxtype"
)

type Token struct {
	Type    Type
	Lexeme  string
	Literal loxtype.Type
	Line    int
}

func NewToken(tokenType Type, lexeme string, literal loxtype.Type, line int) *Token {
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
