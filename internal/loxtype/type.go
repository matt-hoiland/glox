package loxtype

import (
	"fmt"
	"strconv"

	"github.com/matt-hoiland/glox/internal/runes"
)

type Type interface {
	fmt.Stringer
}

type Boolean bool

var _ Type = Boolean(false)

func (b Boolean) String() string {
	if b {
		return "true"
	}
	return "false"
}

type Nil struct{}

var _ Type = Nil{}

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

func (n Number) String() string {
	return strconv.FormatFloat(float64(n), 'f', -1, 64)
}

type String string

var _ Type = String("")

func (s String) String() string {
	return string(s)
}
