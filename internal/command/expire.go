package command

import (
	"strconv"

	"github.com/marianozunino/crapis/internal/resp"
)

func (e *executor) expire(args []resp.Value) resp.Value {
	if len(args) < 2 {
		return resp.NewError("wrong number of arguments for 'expire' command")
	}

	key := args[0].BulkVal
	ttlStr := args[1].BulkVal

	// validate ttl
	ttl, err := strconv.ParseInt(*ttlStr, 10, 32)
	if err != nil {
		return resp.NewError("value is not an integer or out of range")
	}

	if ttl <= 0 {
		return resp.NewError("invalid expire time in 'expire' command")
	}

	expiredCount := e.db.Expire(*key, ttl)

	return resp.NewInteger(expiredCount)
}
