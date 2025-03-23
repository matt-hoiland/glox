package literal

type Nil struct{}

func (n Nil) String() string {
	return "nil"
}
