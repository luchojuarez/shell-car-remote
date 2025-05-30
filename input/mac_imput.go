//go:build darwin

package input

func GetAllDSControllers() (out []any, e error) {
	panic("Not implemented")
	return nil, e
}

func NewDS4Input(inputController any) Input {
	panic("Not implemented")
	return nil
}
func NewKeyboardInput() Input {
	return nil
}
