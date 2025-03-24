package token

//go:generate stringer -type=Type
type Type int

const (
	// Single-character tokens.
	TypeLeftParen Type = iota
	TypeRightParen
	TypeLeftBrace
	TypeRightBrace
	TypeComma
	TypeDot
	TypeMinus
	TypePlus
	TypeSemicolon
	TypeSlash
	TypeStar

	// One or two character tokens.
	TypeBang
	TypeBangEqual
	TypeEqual
	TypeEqualEqual
	TypeGreater
	TypeGreaterEqual
	TypeLess
	TypeLessEqual

	// Literals.
	TypeIdentifier
	TypeString
	TypeNumber

	// Keywords.
	TypeAnd
	TypeClass
	TypeElse
	TypeFalse
	TypeFun
	TypeFor
	TypeIf
	TypeNil
	TypeOr
	TypePrint
	TypeReturn
	TypeSuper
	TypeThis
	TypeTrue
	TypeVar
	TypeWhile

	TypeEOF
)
