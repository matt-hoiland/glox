package scanner

import "github.com/matt-hoiland/glox/internal/token"

func keywords(key string) (token.Type, bool) {
	switch key {
	case "and":
		return token.TypeAnd, true
	case "class":
		return token.TypeClass, true
	case "else":
		return token.TypeElse, true
	case "false":
		return token.TypeFalse, true
	case "for":
		return token.TypeFor, true
	case "fun":
		return token.TypeFun, true
	case "if":
		return token.TypeIf, true
	case "nil":
		return token.TypeNil, true
	case "or":
		return token.TypeOr, true
	case "print":
		return token.TypePrint, true
	case "return":
		return token.TypeReturn, true
	case "super":
		return token.TypeSuper, true
	case "this":
		return token.TypeThis, true
	case "true":
		return token.TypeTrue, true
	case "var":
		return token.TypeVar, true
	case "while":
		return token.TypeWhile, true
	default:
		return token.Type(0), false
	}
}
