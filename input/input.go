package input

// KeyCommand represents the keys using iota for enumeration
type KeyCommand int

//go:generate stringer -type=KeyCommand
const (
	//control keys
	Forward KeyCommand = iota
	Backward
	Right
	Left
	//periferical keys
	Headlights
	Turbo
)

//go:generate stringer -type=ValueCommand
type ValueCommand int

const (
	Release ValueCommand = 0
	Press   ValueCommand = 1
	Hold    ValueCommand = 2
)

//go:generate stringer -type=Command
type Command struct {
	Key   KeyCommand
	Value ValueCommand
}

type Input interface {
	Listen() *chan Command
}
