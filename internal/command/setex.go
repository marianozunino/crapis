package command

import (
	"github.com/marianozunino/crapis/internal/resp"
	"github.com/rs/zerolog/log"
	"strconv"
)

func (e *Executor) setex(args []resp.Value) resp.Value {
	if len(args) != 3 {
		return resp.NewError("wrong number of arguments for 'setex' command")
	}

	key := args[0].BulkVal
	ttlStr := args[1].BulkVal
	val := args[2].BulkVal

	// validate ttl
	ttl, err := strconv.ParseInt(*ttlStr, 10, 32)
	if err != nil {
		return resp.NewError("value is not an integer or out of range")
	}

	if ttl <= 0 {
		return resp.NewError("invalid expire time in 'setex' command")
	}

	log.Debug().Msgf("setex command [%s] = %s ttl=%s", *key, *val, *ttlStr)

	e.db.StoreValueWithTTL(*key, *val, ttl)

	return resp.Ok()
}
