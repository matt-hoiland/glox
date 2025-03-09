package literal

import "strconv"

type Number float64

func ParseNumber(text []rune) (Number, error) {
	f64, err := strconv.ParseFloat(string(text), 64)
	return Number(f64), err
}

func (n Number) String() string {
	return strconv.FormatFloat(float64(n), 'f', -1, 64)
}
