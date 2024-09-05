package command

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidCommand = errors.New("invalid command")

type CommandType string

const (
	PING   CommandType = "PING"
	GET    CommandType = "GET"
	SET    CommandType = "SET"
	SETEX  CommandType = "SETEX"
	DEL    CommandType = "DEL"
	EXPIRE CommandType = "EXPIRE"
)

// commandMap maps command strings to their corresponding CommandType values.
var commandMap = map[string]CommandType{
	string(PING):   PING,
	string(GET):    GET,
	string(SET):    SET,
	string(SETEX):  SETEX,
	string(DEL):    DEL,
	string(EXPIRE): EXPIRE,
}

// ParseCommand parses a string into a CommandType. Returns an error if the command is invalid.
func ParseCommand(in string) (CommandType, error) {
	in = strings.ToUpper(in)
	if cmd, ok := commandMap[in]; ok {
		return cmd, nil
	}
	return "", fmt.Errorf("%q is not a valid command: %w", in, ErrInvalidCommand)
}
