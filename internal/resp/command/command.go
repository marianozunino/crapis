package command

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidCommand = errors.New("invalid command")

type CommandType string

const (
	PING CommandType = "PING"
)

func ParseCommand(in string) (CommandType, error) {
	in = strings.ToUpper(in)
	switch in {
	case string(PING):
		return PING, nil
	}
	return "", fmt.Errorf("%q is not a valid command: %w", in, ErrInvalidCommand)
}
