package input

import (
	"errors"
)

func Parse(args []string) (string, error) {
	if len(args) < 2 {
		return "", errors.New("no command provided")
	}

	return args[1], nil
}
