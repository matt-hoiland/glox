package scanner

import (
	"fmt"

	"github.com/matt-hoiland/glox/internal/scanner/tokentype"
)

type Literal interface {
	fmt.Stringer
}

type Token struct {
	Type    tokentype.TokenType
	Lexeme  string
	Literal Literal
	Line    int
}

func (t *Token) String() string {
	if t.Literal != nil {
		return fmt.Sprintf("%s %s '%s'", t.Type, t.Lexeme, t.Literal)
	}
	return fmt.Sprintf("%s %s", t.Type, t.Lexeme)
}
