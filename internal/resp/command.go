package resp

import (
	"strconv"

	"github.com/marianozunino/crapis/internal/resp/command"
	"github.com/rs/zerolog/log"
)

type handlerFunc = func([]Value) Value

var Handlers map[command.CommandType]handlerFunc
var store = NewStore()

// init registers all the handlers
// TODO: Do I really want to use init?
func init() {
	Handlers = make(map[command.CommandType]handlerFunc)
	Handlers[command.PING] = ping
	Handlers[command.GET] = get
	Handlers[command.SET] = set
	Handlers[command.SETEX] = setex
	Handlers[command.DEL] = del
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

func setex(args []Value) Value {
	if len(args) != 3 {
		return Value{
			kind:   ERROR,
			strVal: "wrong number of arguments for 'setex' command",
		}
	}

	key := args[0].bulkVal
	ttlStr := args[1].bulkVal
	val := args[2].bulkVal

	// validate ttl
	ttl, err := strconv.ParseInt(*ttlStr, 10, 32)
	if err != nil {
		return Value{
			kind:   ERROR,
			strVal: "value is not an integer or out of range",
		}
	}

	if ttl <= 0 {
		return Value{
			kind:   ERROR,
			strVal: "invalid expire time in 'setex' command",
		}
	}

	log.Debug().Msgf("setex command [%s] = %s ttl=%s", *key, *val, *ttlStr)

	store.StoreValueWithTTL(*key, *val, ttl)

	return Value{
		kind:   STRING,
		strVal: "OK",
	}
}

func del(args []Value) Value {
	if len(args) == 0 {
		return Value{
			kind:   ERROR,
			strVal: "wrong number of arguments for 'del' command",
		}
	}
	keys := make([]string, len(args))
	for i, arg := range args {
		keys[i] = *arg.bulkVal
	}
	deletedKeys := store.DeleteKey(keys...)
	return Value{
		kind:   STRING,
		strVal: strconv.Itoa(deletedKeys),
	}
}
