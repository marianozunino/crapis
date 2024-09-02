package command

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidCommand = errors.New("invalid command")

type CommandType string

const (
	PING  CommandType = "PING"
	GET   CommandType = "GET"
	SET   CommandType = "SET"
	SETEX CommandType = "SETEX"
	DEL   CommandType = "DEL"
)

func ParseCommand(in string) (CommandType, error) {
	in = strings.ToUpper(in)
	switch in {
	case string(PING):
		return PING, nil
	case string(SET):
		return SET, nil
	case string(SETEX):
		return SETEX, nil
	case string(GET):
		return GET, nil
	case string(DEL):
		return DEL, nil
	}
	return "", fmt.Errorf("%q is not a valid command: %w", in, ErrInvalidCommand)
}
