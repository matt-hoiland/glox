package literal

import (
	"strconv"

	"github.com/matt-hoiland/glox/internal/scanner/runes"
)

type Number float64

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
