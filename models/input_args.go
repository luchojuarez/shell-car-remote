package models

import (
	"errors"
	"strconv"
)

type InputArgs struct {
	Players int
}

// Expect input parameters like os.Args
func NewInputArgs(args []string) (InputArgs, error) {
	args = args[1:]
	inputs := InputArgs{}
	if len(args)%2 == 1 {
		return inputs, InputError
	}
	for i := 0; i < len(args); i = i + 2 {
		switch args[i] {
		case "-p", "--players":
			n, err := strconv.Atoi(args[i+1])
			if err != nil {
				return inputs, errors.New("players must be a int " + err.Error())
			}
			inputs.Players = n
		}
	}
	return inputs, nil
}
