package loxtype

import (
	"fmt"
	"strconv"

	"github.com/matt-hoiland/glox/internal/runes"
)

type Type interface {
	fmt.Stringer
	Equals(Type) Boolean
	IsTruthy() Boolean
}

type Boolean bool

var (
	_ Type = Boolean(false)
)

func (b Boolean) Equals(other Type) Boolean {
	bo, ok := other.(Boolean)
	if !ok {
		return false
	}
	return b == bo
}

func (b Boolean) IsTruthy() Boolean {
	return b
}

func (b Boolean) Negate() Type {
	return !b
}

func (b Boolean) String() string {
	if b {
		return "true"
	}
	return "false"
}

type Nil struct{}

var _ Type = Nil{}

func (n Nil) Equals(other Type) Boolean {
	if _, ok := other.(Nil); !ok {
		return false
	}
	return true
}

func (Nil) IsTruthy() Boolean {
	return false
}

func (n Nil) String() string {
	return "nil"
}

type Number float64

var _ Type = Number(0)

func ParseNumber(text []runes.Rune) Number {
	f64, err := strconv.ParseFloat(string(text), 64)
	if err != nil {
		return Number(0)
	}
	return Number(f64)
}

func (n Number) Add(right Number) Number {
	return n + right
}

func (n Number) Divide(right Number) Number {
	return n / right
}

func (n Number) Equals(other Type) Boolean {
	no, ok := other.(Number)
	if !ok {
		return false
	}
	return n == no
}

func (n Number) Greater(right Number) Boolean {
	return n > right
}

func (n Number) GreaterEqual(right Number) Boolean {
	return n >= right
}

func (Number) IsTruthy() Boolean {
	return true
}

func (n Number) Less(right Number) Boolean {
	return n < right
}

func (n Number) LessEqual(right Number) Boolean {
	return n <= right
}

func (n Number) Multiply(right Number) Number {
	return n * right
}

func (n Number) Negate() Type {
	return -n
}

func (n Number) String() string {
	return strconv.FormatFloat(float64(n), 'f', -1, 64)
}

func (n Number) Subtract(right Number) Number {
	return n - right
}

type String string

var _ Type = String("")

func (s String) Add(right String) String {
	return s + right
}

func (s String) Equals(other Type) Boolean {
	ns, ok := other.(String)
	if !ok {
		return false
	}
	return s == ns
}

func (String) IsTruthy() Boolean {
	return true
}

func (s String) String() string {
	return string(s)
}
