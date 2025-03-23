package literal

type Boolean bool

func (b Boolean) String() string {
	if b {
		return "true"
	}
	return "false"
}
