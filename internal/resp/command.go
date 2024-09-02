package resp

import (
	"github.com/marianozunino/crapis/internal/resp/command"
)

type handlerFunc = func([]Value) Value

var Handlers map[command.CommandType]handlerFunc

// init registers all the handlers
// TODO: Do I really want to use init?
func init() {
	Handlers = make(map[command.CommandType]handlerFunc)
	Handlers[command.PING] = ping
}

func ping(args []Value) Value {
	return Value{
		kind:   STRING,
		strVal: "PONG",
	}
}
