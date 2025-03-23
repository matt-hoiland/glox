package runes

type Rune rune

func (r Rune) IsAlpha() bool {
	return (r >= 'A' && r <= 'Z') ||
		(r >= 'a' && r <= 'z') ||
		r == '_'
}

func (r Rune) IsAlphaNumeric() bool {
	return r.IsAlpha() || r.IsDigit()
}

func (r Rune) IsDigit() bool {
	return r >= '0' && r <= '9'
}
