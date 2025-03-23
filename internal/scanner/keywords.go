package scanner

import "github.com/matt-hoiland/glox/internal/token/tokentype"

var keywords = map[string]tokentype.TokenType{
	"and":    tokentype.And,
	"class":  tokentype.Class,
	"else":   tokentype.Else,
	"false":  tokentype.False,
	"for":    tokentype.For,
	"fun":    tokentype.Fun,
	"if":     tokentype.If,
	"nil":    tokentype.Nil,
	"or":     tokentype.Or,
	"print":  tokentype.Print,
	"return": tokentype.Return,
	"super":  tokentype.Super,
	"this":   tokentype.This,
	"true":   tokentype.True,
	"var":    tokentype.Var,
	"while":  tokentype.While,
}
