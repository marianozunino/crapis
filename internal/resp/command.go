package resp

import (
	"github.com/marianozunino/crapis/internal/resp/command"
	"github.com/rs/zerolog/log"
)

type handlerFunc = func([]Value) Value

var Handlers map[command.CommandType]handlerFunc
var store = newStore()

// init registers all the handlers
// TODO: Do I really want to use init?
func init() {
	Handlers = make(map[command.CommandType]handlerFunc)
	Handlers[command.PING] = ping
	Handlers[command.SET] = set
	Handlers[command.GET] = get
}

func ping(args []Value) Value {
	if len(args) > 1 {
		return Value{
			kind:   ERROR,
			strVal: "wrong number of arguments for 'ping' command",
		}
	}

	if len(args) == 1 {
		return Value{
			kind:   STRING,
			strVal: *args[0].bulkVal,
		}
	}

	return Value{
		kind:   STRING,
		strVal: "PONG",
	}
}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{
			kind:   ERROR,
			strVal: "wrong number of arguments for 'set' command",
		}
	}

	key := args[0].bulkVal
	val := args[1].bulkVal

	log.Debug().Msgf("set command [%s] = %s", *key, *val)

	store.StoreValue(*key, *val)

	return Value{
		kind:   STRING,
		strVal: "OK",
	}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{
			kind:   ERROR,
			strVal: "wrong number of arguments for 'get' command",
		}
	}

	key := args[0].bulkVal

	log.Debug().Msgf("get command [%s]", *key)

	val := store.ReadVal(*key)

	if val == nil {
		return Value{
			kind:   NULL,
			strVal: "",
		}
	}

	return Value{
		kind:   STRING,
		strVal: *val,
	}

}
