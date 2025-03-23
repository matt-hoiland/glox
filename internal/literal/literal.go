package literal

import (
	"fmt"
	"strconv"

	"github.com/matt-hoiland/glox/internal/runes"
)

type Boolean bool

var _ fmt.Stringer = Boolean(false)

func (b Boolean) String() string {
	if b {
		return "true"
	}
	return "false"
}

type Nil struct{}

var _ fmt.Stringer = Nil{}

func (n Nil) String() string {
	return "nil"
}

type Number float64

var _ fmt.Stringer = Number(0)

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

var _ fmt.Stringer = String("")

func (s String) String() string {
	return string(s)
}
