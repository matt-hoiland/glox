package scanner

import "github.com/matt-hoiland/glox/internal/token"

var keywords = map[string]token.Type{
	"and":    token.TypeAnd,
	"class":  token.TypeClass,
	"else":   token.TypeElse,
	"false":  token.TypeFalse,
	"for":    token.TypeFor,
	"fun":    token.TypeFun,
	"if":     token.TypeIf,
	"nil":    token.TypeNil,
	"or":     token.TypeOr,
	"print":  token.TypePrint,
	"return": token.TypeReturn,
	"super":  token.TypeSuper,
	"this":   token.TypeThis,
	"true":   token.TypeTrue,
	"var":    token.TypeVar,
	"while":  token.TypeWhile,
}
